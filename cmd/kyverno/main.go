package main

import (
	"context"
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/nirmata/kyverno/pkg/checker"
	kyvernoclient "github.com/nirmata/kyverno/pkg/client/clientset/versioned"
	kyvernoinformer "github.com/nirmata/kyverno/pkg/client/informers/externalversions"
	"github.com/nirmata/kyverno/pkg/config"
	dclient "github.com/nirmata/kyverno/pkg/dclient"
	event "github.com/nirmata/kyverno/pkg/event"
	"github.com/nirmata/kyverno/pkg/namespace"
	"github.com/nirmata/kyverno/pkg/policy"
	"github.com/nirmata/kyverno/pkg/policystore"
	"github.com/nirmata/kyverno/pkg/policyviolation"
	"github.com/nirmata/kyverno/pkg/signal"
	"github.com/nirmata/kyverno/pkg/utils"
	"github.com/nirmata/kyverno/pkg/version"
	"github.com/nirmata/kyverno/pkg/webhookconfig"
	"github.com/nirmata/kyverno/pkg/webhooks"
	kubeinformers "k8s.io/client-go/informers"
)

var (
	kubeconfig     string
	serverIP       string
	cpu            bool
	memory         bool
	webhookTimeout int
	//TODO: this has been added to backward support command line arguments
	// will be removed in future and the configuration will be set only via configmaps
	filterK8Resources string
)

func main() {
	defer glog.Flush()
	version.PrintVersionInfo()

	// cleanUp Channel
	cleanUp := make(chan struct{})
	//  handle os signals
	stopCh := signal.SetupSignalHandler()
	// CLIENT CONFIG
	clientConfig, err := config.CreateClientConfig(kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v\n", err)
	}

	// KYVENO CRD CLIENT
	// access CRD resources
	//		- Policy
	//		- PolicyViolation
	pclient, err := kyvernoclient.NewForConfig(clientConfig)
	if err != nil {
		glog.Fatalf("Error creating client: %v\n", err)
	}

	// DYNAMIC CLIENT
	// - client for all registered resources
	// - invalidate local cache of registered resource every 10 seconds
	client, err := dclient.NewClient(clientConfig, 10*time.Second, stopCh)
	if err != nil {
		glog.Fatalf("Error creating client: %v\n", err)
	}
	// CRD CHECK
	// - verify if the CRD for Policy & PolicyViolation are available
	if !utils.CRDInstalled(client.DiscoveryClient) {
		glog.Fatalf("Required CRDs unavailable")
	}
	// KUBERNETES CLIENT
	kubeClient, err := utils.NewKubeClient(clientConfig)
	if err != nil {
		glog.Fatalf("Error creating kubernetes client: %v\n", err)
	}

	// TODO(shuting): To be removed for v1.2.0
	utils.CleanupOldCrd(client)

	// KUBERNETES RESOURCES INFORMER
	// watches namespace resource
	// - cache resync time: 10 seconds
	kubeInformer := kubeinformers.NewSharedInformerFactoryWithOptions(
		kubeClient,
		10*time.Second)

	// WERBHOOK REGISTRATION CLIENT
	webhookRegistrationClient := webhookconfig.NewWebhookRegistrationClient(
		clientConfig,
		client,
		serverIP,
		int32(webhookTimeout))

	// Resource Mutating Webhook Watcher
	lastReqTime := checker.NewLastReqTime()
	rWebhookWatcher := webhookconfig.NewResourceWebhookRegister(
		lastReqTime,
		kubeInformer.Admissionregistration().V1beta1().MutatingWebhookConfigurations(),
		webhookRegistrationClient,
	)

	// KYVERNO CRD INFORMER
	// watches CRD resources:
	//		- Policy
	//		- PolicyVolation
	// - cache resync time: 10 seconds
	pInformer := kyvernoinformer.NewSharedInformerFactoryWithOptions(
		pclient,
		10*time.Second)

	// Configuration Data
	// dyamically load the configuration from configMap
	// - resource filters
	// if the configMap is update, the configuration will be updated :D
	configData := config.NewConfigData(
		kubeClient,
		kubeInformer.Core().V1().ConfigMaps(),
		filterK8Resources)

	// Policy meta-data store
	policyMetaStore := policystore.NewPolicyStore(pInformer.Kyverno().V1().ClusterPolicies())

	// EVENT GENERATOR
	// - generate event with retry mechanism
	egen := event.NewEventGenerator(
		client,
		pInformer.Kyverno().V1().ClusterPolicies())

	// POLICY VIOLATION GENERATOR
	// -- generate policy violation
	pvgen := policyviolation.NewPVGenerator(pclient,
		client,
		pInformer.Kyverno().V1().ClusterPolicyViolations(),
		pInformer.Kyverno().V1().PolicyViolations())

	// POLICY CONTROLLER
	// - reconciliation policy and policy violation
	// - process policy on existing resources
	// - status aggregator: recieves stats when a policy is applied
	//					    & updates the policy status
	pc, err := policy.NewPolicyController(pclient,
		client,
		pInformer.Kyverno().V1().ClusterPolicies(),
		pInformer.Kyverno().V1().ClusterPolicyViolations(),
		pInformer.Kyverno().V1().PolicyViolations(),
		configData,
		egen,
		pvgen,
		policyMetaStore,
		rWebhookWatcher)
	if err != nil {
		glog.Fatalf("error creating policy controller: %v\n", err)
	}

	// GENERATE CONTROLLER
	// - watches for Namespace resource and generates resource based on the policy generate rule
	nsc := namespace.NewNamespaceController(
		pclient,
		client,
		kubeInformer.Core().V1().Namespaces(),
		pInformer.Kyverno().V1().ClusterPolicies(),
		pc.GetPolicyStatusAggregator(),
		egen,
		configData,
		pvgen,
		policyMetaStore)

	// CONFIGURE CERTIFICATES
	tlsPair, err := client.InitTLSPemPair(clientConfig)
	if err != nil {
		glog.Fatalf("Failed to initialize TLS key/certificate pair: %v\n", err)
	}

	// WEBHOOK REGISTRATION
	// - mutating,validatingwebhookconfiguration (Policy)
	// - verifymutatingwebhookconfiguration (Kyverno Deployment)
	// resource webhook confgiuration is generated dynamically in the webhook server and policy controller
	// based on the policy resources created
	if err = webhookRegistrationClient.Register(); err != nil {
		glog.Fatalf("Failed registering Admission Webhooks: %v\n", err)
	}

	// WEBHOOOK
	// - https server to provide endpoints called based on rules defined in Mutating & Validation webhook configuration
	// - reports the results based on the response from the policy engine:
	// -- annotations on resources with update details on mutation JSON patches
	// -- generate policy violation resource
	// -- generate events on policy and resource
	server, err := webhooks.NewWebhookServer(
		pclient,
		client,
		tlsPair,
		pInformer.Kyverno().V1().ClusterPolicies(),
		kubeInformer.Rbac().V1().RoleBindings(),
		kubeInformer.Rbac().V1().ClusterRoleBindings(),
		egen,
		webhookRegistrationClient,
		pc.GetPolicyStatusAggregator(),
		configData,
		policyMetaStore,
		pvgen,
		rWebhookWatcher,
		cleanUp)
	if err != nil {
		glog.Fatalf("Unable to create webhook server: %v\n", err)
	}
	// Start the components
	pInformer.Start(stopCh)
	kubeInformer.Start(stopCh)

	go rWebhookWatcher.Run(stopCh)
	go configData.Run(stopCh)
	go policyMetaStore.Run(stopCh)
	go pc.Run(1, stopCh)
	go egen.Run(1, stopCh)
	go nsc.Run(1, stopCh)
	go pvgen.Run(1, stopCh)

	// verifys if the admission control is enabled and active
	// resync: 60 seconds
	// deadline: 60 seconds (send request)
	// max deadline: deadline*3 (set the deployment annotation as false)
	server.RunAsync(stopCh)

	<-stopCh

	// by default http.Server waits indefinitely for connections to return to idle and then shuts down
	// adding a threshold will handle zombie connections
	// adjust the context deadline to 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	// cleanup webhookconfigurations followed by webhook shutdown
	server.Stop(ctx)
	// resource cleanup
	// remove webhook configurations
	<-cleanUp
	glog.Info("successful shutdown of kyverno controller")
}

func init() {
	flag.StringVar(&filterK8Resources, "filterK8Resources", "", "k8 resource in format [kind,namespace,name] where policy is not evaluated by the admission webhook. example --filterKind \"[Deployment, kyverno, kyverno]\" --filterKind \"[Deployment, kyverno, kyverno],[Events, *, *]\"")
	flag.IntVar(&webhookTimeout, "webhooktimeout", 3, "timeout for webhook configurations")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&serverIP, "serverIP", "", "IP address where Kyverno controller runs. Only required if out-of-cluster.")
	config.LogDefaultFlags()
	flag.Parse()
}

package controller

import (
    "time"
    "fmt"

    "k8s.io/sample-controller/pkg/signals"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/client-go/tools/cache"

    clientset "nirmata/kube-policy/pkg/client/clientset/versioned"
    informers "nirmata/kube-policy/pkg/client/informers/externalversions"
    lister "nirmata/kube-policy/pkg/client/listers/policy/v1alpha1"
    types "nirmata/kube-policy/pkg/apis/policy/v1alpha1"
)

// Controller for CRD
type Controller struct {
    policyInformerFactory informers.SharedInformerFactory
    policyLister lister.PolicyLister
}

// NewController from cmd args
func NewController(masterURL, kubeconfigPath string) (*Controller, error) {
    cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
    if err != nil {
        fmt.Printf("Error building kubeconfig: %v\n", err)
        return nil, err
    }

    policyClientset, err := clientset.NewForConfig(cfg)
    if err != nil {
        fmt.Printf("Error building policy clientset: %v\n", err)
        return nil, err
    }

    policyInformerFactory := informers.NewSharedInformerFactory(policyClientset, time.Second*30)
    policyInformer := policyInformerFactory.Nirmata().V1alpha1().Policies()
    
    controller := &Controller {
        policyInformerFactory: policyInformerFactory,
        policyLister: policyInformer.Lister(),
    }

    policyInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: controller.createPolicyHandler,
        UpdateFunc: controller.updatePolicyHandler,
        DeleteFunc: controller.deletePolicyHandler,
    })

    return controller, nil
}

// Run is main controller thread
func (c *Controller) Run() error {
    stopCh := signals.SetupSignalHandler()
    c.policyInformerFactory.Start(stopCh)

    fmt.Println("Running controller...")
    <-stopCh
    fmt.Println("\nShutting down controller...")

    return nil
}

// GetPolicies retrieves all policy resources
// from cache. Cache is refreshed by informer
func (c *Controller) GetPolicies() ([]*types.Policy, error) {
    // Create nil Selector to grab all the policies
    cachedPolicies, err := c.policyLister.List(labels.NewSelector())

    var policies []*types.Policy

    if err != nil {
        return nil, err
    }

    for _, elem := range cachedPolicies {
        policies = append(policies, elem.DeepCopy())
    }

    return policies, nil
}

func (c *Controller) createPolicyHandler(resource interface{}) {
    key := c.getResourceKey(resource)
    fmt.Printf("Created policy: %s\n", key)
}

func (c *Controller) updatePolicyHandler(oldResource, newResource interface{}) {
    oldKey := c.getResourceKey(oldResource)
    newKey := c.getResourceKey(newResource)

    fmt.Printf("Updated policy from %s to %s\n", oldKey, newKey)
}

func (c *Controller) deletePolicyHandler(resource interface{}) {
    key := c.getResourceKey(resource)
    fmt.Printf("Deleted policy: %s\n", key)
}

func (c *Controller) getResourceKey(resource interface{}) string {    
    if key, err := cache.MetaNamespaceKeyFunc(resource); err != nil {
        fmt.Printf("Error retrieving policy key: %v\n", err)
        return ""
    } else {
        return key
    }
}
package policy

import (
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/minio/minio/pkg/wildcard"
	kyverno "github.com/nirmata/kyverno/pkg/api/kyverno/v1alpha1"
	client "github.com/nirmata/kyverno/pkg/dclient"
	"github.com/nirmata/kyverno/pkg/info"
	"github.com/nirmata/kyverno/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (pc *PolicyController) processExistingResources(policy kyverno.Policy) {
	// Parse through all the resources
	// drops the cache after configured rebuild time
	pc.rm.Drop()

	// get resource that are satisfy the resource description defined in the rules
	resourceMap := listResources(pc.client, policy, pc.filterK8Resources)
	for _, resource := range resourceMap {
		// pre-processing, check if the policy and resource version has been processed before
		if !pc.rm.ProcessResource(policy.Name, policy.ResourceVersion, resource.GetKind(), resource.GetNamespace(), resource.GetName(), resource.GetResourceVersion()) {
			glog.V(4).Infof("policy %s with resource version %s already processed on resource %s/%s/%s with resource version %s", policy.Name, policy.ResourceVersion, resource.GetKind(), resource.GetNamespace(), resource.GetName(), resource.GetResourceVersion())
			continue
		}
		// apply the policy on each
		glog.V(4).Infof("apply policy %s with resource version %s on resource %s/%s/%s with resource version %s", policy.Name, policy.ResourceVersion, resource.GetKind(), resource.GetNamespace(), resource.GetName(), resource.GetResourceVersion())
		applyPolicyOnResource(policy, resource)
		// post-processing, register the resource as processed
		pc.rm.RegisterResource(policy.GetName(), policy.GetResourceVersion(), resource.GetKind(), resource.GetNamespace(), resource.GetName(), resource.GetResourceVersion())
	}
}

func applyPolicyOnResource(policy kyverno.Policy, resource unstructured.Unstructured) *info.PolicyInfo {
	policyInfo, err := applyPolicy(policy, resource)
	if err != nil {
		glog.V(4).Infof("failed to process policy %s on resource %s/%s/%s: %v", policy.GetName(), resource.GetKind(), resource.GetNamespace(), resource.GetName(), err)
		return nil
	}
	return &policyInfo
}

func listResources(client *client.Client, policy kyverno.Policy, filterK8Resources []utils.K8Resource) map[string]unstructured.Unstructured {
	// key uid
	resourceMap := map[string]unstructured.Unstructured{}

	for _, rule := range policy.Spec.Rules {
		// resources that match
		for _, k := range rule.MatchResources.Kinds {
			if kindIsExcluded(k, rule.ExcludeResources.Kinds) {
				glog.V(4).Infof("processing policy %s rule %s: kind %s is exluded", policy.Name, rule.Name, k)
				continue
			}
			var namespaces []string
			if k == "Namespace" {
				// TODO
				// this is handled by generator controller
				glog.V(4).Infof("skipping processing policy %s rule %s for kind Namespace", policy.Name, rule.Name)
				continue
			}
			//TODO: if namespace is not define can we default to *
			if rule.MatchResources.Namespace != "" {
				namespaces = append(namespaces, rule.MatchResources.Namespace)
			} else {
				glog.V(4).Infof("processing policy %s rule %s, namespace not defined, getting all namespaces ", policy.Name, rule.Name)
				// get all namespaces
				namespaces = getAllNamespaces(client)
			}
			// check if exclude namespace is not clashing
			namespaces = excludeNamespaces(namespaces, rule.ExcludeResources.Namespace)

			// get resources in the namespaces
			for _, ns := range namespaces {
				rMap := getResourcesPerNamespace(k, client, ns, rule, filterK8Resources)
				mergeresources(resourceMap, rMap)
			}

		}
	}
	return resourceMap
}

func getResourcesPerNamespace(kind string, client *client.Client, namespace string, rule kyverno.Rule, filterK8Resources []utils.K8Resource) map[string]unstructured.Unstructured {
	resourceMap := map[string]unstructured.Unstructured{}
	// merge include and exclude label selector values
	ls := mergeLabelSectors(rule.MatchResources.Selector, rule.ExcludeResources.Selector)
	// list resources
	glog.V(4).Infof("get resources for kind %s, namespace %s, selector %v", kind, namespace, rule.MatchResources.Selector)
	list, err := client.ListResource(kind, namespace, ls)
	if err != nil {
		glog.Infof("unable to get resources: err %v", err)
		return nil
	}
	// filter based on name
	for _, r := range list.Items {
		// match name
		if rule.MatchResources.Name != "" {
			if !wildcard.Match(rule.MatchResources.Name, r.GetName()) {
				glog.V(4).Infof("skipping resource %s/%s due to include condition name=%s mistatch", r.GetNamespace(), r.GetName(), rule.MatchResources.Name)
				continue
			}
		}
		// exclude name
		if rule.ExcludeResources.Name != "" {
			if wildcard.Match(rule.ExcludeResources.Name, r.GetName()) {
				glog.V(4).Infof("skipping resource %s/%s due to exclude condition name=%s mistatch", r.GetNamespace(), r.GetName(), rule.MatchResources.Name)
				continue
			}
		}
		// Skip the filtered resources
		if utils.SkipFilteredResources(r.GetKind(), r.GetNamespace(), r.GetName(), filterK8Resources) {
			continue
		}

		//TODO check if the group version kind is present or not
		resourceMap[string(r.GetUID())] = r
	}
	return resourceMap
}

// merge b into a map
func mergeresources(a, b map[string]unstructured.Unstructured) {
	for k, v := range b {
		a[k] = v
	}
}
func mergeLabelSectors(include, exclude *metav1.LabelSelector) *metav1.LabelSelector {
	if exclude == nil {
		return include
	}
	// negate the exclude information
	// copy the label selector
	//TODO: support exclude expressions in exclude
	ls := include.DeepCopy()
	for k, v := range exclude.MatchLabels {
		lsreq := metav1.LabelSelectorRequirement{
			Key:      k,
			Operator: metav1.LabelSelectorOpNotIn,
			Values:   []string{v},
		}
		ls.MatchExpressions = append(ls.MatchExpressions, lsreq)
	}
	return ls
}

func kindIsExcluded(kind string, list []string) bool {
	for _, b := range list {
		if b == kind {
			return true
		}
	}
	return false
}

func excludeNamespaces(namespaces []string, excludeNs string) []string {
	if excludeNs == "" {
		return namespaces
	}
	filteredNamespaces := []string{}
	for _, n := range namespaces {
		if n == excludeNs {
			continue
		}
		filteredNamespaces = append(filteredNamespaces, n)
	}
	return filteredNamespaces
}

func getAllNamespaces(client *client.Client) []string {
	var namespaces []string
	// get all namespaces
	nsList, err := client.ListResource("Namespace", "", nil)
	if err != nil {
		glog.Error(err)
		return namespaces
	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.GetName())
	}
	return namespaces
}

func NewResourceManager(rebuildTime int64) *ResourceManager {
	rm := ResourceManager{
		data:        make(map[string]interface{}),
		time:        time.Now(),
		rebuildTime: rebuildTime,
	}
	// set time it was built
	return &rm
}

// ResourceManager
type ResourceManager struct {
	// we drop and re-build the cache
	// based on the memory consumer of by the map
	data        map[string]interface{}
	mux         sync.RWMutex
	time        time.Time
	rebuildTime int64 // after how many seconds should we rebuild the cache
}

type resourceManager interface {
	ProcessResource(policy, pv, kind, ns, name, rv string) bool
	//TODO	removeResource(kind, ns, name string) error
	RegisterResource(policy, pv, kind, ns, name, rv string)
	// reload
	Drop()
}

//Drop drop the cache after every rebuild interval mins
//TODO: or drop based on the size
func (rm *ResourceManager) Drop() {
	timeSince := time.Since(rm.time)
	glog.V(4).Infof("time since last cache reset time %v is %v", rm.time, timeSince)
	glog.V(4).Infof("cache rebuild time %v", time.Duration(rm.rebuildTime)*time.Second)
	if timeSince > time.Duration(rm.rebuildTime)*time.Second {
		rm.mux.Lock()
		defer rm.mux.Unlock()
		rm.data = map[string]interface{}{}
		rm.time = time.Now()
		glog.V(4).Infof("dropping cache at time %v", rm.time)
	}
}

var empty struct{}

//RegisterResource stores if the policy is processed on this resource version
func (rm *ResourceManager) RegisterResource(policy, pv, kind, ns, name, rv string) {
	rm.mux.Lock()
	defer rm.mux.Unlock()
	// add the resource
	key := buildKey(policy, pv, kind, ns, name, rv)
	rm.data[key] = empty
}

//ProcessResource returns true if the policy was not applied on the resource
func (rm *ResourceManager) ProcessResource(policy, pv, kind, ns, name, rv string) bool {
	rm.mux.RLock()
	defer rm.mux.RUnlock()

	key := buildKey(policy, pv, kind, ns, name, rv)
	_, ok := rm.data[key]
	return ok == false
}

func buildKey(policy, pv, kind, ns, name, rv string) string {
	return policy + "/" + pv + "/" + kind + "/" + ns + "/" + name + "/" + rv
}

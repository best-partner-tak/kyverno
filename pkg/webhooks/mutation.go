package webhooks

import (
	"github.com/golang/glog"
	engine "github.com/nirmata/kyverno/pkg/engine"
	policyctr "github.com/nirmata/kyverno/pkg/policy"
	"github.com/nirmata/kyverno/pkg/utils"
	v1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// HandleMutation handles mutating webhook admission request
func (ws *WebhookServer) HandleMutation(request *v1beta1.AdmissionRequest) (bool, []byte, string) {
	glog.V(4).Infof("Receive request in mutating webhook: Kind=%s, Namespace=%s Name=%s UID=%s patchOperation=%s",
		request.Kind.Kind, request.Namespace, request.Name, request.UID, request.Operation)

	var patches [][]byte
	var policyStats []policyctr.PolicyStat

	// gather stats from the engine response
	gatherStat := func(policyName string, policyResponse engine.PolicyResponse) {
		ps := policyctr.PolicyStat{}
		ps.PolicyName = policyName
		ps.Stats.MutationExecutionTime = policyResponse.ProcessingTime
		ps.Stats.RulesAppliedCount = policyResponse.RulesAppliedCount
		policyStats = append(policyStats, ps)
	}
	// send stats for aggregation
	sendStat := func(blocked bool) {
		for _, stat := range policyStats {
			stat.Stats.ResourceBlocked = utils.Btoi(blocked)
			//SEND
			ws.policyStatus.SendStat(stat)
		}
	}
	// convert RAW to unstructured
	resource, err := engine.ConvertToUnstructured(request.Object.Raw)
	if err != nil {
		//TODO: skip applying the amiddions control ?
		glog.Errorf("unable to convert raw resource to unstructured: %v", err)
		return true, nil, ""
	}

	//TODO: check if resource gvk is available in raw resource,
	//TODO: check if the name and namespace is also passed right in the resource?
	// if not then set it from the api request
	resource.SetGroupVersionKind(schema.GroupVersionKind{Group: request.Kind.Group, Version: request.Kind.Version, Kind: request.Kind.Kind})
	policies, err := ws.pLister.List(labels.NewSelector())
	if err != nil {
		//TODO check if the CRD is created ?
		// Unable to connect to policy Lister to access policies
		glog.Errorln("Unable to connect to policy controller to access policies. Mutation Rules are NOT being applied")
		glog.Warning(err)
		return true, nil, ""
	}

	var engineResponses []engine.EngineResponseNew
	for _, policy := range policies {

		// check if policy has a rule for the admission request kind
		if !utils.Contains(getApplicableKindsForPolicy(policy), request.Kind.Kind) {
			continue
		}

		glog.V(4).Infof("Handling mutation for Kind=%s, Namespace=%s Name=%s UID=%s patchOperation=%s",
			resource.GetKind(), resource.GetNamespace(), resource.GetName(), request.UID, request.Operation)
		// TODO: this can be
		engineResponse := engine.MutateNew(*policy, *resource)
		engineResponses = append(engineResponses, engineResponse)
		// Gather policy application statistics
		gatherStat(policy.Name, engineResponse.PolicyResponse)
		if !engineResponse.IsSuccesful() {
			glog.V(4).Infof("Failed to apply policy %s on resource %s/%s\n", policy.Name, resource.GetNamespace(), resource.GetName())
			continue
		}
		// gather patches
		patches = append(patches, engineResponse.GetPatches())
		// generate annotations
		if annPatches := generateAnnotationPatches(resource.GetAnnotations(), engineResponse.PolicyResponse); annPatches != nil {
			patches = append(patches, annPatches)
		}
		glog.V(4).Infof("Mutation from policy %s has applied succesfully to %s %s/%s", policy.Name, request.Kind.Kind, resource.GetNamespace(), resource.GetName())
		//TODO: check if there is an order to policy application on resource
		// resource = &engineResponse.PatchedResource
	}

	// ADD EVENTS
	events := generateEvents(engineResponses, (request.Operation == v1beta1.Update))
	ws.eventGen.Add(events...)

	if isResponseSuccesful(engineResponses) {
		sendStat(false)
		patch := engine.JoinPatches(patches)
		return true, patch, ""
	}

	sendStat(true)
	glog.Errorf("Failed to mutate the resource\n")
	return false, nil, getErrorMsg(engineResponses)
}

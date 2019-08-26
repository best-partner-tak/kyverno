package webhooks

import (
	"github.com/golang/glog"
	engine "github.com/nirmata/kyverno/pkg/engine"
	policyctr "github.com/nirmata/kyverno/pkg/policy"
	"github.com/nirmata/kyverno/pkg/policyviolation"
	"github.com/nirmata/kyverno/pkg/utils"
	v1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// HandleValidation handles validating webhook admission request
// If there are no errors in validating rule we apply generation rules
// patchedResource is the (resource + patches) after applying mutation rules
func (ws *WebhookServer) HandleValidation(request *v1beta1.AdmissionRequest, patchedResource []byte) (bool, string) {
	glog.V(4).Infof("Receive request in validating webhook: Kind=%s, Namespace=%s Name=%s UID=%s patchOperation=%s",
		request.Kind.Kind, request.Namespace, request.Name, request.UID, request.Operation)

	var policyStats []policyctr.PolicyStat

	// gather stats from the engine response
	gatherStat := func(policyName string, policyResponse engine.PolicyResponse) {
		ps := policyctr.PolicyStat{}
		ps.PolicyName = policyName
		ps.Stats.ValidationExecutionTime = policyResponse.ProcessingTime
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

	resourceRaw := request.Object.Raw
	if patchedResource != nil {
		glog.V(4).Info("using patched resource from mutation to process validation rules")
		resourceRaw = patchedResource
	}
	// convert RAW to unstructured
	resource, err := engine.ConvertToUnstructured(resourceRaw)
	if err != nil {
		//TODO: skip applying the amiddions control ?
		glog.Errorf("unable to convert raw resource to unstructured: %v", err)
		return true, ""
	}
	//TODO: check if resource gvk is available in raw resource,
	// if not then set it from the api request
	resource.SetGroupVersionKind(schema.GroupVersionKind{Group: request.Kind.Group, Version: request.Kind.Version, Kind: request.Kind.Kind})
	//TODO: check if the name and namespace is also passed right in the resource?
	// all the patches to be applied on the resource

	policies, err := ws.pLister.List(labels.NewSelector())
	if err != nil {
		//TODO check if the CRD is created ?
		// Unable to connect to policy Lister to access policies
		glog.Error("Unable to connect to policy controller to access policies. Validation Rules are NOT being applied")
		glog.Warning(err)
		return true, ""
	}

	var engineResponses []engine.EngineResponseNew
	for _, policy := range policies {

		if !utils.Contains(getApplicableKindsForPolicy(policy), request.Kind.Kind) {
			continue
		}

		glog.V(4).Infof("Handling validation for Kind=%s, Namespace=%s Name=%s UID=%s patchOperation=%s",
			resource.GetKind(), resource.GetNamespace(), resource.GetName(), request.UID, request.Operation)

		// glog.V(4).Infof("Validating resource %s/%s/%s with policy %s with %d rules\n", resource.GetKind(), resource.GetNamespace(), resource.GetName(), policy.ObjectMeta.Name, len(policy.Spec.Rules))

		engineResponse := engine.ValidateNew(*policy, *resource)
		engineResponses = append(engineResponses, engineResponse)
		// Gather policy application statistics
		gatherStat(policy.Name, engineResponse.PolicyResponse)
		if !engineResponse.IsSuccesful() {
			glog.V(4).Infof("Failed to apply policy %s on resource %s/%s\n", policy.Name, resource.GetNamespace(), resource.GetName())
			continue
		}
	}
	// ADD EVENTS
	events := generateEvents(engineResponses, (request.Operation == v1beta1.Update))
	ws.eventGen.Add(events...)

	// If Validation fails then reject the request
	// violations are created if "audit" flag is set
	// and if there are any then we dont block the resource creation
	// Even if one the policy being applied
	if !isResponseSuccesful(engineResponses) && toBlockResource(engineResponses) {
		sendStat(true)
		return false, getErrorMsg(engineResponses)
	}

	// ADD POLICY VIOLATIONS
	policyviolation.CreatePV(ws.pvLister, ws.kyvernoClient, engineResponses)

	sendStat(false)
	return true, ""
}

package engine

import (
	"github.com/golang/glog"
	kubepolicy "github.com/nirmata/kyverno/pkg/apis/policy/v1alpha1"
	"github.com/nirmata/kyverno/pkg/info"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate performs mutation. Overlay first and then mutation patches
func Mutate(policy kubepolicy.Policy, rawResource []byte, gvk metav1.GroupVersionKind) ([][]byte, []byte, []*info.RuleInfo) {
	var allPatches, rulePatches [][]byte
	var err error
	var errs []error
	patchedDocument := rawResource
	ris := []*info.RuleInfo{}

	for _, rule := range policy.Spec.Rules {
		if rule.Mutation == nil {
			continue
		}
		ri := info.NewRuleInfo(rule.Name, info.Mutation)

		ok := ResourceMeetsDescription(rawResource, rule.MatchResources.ResourceDescription, rule.ExcludeResources.ResourceDescription, gvk)
		if !ok {
			glog.V(3).Infof("Not applicable on specified resource kind%s", gvk.Kind)
			continue
		}
		// Process Overlay
		if rule.Mutation.Overlay != nil {
			rulePatches, err = ProcessOverlay(rule, patchedDocument, gvk)
			if err == nil {
				if len(rulePatches) == 0 {
					// if array elements dont match then we skip(nil patch, no error)
					// or if acnohor is defined and doenst match
					// policy is not applicable
					continue
				}
				ri.Addf("Rule %s: Overlay succesfully applied.", rule.Name)
				// merge the json patches
				patch := JoinPatches(rulePatches)
				// strip slashes from string
				ri.Changes = string(patch)
				allPatches = append(allPatches, rulePatches...)
			} else {
				ri.Fail()
				ri.Addf("overlay application has failed, err %v.", err)
			}
		}

		// Process Patches
		if len(rule.Mutation.Patches) != 0 {
			rulePatches, errs = ProcessPatches(rule, patchedDocument)
			if len(errs) > 0 {
				ri.Fail()
				for _, err := range errs {
					ri.Addf("patches application has failed, err %v.", err)
				}
			} else {
				ri.Addf("Rule %s: Patches succesfully applied.", rule.Name)
				allPatches = append(allPatches, rulePatches...)
			}
		}

		patchedDocument, err = ApplyPatches(rawResource, rulePatches)
		if err != nil {
			glog.Errorf("Failed to apply patches on ruleName=%s, err%v\n:", rule.Name, err)
		}
		ris = append(ris, ri)
	}

	return allPatches, patchedDocument, ris
}

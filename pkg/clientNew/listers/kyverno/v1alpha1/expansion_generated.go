/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"fmt"

	kyverno "github.com/nirmata/kyverno/pkg/api/kyverno/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// PolicyListerExpansion allows custom methods to be added to
// PolicyLister.
type PolicyListerExpansion interface {
	GetPolicyForPolicyViolation(pv *kyverno.PolicyViolation) ([]*kyverno.Policy, error)
}

// PolicyViolationListerExpansion allows custom methods to be added to
// PolicyViolationLister.
type PolicyViolationListerExpansion interface{}

func (pl *policyLister) GetPolicyForPolicyViolation(pv *kyverno.PolicyViolation) ([]*kyverno.Policy, error) {
	if len(pv.Labels) == 0 {
		return nil, fmt.Errorf("no Policy found for PolicyViolation %v because it has no labels", pv.Name)
	}

	pList, err := pl.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var policies []*kyverno.Policy
	for _, p := range pList {
		policyLabelmap := map[string]string{"policy": p.Name}

		ls := &metav1.LabelSelector{}
		err = metav1.Convert_Map_string_To_string_To_v1_LabelSelector(&policyLabelmap, ls, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate label sector of Policy name %s: %v", p.Name, err)
		}
		selector, err := metav1.LabelSelectorAsSelector(ls)
		if err != nil {
			return nil, fmt.Errorf("invalid label selector: %v", err)
		}
		// If a policy with a nil or empty selector creeps in, it should match nothing, not everything.
		if selector.Empty() || !selector.Matches(labels.Set(pv.Labels)) {
			continue
		}
		policies = append(policies, p)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("could not find Policy set for PolicyViolation %s with labels: %v", pv.Name, pv.Labels)
	}

	return policies, nil

}

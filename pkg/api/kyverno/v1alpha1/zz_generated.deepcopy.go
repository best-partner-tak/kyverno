// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloneFrom) DeepCopyInto(out *CloneFrom) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloneFrom.
func (in *CloneFrom) DeepCopy() *CloneFrom {
	if in == nil {
		return nil
	}
	out := new(CloneFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPolicy) DeepCopyInto(out *ClusterPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPolicy.
func (in *ClusterPolicy) DeepCopy() *ClusterPolicy {
	if in == nil {
		return nil
	}
	out := new(ClusterPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPolicyList) DeepCopyInto(out *ClusterPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPolicyList.
func (in *ClusterPolicyList) DeepCopy() *ClusterPolicyList {
	if in == nil {
		return nil
	}
	out := new(ClusterPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPolicyViolation) DeepCopyInto(out *ClusterPolicyViolation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPolicyViolation.
func (in *ClusterPolicyViolation) DeepCopy() *ClusterPolicyViolation {
	if in == nil {
		return nil
	}
	out := new(ClusterPolicyViolation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPolicyViolation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPolicyViolationList) DeepCopyInto(out *ClusterPolicyViolationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterPolicyViolation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPolicyViolationList.
func (in *ClusterPolicyViolationList) DeepCopy() *ClusterPolicyViolationList {
	if in == nil {
		return nil
	}
	out := new(ClusterPolicyViolationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPolicyViolationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExcludeResources) DeepCopyInto(out *ExcludeResources) {
	*out = *in
	in.ResourceDescription.DeepCopyInto(&out.ResourceDescription)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExcludeResources.
func (in *ExcludeResources) DeepCopy() *ExcludeResources {
	if in == nil {
		return nil
	}
	out := new(ExcludeResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Generation.
func (in *Generation) DeepCopy() *Generation {
	if in == nil {
		return nil
	}
	out := new(Generation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchResources) DeepCopyInto(out *MatchResources) {
	*out = *in
	in.ResourceDescription.DeepCopyInto(&out.ResourceDescription)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchResources.
func (in *MatchResources) DeepCopy() *MatchResources {
	if in == nil {
		return nil
	}
	out := new(MatchResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mutation.
func (in *Mutation) DeepCopy() *Mutation {
	if in == nil {
		return nil
	}
	out := new(Mutation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Policy) DeepCopyInto(out *Policy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Policy.
func (in *Policy) DeepCopy() *Policy {
	if in == nil {
		return nil
	}
	out := new(Policy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyStatus) DeepCopyInto(out *PolicyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyStatus.
func (in *PolicyStatus) DeepCopy() *PolicyStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyViolation) DeepCopyInto(out *PolicyViolation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyViolation.
func (in *PolicyViolation) DeepCopy() *PolicyViolation {
	if in == nil {
		return nil
	}
	out := new(PolicyViolation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyViolationSpec) DeepCopyInto(out *PolicyViolationSpec) {
	*out = *in
	out.ResourceSpec = in.ResourceSpec
	if in.ViolatedRules != nil {
		in, out := &in.ViolatedRules, &out.ViolatedRules
		*out = make([]ViolatedRule, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyViolationSpec.
func (in *PolicyViolationSpec) DeepCopy() *PolicyViolationSpec {
	if in == nil {
		return nil
	}
	out := new(PolicyViolationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyViolationStatus) DeepCopyInto(out *PolicyViolationStatus) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyViolationStatus.
func (in *PolicyViolationStatus) DeepCopy() *PolicyViolationStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyViolationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceDescription) DeepCopyInto(out *ResourceDescription) {
	*out = *in
	if in.Kinds != nil {
		in, out := &in.Kinds, &out.Kinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceDescription.
func (in *ResourceDescription) DeepCopy() *ResourceDescription {
	if in == nil {
		return nil
	}
	out := new(ResourceDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSpec) DeepCopyInto(out *ResourceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSpec.
func (in *ResourceSpec) DeepCopy() *ResourceSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	in.MatchResources.DeepCopyInto(&out.MatchResources)
	in.ExcludeResources.DeepCopyInto(&out.ExcludeResources)
	in.Mutation.DeepCopyInto(&out.Mutation)
	in.Validation.DeepCopyInto(&out.Validation)
	in.Generation.DeepCopyInto(&out.Generation)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Spec) DeepCopyInto(out *Spec) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]Rule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Spec.
func (in *Spec) DeepCopy() *Spec {
	if in == nil {
		return nil
	}
	out := new(Spec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Validation.
func (in *Validation) DeepCopy() *Validation {
	if in == nil {
		return nil
	}
	out := new(Validation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ViolatedRule) DeepCopyInto(out *ViolatedRule) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ViolatedRule.
func (in *ViolatedRule) DeepCopy() *ViolatedRule {
	if in == nil {
		return nil
	}
	out := new(ViolatedRule)
	in.DeepCopyInto(out)
	return out
}

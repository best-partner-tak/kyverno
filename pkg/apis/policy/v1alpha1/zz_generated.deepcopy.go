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

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Generation.
func (in *Generation) DeepCopy() *Generation {
	if in == nil {
		return nil
	}
	out := new(Generation)
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
	in.Status.DeepCopyInto(&out.Status)
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

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Policy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyList) DeepCopyInto(out *PolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Policy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyList.
func (in *PolicyList) DeepCopy() *PolicyList {
	if in == nil {
		return nil
	}
	out := new(PolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceDescription) DeepCopyInto(out *ResourceDescription) {
	*out = *in
	if in.Kinds != nil {
		in, out := &in.Kinds, &out.Kinds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
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
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	in.ResourceDescription.DeepCopyInto(&out.ResourceDescription)
	if in.Mutation != nil {
		in, out := &in.Mutation, &out.Mutation
		*out = (*in).DeepCopy()
	}
	if in.Validation != nil {
		in, out := &in.Validation, &out.Validation
		*out = (*in).DeepCopy()
	}
	if in.Generation != nil {
		in, out := &in.Generation, &out.Generation
		*out = (*in).DeepCopy()
	}
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	if in.Violations != nil {
		in, out := &in.Violations, &out.Violations
		*out = make([]Violation, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
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
func (in *Violation) DeepCopyInto(out *Violation) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Violation.
func (in *Violation) DeepCopy() *Violation {
	if in == nil {
		return nil
	}
	out := new(Violation)
	in.DeepCopyInto(out)
	return out
}

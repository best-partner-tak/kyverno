package v1alpha1

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Validate checks if rule is not empty and all substructures are valid
func (r *Rule) Validate() error {
	err := r.ResourceDescription.Validate()
	if err != nil {
		return err
	}

	if r.Mutation == nil && r.Validation == nil && r.Generation == nil {
		return errors.New("The rule is empty")
	}

	return nil
}

// Validate checks if all necesarry fields are present and have values. Also checks a Selector.
// Returns error if
// - kinds is not defined
func (pr *ResourceDescription) Validate() error {
	if len(pr.Kinds) == 0 {
		return errors.New("The Kind is not specified")
	}

	if pr.Selector != nil {
		selector, err := metav1.LabelSelectorAsSelector(pr.Selector)
		if err != nil {
			return err
		}
		requirements, _ := selector.Requirements()
		if len(requirements) == 0 {
			return errors.New("The requirements are not specified in selector")
		}
	}

	return nil
}

// Validate if all mandatory PolicyPatch fields are set
func (pp *Patch) Validate() error {
	if pp.Path == "" {
		return errors.New("JSONPatch field 'path' is mandatory")
	}

	if pp.Operation == "add" || pp.Operation == "replace" {
		if pp.Value == nil {
			return fmt.Errorf("JSONPatch field 'value' is mandatory for operation '%s'", pp.Operation)
		}

		return nil
	} else if pp.Operation == "remove" {
		return nil
	}

	return fmt.Errorf("Unsupported JSONPatch operation '%s'", pp.Operation)
}

// Validate returns error if generator is configured incompletely
func (pcg *Generation) Validate() error {
	if len(pcg.Data) == 0 && pcg.CopyFrom == nil {
		return fmt.Errorf("Neither Data nor CopyFrom (source) of %s/%s is specified", pcg.Kind, pcg.Name)
	}
	return nil
}

// DeepCopyInto is declared because k8s:deepcopy-gen is
// not able to generate this method for interface{} member
func (in *Mutation) DeepCopyInto(out *Mutation) {
	if out != nil {
		*out = *in
	}
}

// DeepCopyInto is declared because k8s:deepcopy-gen is
// not able to generate this method for interface{} member
func (pp *Patch) DeepCopyInto(out *Patch) {
	if out != nil {
		*out = *pp
	}
}

// DeepCopyInto is declared because k8s:deepcopy-gen is
// not able to generate this method for interface{} member
func (in *Validation) DeepCopyInto(out *Validation) {
	if out != nil {
		*out = *in
	}
}

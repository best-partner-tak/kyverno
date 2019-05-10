package event

import (
	"fmt"
	"regexp"
)

//Key to describe the event
type EventMsg int

const (
	FResourcePolcy EventMsg = iota
	FProcessRule
	SPolicyApply
	SRuleApply
	FPolicyApplyBlockCreate
	FPolicyApplyBlockUpdate
	FPolicyApplyBlockUpdateRule
)

func (k EventMsg) String() string {
	return [...]string{
		"Failed to satisfy policy on resource %s.The following rules %s failed to apply. Created Policy Violation",
		"Failed to process rule %s of policy %s. Created Policy Violation %s",
		"Policy applied successfully on the resource %s",
		"Rule %s of Policy %s applied successfull",
		"Failed to apply policy, blocked creation of resource %s. The following rules %s failed to apply",
		"Failed to apply rule %s of policy %s Blocked update of the resource",
		"Failed to apply policy on resource %s.Blocked update of the resource. The following rules %s failed to apply",
	}[k]
}

const argRegex = "%[s,d,v]"

//GetEventMsg return the application message based on the message id and the arguments,
// if the number of arguments passed to the message are incorrect generate an error
func getEventMsg(key EventMsg, args ...interface{}) (string, error) {
	// Verify the number of arguments
	re := regexp.MustCompile(argRegex)
	argsCount := len(re.FindAllString(key.String(), -1))
	if argsCount != len(args) {
		return "", fmt.Errorf("message expects %d arguments, but %d arguments passed", argsCount, len(args))
	}
	return fmt.Sprintf(key.String(), args...), nil
}

package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/golang/glog"

	jsonpatch "github.com/evanphx/json-patch"
	kubepolicy "github.com/nirmata/kyverno/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProcessOverlay handles validating admission request
// Checks the target resources for rules defined in the policy
func ProcessOverlay(rule kubepolicy.Rule, rawResource []byte, gvk metav1.GroupVersionKind) ([]PatchBytes, error) {

	var resource interface{}
	var appliedPatches []PatchBytes
	err := json.Unmarshal(rawResource, &resource)
	if err != nil {
		return nil, err
	}

	patches, err := mutateResourceWithOverlay(resource, *rule.Mutation.Overlay)
	if err != nil {
		return nil, err
	}
	appliedPatches = append(appliedPatches, patches...)

	return appliedPatches, err
}

// mutateResourceWithOverlay is a start of overlaying process
func mutateResourceWithOverlay(resource, pattern interface{}) ([]PatchBytes, error) {
	// It assumes that mutation is started from root, so "/" is passed
	return applyOverlay(resource, pattern, "/")
}

// applyOverlay detects type of current item and goes down through overlay and resource trees applying overlay
func applyOverlay(resource, overlay interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes
	// resource item exists but has different type - replace
	// all subtree within this path by overlay
	if reflect.TypeOf(resource) != reflect.TypeOf(overlay) {
		patch, err := replaceSubtree(overlay, path)
		if err != nil {
			return nil, err
		}

		appliedPatches = append(appliedPatches, patch)
	}
	return applyOverlayForSameTypes(resource, overlay, path)
}

// applyOverlayForSameTypes is applyOverlay for cases when TypeOf(resource) == TypeOf(overlay)
func applyOverlayForSameTypes(resource, overlay interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	// detect the type of resource and overlay and select corresponding handler
	switch typedOverlay := overlay.(type) {
	// map
	case map[string]interface{}:
		typedResource := resource.(map[string]interface{})
		patches, err := applyOverlayToMap(typedResource, typedOverlay, path)
		if err != nil {
			return nil, err
		}
		appliedPatches = append(appliedPatches, patches...)
	// array
	case []interface{}:
		typedResource := resource.([]interface{})
		patches, err := applyOverlayToArray(typedResource, typedOverlay, path)
		if err != nil {
			return nil, err
		}
		appliedPatches = append(appliedPatches, patches...)
	// elementary types
	case string, float64, int64, bool:
		patch, err := replaceSubtree(overlay, path)
		if err != nil {
			return nil, err
		}
		appliedPatches = append(appliedPatches, patch)
	default:
		return nil, fmt.Errorf("Overlay has unsupported type: %T", overlay)
	}

	return appliedPatches, nil
}

// for each overlay and resource map elements applies overlay
func applyOverlayToMap(resourceMap, overlayMap map[string]interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	for key, value := range overlayMap {
		// skip anchor element because it has condition, not
		// the value that must replace resource value
		if isConditionAnchor(key) {
			continue
		}

		noAnchorKey := removeAnchor(key)
		currentPath := path + noAnchorKey + "/"
		resourcePart, ok := resourceMap[noAnchorKey]

		if ok && !isAddingAnchor(key) {
			// Key exists - go down through the overlay and resource trees
			patches, err := applyOverlay(resourcePart, value, currentPath)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patches...)
		}

		if !ok {
			// Key does not exist - insert entire overlay subtree
			patch, err := insertSubtree(value, currentPath)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patch)
		}
	}

	return appliedPatches, nil
}

// for each overlay and resource array elements applies overlay
func applyOverlayToArray(resource, overlay []interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	if 0 == len(overlay) {
		return nil, errors.New("Empty array detected in the overlay")
	}

	if 0 == len(resource) {
		// If array resource is empty, insert part from overlay
		patch, err := insertSubtree(overlay, path)
		if err != nil {
			return nil, err
		}
		appliedPatches = append(appliedPatches, patch)

		return appliedPatches, nil
	}

	if reflect.TypeOf(resource[0]) != reflect.TypeOf(overlay[0]) {
		return nil, fmt.Errorf("Overlay array and resource array have elements of different types: %T and %T", overlay[0], resource[0])
	}

	return applyOverlayToArrayOfSameTypes(resource, overlay, path)
}

// applyOverlayToArrayOfSameTypes applies overlay to array elements if they (resource and overlay elements) have same type
func applyOverlayToArrayOfSameTypes(resource, overlay []interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	switch overlay[0].(type) {
	case map[string]interface{}:
		return applyOverlayToArrayOfMaps(resource, overlay, path)
	default:
		lastElementIdx := len(resource)

		// Add elements to the end
		for i, value := range overlay {
			currentPath := path + strconv.Itoa(lastElementIdx+i) + "/"
			// currentPath example: /spec/template/spec/containers/3/
			patch, err := insertSubtree(value, currentPath)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patch)
		}
	}

	return appliedPatches, nil
}

// Array of maps needs special handling as far as it can have anchors.
func applyOverlayToArrayOfMaps(resource, overlay []interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	lastElementIdx := len(resource)
	for i, overlayElement := range overlay {
		typedOverlay := overlayElement.(map[string]interface{})
		anchors := getAnchorsFromMap(typedOverlay)

		if len(anchors) > 0 {
			// If we have anchors - choose corresponding resource element and mutate it
			patches, err := applyOverlayWithAnchors(resource, overlayElement, anchors, path)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patches...)
		} else if hasNestedAnchors(overlayElement) {
			// If we have anchors on the lower level - continue traversing overlay and resource trees
			for j, resourceElement := range resource {
				currentPath := path + strconv.Itoa(j) + "/"
				// currentPath example: /spec/template/spec/containers/3/
				patches, err := applyOverlay(resourceElement, overlayElement, currentPath)
				if err != nil {
					return nil, err
				}
				appliedPatches = append(appliedPatches, patches...)
			}
		} else {
			// Overlay subtree has no anchors - insert new element
			currentPath := path + strconv.Itoa(lastElementIdx+i) + "/"
			// currentPath example: /spec/template/spec/containers/3/
			patch, err := insertSubtree(overlayElement, currentPath)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patch)
		}
	}

	return appliedPatches, nil
}

func applyOverlayWithAnchors(resource []interface{}, overlay interface{}, anchors map[string]interface{}, path string) ([]PatchBytes, error) {
	var appliedPatches []PatchBytes

	for i, resourceElement := range resource {
		typedResource := resourceElement.(map[string]interface{})

		currentPath := path + strconv.Itoa(i) + "/"
		// currentPath example: /spec/template/spec/containers/3/
		if !skipArrayObject(typedResource, anchors) {
			patches, err := applyOverlay(resourceElement, overlay, currentPath)
			if err != nil {
				return nil, err
			}
			appliedPatches = append(appliedPatches, patches...)
		}
	}

	return appliedPatches, nil
}

func insertSubtree(overlay interface{}, path string) (PatchBytes, error) {
	return processSubtree(overlay, path, "add")
}

func replaceSubtree(overlay interface{}, path string) (PatchBytes, error) {
	return processSubtree(overlay, path, "replace")
}

func processSubtree(overlay interface{}, path string, op string) (PatchBytes, error) {

	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if path == "" {
		path = "/"
	}

	value := prepareJSONValue(overlay)
	patchStr := fmt.Sprintf(`{ "op": "%s", "path": "%s", "value": %s }`, op, path, value)

	// check the patch
	_, err := jsonpatch.DecodePatch([]byte("[" + patchStr + "]"))
	if err != nil {
		glog.V(3).Info(err)
		return nil, fmt.Errorf("Failed to make '%s' patch from an overlay '%s' for path %s", op, value, path)
	}

	return PatchBytes(patchStr), nil
}

// converts overlay to JSON string to be inserted into the JSON Patch
func prepareJSONValue(overlay interface{}) string {
	jsonOverlay, err := json.Marshal(overlay)

	if err != nil || hasOnlyAnchors(overlay) {
		glog.V(3).Info(err)
		return ""
	}

	return string(jsonOverlay)
}

// Anchor has pattern value, so resource shouldn't be mutated with it
// If entire subtree has only anchor keys - we should skip inserting it
func hasOnlyAnchors(overlay interface{}) bool {
	switch typed := overlay.(type) {
	case map[string]interface{}:
		if anchors := getAnchorsFromMap(typed); len(anchors) == len(typed) {
			return true
		}

		for _, value := range typed {
			if !hasOnlyAnchors(value) {
				return false
			}
		}
	case []interface{}:
		for _, value := range typed {
			if !hasOnlyAnchors(value) {
				return false
			}
		}
	default:
		return false
	}

	return true
}

// Checks if subtree has anchors
func hasNestedAnchors(overlay interface{}) bool {
	switch typed := overlay.(type) {
	case map[string]interface{}:
		if anchors := getAnchorsFromMap(typed); len(anchors) > 0 {
			return true
		}

		for _, value := range typed {
			if hasNestedAnchors(value) {
				return true
			}
		}
		return false
	case []interface{}:
		for _, value := range typed {
			if hasNestedAnchors(value) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

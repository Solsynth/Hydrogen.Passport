package services

import (
	"reflect"
	"regexp"
	"strings"
)

func HasPermNode(held any, required any) bool {
	heldValue := reflect.ValueOf(held)
	requiredValue := reflect.ValueOf(required)

	switch heldValue.Kind() {
	case reflect.Int, reflect.Float64:
		if heldValue.Float() >= requiredValue.Float() {
			return true
		}
	case reflect.String:
		if heldValue.String() == requiredValue.String() {
			return true
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < heldValue.Len(); i++ {
			if reflect.DeepEqual(heldValue.Index(i).Interface(), required) {
				return true
			}
		}
	default:
		if reflect.DeepEqual(held, required) {
			return true
		}
	}

	return false
}

func FilterPermNodes(tree map[string]any, claims []string) map[string]any {
	filteredTree := make(map[string]any)

	match := func(claim, permission string) bool {
		regex := strings.Replace(permission, "*", ".*", -1)
		match, _ := regexp.MatchString("^"+regex+"$", claim)
		return match
	}

	for _, claim := range claims {
		for key, value := range tree {
			if match(claim, key) {
				filteredTree[key] = value
			}
		}
	}

	return filteredTree
}

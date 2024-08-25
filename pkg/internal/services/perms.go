package services

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func HasPermNode(perms map[string]any, requiredKey string, requiredValue any) bool {
	if heldValue, ok := perms[requiredKey]; ok {
		return ComparePermNode(heldValue, requiredValue)
	}
	return false
}

func HasPermNodeWithDefault(perms map[string]any, requiredKey string, requiredValue any, defaultValue any) bool {
	if heldValue, ok := perms[requiredKey]; ok {
		return ComparePermNode(heldValue, requiredValue)
	}
	return ComparePermNode(defaultValue, requiredValue)
}

func ComparePermNode(held any, required any) bool {
	heldValue := reflect.ValueOf(held)
	requiredValue := reflect.ValueOf(required)

	switch heldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if heldValue.Int() >= requiredValue.Int() {
			return true
		}
	case reflect.Float32, reflect.Float64:
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
		regex := strings.ReplaceAll(claim, "*", ".*")
		match, _ := regexp.MatchString(fmt.Sprintf("^%s$", regex), permission)
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

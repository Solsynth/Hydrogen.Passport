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
	isNumeric := func(val reflect.Value) bool {
		kind := val.Kind()
		return kind >= reflect.Int && kind <= reflect.Uint64 || kind >= reflect.Float32 && kind <= reflect.Float64
	}

	toFloat64 := func(val reflect.Value) float64 {
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(val.Uint())
		case reflect.Float32, reflect.Float64:
			return val.Float()
		default:
			panic(fmt.Sprintf("non-numeric value of kind %s", val.Kind()))
		}
	}

	heldValue := reflect.ValueOf(held)
	requiredValue := reflect.ValueOf(required)

	switch heldValue.Kind() {
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
		if isNumeric(heldValue) && isNumeric(requiredValue) {
			return toFloat64(heldValue) >= toFloat64(requiredValue)
		}

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

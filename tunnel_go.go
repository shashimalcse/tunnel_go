package tunnel_go

import (
	"encoding/json"
	"strings"
)

// Attribute represents an attribute with a name and a type.
type Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Property represents a property with an attribute, operator, and value.
type Property struct {
	Attribute Attribute `json:"attribute"`
	Operator  string    `json:"operator"`
	Value     []string  `json:"value"`
}

// Paths is a slice of slices of Property.
type Paths [][]Property

// ValidateTunnelPolicy validates the tunnel policy.
func ValidateTunnelPolicy(policyString string, inputString string) bool {
	var paths Paths
	err := json.Unmarshal([]byte(policyString), &paths)
	if err != nil {
		return false
	}

	var input map[string]interface{}
	err = json.Unmarshal([]byte(inputString), &input)
	if err != nil {
		return false
	}

	for _, path := range paths {
		allConditionsMet := true
		for _, property := range path {
			value := getValueFromInput(input, property.Attribute.Name)
			switch v := value.(type) {
			case string:
				if !evaluateSingleValue(v, property) {
					allConditionsMet = false
				}
			case []interface{}:
				stringSlice := make([]string, len(v))
				for i, iv := range v {
					if str, ok := iv.(string); ok {
						stringSlice[i] = str
					} else {
						allConditionsMet = false
						break
					}
				}
				if !evaluateMultipleValues(stringSlice, property) {
					allConditionsMet = false
				}
			default:
				allConditionsMet = false
			}
		}
		if allConditionsMet {
			return true
		}
	}
	return false
}

// GetValueFromInput retrieves a value from the input based on the attribute name.
func getValueFromInput(input map[string]interface{}, attributeName string) interface{} {
	keys := strings.Split(attributeName, ".")
	var current interface{} = input
	for _, key := range keys {
		if currentMap, ok := current.(map[string]interface{}); ok {
			current = currentMap[key]
		} else {
			return nil
		}
	}
	return current
}

// evaluateSingleValue evaluates the condition for a single string value.
func evaluateSingleValue(value string, property Property) bool {
	switch property.Operator {
	case "equal":
		return value == property.Value[0]
	case "not_equal":
		return value != property.Value[0]
	}
	return false
}

// evaluateMultipleValues evaluates the condition for a slice of string values.
func evaluateMultipleValues(values []string, property Property) bool {
	switch property.Operator {
	case "contains":
		return contains(values, property.Value[0])
	case "not_contains":
		return !contains(values, property.Value[0])
	case "contain_at_least_one":
		for _, val := range property.Value {
			if contains(values, val) {
				return true
			}
		}
	case "not_contain_at_least_one":
		for _, val := range property.Value {
			if !contains(values, val) {
				return true
			}
		}
	}
	return false
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

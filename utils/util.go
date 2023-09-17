package utils

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseInterface[T any](i interface{}, v *T) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}

func Contains[T comparable](arr []T, el T) bool {
	for _, element := range arr {
		if element == el {
			return true
		}
	}
	return false
}

func ToBytes[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

func CompareInterfaces(a, b interface{}) bool {
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)

	if a == nil && b == nil {
		return true
	}

	if typeA != typeB {
		return false
	}

	switch a.(type) {
	case int, string, bool, float64:
		return a == b
	case []interface{}:
		sliceA := a.([]interface{})
		sliceB := b.([]interface{})

		if len(sliceA) != len(sliceB) {
			return false
		}

		for i := 0; i < len(sliceA); i++ {
			if !CompareInterfaces(sliceA[i], sliceB[i]) {
				return false
			}
		}

		return true
	case map[string]interface{}:
		mapA := a.(map[string]interface{})
		mapB := b.(map[string]interface{})

		if len(mapA) != len(mapB) {
			return false
		}

		for key, valueA := range mapA {
			if valueA != nil && valueA.(string) == "any" {
				continue // Skip keys and values with "any"
			}

			if valueB, exists := mapB[key]; !exists || !CompareInterfaces(valueA, valueB) {
				return false
			}
		}

		return true
	default:
		// Handle other custom types if needed
		return false
	}
}

func CompareStrings(a *string, e *string) error {
	if *e == "" {
		return nil
	}

	var actual = parseString(a)
	var expected = parseString(e)

	if !CompareInterfaces(expected, actual) {
		return errors.New("Not matching")
	}

	return nil
}

func parseString(a *string) interface{} {
	var val interface{}
	if err := json.Unmarshal([]byte(*a), &val); err == nil && val != nil {
		switch val.(type) {
		case []interface{}:
			val = val.([]interface{})
		case map[string]interface{}:
			val = val.(interface{})
		default:
			val = val.(string)
		}
	} else {
		val = *a
	}

	return val
}

func MakeObjectIds(ids []string) ([]primitive.ObjectID, error) {
	objectIds := make([]primitive.ObjectID, len(ids))
	for i, oid := range ids {
		id, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return nil, errors.New("Invalid object id passed -> " + oid)
		}
		objectIds[i] = id
	}

	return objectIds, nil
}
package utils

import "encoding/json"

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

package dataformat

import (
	"encoding/json"
)

// JSON provides Marshal (employs json.Marshal) / Unmarshal (employs json.Unmarshal) functions.
type JSON struct{}

func (JSON) Unmarshal(data []byte, targetType any) (any, error) {
	var target = newInstanceOfType(targetType)

	err := json.Unmarshal(data, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (JSON) Marshal(data any) (string, error) {
	v, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

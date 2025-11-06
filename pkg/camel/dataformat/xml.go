package dataformat

import (
	"encoding/xml"
)

// XML provides Marshal (employs xml.Marshal) / Unmarshal (employs xml.Unmarshal) functions.
type XML struct{}

func (XML) Unmarshal(data []byte, targetType any) (any, error) {
	var target = newInstanceOfType(targetType)

	err := xml.Unmarshal(data, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (XML) Marshal(data any) (string, error) {
	v, err := xml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

package dataformat

import (
	"bytes"
	"encoding/json"
)

// JSONFormat implementation for JSON using the standard encoding/json library.
// For maximum efficiency, we use Decoder/Encoder with buffers instead of Marshal/Unmarshal,
// to minimize allocations (especially in Read/Marshal on large data).
type JSONFormat struct{}

// Unmarshal deserializes []byte into new instance of targetType and returns data or error.
// If targetType is not a pointer, creates a new pointer.
// Using json.Decoder for stream-parsing (more effective for big data).
func (JSONFormat) Unmarshal(data []byte, targetType any) (any, error) {
	var target = newInstanceOfType(targetType)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber() // avoid float precision loss
	err := dec.Decode(target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// Marshal serializes data to JSON.
// Using json.Encoder for stream-serialization (more effective than Marshal for big data).
func (JSONFormat) Marshal(data any) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024)) // initial buffer for reducing realloc
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // disable escape
	err := enc.Encode(data)
	return "", err
}

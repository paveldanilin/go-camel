package dataformat

import (
	"bytes"
	"encoding/xml"
)

// XMLFormat implementation for XML using the standard encoding/xml library.
// For efficiency, we use Decoder/Encoder with buffers (stream parsing).
type XMLFormat struct{}

// Unmarshal deserializes []byte into a new instance of targetType and returns it.
// Similar to JSON: create a pointer if needed. Use xml.Decoder.
func (XMLFormat) Unmarshal(data []byte, targetType any) (any, error) {
	var target = newInstanceOfType(targetType)

	dec := xml.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// Marshal serializes data into XML, discarding []byte to match the signature.
// We use xml.Encoder for stream serialization.
func (XMLFormat) Marshal(data any) error {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	enc := xml.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return err
	}
	return enc.Flush()
}

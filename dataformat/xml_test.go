package dataformat

import (
	"reflect"
	"testing"
)

func TestXMLFormat_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		targetType  any
		wantErr     bool
		wantNonZero bool
	}{
		{
			name:        "Valid struct",
			data:        []byte(`{"name":"Alice","age":30}`),
			targetType:  person{},
			wantErr:     false,
			wantNonZero: true,
		},
		{
			name:        "Valid with pointer",
			data:        []byte(`{"name":"Bob","age":25}`),
			targetType:  &person{},
			wantErr:     false,
			wantNonZero: true,
		},
		{
			name:       "Invalid JSON",
			data:       []byte(`invalid`),
			targetType: person{},
			wantErr:    true,
		},
		{
			name:       "Empty data",
			data:       []byte{},
			targetType: person{},
			wantErr:    true,
		},
		{
			name:        "Complex data",
			data:        []byte(`{"name":"Charlie","age":40,"friends":["Dave"],"nestedList":[{"item":"first"}]}`),
			targetType:  person{},
			wantErr:     false,
			wantNonZero: true,
		},
	}

	jsonFmt := &JSONFormat{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonFmt.Unmarshal(tt.data, tt.targetType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.wantNonZero && reflect.DeepEqual(got, reflect.Zero(reflect.TypeOf(tt.targetType)).Interface()) {
				t.Errorf("Unmarshal() returned zero value, want non-zero")
			}
		})
	}
}

func TestXMLFormat_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name:    "Valid struct",
			data:    person{Name: "Alice", Age: 30},
			wantErr: false,
		},
		{
			name:    "Nil data",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "Invalid data (unserializable)",
			data:    func() {},
			wantErr: true,
		},
		{
			name:    "Complex data",
			data:    person{Name: "Bob", Friends: []string{"Dave"}, NestedList: []item{{Item: "first"}}},
			wantErr: false,
		},
	}

	xmlFmt := &XMLFormat{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := xmlFmt.Marshal(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

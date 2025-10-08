package dataformat

import (
	"reflect"
	"testing"
	"time"
)

type person struct {
	Name       string    `json:"name" xml:"Name"`
	Age        int       `json:"age" xml:"Age"`
	Friends    []string  `json:"friends" xml:"Friends>Friend"`
	NestedList []item    `json:"nestedList" xml:"NestedList>Item"`
	Birthday   time.Time `json:"birthday" xml:"Birthday"`
}

type item struct {
	Item string `json:"item" xml:"Item"`
}

func TestJSONFormat_Unmarshal(t *testing.T) {
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

func TestJSONFormat_Marshal(t *testing.T) {
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
			wantErr: false, // json.Encode nil to null without error
		},
		{
			name:    "Invalid data (unserializable)",
			data:    func() {}, // Func not serializable
			wantErr: true,
		},
		{
			name:    "Complex data",
			data:    person{Name: "Bob", Friends: []string{"Dave"}, NestedList: []item{{Item: "first"}}},
			wantErr: false,
		},
	}

	jsonFmt := &JSONFormat{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsonFmt.Marshal(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

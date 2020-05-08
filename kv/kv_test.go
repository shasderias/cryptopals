package kv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshal(t *testing.T) {
	type simpleStruct struct {
		Foo string
		Baz string
		Zap string
	}

	type mixedIntAndStringStruct struct {
		Foo string
		Baz int
	}

	tests := []struct {
		name      string
		data      []byte
		v         interface{}
		expectedV interface{}
		wantErr   bool
	}{
		{
			name: "simple",
			data: []byte("foo=bar&baz=qux&zap=zazzle"),
			v:    &simpleStruct{},
			expectedV: &simpleStruct{
				Foo: "bar",
				Baz: "qux",
				Zap: "zazzle",
			},
			wantErr: false,
		},
		{
			name: "mixedIntAndString",
			data: []byte("foo=bar&baz=3"),
			v:    &mixedIntAndStringStruct{},
			expectedV: &mixedIntAndStringStruct{
				Foo: "bar",
				Baz: 3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.data, tt.v)
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.v, tt.expectedV); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "simple",
			v: struct {
				foo string
				baz int
			}{
				foo: "bar",
				baz: 3,
			},
			want:    "foo=bar&baz=3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("got %s; want %s", string(got), tt.want)
			}
		})
	}
}

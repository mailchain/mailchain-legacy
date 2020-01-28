package output

import (
	"reflect"
	"testing"
)

func TestElement_IsDefault(t *testing.T) {
	type fields struct {
		Name       string
		Attributes []Attribute
		Elements   []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"non-default-elements",
			fields{
				"name",
				[]Attribute{{}, {}},
				[]Element{
					{}, {Attributes: []Attribute{{}, {IsDefault: true}}},
				},
			},
			false,
		},
		{
			"non-default-attributes",
			fields{
				"name",
				[]Attribute{{}, {IsDefault: true}},
				[]Element{},
			},
			false,
		},
		{
			"default-attributes",
			fields{
				"name",
				[]Attribute{{IsDefault: true}},
				[]Element{{}},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Element{
				FullName:   tt.fields.Name,
				Attributes: tt.fields.Attributes,
				Elements:   tt.fields.Elements,
			}
			if got := e.IsDefault(); got != tt.want {
				t.Errorf("Element.IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttribute_ShortName(t *testing.T) {
	type fields struct {
		FullName          string
		IsDefault         bool
		AdditionalComment string
		Value             interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"short",
			fields{
				FullName: "short",
			},
			"short",
		},
		{
			"multi",
			fields{
				FullName: "key.sub-key.sub-sub-key",
			},
			"sub-sub-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attribute{
				FullName:          tt.fields.FullName,
				IsDefault:         tt.fields.IsDefault,
				AdditionalComment: tt.fields.AdditionalComment,
				Value:             tt.fields.Value,
			}
			if got := a.ShortName(); got != tt.want {
				t.Errorf("Attribute.ShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_ShortName(t *testing.T) {
	type fields struct {
		FullName   string
		Attributes []Attribute
		Elements   []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"simple",
			fields{FullName: "simple"},
			"simple",
		},
		{
			"sub",
			fields{FullName: "root.sub-key.sub-sub-key"},
			"sub-sub-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Element{
				FullName:   tt.fields.FullName,
				Attributes: tt.fields.Attributes,
				Elements:   tt.fields.Elements,
			}
			if got := e.ShortName(); got != tt.want {
				t.Errorf("Element.ShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElement_SortedAttributes(t *testing.T) {
	type fields struct {
		FullName   string
		Attributes []Attribute
		Elements   []Element
	}
	tests := []struct {
		name   string
		fields fields
		want   []Attribute
	}{
		{
			"success",
			fields{
				"simple",
				[]Attribute{{FullName: "D"}, {FullName: "F"}, {FullName: "A"}},
				nil,
			},
			[]Attribute{{FullName: "A"}, {FullName: "D"}, {FullName: "F"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Element{
				FullName:   tt.fields.FullName,
				Attributes: tt.fields.Attributes,
				Elements:   tt.fields.Elements,
			}
			if got := e.SortedAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Element.SortedAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

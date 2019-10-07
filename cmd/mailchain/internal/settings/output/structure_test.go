package output

import "testing"

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
				[]Attribute{Attribute{}, Attribute{}},
				[]Element{
					Element{}, Element{Attributes: []Attribute{Attribute{}, Attribute{IsDefault: true}}},
				},
			},
			false,
		},
		{
			"non-default-attributes",
			fields{
				"name",
				[]Attribute{Attribute{}, Attribute{IsDefault: true}},
				[]Element{},
			},
			false,
		},
		{
			"default-attributes",
			fields{
				"name",
				[]Attribute{Attribute{}},
				[]Element{Element{}},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Element{
				Name:       tt.fields.Name,
				Attributes: tt.fields.Attributes,
				Elements:   tt.fields.Elements,
			}
			if got := e.IsDefault(); got != tt.want {
				t.Errorf("Element.IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

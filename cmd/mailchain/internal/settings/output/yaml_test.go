package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNamePortion(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name            string
		isDefault       bool
		trailingNewLine bool
		commentDefaults bool
		excludeDefaults bool
		tabsize         int
		indent          int
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"new-line",
			args{
				"property",
				false,
				true,
				true,
				false,
				2,
				1,
			},
			"  property:\n",
		},
		{
			"no-new-line",
			args{
				"property",
				false,
				false,
				true,
				false,
				2,
				1,
			},
			"  property:",
		},
		{
			"default-include-defaults",
			args{
				"property",
				true,
				false,
				false,
				false,
				2,
				1,
			},
			"  property:",
		},
		{
			"default-comment-defaults",
			args{
				"property",
				true,
				false,
				true,
				false,
				2,
				1,
			},
			"#   property:",
		},
		{
			"default-excluded",
			args{
				"property",
				true,
				false,
				true,
				true,
				2,
				1,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			createNamePortion(tt.args.name, tt.args.isDefault, tt.args.trailingNewLine, tt.args.commentDefaults, tt.args.excludeDefaults, out, tt.args.tabsize, tt.args.indent)
			if gotOut := out.String(); !assert.Equal(tt.wantOut, gotOut) {
				t.Errorf("createNamePortion() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_yamlAttributeDefault(t *testing.T) {
	type args struct {
		a               Attribute
		tabsize         int
		indent          int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"value-set-string",
			args{
				Attribute{FullName: "key.sub-key", IsDefault: false, Value: "test-value"},
				2,
				2,
				true,
				false,
			},
			"    sub-key: \"test-value\"\n",
		},
		{
			"value-set-int",
			args{
				Attribute{FullName: "key.sub-key", IsDefault: false, Value: 100},
				2,
				2,
				true,
				false,
			},
			"    sub-key: 100\n",
		},
		{
			"comment-default-value",
			args{
				Attribute{FullName: "key.sub-key", IsDefault: true, Value: "test-value"},
				2,
				2,
				true,
				false,
			},
			"#     sub-key: \"test-value\"\n",
		},
		{
			"exclude-default-value",
			args{
				Attribute{FullName: "key.sub-key", IsDefault: true, Value: "test-value"},
				2,
				2,
				true,
				true,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			yamlAttributeDefault(tt.args.a, out, tt.args.tabsize, tt.args.indent, tt.args.commentDefaults, tt.args.excludeDefaults)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("yamlAttributeDefault() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_yamlValueFormat(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"string",
			args{"value"},
			" %q",
		},
		{
			"string slice",
			args{[]string{"value1", "value2"}},
			" %v",
		},
		{
			"int",
			args{2},
			" %v",
		},
		{
			"bool",
			args{true},
			" %v",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yamlValueFormat(tt.args.val); got != tt.want {
				t.Errorf("yamlValueFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_yamlAttributeStringSlice(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		fullName        string
		val             []string
		isDefault       bool
		tabsize         int
		indent          int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"empty",
			args{
				"key.subkey",
				[]string{},
				false,
				2,
				2,
				true,
				false,
			},
			"    subkey: []\n",
		},
		{
			"exclude-default",
			args{
				"key.subkey",
				[]string{"value1", "value2"},
				true,
				2,
				2,
				true,
				true,
			},
			"",
		},
		{
			"set-value",
			args{
				"key.subkey",
				[]string{"value1", "value2"},
				false,
				2,
				2,
				true,
				false,
			},
			"    subkey:\n      - \"value1\"\n      - \"value2\"\n",
		},
		{
			"comment-default-value",
			args{
				"key.subkey",
				[]string{"value1", "value2"},
				true,
				2,
				2,
				true,
				false,
			},
			"#     subkey:\n#       - \"value1\"\n#       - \"value2\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			yamlAttributeStringSlice(tt.args.fullName, tt.args.val, tt.args.isDefault, out, tt.args.tabsize, tt.args.indent, tt.args.commentDefaults, tt.args.excludeDefaults)
			if gotOut := out.String(); !assert.Equal(tt.wantOut, gotOut) {
				t.Errorf("yamlAttributeStringSlice() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_yamlAttribute(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		a               Attribute
		tabsize         int
		indent          int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"string-slice",
			args{
				Attribute{FullName: "key.subkey", IsDefault: false, Value: []string{"value1", "value2"}},
				2,
				0,
				true,
				true,
			},
			"subkey:\n  - \"value1\"\n  - \"value2\"\n",
		},
		{
			"string",
			args{
				Attribute{FullName: "key.subkey", IsDefault: false, Value: "value"},
				2,
				0,
				true,
				true,
			},
			"subkey: \"value\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			yamlAttribute(tt.args.a, out, tt.args.tabsize, tt.args.indent, tt.args.commentDefaults, tt.args.excludeDefaults)
			if gotOut := out.String(); !assert.Equal(tt.wantOut, gotOut) {
				t.Errorf("yamlAttribute() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_yamlElement(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		e               Element
		tabsize         int
		indent          int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"full",
			args{
				Element{
					"element.parent-element",
					[]Attribute{
						{FullName: "key1", IsDefault: false, Value: []string{"value1", "value2"}},
						{FullName: "key2", IsDefault: false, Value: "value"},
					},
					[]Element{
						{
							FullName: "sub-element",
							Attributes: []Attribute{
								{FullName: "subkey2", IsDefault: false, Value: "value"},
							},
						},
					},
				},
				2, 0, true, false,
			},
			"parent-element:\n  key1:\n    - \"value1\"\n    - \"value2\"\n  key2: \"value\"\n  sub-element:\n    subkey2: \"value\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			yamlElement(tt.args.e, out, tt.args.tabsize, tt.args.indent, tt.args.commentDefaults, tt.args.excludeDefaults)
			if gotOut := out.String(); !assert.Equal(tt.wantOut, gotOut) {
				t.Errorf("yamlElement() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestToYaml(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		root            Root
		tabsize         int
		commentDefaults bool
		excludeDefaults bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			"full",
			args{
				Root{
					[]Element{
						{
							FullName: "sub-element",
							Attributes: []Attribute{
								{FullName: "subkey2", IsDefault: false, Value: "value"},
							},
						},
					},
				},
				2, true, false,
			},
			"sub-element:\n  subkey2: \"value\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			ToYaml(tt.args.root, out, tt.args.tabsize, tt.args.commentDefaults, tt.args.excludeDefaults)
			if gotOut := out.String(); !assert.Equal(tt.wantOut, gotOut) {
				t.Errorf("ToYaml() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

package output

import (
	"fmt"
	"io"
	"strings"
)

// ToYaml create a yaml output of the configuration. Tab size, comment out default values, or exclude default values are configurable options.
func ToYaml(root Root, out io.Writer, tabsize int, commentDefaults, excludeDefaults bool) {
	for _, em := range root.Elements {
		yamlElement(em, out, tabsize, 0, commentDefaults, excludeDefaults)
	}
}

func yamlElement(e Element, out io.Writer, tabsize, indent int, commentDefaults, excludeDefaults bool) {
	dots := strings.Split(e.FullName, ".")
	shortKey := dots[len(dots)-1]
	createNamePortion(shortKey, e.IsDefault(), true, commentDefaults, excludeDefaults, out, tabsize, indent)

	for _, a := range e.SortedAttributes() {
		yamlAttribute(a, out, tabsize, indent+1, commentDefaults, excludeDefaults)
	}

	for _, em := range e.SortedElements() {
		yamlElement(em, out, tabsize, indent+1, commentDefaults, excludeDefaults)
	}
}

func yamlAttribute(a Attribute, out io.Writer, tabsize, indent int, commentDefaults, excludeDefaults bool) {
	switch value := a.Value.(type) {
	case []string:
		yamlAttributeStringSlice(a.ShortName(), value, a.IsDefault, out, tabsize, indent, commentDefaults, excludeDefaults)
	default:
		yamlAttributeDefault(a, out, tabsize, indent, commentDefaults, excludeDefaults)
	}
}

func yamlAttributeDefault(a Attribute, out io.Writer, tabsize, indent int, commentDefaults, excludeDefaults bool) {
	if excludeDefaults && a.IsDefault {
		return
	}

	createNamePortion(a.ShortName(), a.IsDefault, false, commentDefaults, excludeDefaults, out, tabsize, indent)
	fmt.Fprintf(out, yamlValueFormat(a.Value), a.Value)
	fmt.Fprint(out, "\n")
}

func yamlAttributeStringSlice(fullName string, val []string, isDefault bool, out io.Writer, tabsize, indent int, commentDefaults, excludeDefaults bool) {
	dots := strings.Split(fullName, ".")
	shortKey := dots[len(dots)-1]

	if isDefault && excludeDefaults {
		return
	}

	if len(val) == 0 {
		createNamePortion(shortKey, isDefault, false, commentDefaults, excludeDefaults, out, tabsize, indent)
		fmt.Fprint(out, " []\n")

		return
	}

	createNamePortion(shortKey, isDefault, true, commentDefaults, excludeDefaults, out, tabsize, indent)

	for _, item := range val {
		itemFormat := "%s- %q\n"
		if commentDefaults && isDefault {
			itemFormat = "# " + itemFormat
		}

		fmt.Fprintf(out, itemFormat, strings.Repeat(" ", tabsize*(indent+1)), item)
	}
}

func createNamePortion(name string, isDefault, trailingNewLine, commentDefaults, excludeDefaults bool, out io.Writer, tabsize, indent int) {
	if isDefault && excludeDefaults {
		return
	}

	if isDefault && commentDefaults {
		fmt.Fprintf(out, "# ")
	}

	fmt.Fprintf(out, "%s%s:", strings.Repeat(" ", tabsize*indent), name)

	if trailingNewLine {
		fmt.Fprintf(out, "\n")
	}
}

func yamlValueFormat(val interface{}) string {
	switch val.(type) {
	case string:
		return " %q"
	default:
		return " %v"
	}
}

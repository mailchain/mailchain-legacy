package output

import (
	"sort"
	"strings"
)

// Root element for that contains only `[]Element` structs.
type Root struct {
	Elements []Element
}

// Element groups settings and child groups of settings.
type Element struct {
	FullName   string
	Attributes []Attribute
	Elements   []Element
}

// SortedElements to ensure that `[]Element` slice is sorted consistently.
func (e Element) SortedElements() []Element {
	sort.Slice(e.Elements, func(i, j int) bool {
		return e.Elements[i].FullName < e.Elements[j].FullName
	})

	return e.Elements
}

// SortedAttributes to ensure that `[]Attribute` slice is sorted consistently.
func (e Element) SortedAttributes() []Attribute {
	sort.Slice(e.Attributes, func(i, j int) bool {
		return e.Attributes[i].FullName < e.Attributes[j].FullName
	})

	return e.Attributes
}

// ShortName of the element split by ".".
func (e Element) ShortName() string {
	dots := strings.Split(e.FullName, ".")

	return dots[len(dots)-1]
}

// IsDefault checks the entire Element and child Elements and Attributes to deterimine if all values are default.
func (e Element) IsDefault() bool {
	for _, i := range e.Elements {
		if !i.IsDefault() {
			return false
		}
	}

	for _, i := range e.Attributes {
		if !i.IsDefault {
			return false
		}
	}

	return true
}

// Attribute contains setting of a single value.
type Attribute struct {
	FullName          string
	IsDefault         bool
	AdditionalComment string
	Value             interface{}
}

// ShortName of an attributes FullName split by ".".
func (a Attribute) ShortName() string {
	dots := strings.Split(a.FullName, ".")
	return dots[len(dots)-1]
}

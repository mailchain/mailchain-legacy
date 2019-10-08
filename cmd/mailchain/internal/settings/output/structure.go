package output

import (
	"sort"
	"strings"
)

type Root struct {
	Elements []Element
}

type Element struct {
	FullName   string
	Attributes []Attribute
	Elements   []Element
}

func (e Element) SortedElements() []Element {
	sort.Slice(e.Elements, func(i, j int) bool {
		return e.Elements[i].FullName < e.Elements[j].FullName
	})
	return e.Elements
}

func (e Element) SortedAttributes() []Attribute {
	sort.Slice(e.Attributes, func(i, j int) bool {
		return e.Attributes[i].FullName < e.Attributes[j].FullName
	})
	return e.Attributes
}

func (e Element) ShortName() string {
	dots := strings.Split(e.FullName, ".")
	return dots[len(dots)-1]
}

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

type Attribute struct {
	FullName          string
	IsDefault         bool
	AdditionalComment string
	Value             interface{}
}

func (a Attribute) ShortName() string {
	dots := strings.Split(a.FullName, ".")
	return dots[len(dots)-1]
}

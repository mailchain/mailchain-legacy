package output

type Root struct {
	Elements []Element
}

type Element struct {
	Name       string
	Attributes []Attribute
	Elements   []Element
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
	Name              string
	IsDefault         bool
	AdditionalComment string
	Value             interface{}
}

package gousmap

import "fmt"

type UsmapSchema struct {
	Name       string
	SuperType  *string
	PropCount  uint16
	Properties []*UsmapProperty
}

// Returns UsmapSchema representation as string.
func (cls *UsmapSchema) ToString() string {
	name := cls.Name
	if cls.SuperType != nil {
		name += ": " + *cls.SuperType
	}

	return fmt.Sprintf("%s | %d properties", name, len(cls.Properties))
}

// Returns a list containing the properties of the current UsmapSchema.
func (cls *UsmapSchema) GetProps() []string {
	props := make([]string, len(cls.Properties))
	for i := range cls.Properties {
		props[i] = cls.Properties[i].ToString()
	}

	return props
}

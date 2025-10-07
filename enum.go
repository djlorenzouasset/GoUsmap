package gousmap

import (
	"fmt"
	"strconv"
)

type UsmapEnum struct {
	Name  string
	Names []string
}

// Returns UsmapEnum representation as string.
func (cls *UsmapEnum) ToString() string {
	return cls.Name + " | " + strconv.Itoa(len(cls.Names))
}

// Returns a list of enum members of the current UsmapEnum.
func (cls *UsmapEnum) GetValues() []string {
	values := make([]string, len(cls.Names))
	for i := range cls.Names {
		values[i] = fmt.Sprintf("%s::%s", cls.Name, cls.Names[i])
	}

	return values
}
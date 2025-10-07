package gousmap

import "fmt"

type UsmapProperty struct {
	Name      string
	Data      *UsmapPropertyData
	SchemaIdx uint16
	ArraySize byte
}

// Returns UsmapProperty representation as string.
func (cls *UsmapProperty) ToString() string {
	return fmt.Sprintf("%s | %v", cls.Name, cls.Data.Type.ToString())
}
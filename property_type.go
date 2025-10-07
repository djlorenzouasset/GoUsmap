package gousmap

import "fmt"

type EUsmapPropertyType uint8
const (
	ByteProperty EUsmapPropertyType = iota
	BoolProperty
	IntProperty
	FloatProperty
	ObjectProperty
	NameProperty
	DelegateProperty
	DoubleProperty
	ArrayProperty
	StructProperty
	StrProperty
	TextProperty
	InterfaceProperty
	MulticastDelegateProperty
	WeakObjectProperty
	LazyObjectProperty
	AssetObjectProperty
	SoftObjectProperty
	UInt64Property
	UInt32Property
	UInt16Property
	Int64Property
	Int16Property
	Int8Property
	MapProperty
	SetProperty
	EnumProperty
	FieldPathProperty
	OptionalProperty
	Utf8StrProperty
	AnsiStrProperty
	Unknown EUsmapPropertyType = 255
)

// Returns a string representation of EUsmapPropertyType.
func (t EUsmapPropertyType) ToString() string {
	switch t {
	case ByteProperty:
		return "ByteProperty"
	case BoolProperty:
		return "BoolProperty"
	case IntProperty:
		return "IntProperty"
	case FloatProperty:
		return "FloatProperty"
	case ObjectProperty:
		return "ObjectProperty"
	case NameProperty:
		return "NameProperty"
	case DelegateProperty:
		return "DelegateProperty"
	case DoubleProperty:
		return "DoubleProperty"
	case ArrayProperty:
		return "ArrayProperty"
	case StructProperty:
		return "StructProperty"
	case StrProperty:
		return "StrProperty"
	case TextProperty:
		return "TextProperty"
	case InterfaceProperty:
		return "InterfaceProperty"
	case MulticastDelegateProperty:
		return "MulticastDelegateProperty"
	case WeakObjectProperty:
		return "WeakObjectProperty"
	case LazyObjectProperty:
		return "LazyObjectProperty"
	case AssetObjectProperty:
		return "AssetObjectProperty"
	case SoftObjectProperty:
		return "SoftObjectProperty"
	case UInt64Property:
		return "UInt64Property"
	case UInt32Property:
		return "UInt32Property"
	case UInt16Property:
		return "UInt16Property"
	case Int64Property:
		return "Int64Property"
	case Int16Property:
		return "Int16Property"
	case Int8Property:
		return "Int8Property"
	case MapProperty:
		return "MapProperty"
	case SetProperty:
		return "SetProperty"
	case EnumProperty:
		return "EnumProperty"
	case FieldPathProperty:
		return "FieldPathProperty"
	case OptionalProperty:
		return "OptionalProperty"
	case Utf8StrProperty:
		return "Utf8StrProperty"
	case AnsiStrProperty:
		return "AnsiStrProperty"
	case Unknown:
		return "Unknown"
	default:
		return fmt.Sprintf("EUsmapPropertyType(%d)", t)
	}
}

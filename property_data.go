package gousmap

type UsmapPropertyData struct {
	Type       EUsmapPropertyType
	EnumName   string
	StructType string
	InnerType  *UsmapPropertyData
	ValueType  *UsmapPropertyData
}

func Deserialize(reader *UsmapReader, names []string) *UsmapPropertyData {
	propTypeByte, _ := reader.ReadByte()
	propType := EUsmapPropertyType(propTypeByte)

	prop := &UsmapPropertyData{ Type: propType }
	switch propType {
	case EnumProperty:
		prop.InnerType = Deserialize(reader, names)
		nameIdx, _ := reader.ReadUint32()
		prop.EnumName = names[nameIdx]
	case StructProperty:
		structTypeIdx, _ := reader.ReadUint32()
		prop.StructType = names[structTypeIdx]
	case SetProperty, ArrayProperty, OptionalProperty:
		prop.InnerType = Deserialize(reader, names)
	case MapProperty:
		prop.InnerType = Deserialize(reader, names)
		prop.ValueType = Deserialize(reader, names)
	}

	return prop
}
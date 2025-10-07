package gousmap

type EUsmapVersion byte
const (
	EUsmapVersionInitial EUsmapVersion = iota
	// Adds package versioning to aid with compatibility
	EUsmapVersionPackageVersioning
	// 16-bit wide name-lengths (ushort/uint16)
	EUsmapVersionLongFName
	// Enums with more than 255 values
	EUsmapVersionLargeEnums
	// Support for explicit enum values
	EUsmapExplicitEnumValues
	// Support for Utf8StrProperty/AnsiStrProperty
	EUsmapUtf8AndAnsiStrProps
	EUsmapVersionLatest
)

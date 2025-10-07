package gousmap

import "fmt"

type ECompressionMethod byte
const (
	// No compression method
	ECompressionMethodNone ECompressionMethod = iota
	// Data is Oodle compressed
	ECompressionMethodOodle
	// Data is Brotli compressed
	ECompressionMethodBrotli
	// Data is ZStandard compressed
	ECompressionMethodZStandard
	ECompressionMethodMax
)

// Returns a string representation of ECompressionMethod.
func (t ECompressionMethod) ToString() string {
	switch t {
	case ECompressionMethodNone:
		return "None"
	case ECompressionMethodOodle:
		return "Oodle"
	case ECompressionMethodBrotli:
		return "Brotli"
	case ECompressionMethodZStandard:
		return "ZStandard"
	case ECompressionMethodMax:
		return fmt.Sprintf("ECompressionMethod_MAX(%d)", t)
	default:
		return fmt.Sprintf("ECompressionMethod(%d)", t)
	}
}
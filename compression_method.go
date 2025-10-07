package gousmap

type ECompressionMethod byte
const (
	ECompressionMethodNone ECompressionMethod = iota
	ECompressionMethodOodle
	ECompressionMethodBrotli
	ECompressionMethodZStandard
	ECompressionMethodMax
)

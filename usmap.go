// GoUsmap provides utilities for parsing .usmap files into
// structured and usable Go data types.
//
// This package was created for learning purposes and to explore
// both the Go language and the .usmap file format.
package gousmap

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// Usmap file magic.
const UsmapMagic uint16 = 0x30C4

type Usmap struct {
	CompressionMethod ECompressionMethod

	Names   []string
	Enums   []*UsmapEnum
	Schemas []*UsmapSchema

	reader   *UsmapReader
	oodleDec *Oodle
}

func (ump *Usmap) parseInternal(longFNames bool, largeEnums bool) error {
	size, _ := ump.reader.ReadUint32()
	ump.Names = make([]string, size)
	for i := 0; i < int(size); i++ {
		var nameLength int
		if longFNames {
			length, _ := ump.reader.ReadUint16()
			nameLength = int(length)
		} else {
			length, _ := ump.reader.ReadByte()
			nameLength = int(length)
		}

		name, _ := ump.reader.ReadString(nameLength)
		ump.Names[i] = name
	}

	size, _ = ump.reader.ReadUint32()
	ump.Enums = make([]*UsmapEnum, size)
	for i := 0; i < int(size); i++ {
		idx, _ := ump.reader.ReadUint32()
		enumName := ump.Names[idx]

		var enumNamesSize int
		if largeEnums {
			length, _ := ump.reader.ReadUint16()
			enumNamesSize = int(length)
		} else {
			length, _ := ump.reader.ReadByte()
			enumNamesSize = int(length)
		}

		enumNames := make([]string, enumNamesSize)
		for n := 0; n < enumNamesSize; n++ {
			nameIdx, _ := ump.reader.ReadUint32()
			enumNames[n] = ump.Names[nameIdx]
		}

		ump.Enums[i] = &UsmapEnum{
			Name:  enumName,
			Names: enumNames,
		}
	}

	size, _ = ump.reader.ReadUint32()
	ump.Schemas = make([]*UsmapSchema, size)
	for i := 0; i < int(size); i++ {
		idx, _ := ump.reader.ReadUint32()
		schemaName := ump.Names[idx]

		var schemaSuperType *string
		superIdx, _ := ump.reader.ReadUint32()
		if superIdx == math.MaxUint32 {
			schemaSuperType = nil
		} else {
			schemaSuperType = &ump.Names[superIdx]
		}

		propCount, _ := ump.reader.ReadUint16()
		serializablePropCount, _ := ump.reader.ReadUint16()

		props := make([]*UsmapProperty, serializablePropCount)
		for p := 0; p < int(serializablePropCount); p++ {
			schemaIdx, _ := ump.reader.ReadUint16()
			arraySize, _ := ump.reader.ReadByte()
			nameIdx, _ := ump.reader.ReadUint32()

			prop := Deserialize(ump.reader, ump.Names)
			props[p] = &UsmapProperty{
				Name:      ump.Names[nameIdx],
				Data:      prop,
				SchemaIdx: schemaIdx,
				ArraySize: arraySize,
			}
		}

		ump.Schemas[i] = &UsmapSchema{
			Name:       schemaName,
			SuperType:  schemaSuperType,
			PropCount:  propCount,
			Properties: props,
		}
	}

	return nil
}

// Returns a summary of the Usmap contents as a formatted string.
func (ump *Usmap) ToString() string {
	return fmt.Sprintf(
		"Names: %d | Enums: %d | Schemas: %d",
		len(ump.Names), len(ump.Enums), len(ump.Schemas),
	)
}

// Parses a .usmap file from a raw byte slice.
//
// The function reads and validates the .usmap header, determines the compression
// method, and decompresses the data if necessary. It supports Oodle, Brotli, and
// Zstandard compression methods.
//
// If an Oodle-compressed file is encountered, oodleInstance must not be nil.
// It returns a pointer to a Usmap structure containing the parsed data, or an
// error if the file is invalid or decompression fails.
func ParseFromBytes(buffer []byte, oodleInstance *Oodle) (*Usmap, error) {
	reader, err := CreateReader(buffer)
	if err != nil {
		return nil, err
	}

	usmap := &Usmap{}
	if oodleInstance != nil {
		usmap.oodleDec = oodleInstance
	}

	magic, _ := reader.ReadUint16()
	if magic != UsmapMagic {
		return nil, fmt.Errorf("invalid .usmap magic: 0x%04x, requested 0x%04x", magic, UsmapMagic)
	}

	versionByte, _ := reader.ReadByte()
	version := EUsmapVersion(versionByte)
	if version > EUsmapVersionLatest {
		return nil, fmt.Errorf("unsupported .usmap file: %d", int(version))
	}

	bHasVersioning, _ := reader.ReadInt32()
	if version > EUsmapVersionPackageVersioning && bHasVersioning != 0 {
		reader.Position += 4 * 2 // FPackageFileVersion

		versionsLength, _ := reader.ReadInt32() // FCustomVersionContainer
		reader.Position += int(versionsLength * (16 /* FGuid */ + 4))
	}

	compressionMethodByte, _ := reader.ReadByte()
	usmap.CompressionMethod = ECompressionMethod(compressionMethodByte)

	compressedSize, _ := reader.ReadInt32()
	uncompressedSize, _ := reader.ReadInt32()

	if usmap.CompressionMethod > ECompressionMethodMax {
		return nil, fmt.Errorf("unsupported compression method: %s", usmap.CompressionMethod.ToString())
	}

	if len(reader.Buffer)-reader.Position < int(compressedSize) {
		return nil, fmt.Errorf("there is not enough data in the .usmap file")
	}

	compressedData, _ := reader.ReadBytes(int(compressedSize))

	var uncompressedReader *UsmapReader
	decompressedData := make([]byte, uncompressedSize)

	switch usmap.CompressionMethod {
	case ECompressionMethodNone:
		uncompressedReader = reader

	case ECompressionMethodOodle:
		if usmap.oodleDec == nil {
			return nil, fmt.Errorf("oodle compression used but no oodle instance provided")
		}

		result, err := usmap.oodleDec.Decompress(compressedData, decompressedData, int(uncompressedSize))
		if err != nil {
			return nil, err
		}
		if result != int(uncompressedSize) {
			return nil, fmt.Errorf("invalid .usmap decompress result: %d/%d", result, uncompressedSize)
		}

		uncompressedReader, err = CreateReader(decompressedData)
		if err != nil {
			return nil, err
		}

	case ECompressionMethodBrotli:
		brotliDec := brotli.NewReader(bytes.NewReader(compressedData))
		n, err := io.ReadFull(brotliDec, decompressedData)
		if err != nil {
			return nil, err
		}
		if n != int(uncompressedSize) {
			return nil, fmt.Errorf("brotli: decompressed size mismatch: got %d, expected %d", n, uncompressedSize)
		}

		uncompressedReader, err = CreateReader(decompressedData)
		if err != nil {
			return nil, err
		}

	case ECompressionMethodZStandard:
		zstdDec, err := zstd.NewReader(nil)
		if err != nil {
			return nil, err
		}
		defer zstdDec.Close()

		decompressed, err := zstdDec.DecodeAll(compressedData, make([]byte, 0, uncompressedSize))
		if err != nil {
			return nil, err
		}
		if len(decompressed) != int(uncompressedSize) {
			return nil, fmt.Errorf("zstandard: decompressed size mismatch: got %d, expected %d", len(decompressed), uncompressedSize)
		}

		uncompressedReader, err = CreateReader(decompressed)
		if err != nil {
			return nil, err
		}
	}

	usmap.reader = uncompressedReader
	usmap.parseInternal(version >= EUsmapVersionLongFName, version >= EUsmapVersionLargeEnums)

	return usmap, nil
}

// Parses a .usmap file from the specified file path.
//
// The function reads the file from disk, validates the .usmap header,
// determines the compression method, and decompresses the data if necessary.
// It supports Oodle, Brotli, and Zstandard compression methods.
//
// If an Oodle-compressed file is encountered, oodleInstance must not be nil.
// It returns a pointer to a Usmap structure containing the parsed data, or an
// error if the file is invalid or decompression fails.
func ParseFromFile(filePath string, oodleInstance *Oodle) (*Usmap, error) {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParseFromBytes(f, oodleInstance)
}

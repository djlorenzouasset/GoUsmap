package gousmap

import (
	"io"
	"fmt"
	"encoding/binary"
)

type UsmapReader struct {
	Buffer   []byte
	Position int
}

func CreateReader(buff []byte) (*UsmapReader, error) {
	if len(buff) == 0 {
		return nil, fmt.Errorf("the buffer is empty")
	}

	reader := &UsmapReader{
		Buffer: buff,
		Position: 0,
	}

	return reader, nil
}

func (r *UsmapReader) ensure(n int) error {
	if n < 0 {
		return fmt.Errorf("negative length: %d", n)
	}
	if r.Position < 0 {
		return fmt.Errorf("negative reader position: %d", r.Position)
	}
	if r.Position + n > len(r.Buffer) {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func (r *UsmapReader) ReadBytes(length int) ([]byte, error) {
	if err := r.ensure(length); err != nil {
		return nil, err
	}
	start := r.Position
	r.Position += length
	out := make([]byte, length)
	copy(out, r.Buffer[start:r.Position])
	return out, nil
}

func (r *UsmapReader) ReadByte() (byte, error) { 
	bytes, err := r.ReadBytes(1)
	if err != nil {
		return 0x0, err
	}
	return bytes[0], nil
}

func (r *UsmapReader) ReadUint8() (uint8, error) {
	v, err := r.ReadByte()
	return uint8(v), err
}

func (r *UsmapReader) ReadInt8() (int8, error) {
	v, err := r.ReadByte()
	return int8(v), err
}

func (r *UsmapReader) ReadUint16() (uint16, error) {
	if err := r.ensure(2); err != nil {
		return 0, err
	}
	b := r.Buffer[r.Position : r.Position+2]
	r.Position += 2
	return binary.LittleEndian.Uint16(b), nil
}

func (r *UsmapReader) ReadInt16() (int16, error) {
	v, err := r.ReadUint16()
	return int16(v), err
}

func (r *UsmapReader) ReadUint32() (uint32, error) {
	if err := r.ensure(4); err != nil {
		return 0, err
	}
	b := r.Buffer[r.Position : r.Position+4]
	r.Position += 4
	return binary.LittleEndian.Uint32(b), nil
}

func (r *UsmapReader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

func (r *UsmapReader) ReadUint64() (uint64, error) {
	if err := r.ensure(8); err != nil {
		return 0, err
	}
	b := r.Buffer[r.Position : r.Position+8]
	r.Position += 8
	return binary.LittleEndian.Uint64(b), nil
}

func (r *UsmapReader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

func (r *UsmapReader) ReadBool() (bool, error) {
	v, err := r.ReadUint8()
	if err != nil {
		return false, err
	}
	return v != 0, nil
}

func (r *UsmapReader) ReadString(length int) (string, error) {
	v, err := r.ReadBytes(length)
	if err != nil {
		return "", err
	}
	return string(v), nil // assume its UTF8 ??
}
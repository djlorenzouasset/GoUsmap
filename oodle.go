package gousmap

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type OodleFuzzSafe int
const (
	OodleFuzzSafe_No  OodleFuzzSafe = 0
	OodleFuzzSafe_Yes OodleFuzzSafe = 1
)

type OodleCheckCrc int
const (
	OodleCheckCrc_No      OodleCheckCrc = 0
	OodleCheckCrc_Yes     OodleCheckCrc = 1
	OodleCheckCrc_Force32 OodleCheckCrc = 0x40000000
)

type OodleVerbosity int
const (
	OodleVerbosity_None    OodleVerbosity = 0
	OodleVerbosity_Minimal OodleVerbosity = 1
	OodleVerbosity_Some    OodleVerbosity = 2
	OodleVerbosity_Lots    OodleVerbosity = 3
	OodleVerbosity_Force32 OodleVerbosity = 0x40000000
)

type OodleDecodeThreadPhase int
const (
	OodleDecodeThreadPhase_Phase1     OodleDecodeThreadPhase = 1
	OodleDecodeThreadPhase_Phase2     OodleDecodeThreadPhase = 2
	OodleDecodeThreadPhase_All        OodleDecodeThreadPhase = 3
	OodleDecodeThreadPhase_Unthreaded OodleDecodeThreadPhase = OodleDecodeThreadPhase_All
)

type Oodle struct {
	oodleInstance  *syscall.DLL
	procDecompress *syscall.Proc
}

// Creates a new Oodle instance from the specified DLL path.
//
// The function loads the Oodle DLL located at oodlePath and retrieves the
// OodleLZ_Decompress procedure. It returns a pointer to an Oodle instance
// if successful, or an error if the DLL cannot be found or loaded.
func CreateOodleInstance(oodlePath string) (*Oodle, error) {
	_, err := os.Stat(oodlePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("oodle dll path not found: %s", oodlePath)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating instance: %s", err)
	}

	dll, err := syscall.LoadDLL(oodlePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load Oodle DLL: %w", err)
	}

	decompressProc, err := dll.FindProc("OodleLZ_Decompress")
	if err != nil {
		return nil, fmt.Errorf("failed to find 'OodleLZ_Decompress': %w", err)
	}

	return &Oodle{
		oodleInstance:  dll,
		procDecompress: decompressProc,
	}, nil
}

// Decompresses a buffer previously compressed with Oodle.
//
// It takes a compressed byte slice (compBuf) and writes the decompressed
// data into rawBuf, which must be large enough to hold the uncompressed data.
// rawLen specifies the expected uncompressed length.
//
// The function returns the number of bytes written to rawBuf and any error
// that occurred during decompression.
func (o *Oodle) Decompress(compBuf []byte, rawBuf []byte, rawLen int) (int, error) {
	if o.procDecompress == nil {
		return 0, fmt.Errorf("procedure 'OodleLZ_Decompress' not found")
	}

	compPtr := uintptr(unsafe.Pointer(&compBuf[0]))
	compSize := uintptr(len(compBuf))
	rawPtr := uintptr(unsafe.Pointer(&rawBuf[0]))
	rawSize := uintptr(rawLen)

	ret, _, err := o.procDecompress.Call(
		compPtr, compSize, rawPtr, rawSize,
		uintptr(OodleFuzzSafe_Yes),
		uintptr(OodleCheckCrc_No),
		uintptr(OodleVerbosity_None),
		0, // decBufBase (NULL)
		0, // decBufSize
		0, // fpCallback (NULL)
		0, // callbackUserData (NULL)
		0, // decoderMemory (NULL)
		0, // decoderMemorySize
		uintptr(OodleDecodeThreadPhase_Unthreaded),
	)

	if ret == 0 || int(ret) < 0 {
		return 0, fmt.Errorf("decompression failed: %v", err)
	}

	return int(ret), nil
}

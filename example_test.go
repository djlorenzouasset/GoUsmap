package gousmap_test

import (
	"fmt"
	"testing"

	"github.com/djlorenzouasset/GoUsmap"
)

func TestParseFromFile_Basic(t *testing.T) {
	// Parse a mapping file using ZStandard compression.
	usmap, err := gousmap.ParseFromFile("4.27.1-197800+release-4-1-OPP.usmap", nil)
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}
	fmt.Printf("ZStandard-compressed mapping: %s\n", usmap.ToString())

	// Parse a mapping file using Oodle compression.
	oodleInst, err := gousmap.CreateOodleInstance("oo2core_9_win64.dll")
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}

	usmap, err = gousmap.ParseFromFile("++Fortnite+Release-37.40-CL-46295673-Windows_oo.usmap", oodleInst)
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}

	fmt.Printf("Oodle-compressed mapping: %s\n", usmap.ToString())
}

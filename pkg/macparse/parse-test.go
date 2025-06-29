package macparse

import (
	"testing"
)

func TestMacParse(t *testing.T) {

	macAddr, err := ParseThreeSectionMAC("0011-326d-bfb3")
	if err != nil {
		t.Errorf("field db parse failed")
		return
	}
	if macAddr != "00:11:32:6d:bf:b3" {
		t.Errorf("field db parse failed")
		return
	}
}

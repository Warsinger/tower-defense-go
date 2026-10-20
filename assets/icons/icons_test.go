package icons

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"
	"testing"
)

func TestPlatformIconContainers(t *testing.T) {
	images, err := Images()
	if err != nil {
		t.Fatal(err)
	}
	for i, size := range []int{32, 48, 256} {
		if images[i].Bounds().Dx() != size || images[i].Bounds().Dy() != size {
			t.Fatalf("window icon %d has wrong size", size)
		}
		r, g, b, a := images[i].At(0, 0).RGBA()
		if r != 219*257 || g != 229*257 || b != 235*257 || a != 65535 {
			t.Fatal("approved background changed")
		}
	}
	ico, err := os.ReadFile("favicon.ico")
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(ico[2:4]) != 1 || binary.LittleEndian.Uint16(ico[4:6]) != 6 {
		t.Fatal("invalid ICO header")
	}
	for i, size := range []int{16, 32, 48, 64, 128, 256} {
		entry := ico[6+i*16 : 6+(i+1)*16]
		length, offset := binary.LittleEndian.Uint32(entry[8:12]), binary.LittleEndian.Uint32(entry[12:16])
		if uint64(offset)+uint64(length) > uint64(len(ico)) {
			t.Fatal("ICO payload out of bounds")
		}
		img, err := png.Decode(bytes.NewReader(ico[offset : offset+length]))
		if err != nil {
			t.Fatal(err)
		}
		if img.Bounds().Dx() != size || img.Bounds().Dy() != size {
			t.Fatal("ICO dimensions mismatch")
		}
	}
	icns, err := os.ReadFile("tower-defense.icns")
	if err != nil {
		t.Fatal(err)
	}
	if string(icns[:4]) != "icns" || int(binary.BigEndian.Uint32(icns[4:8])) != len(icns) {
		t.Fatal("invalid ICNS header")
	}
	expected := map[string]int{"icp4": 16, "icp5": 32, "icp6": 64, "ic07": 128, "ic08": 256, "ic09": 512, "ic10": 1024}
	for offset := 8; offset < len(icns); {
		if offset+8 > len(icns) {
			t.Fatal("truncated ICNS header")
		}
		kind := string(icns[offset : offset+4])
		length := int(binary.BigEndian.Uint32(icns[offset+4 : offset+8]))
		if length < 8 || offset+length > len(icns) {
			t.Fatal("ICNS payload out of bounds")
		}
		img, err := png.Decode(bytes.NewReader(icns[offset+8 : offset+length]))
		if err != nil {
			t.Fatal(err)
		}
		size, ok := expected[kind]
		if !ok || img.Bounds().Dx() != size || img.Bounds().Dy() != size {
			t.Fatal("ICNS representation mismatch")
		}
		delete(expected, kind)
		offset += length
	}
	if len(expected) != 0 {
		t.Fatal("missing ICNS representations")
	}
}

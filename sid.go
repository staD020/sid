// Package sid contains a pure-go implementation of the Commodore 64's .sid music format.
package sid

import (
	"fmt"
	"io"
	"os"
)

type (
	SID      []byte
	Word     uint16
	Version  uint16
	LongWord uint32
)

func New(r io.Reader) (SID, error) {
	bin, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	s := SID(bin)
	return s, s.Validate()
}

func LoadSID(path string) (SID, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return New(f)
}

func (s SID) Validate() error {
	if len(s) < 0x7c {
		return fmt.Errorf("input too short to be a .sid, length: %d bytes", len(s))
	}
	if err := s.headerMarkerOK(); err != nil {
		return err
	}
	v := s.Version()
	if v < 1 || v > 4 {
		return fmt.Errorf("incorrect version: %s", v)
	}
	o := s.dataOffset()
	if o != 0x76 && o != 0x7c {
		return fmt.Errorf("incorrect dataOffset %s", o)
	}
	return nil
}

func bytesToWord(bHi, bLo byte) Word {
	return Word(uint16(bHi)<<8 + uint16(bLo))
}

func (s SID) Version() Version {
	return Version(bytesToWord(s[4], s[5]))
}

func (s SID) dataOffset() Word {
	return bytesToWord(s[6], s[7])
}

func (s SID) LoadAddress() Word {
	if a := bytesToWord(s[8], s[9]); a > 0 {
		return a
	}
	offset := s.dataOffset()
	return bytesToWord(s[offset+1], s[offset])
}

func (s SID) Bytes() []byte {
	offset := s.dataOffset()
	if loadTo := bytesToWord(s[8], s[9]); loadTo == 0 {
		return s[offset:]
	}
	buf := []byte{s[8], s[9]}
	return append(buf, s[offset:]...)
}

func (s SID) RawBytes() []byte {
	offset := s.dataOffset()
	if loadTo := bytesToWord(s[8], s[9]); loadTo == 0 {
		return s[offset+2:]
	}
	return s[offset:]
}

func (s SID) InitAddress() Word {
	return bytesToWord(s[0xa], s[0xb])
}

func (s SID) PlayAddress() Word {
	return bytesToWord(s[0xc], s[0xd])
}

func (s SID) Songs() Word {
	return bytesToWord(s[0xe], s[0xf])
}

func (s SID) StartSong() Word {
	return bytesToWord(s[0x10], s[0x11])
}

func (s SID) Speed() LongWord {
	return LongWord(s[0x12])<<24 + LongWord(s[0x13])<<16 + LongWord(s[0x14])<<8 + LongWord(s[0x15])
}

func (s SID) Speed50Herz() bool {
	return s.Speed()&1 == 0
}

func chopString(in string) (out string) {
	for _, c := range in {
		if byte(c) == 0 {
			return out
		}
		out += string(c)
	}
	return out
}

func (s SID) Name() string {
	return chopString(string(s[0x16:0x35]))
}

func (s SID) Author() string {
	return chopString(string(s[0x36:0x55]))
}

func (s SID) Released() string {
	return chopString(string(s[0x56:0x75]))
}

func (s SID) String() string {
	return fmt.Sprintf("%q by %s (c) %s (%s-%s)", s.Name(), s.Author(), s.Released(), s.LoadAddress(), s.LoadAddress()+Word(len(s.RawBytes())))
}

func (v Version) String() string {
	switch v {
	case 1:
		return "PSID, 0x0001"
	case 2:
		return "PSID V2NG, RSID, 0x0002"
	case 3:
		return "PSID V2NG, RSID, 0x0003"
	case 4:
		return "PSID V2NG, RSID, 0x0004"
	}
	return fmt.Sprintf("unknown version %s", Word(v))
}

func (w Word) String() string {
	return fmt.Sprintf("0x%04x", uint16(w))
}

func (w Word) LowByte() byte {
	return byte(w & 0xff)
}

func (w Word) HighByte() byte {
	return byte(w >> 8)
}

func (w LongWord) String() string {
	return fmt.Sprintf("0x%08x", uint32(w))
}

// headerMarkerOK returns an error if the header of s does not match.
func (s SID) headerMarkerOK() error {
	if s[0] != 'P' && s[0] != 'R' {
		return fmt.Errorf("incorrect PSID/RSID header marker: first byte incorrect: %q", string(s[0:3]))
	}
	const postfix = "SID"
	for i, c := range postfix {
		if s[i+1] != byte(c) {
			return fmt.Errorf("incorrect PSID/RSID header marker: %q", string(s[0:3]))
		}
	}
	return nil
}

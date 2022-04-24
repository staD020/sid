// Package sid contains a pure-go implementation of the Commodore 64's .sid music format.
package sid

import (
	"fmt"
	"os"
)

type (
	SID struct {
		bin []byte
	}
	Word     uint16
	Version  uint16
	LongWord uint32
)

func LoadSID(path string) (*SID, error) {
	bin, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := &SID{bin: bin}
	err = s.Validate()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SID) Validate() error {
	if err := s.headerMarkerOK(); err != nil {
		return err
	}
	v := s.Version()
	if v < 1 || v > 4 {
		return fmt.Errorf("incorrect version number 0x%04x: %s", v, v)
	}
	o := s.dataOffset()
	if o != 0x76 && o != 0x7c {
		return fmt.Errorf("incorrect dataOffset 0x%02x", o)
	}
	return nil
}

func bytesToWord(bHi, bLo byte) Word {
	return Word(uint16(bHi)<<8 + uint16(bLo))
}

func (s *SID) Version() Version {
	return Version(bytesToWord(s.bin[4], s.bin[5]))
}

func (s *SID) dataOffset() Word {
	return bytesToWord(s.bin[6], s.bin[7])
}

func (s *SID) LoadAddress() Word {
	if a := bytesToWord(s.bin[8], s.bin[9]); a > 0 {
		return a
	}
	offset := s.dataOffset()
	return bytesToWord(s.bin[offset+1], s.bin[offset])
}

func (s *SID) Bytes() []byte {
	offset := s.dataOffset()
	if loadTo := bytesToWord(s.bin[8], s.bin[9]); loadTo == 0 {
		return s.bin[offset:]
	}
	buf := []byte{s.bin[8], s.bin[9]}
	return append(buf, s.bin[offset:]...)
}

func (s *SID) RawBytes() []byte {
	offset := s.dataOffset()
	if loadTo := bytesToWord(s.bin[8], s.bin[9]); loadTo == 0 {
		return s.bin[offset+2:]
	}
	return s.bin[offset:]
}

func (s *SID) String() string {
	return fmt.Sprintf("%q by %q - %q", s.Name(), s.Author(), s.Released())
}

func (s *SID) InitAddress() Word {
	return bytesToWord(s.bin[0xa], s.bin[0xb])
}

func (s *SID) PlayAddress() Word {
	return bytesToWord(s.bin[0xc], s.bin[0xd])
}

func (s *SID) Songs() Word {
	return bytesToWord(s.bin[0xe], s.bin[0xf])
}

func (s *SID) StartSong() Word {
	return bytesToWord(s.bin[0x10], s.bin[0x11])
}

func (s *SID) Speed() LongWord {
	return LongWord(s.bin[0x12])<<24 + LongWord(s.bin[0x13])<<16 + LongWord(s.bin[0x14])<<8 + LongWord(s.bin[0x15])
}

func chopString(in string) (out string) {
	for _, c := range in {
		if byte(c) == 0 {
			return out
		}
		out = out + string(c)
	}
	return out
}

func (s *SID) Name() string {
	return chopString(string(s.bin[0x16:0x35]))
}

func (s *SID) Author() string {
	return chopString(string(s.bin[0x36:0x55]))
}

func (s *SID) Released() string {
	return chopString(string(s.bin[0x56:0x75]))
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
	default:
		return "unknown version"
	}
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

func (s *SID) headerMarkerOK() error {
	if s.bin[0] != 'P' && s.bin[0] != 'R' {
		return fmt.Errorf("first byte incorrect")
	}
	const postfix = "SID"
	for i, c := range postfix {
		if s.bin[i+1] != byte(c) {
			return fmt.Errorf("incorrect PSID/RSID header marker")
		}
	}
	return nil
}

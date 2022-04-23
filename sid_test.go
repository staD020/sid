package sid

import (
	"testing"
)

const testSID = "testdata/Rivalry_tune_5.sid"

func TestLoadSID(t *testing.T) {
	s, err := LoadSID(testSID)
	if err != nil {
		t.Fatalf("LoadSID %q error: %v", testSID, err)
	}
	if err := s.Validate(); err != nil {
		t.Errorf("s.Validate() %q error: %v", testSID, err)
	}

	cases := []struct {
		got  map[string]string
		want map[string]string
	}{
		{
			got: map[string]string{
				"Name":        s.Name(),
				"Author":      s.Author(),
				"Released":    s.Released(),
				"Version":     s.Version().String(),
				"Speed":       s.Speed().String(),
				"dataOffset":  s.dataOffset().String(),
				"LoadAddress": s.LoadAddress().String(),
				"InitAddress": s.InitAddress().String(),
				"PlayAddress": s.PlayAddress().String(),
			},
			want: map[string]string{
				"Name":        "Rivalry (tune 5)",
				"Author":      "Thomas E. Petersen (Laxity)",
				"Released":    "2019 Seniors",
				"Version":     "PSID V2NG, RSID, 0x0002",
				"Speed":       "0x00000000",
				"dataOffset":  "0x007c",
				"LoadAddress": "0x1000",
				"InitAddress": "0x1000",
				"PlayAddress": "0x1009",
			},
		},
	}

	for _, c := range cases {
		for k, got := range c.got {
			if got != c.want[k] {
				t.Errorf("s.%s() mismatch got: %q want: %q", k, got, c.want[k])
			}
		}
	}

	/*
		fmt.Println("name:", s.Name())
		fmt.Println("author:", s.Author())
		fmt.Println("released:", s.Released())
		fmt.Println("version:", s.Version())
		fmt.Println("speed:", s.Speed())
		fmt.Println("dataOffset:", s.dataOffset())
		fmt.Println("LoadAddress:", s.LoadAddress())
		fmt.Println("InitAddress:", s.InitAddress())
		fmt.Println("PlayAddress:", s.PlayAddress())
	*/
}

func TestBytesToWord(t *testing.T) {
	cases := []struct {
		in   [2]byte
		want Word
	}{
		{[2]byte{0, 0}, 0},
		{[2]byte{16, 0}, 4096},
		{[2]byte{16, 1}, 4097},
		{[2]byte{64, 65}, 16449},
		{[2]byte{255, 255}, 65535},
	}
	for _, c := range cases {
		got := bytesToWord(c.in[0], c.in[1])
		if got != c.want {
			t.Errorf("bytesToWord(%d, %d) == %s, want %s", c.in[0], c.in[1], got, c.want)
		}
	}

}

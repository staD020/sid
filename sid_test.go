package sid

import (
	"bytes"
	"testing"
)

// const testSID = "testdata/jasonpage_eighth.sid"
const testSID = "testdata/Rivalry_tune_5.sid"

func TestNew(t *testing.T) {
	t.Parallel()
	s, err := New(&bytes.Buffer{})
	if err == nil {
		t.Errorf("New empty buf succeeded while it should have failed")
	}
	if len(s) != 0 {
		t.Errorf("New empty buf not empty")
	}
}

func TestLoadSID(t *testing.T) {
	t.Parallel()
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
				"String":      s.String(),
				"Songs":       s.Songs().String(),
				"StartSong":   s.StartSong().String(),
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
				"String":      `"Rivalry (tune 5)" by Thomas E. Petersen (Laxity) (c) 2019 Seniors (0x1000-0x1fec)`,
				"Songs":       "0x0001",
				"StartSong":   "0x0001",
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
}

func TestSpeed50Herz(t *testing.T) {
	t.Parallel()
	s, err := LoadSID(testSID)
	if err != nil {
		t.Fatalf("LoadSID %q error: %v", testSID, err)
	}

	got := s.Speed50Herz()
	want := true
	if got != want {
		t.Errorf("s.Speed50Herz() mismatch got: %v want: %v", got, want)
	}
}

func TestBytes(t *testing.T) {
	t.Parallel()
	s, err := LoadSID(testSID)
	if err != nil {
		t.Fatalf("LoadSID %q error: %v", testSID, err)
	}

	got := len(s.Bytes())
	want := 4078
	if got != want {
		t.Errorf("len(s.Bytes()) mismatch got: %d want: %d", got, want)
	}
}

func TestRawBytes(t *testing.T) {
	t.Parallel()
	s, err := LoadSID(testSID)
	if err != nil {
		t.Fatalf("LoadSID %q error: %v", testSID, err)
	}

	got := len(s.RawBytes())
	want := 4076
	if got != want {
		t.Errorf("len(s.RawBytes()) mismatch got: %d want: %d", got, want)
	}
}

func TestBytesToWord(t *testing.T) {
	t.Parallel()
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

func TestChopString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"foo", "foo"},
		{"foo ", "foo "},
		{"foo\x00", "foo"},
		{"foo\x00bar", "foo"},
	}
	for _, c := range cases {
		got := chopString(c.in)
		if got != c.want {
			t.Errorf("chopString(%s) == %s, want %s", c.in, got, c.want)
		}
	}
}

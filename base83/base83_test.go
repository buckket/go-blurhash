package base83_test

import (
	"fmt"
	"github.com/buckket/go-blurhash/base83"
	"testing"
)

func TestDecode(t *testing.T) {
	var testCasesValid = []struct {
		in  string
		out int
	}{
		{"", 0},
		{"foobar", 163902429697},
		{"LFE.@D9F01_2%L%MIVD*9Goe-;WB", -1597651267176502418},
	}

	for _, tc := range testCasesValid {
		t.Run(fmt.Sprintf("%q", tc.in), func(t *testing.T) {
			out, err := base83.Decode(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.out {
				t.Fatalf("got %d, wanted %d", out, tc.out)
			}
		})
	}

	_, err := base83.Decode("LFE.@D9F01_2%L%MIVD*9Goe-;Wµ")
	if err == nil {
		t.Fatal("should have failed")
	}

	err, ok := err.(base83.InvalidCharacterError)
	if !ok {
		t.Fatal("wrong error type")
	}
	_ = err.Error()
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = base83.Decode("~$")
	}
}

func TestEncode(t *testing.T) {
	var testCasesValid = []struct {
		in     int
		length int
		out    string
	}{
		{0, 0, ""},
		{163902429697, 6, "foobar"},
		{100, 2, "1H"},
	}

	for _, tc := range testCasesValid {
		t.Run(fmt.Sprintf("%d", tc.in), func(t *testing.T) {
			out, err := base83.Encode(tc.in, tc.length)
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.out {
				t.Fatalf("got %q, wanted %q", out, tc.out)
			}
		})
	}

	var testCasesInvalid = []struct {
		in     int
		length int
	}{
		{84, 1},
		{-1597651267176502418, 16},
		{163902429697, 10},
	}

	for _, tc := range testCasesInvalid {
		t.Run(fmt.Sprintf("%d", tc.in), func(t *testing.T) {
			_, err := base83.Encode(tc.in, tc.length)
			if err == nil {
				t.Fatal("should have failed")
			}

			err, ok := err.(base83.InvalidLengthError)
			if !ok {
				t.Fatal("wrong error type")
			}
			_ = err.Error()
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = base83.Encode(6869, 2)
	}
}

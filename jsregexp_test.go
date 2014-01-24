package jsregexp_test

import (
	"."
	"fmt"
	"testing"
)

func ExampleTranslate() {
	fmt.Println(jsregexp.Translate("/asdf (.+)/i"))
	// Output: (?i:asdf (.+))
}

func TestForwardSlashes(t *testing.T) {
	expected := "(?i:\\/r\\/(.*))"
	translated := jsregexp.Translate("/\\/r\\/(.*)/i")

	if translated != expected {
		t.Errorf("Expected " + expected + " but got " + translated)
	}
}

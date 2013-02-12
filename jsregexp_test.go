package jsregexp_test

import (
	"."
	"fmt"
)

func ExampleTranslate() {
	fmt.Println(jsregexp.Translate("/asdf (.+)/i"))
	// Output: (?i:asdf (.+))
}

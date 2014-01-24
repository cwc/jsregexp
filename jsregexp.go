// Package jsregexp provides the Translate function for using JavaScript regex
// strings with the Go 'regexp' package.
package jsregexp

import (
	"bytes"
	"fmt"
	"strings"

	re "regexp"
)

// Translate reformats a JavaScript regex pattern into one compatible with
// Go's 'regexp' package.
func Translate(pattern string) (goPattern string) {
	re2flags := ""

	// Last group after '/' will be the flags
	split := strings.Split(pattern, "/")
	flags := split[len(split)-1]

	// Parse flags
	for _, rune := range flags {
		switch rune {
		case 'i':
			re2flags += "i"
			break
		case 'm':
			re2flags += "m"
			break
		}
	}

	// Rebuild pattern without the flags
	pattern = strings.Join(split[:len(split)-1], "/")

	goPattern = transformRegExp(pattern[1:len(pattern)])
	if len(re2flags) > 0 {
		goPattern = fmt.Sprintf("(?%s:%s)", re2flags, goPattern)
	}

	return
}

// 0031,0032,0033,0034,0035,0036,0037,0038,0039 // 1 - 9
// 0043,0045,0046,0047,0048,0049,004A,004B,004C,004D,004E,004F
// 0050,0052,0054,0055,0056,0058,0059,005A
// 0063,0065,0067,0068,0069,006A,006B,006C,006D,006F
// 0070,0071,0075,0078,0079
// 0080,0081,0082,0083,0084,0085,0086,0087,0088,0089,008A,008B,008C,008D,008E,008F
// 0090,0091,0092,0093,0094,0095,0096,0097,0098,0099,009A,009B,009C,009D,009E,009F
// 00A0,00A1,00A2,00A3,00A4,00A5,00A6,00A7,00A8,00A9,00AA,00AB,00AC,00AD,00AE,00AF
// 00B0,00B1,00B2,00B3,00B4,00B5,00B6,00B7,00B8,00B9,00BA,00BB,00BC,00BD,00BE,00BF
// 00C0,00C1,00C2,00C3,00C4,00C5,00C6,00C7,00C8,00C9,00CA,00CB,00CC,00CD,00CE,00CF
// ...
// c = 63* c[A-Z]
// p = 70
// u = 75* u[:xdigit:]{4}
// x = 78* x[:xdigit:]{2}
//\x{0031}-\x{0039}

var transformRegExp_matchSlashU = re.MustCompile(`\\u([[:xdigit:]]{1,4})`)
var transformRegExp_escape_c = re.MustCompile(`\\c([A-Za-z])`)
var transformRegExp_unescape_c = re.MustCompile(`\\c`)
var transformRegExp_unescape = []*re.Regexp{
	re.MustCompile(strings.NewReplacer("\n", "", "\t", "", " ", "").Replace(`
		\\(
		[
			\x{0043}\x{0045}-\x{004F}
			\x{0050}\x{0052}\x{0054}-\x{0056}\x{0058}-\x{005A}
			\x{0065}\x{0067}-\x{006D}\x{006F}
			\x{0070}\x{0071}\x{0079}
			\x{0080}-\x{FFFF}
		]
		)()
	`)),
	re.MustCompile(`\\(u)([^[:xdigit:]])`),
	re.MustCompile(`\\(u)([[:xdigit:]][^[:xdigit:]])`),
	re.MustCompile(`\\(u)([[:xdigit:]][[:xdigit:]][^[:xdigit:]])`),
	re.MustCompile(`\\(u)([[:xdigit:]][[:xdigit:]][[:xdigit:]][^[:xdigit:]])`),
	re.MustCompile(`\\(x)([^[:xdigit:]])`),
	re.MustCompile(`\\(x)([[:xdigit:]][^[:xdigit:]])`),
}

var transformRegExp_unescapeDollar = re.MustCompile(`\\([cux])$`)

// TODO Go "re" bug? Can't do: (?:)|(?:$)

func transformRegExp(ecmaRegExp string) (goRegExp string) {
	// https://bugzilla.mozilla.org/show_bug.cgi/show_bug.cgi?id=334158
	tmp := []byte(ecmaRegExp)
	for _, value := range transformRegExp_unescape {
		tmp = value.ReplaceAll(tmp, []byte(`$1$2`))
	}
	tmp = transformRegExp_escape_c.ReplaceAllFunc(tmp, func(in []byte) []byte {
		in = bytes.ToUpper(in)
		return []byte(fmt.Sprintf("\\%o", in[0]-64)) // \cA => \001 (A == 65)
	})
	tmp = transformRegExp_unescape_c.ReplaceAll(tmp, []byte(`c`))
	tmp = transformRegExp_unescapeDollar.ReplaceAll(tmp, []byte(`$1`))
	tmp = transformRegExp_matchSlashU.ReplaceAll(tmp, []byte(`\x{$1}`))
	return string(tmp)
}

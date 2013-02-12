jsregexp
========

A package for using JavaScript regexes with Go

Example
-------
```
func ExampleTranslate() {
        fmt.Println(jsregexp.Translate("/asdf (.+)/i"))
        // Output: (?i:asdf (.+))
}
```

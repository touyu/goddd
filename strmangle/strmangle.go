package strmangle

import (
	"strings"
	"sync"
)

var uppercaseWords = map[string]struct{}{
	"acl":   {},
	"api":   {},
	"ascii": {},
	"cpu":   {},
	"eof":   {},
	"guid":  {},
	"id":    {},
	"ip":    {},
	"json":  {},
	"ram":   {},
	"sla":   {},
	"udp":   {},
	"ui":    {},
	"uid":   {},
	"uuid":  {},
	"uri":   {},
	"url":   {},
	"utf8":  {},
}

// titleCaseCache holds the mapping of title cases.
// Example: map["MyWord"] == "my_word"
var (
	mut            sync.RWMutex
	titleCaseCache = map[string]string{}
)

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
//
// Note: This method is ugly because it has been highly optimized,
// we found that it was a fairly large bottleneck when we were using regexp.
func TitleCase(n string) string {
	// Attempt to fetch from cache
	mut.RLock()
	val, ok := titleCaseCache[n]
	mut.RUnlock()
	if ok {
		return val
	}

	ln := len(n)
	name := []byte(n)
	buf := GetBuffer()

	start := 0
	end := 0
	for start < ln {
		// Find the start and end of the underscores to account
		// for the possibility of being multiple underscores in a row.
		if end < ln {
			if name[start] == '_' {
				start++
				end++
				continue
				// Once we have found the end of the underscores, we can
				// find the end of the first full word.
			} else if name[end] != '_' {
				end++
				continue
			}
		}

		word := name[start:end]
		wordLen := len(word)
		var vowels bool

		numStart := wordLen
		for i, c := range word {
			vowels = vowels || (c == 97 || c == 101 || c == 105 || c == 111 || c == 117 || c == 121)

			if c > 47 && c < 58 && numStart == wordLen {
				numStart = i
			}
		}

		_, match := uppercaseWords[string(word[:numStart])]

		if match || !vowels {
			// Uppercase all a-z characters
			for _, c := range word {
				if c > 96 && c < 123 {
					buf.WriteByte(c - 32)
				} else {
					buf.WriteByte(c)
				}
			}
		} else {
			if c := word[0]; c > 96 && c < 123 {
				buf.WriteByte(word[0] - 32)
				buf.Write(word[1:])
			} else {
				buf.Write(word)
			}
		}

		start = end + 1
		end = start
	}

	ret := buf.String()
	PutBuffer(buf)

	// Cache the title case result
	mut.Lock()
	titleCaseCache[n] = ret
	mut.Unlock()

	return ret
}

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID". It will also lowercase the first letter
// of the name in the case where it's fed something that starts with uppercase.
func CamelCase(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	// Discard all leading '_'
	index := -1
	for i := 0; i < len(name); i++ {
		if name[i] != '_' {
			index = i
			break
		}
	}

	if index != -1 {
		name = name[index:]
	} else {
		return ""
	}

	index = -1
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			index = i
			break
		}
	}

	if index == -1 {
		buf.WriteString(strings.ToLower(string(name[0])))
		if len(name) > 1 {
			buf.WriteString(name[1:])
		}
	} else {
		buf.WriteString(strings.ToLower(string(name[0])))
		if len(name) > 1 {
			buf.WriteString(name[1:index])
			buf.WriteString(TitleCase(name[index+1:]))
		}
	}

	return buf.String()
}

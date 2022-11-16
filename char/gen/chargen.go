// Package char/gen is a library to generate lists of characters.
package gen

import (
	"errors"
	"regexp/syntax"
)

func CodePt(p int32) string {
	return string(p)
}

// Returns unicode characters from start to end (inclusive)
func Range(start int32, end int32) (string, error) {
	var chars string
	if start < 0 || end < 0 || start >= end {
		return "", errors.New("RangeChars(" + string(start) + ", " + string(end) + "): start must be smaller than end, and both must be positive")
	}
	
	for n := start; n <= end; n++ {
		chars += string(n)
	}
	return chars, nil
}


// Make a regex program with the input.
func makeRE(input string) (*syntax.Prog, error) {
	var re *syntax.Prog
	var rp *syntax.Regexp
	var e error
	if rp, e = syntax.Parse(input, syntax.Simple);
	 e != nil {
		return nil, e
	}
	rp.Simplify()
	if re, e = syntax.Compile(rp); e != nil {
		return nil, e
	}
	// de-dupe each Inst.Rune -- it is not clear why, but duplicate values for the space character were observed
	for i := 0; i < len(re.Inst); i++ {
		re.Inst[i].Rune = dedupe(re.Inst[i].Rune)
	}
	return re, nil
}

// Walk regex program values to get matching characters.
func inspectRE(re *syntax.Prog) (string, error) {
	var base string
	
	// walk all RE instructions
	for i := re.Start; i < len(re.Inst); i++ {
		inst := re.Inst[i]
		runeLen := len(inst.Rune)
		runeCap := cap(inst.Rune)
		if runeLen == 0 { // no contents
			continue
		} else if runeLen == runeCap { // treat as start-end pairs
			for k := 0; k < runeLen; k++ {
				start := inst.Rune[k]
				k++
				end := inst.Rune[k]

				if s, e := Range(int32(start), int32(end)); e != nil {
					return "", e
				} else {
					base += s
				}
			}
		} else { //treat as individual characters
			for k := 0; k < runeLen; k++ {
				base += string(inst.Rune[k])
			}
		}

	}
	
	return base, nil
}

// Get a string of characters (grapheme clusters) from regex input.
func Regex(input string) (string, error) {
	if re, e := makeRE(input); e!= nil {
		return "",e
	} else if s, e := inspectRE(re); e != nil {
		return "",e
	} else {
		return s, nil
	}
}

type Numeric interface {
	~uint       |
	~uint8      |
	~uint16     |
	~uint32     |
	~uint64     |
	~uintptr    |
	~int        |
	~int8       |
	~int16      |
	~int32      |
	~int64      |
	~float32    |
	~float64    |
	~complex64  |
	~complex128 
}

type Comparable interface {
	Numeric | ~string
}

// TODO: Move to a more common util location
func dedupe[T Comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

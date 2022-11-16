package sub

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

const (
	cpSeparator string = "|"
)

// High level container for elements used to shift the characters of a text.
type CharSub struct {
	includes []*CharSubMap
	excludes []*CharSubMap
}

/* CharSub methods */

func (o *CharSub) Includes() []*CharSubMap {
	return o.includes
}
func (o *CharSub) Excludes() []*CharSubMap {
	return o.excludes
}
func (o *CharSub) AddInclude(baseline string, substitution string) (*CharSub, error) {
	if csm, e := makeCharSubMap(baseline, substitution); e != nil {
		return o, e
	} else {
		o.includes = append(o.includes, csm)
	}
	return o, nil
}
func (o *CharSub) AddExclude(baseline string) (*CharSub, error) {
	if csm, e := makeCharSubMap(baseline, baseline); e != nil {
		return o, e
	} else {
		o.excludes = append(o.excludes, csm)
	}
	return o, nil
}
func (o *CharSub) Do(s string) (string, error) {
	return do(o, s)
}

/* CharSub functions */

// Map character substitutions.
func makeCharSubMap(base string, sub string) (*CharSubMap, error) {
	if csm, e := CreateCharSubMap(base, sub); e != nil {
		return nil, e
	} else {
		return csm, nil
	}
}

// Shift the characters of a string based on the character sets.
func do(o *CharSub, s string) (string, error) {
	// check for required character sets
	if len(o.includes) < 1 {
		return "", errors.New("includes must be defined before text can be shifted")
	}
	if len(s) < 1 {
		return "", errors.New("empty strings cannot be shifted")
	}

	var output string
	// Read each character and look for it in excludes and then includes, stopping at first match.
	g := uniseg.NewGraphemes(s)
next:
	for g.Next() {
		// get the code points
		runes := g.Runes()
		var strs []string
		for i := 0; i < len(runes); i++ {
			strs = append(strs, strconv.Itoa(int(runes[i])))
		}
		// concatenate all the code points
		mapKey := strings.Join(strs, cpSeparator)

		var included, excluded bool

	excludes:
		for i := 0; i < len(o.excludes); i++ {
			exclude := o.excludes[i]
			if _, ok := exclude.Map()[mapKey]; ok {
				excluded = true
				output += g.Str()
				break excludes
			}
		}
		if excluded {
			continue next
		}
	includes:
		for i := 0; i < len(o.includes); i++ {
			include := o.includes[i]
			if v, ok := include.Map()[mapKey]; ok {
				included = true
				// parse map value back into a character string
				v := strings.Split(v, cpSeparator)
				for n := 0; n < len(v); n++ {
					tmpInt, e := strconv.Atoi(v[n])
					if e != nil {
						return "", e
					}
					output += string(int32(tmpInt))
				}
				break includes
			}
		}
		if included {
			continue next
		}
		// handle an unrecognized character
		fmt.Println("warning: encountered unrecognized character:  " + g.Str())
		fmt.Println("\t unrecognized characters are excluded from substitution")
		output += g.Str()
	}

	return output, nil
}



// Create a new CharSub.
func Blank() *CharSub {
	return &CharSub{}
}

// Create a new CharSub with inputs.
//
// variadic function to allow excludes to be optional
// arguments should be ordered bases, subs, excludes
// bases correspond to subs, i.e. their lengths must be the same.
// excludes have no subs.
func Create(lists ...[]string) (*CharSub, error) {
	o := Blank()
	switch n := len(lists); {
	case n == 0:
		return o, nil
	case n == 1:
		return nil, errors.New("list of base characters provided without substitutions")
	case 3 < n:
		return nil, errors.New("too many (" + strconv.Itoa(n) + ") lists provided")
	}

	bases := lists[0]
	subs := lists[1]
	var excludes []string
	if len(lists) == 3 {
		excludes = lists[2]
	}

	if len(bases) != len(subs) {
		return nil, errors.New("bases length different from subs length")
	}

	for i := 0; i < len(bases); i++ {
		if _, e := o.AddInclude(bases[i], subs[i]); e != nil {
			return nil, e
		}
	}

	for i := 0; i < len(excludes); i++ {
		if _, e := o.AddExclude(excludes[i]); e != nil {
			return nil, e
		}
	}

	return o, nil
}

package sub

import (
	"errors"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

// Character sets used to align character substitution.
type CharSubMap struct {
	base        string
	baseList    []string
	baseUnicode []string
	sub         string
	subList     []string
	subUnicode  []string
	unicodeMap  map[string]string
	count       int
}

/* CharSubMap methods */

// Get the base string.
func (c *CharSubMap) Base() string {
	return c.base
}

// Get the sub string.
func (c *CharSubMap) Sub() string {
	return c.sub
}

// Get the base to sub character mapping.
func (c *CharSubMap) Map() map[string]string {
	return c.unicodeMap
}

// Count base grapheme clusters, i.e. user-perceived characters.
func (c *CharSubMap) Count() int {
	c.count = CharCount(c.base)
	return c.count
}

func (c *CharSubMap) Build() (*CharSubMap, error) {
	var e error = nil
	if c.baseList,
		c.baseUnicode,
		c.subList,
		c.subUnicode,
		c.unicodeMap,
		e =
		buildCharSubMap(c.base, c.sub); e != nil {
		return c, e
	} else {
		return c, nil
	}
}

// Process the input and offset to generate the other field values.
func buildCharSubMap(base string, sub string) ([]string, []string, []string, []string, map[string]string, error) {
	var e error = nil
	var baseList, baseUnicode, subList, subUnicode []string
	unicodeMap := make(map[string]string)

	if baseList, baseUnicode, e = parseCharSet(base); e != nil {
		return nil, nil, nil, nil, nil, e
	}

	if subList, subUnicode, e = parseCharSet(sub); e != nil {
		return nil, nil, nil, nil, nil, e
	}

	if unicodeMap, e = mapLists(baseUnicode, subUnicode); e != nil {
		return nil, nil, nil, nil, nil, e
	}

	return baseList, baseUnicode, subList, subUnicode, unicodeMap, nil
}

// Return an empty CharSubMap
func BlankCharSubMap() CharSubMap {
	return CharSubMap{}
}

// Create a CharSubMap with initial values.
func CreateCharSubMap(base string, sub string) (*CharSubMap, error) {
	c := BlankCharSubMap()
	c.base = base
	c.sub = sub
	c.Count()
	if _, e := c.Build(); e != nil {
		return nil, e
	}
	return &c, nil
}

/* Utility functions */


// TODO: This should be abstracted a little... into a list of grapheme cluster groups
// Separate base by character (grapheme cluster).
func parseCharSet(base string) ([]string, []string, error) {
	var l []string
	var u []string
	// process the graphemes (user perceived characters) within the character range
	g := uniseg.NewGraphemes(base)
	for g.Next() {
		l = append(l, g.Str()) // store the character
		// get the code points
		runes := g.Runes()
		// get the string equivalent of the code points e.g. 3 => "3"
		var s []string
		for i := 0; i < len(runes); i++ {
			s = append(s, strconv.Itoa(int(runes[i])))
		}
		// concatenate all the code points
		u = append(u, strings.Join(s, cpSeparator))
	}
	return l, u, nil
}

// Returns the count of grapheme clusters, i.e. user-perceived characters, including whitespace.
// "\t" == 1    "ab cd" == 5
func CharCount(s string) int {
	return uniseg.GraphemeClusterCount(s)
}

// TODO: Make this more generic and move it to a util location
// Map base characters to sub characters.
func mapLists(a []string, b []string) (map[string]string, error) {
	if len(a) != len(b) {
		return nil, errors.New("lists must be equal length to be mapped")
	}
	m := make(map[string]string)

	// map baseUnicode to subUnicode
	for i := 0; i < len(a); i++ {
		m[a[i]] = b[i]
	}
	return m, nil
}

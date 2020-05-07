package types

import (
	"bytes"
	"fmt"
)

type HashSet struct {
	m map[interface{}]bool
}

func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

func (set *HashSet) Add(e interface{}) (b bool) {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}

	if set.Len() != other.Len() {
		return false
	}

	for k, _ := range set.m {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

func (set *HashSet) Elements() []interface{} {
	initlen := len(set.m)

	snaphot := make([]interface{}, initlen)

	actuallen := 0

	for k, _ := range set.m {
		if actuallen < initlen {
			snaphot[actuallen] = k
		} else {
			snaphot = append(snaphot, k)
		}
		actuallen++
	}

	if actuallen < initlen {
		snaphot = snaphot[:actuallen]
	}

	return snaphot
}

func (set *HashSet) String() string {
	var buf bytes.Buffer

	buf.WriteString("set{")

	first := true

	for k, _ := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(fmt.Sprintf("%v", k))
	}

	buf.WriteString("}")

	return buf.String()
}
func (set *HashSet) AddSlice(as []interface{}) {
	for _, k := range as {
		set.Add(k)
	}
}
func (set *HashSet) AddStringSlice(as []string) {
	for _, k := range as {
		set.Add(k)
	}
}
func (set *HashSet) ToSlice() []interface{} {
	s := make([]interface{}, len(set.m))
	i := 0
	for k, _ := range set.m {
		s[i] = k
		i++
	}
	return s
}
func (set *HashSet) ToStringSlice() []string {
	s := make([]string, set.Len())
	i := 0
	for k, _ := range set.m {
		s[i] = fmt.Sprintf("%v", k)
		i++
	}
	return s
}

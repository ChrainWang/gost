package gost

import (
	"testing"
)

func TestPathSplit(t *testing.T) {
	pathExamples := []string{
		"aaa/bbbb/cccc/dddd",
	}

	for _, path := range pathExamples {
		ps := NewSpliter([]byte(path), []rune{'/', '\\'})
		ps.Split()
		for {
			if section, err := ps.Next(); err == nil {
				println(string(section))
				if section == nil {
					break
				}
			} else {
				break
			}
		}
	}
}

package xmd

import (
	"fmt"
	"strings"
)

func SpaceFn(cache *Cache) map[int]int {
	spaces := make(map[int]int)

	for i, item := range cache.histories {
		if _, ok := spaces[item.result]; ok {
			continue
		}

		spaces[item.result] = i + 1
	}

	return spaces
}

func fmtIntSlice(s []int) string {
	s0 := make([]string, 0, len(s))

	for _, i := range s {
		s0 = append(s0, fmt.Sprintf("%02d", i))
	}

	return strings.Join(s0, ",")
}

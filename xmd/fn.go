package xmd

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

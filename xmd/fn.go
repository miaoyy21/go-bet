package xmd

func SpaceFn(cache *Cache) map[int]int {
	spaces := make(map[int]int)

	for i := len(cache.histories) - 1; i >= 0; i-- {
		if _, ok := spaces[cache.histories[i].result]; ok {
			continue
		}

		spaces[cache.histories[i].result] = len(cache.histories) - i
	}

	return spaces
}

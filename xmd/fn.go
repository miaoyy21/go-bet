package xmd

import (
	"fmt"
	"log"
	"math"
	"sort"
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

func m2sFn(as map[int]int) []int {
	ss := make([]int, 0, len(as))

	for i := range as {
		ss = append(ss, i)
	}
	sort.Ints(ss)

	return ss
}

func fmtIntSlice(s []int) string {
	s0 := make([]string, 0, len(s))

	for _, i := range s {
		s0 = append(s0, fmt.Sprintf("%02d", i))
	}

	return strings.Join(s0, ",")
}

func modeFn(bets map[int]float64, md int) (int, string) {
	var m1, m2, m3, m4, m5, m6, m7, m8 int
	for result, rx := range bets {
		if rx < 1.0 {
			continue
		}

		if result >= 14 {
			m1 += stds[result]
		} else {
			m2 += stds[result]
		}

		if result%2 == 1 {
			m3 += stds[result]
		} else {
			m4 += stds[result]
		}

		if result >= 10 && result <= 17 {
			m5 += stds[result]
		} else {
			m6 += stds[result]
		}

		if result%10 >= 5 && result%10 <= 9 {
			m7 += stds[result]
		} else {
			m8 += stds[result]
		}
	}

	m5 = int(float64(m5) * 44.0 / 56.0)
	m6 = int(float64(m6) * 56.0 / 44.0)
	log.Printf("模式权重：大数【%d】, 小数【%d】, 奇数【%d】, 偶数【%d】, 中数【%d】, 边数【%d】, 大尾数【%d】, 小尾数【%d】 \n", m1, m2, m3, m4, m5, m6, m7, m8)

	if m1 >= md && m1 >= m2 && m1 >= m3 && m1 >= m4 && m1 >= m5 && m1 >= m6 && m1 >= m7 && m1 >= m8 {
		return 1, "大数"
	}

	if m2 >= md && m2 >= m1 && m2 >= m3 && m2 >= m4 && m2 >= m5 && m2 >= m6 && m2 >= m7 && m2 >= m8 {
		return 2, "小数"
	}

	if m3 >= md && m3 >= m1 && m3 >= m2 && m3 >= m4 && m3 >= m5 && m3 >= m6 && m3 >= m7 && m3 >= m8 {
		return 3, "奇数"
	}

	if m4 >= md && m4 >= m1 && m4 >= m2 && m4 >= m3 && m4 >= m5 && m4 >= m6 && m4 >= m7 && m4 >= m8 {
		return 4, "偶数"
	}

	if m5 >= md && m5 >= m1 && m5 >= m2 && m5 >= m3 && m5 >= m4 && m5 >= m6 && m5 >= m7 && m5 >= m8 {
		return 5, "中数"
	}

	if m6 >= md && m6 >= m1 && m6 >= m2 && m6 >= m3 && m6 >= m4 && m6 >= m5 && m6 >= m7 && m6 >= m8 {
		return 6, "边数"
	}

	if m7 >= md && m7 >= m1 && m7 >= m2 && m7 >= m3 && m7 >= m4 && m7 >= m5 && m7 >= m6 && m7 >= m8 {
		return 7, "大尾数"
	}

	if m8 >= md && m8 >= m1 && m8 >= m2 && m8 >= m3 && m8 >= m4 && m8 >= m5 && m8 >= m6 && m8 >= m7 {
		return 8, "小尾数"
	}

	return 0, "未知"
}

func extraFn(modeId int, mGold int, bets map[int]float64) (map[int]struct{}, map[int]int) {
	as := make([]int, 0)

	switch modeId {
	case 1:
		as = append(as, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27)
	case 2:
		as = append(as, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	case 3:
		as = append(as, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27)
	case 4:
		as = append(as, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26)
	case 5:
		as = append(as, 10, 11, 12, 13, 14, 15, 16, 17)
	case 6:
		as = append(as, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27)
	case 7:
		as = append(as, 5, 6, 7, 8, 9, 15, 16, 17, 18, 19, 25, 26, 27)
	case 8:
		as = append(as, 0, 1, 2, 3, 4, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24)
	default:
		as = make([]int, 0)
	}

	ams := make(map[int]struct{})
	for _, a := range as {
		ams[a] = struct{}{}
	}

	extras := make(map[int]int)
	for result, rx := range bets {
		if _, ok := ams[result]; ok {
			continue
		}

		betGold := rx * float64(2*mGold) * float64(stds[result]) / 1000
		iGold := int(math.Floor(betGold/500.0) * 500)
		if iGold > 0 {
			extras[result] = iGold
		}
	}

	return ams, extras
}

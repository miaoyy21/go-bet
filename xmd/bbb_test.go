package xmd

import (
	"log"
	"sort"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	extras := make(map[int]int)
	extras[13] = 12500
	extras[14] = 14500
	extras[15] = 16500

	stdBets := []int{50000, 10000, 5000, 2000, 1000, 500}
	betMaps := make(map[int][]int)

	for _, stdBet := range stdBets {
		betSlice, ok := betMaps[stdBet]
		if !ok {
			betSlice = make([]int, 0)
		}

		for result, betGold := range extras {
			qn := betGold / stdBet
			if qn > 0 {
				for i := 0; i < qn; i++ {
					betSlice = append(betSlice, result)
				}

				extras[result] = betGold - qn*stdBet
			}
		}

		sort.Ints(betSlice)
		betMaps[stdBet] = betSlice
	}

	for _, stdBet := range stdBets {
		log.Printf("%d :: %v \n", stdBet, betMaps[stdBet])

		//if err := hBetting1(issue, betGold, result, cache.user); err != nil {
		//	return err
		//}
	}
	//for stdBet, betSlice := range betMaps {
	//	log.Printf("%d :: %v \n", stdBet, betSlice)
	//}
}

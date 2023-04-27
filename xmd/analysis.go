package xmd

import (
	"log"
	"sort"
	"strconv"
)

var latest = make(map[int]struct{})

var wins int
var fails int

func analysis(cache *Cache) error {
	if err := cache.Sync(200); err != nil {
		return err
	}

	nextIssue := strconv.Itoa(cache.issue + 1)

	// 当前账户余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	// 输出
	if len(latest) == 0 {
		log.Printf("第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0

			log.Printf("第【👍 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
		} else {
			wins = 0
			fails++

			log.Printf("第【👀 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
		}
	}

	var nn int

	spaces := SpaceFn(cache)
	for i := len(cache.histories) - 1; i >= len(cache.histories)-20; i-- {
		result := cache.histories[i].result
		if result <= 4 || result >= 23 {
			nn++
			continue
		}
	}

	if nn <= 1 {
		for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
			result := cache.histories[i].result
			if result <= 5 || result >= 22 {
				latest = make(map[int]struct{})
				return nil
			}
		}
	} else {
		if bets, err := bet28(cache, nextIssue, surplus, SN10, spaces, float64(cache.user.gold)); err != nil {
			return err
		} else {
			latest = bets
		}

		return nil
	}

	size := len(cache.histories)
	r1 := cache.histories[size-1].result
	r2 := cache.histories[size-2].result

	if r1 < 10 || r1 > 17 {
		latest = make(map[int]struct{})
		return nil
	}

	if r2 >= 10 && r2 <= 17 {
		latest = make(map[int]struct{})
		return nil
	}

	var total, coverage int

	latest = make(map[int]struct{})
	target := getTarget(spaces)
	for _, result := range SN28 {
		if _, ok := target[result]; !ok {
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【❤️ %02d】，标准赔率【%-7.2f】，间隔次数【%-4d】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), spaces[result], betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus-total, float64(coverage)/10)

	return nil
}

func getTarget(spaces map[int]int) map[int]struct{} {
	type Space struct {
		Result int
		Space  int
		Rate   float64
	}

	newSpaces := make([]Space, 0)
	for result, space := range spaces {
		rate := float64(space) / (1000 / float64(stds[result]))
		newSpaces = append(newSpaces, Space{Result: result, Space: space, Rate: rate})
	}

	sort.Slice(newSpaces, func(i, j int) bool {
		return newSpaces[i].Rate > newSpaces[j].Rate
	})

	target := make(map[int]struct{})
	for _, newSpace := range newSpaces {
		if newSpace.Result <= 5 || newSpace.Result >= 22 {
			continue
		}

		if newSpace.Rate < 2.0 {
			target[newSpace.Result] = struct{}{}
		}
	}

	return target
}

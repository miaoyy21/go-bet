package xmd

import (
	"log"
	"sort"
	"strconv"
)

var latest = make(map[int]struct{})

var wins int
var fails int
var zWins int
var zFails int

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
		log.Printf("第【✊ %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0

			//rate = rate - 0.25
			//if rate < 1.0 {
			//	rate = 1.0
			//}

			zWins++
			log.Printf("第【👍 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus)
		} else {
			wins = 0
			fails++

			//if rate < 5.0 {
			//	rate = rate + 1
			//}

			zFails++
			log.Printf("第【👀 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus)
		}
	}

	// 12期 00~04、23～27
	for i := len(cache.histories) - 1; i >= len(cache.histories)-12; i-- {
		result := cache.histories[i].result
		if result <= 4 || result >= 23 {
			if ns, err := bet28(cache, nextIssue, surplus, SN8, float64(cache.user.gold)); err != nil {
				return err
			} else {
				latest = ns
			}

			return nil
		}
	}

	for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
		result := cache.histories[i].result
		if result <= 5 || result >= 22 {
			if _, err := bet28(cache, nextIssue, surplus, SN28, 2000); err != nil {
				return err
			}
			latest = make(map[int]struct{})

			return nil
		}
	}

	for i := len(cache.histories) - 1; i >= len(cache.histories)-4; i-- {
		result := cache.histories[i].result
		if result <= 6 || result >= 21 {
			if _, err := bet28(cache, nextIssue, surplus, SN28, 1000); err != nil {
				return err
			}
			latest = make(map[int]struct{})

			return nil
		}
	}

	var total, coverage int

	latest = make(map[int]struct{})
	target := getTarget(cache)
	for _, result := range SN28 {
		if _, ok := target[result]; !ok {
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【❤️ %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus-total, float64(coverage)/10)

	return nil
}

func getTarget(cache *Cache) map[int]struct{} {
	type Space struct {
		Result int
		Space  int

		Rate float64
	}

	spaces := make(map[int]int)
	for i := len(cache.histories) - 1; i >= 0; i-- {
		if _, ok := spaces[cache.histories[i].result]; ok {
			continue
		}

		spaces[cache.histories[i].result] = len(cache.histories) - i
	}

	newSpaces := make([]Space, 0, len(spaces))
	for result, space := range spaces {
		rate := float64(space) / (1000 / float64(stds[result]))
		newSpaces = append(newSpaces, Space{Result: result, Space: space, Rate: rate})
	}
	sort.Slice(newSpaces, func(i, j int) bool {
		return newSpaces[i].Rate > newSpaces[j].Rate
	})

	var n1, n2, n3 int
	target := make(map[int]struct{})
	for _, newSpace := range newSpaces {
		if newSpace.Result >= 10 && newSpace.Result <= 17 {
			if n1 < 8 && newSpace.Rate > 2.0 {
				n1++
				continue
			}
		} else if newSpace.Result <= 6 || newSpace.Result >= 21 {
			if n2 < 14 {
				n2++
				continue
			}
		} else {
			if n3 < 6 && newSpace.Rate > 2.0 {
				n3++
				continue
			}
		}

		target[newSpace.Result] = struct{}{}
	}

	return target
}

package xmd

import (
	"log"
	"sort"
	"strconv"
)

var latest = make(map[int]struct{})

var rate = 1.0
var wins int
var fails int
var zWins int
var zFails int

func analysis(cache *Cache) error {
	if err := cache.Sync(2000); err != nil {
		return err
	}

	nextIssue := strconv.Itoa(cache.issue + 1)

	// 当前账户余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	size := len(cache.histories)
	r1 := cache.histories[size-1].result
	r2 := cache.histories[size-2].result

	if r1 < 10 || r1 > 17 {
		latest = make(map[int]struct{})
		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，不符合投注条件A ...\n", cache.issue, cache.result, surplus)
		return nil
	}

	if r2 >= 10 && r2 <= 17 {
		latest = make(map[int]struct{})
		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，不符合投注条件B ...\n", cache.issue, cache.result, surplus)
		return nil
	}

	// 输出
	rate0 := 1.0
	if len(latest) == 0 {
		log.Printf("第【✊ %d %03d/%03d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0

			if wins >= 2 {
				if wins%3 == 2 {
					if rate <= 2.5 {
						rate = rate + 0.50
					}
				} else {
					rate = rate - 0.25
					if rate < 1.0 {
						rate = 1.0
					}
				}

				rate0 = rate
			}

			zWins++
			log.Printf("第【👍 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，倍率【%.4f -> %.4f】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus, rate, rate0)
		} else {
			wins = 0
			fails++

			zFails++
			log.Printf("第【👀 %d %03d/%03d】期：开奖结果【%d】，余额【%d】，倍率【%.4f -> %.4f】，开始执行分析 ...\n", cache.issue, zWins, zFails, cache.result, surplus, rate, rate0)
		}
	}

	var total, coverage int

	latest = make(map[int]struct{})
	for result := range getTarget(cache) {
		betGold := int(rate0 * float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)

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
			// [10,17]
			if n1 < 2 {
				n1++
				continue
			}
		} else if newSpace.Result <= 5 || newSpace.Result >= 22 {
			// [00,05] [22,27]
			if n2 < 12 {
				n2++
				continue
			}
		} else {
			// [06,09] [18,21]
			if n3 < 2 {
				n3++
				continue
			}
		}

		target[newSpace.Result] = struct{}{}
	}

	return target
}

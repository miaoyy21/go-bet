package xmd

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var latest = make(map[int]struct{})
var times = 1

var rate float64
var wins int
var fails int
var maxFails int

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
		rate = 1.0
		log.Printf("【%-4d  %d】第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", times, maxFails, cache.issue, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			fails = 0

			rate = 1.0
			log.Printf("【%-4d  %d】第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, maxFails, cache.issue, wins, cache.result, surplus, rate)
		} else {
			wins = 0
			fails++
			if fails > maxFails {
				maxFails = fails
			}

			// 最大失败次数
			if fails >= 6 {
				rate = 1.0
			}

			rate = rate * (1.375 + math.Pow(0.75, float64(fails)-1))
			log.Printf("【%-4d  %d】第【👀 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, maxFails, cache.issue, fails, cache.result, surplus, rate)
		}
	}

	p50s, sp50s, coverage := getP50()
	log.Printf("第【%s】期：随机数字【🧐 %s】\n", nextIssue, strings.Join(sp50s, ","))

	var total int

	for i := 0; i <= 27; i++ {
		if _, ok := p50s[i]; !ok {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，投注金额【    -】\n", nextIssue, i, 1000.0/float64(stds[i]))
			continue
		}

		betGold := int(rate * float64(cache.user.gold) * float64(stds[i]) / 1000)
		if err := hPostBet(nextIssue, betGold, i, cache.user); err != nil {
			return err
		}

		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, i, 1000.0/float64(stds[i]), betGold)
		total = total + betGold
	}
	latest = p50s

	times++
	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	return nil
}

func getP50() (map[int]struct{}, []string, int) {
	src := rand.NewSource(time.Now().UnixNano())

	coverage := 0
	p50s, sp50s := make(map[int]struct{}), make([]string, 0)
	for {
		d1 := rand.New(src).Intn(10)
		d2 := rand.New(src).Intn(10)
		d3 := rand.New(src).Intn(10)

		d := d1 + d2 + d3
		if _, ok := p50s[d]; ok {
			continue
		}

		p50s[d] = struct{}{}
		sp50s = append(sp50s, fmt.Sprintf("%02d", d))
		coverage = coverage + stds[d]
		if coverage > 500 {
			break
		}
	}

	sort.Strings(sp50s)
	return p50s, sp50s, coverage
}

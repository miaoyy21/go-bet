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

var source rand.Source

var rate float64
var wins int
var fails int
var mWins int
var mFails int
var xWins int
var xFails int
var nFails map[int]int

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

	if len(latest) == 0 {
		wins, fails, rate = 0, 0, 1.0

		seed := time.Now().UnixNano()
		source = rand.NewSource(seed)
		nFails = make(map[int]int)
		log.Printf("【%-4d】第【✊ %d】期：开奖结果【%d】，余额【%d】，初始化随机种子【%d】，开始执行分析 ...\n", times, cache.issue, cache.result, surplus, seed)
	} else {
		if _, exists := latest[cache.result]; exists {
			wins++
			if wins > mWins {
				mWins = wins
			}

			if fails > 0 {
				nFails[fails]++
			}

			if rate >= 2.0 {
				rate = rate / 2.0
			} else {
				rate = math.Pow(1.25, float64(fails))
			}

			fails = 0
			xWins++
			log.Printf("【%-4d W(%d,%d) F(%d,%d)】第【👍 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, xWins, mWins, xFails, mFails, cache.issue, wins, cache.result, surplus, rate)
		} else {
			if fails <= 3 {
				for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
					result := cache.histories[i].result
					if result <= 5 || result >= 22 {
						latest = make(map[int]struct{})
						log.Printf("【%-4d W(%d,%d) F(%d,%d)】第【✊ %d %02d】期：开奖结果【%d】，余额【%d】，重置竞猜状态 ...\n", times, xWins, mWins, xFails, mFails, cache.issue, fails, cache.result, surplus)
						return nil
					}
				}

				if fails == 0 {
					seed := time.Now().UnixNano()
					source = rand.NewSource(seed) // 重新初始化随机种子

					log.Printf("【%-4d】第【%d】期：重新初始化随机种子【%d】 ...\n", times, cache.issue, seed)
				}
			} else {
				rate = rate * 0.75
			}

			fails++
			if fails > mFails {
				mFails = fails
			}

			rate = rate * (1.35 + math.Pow(0.65, float64(fails)-1))

			wins = 0
			xFails++
			log.Printf("【%-4d W(%d,%d) F(%d,%d)】第【👀 %d %02d】期：开奖结果【%d】，余额【%d】，投注倍率【%.3f】，开始执行分析 ...\n", times, xWins, mWins, xFails, mFails, cache.issue, fails, cache.result, surplus, rate)
		}
	}

	p50s, sp50s, coverage := getP50()
	nfs := make([]string, 0)
	for i := 1; ; i++ {
		if len(nfs) == len(nFails) {
			break
		}

		if n, ok := nFails[i]; ok {
			nfs = append(nfs, fmt.Sprintf("%d:%d", i, n))
		}
	}
	log.Printf("第【%s】期：随机数字【🧐 %s】，分布情况【%s】......\n", nextIssue, strings.Join(sp50s, ","), strings.Join(nfs, " , "))

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
	coverage := 0
	p50s, sp50s := make(map[int]struct{}), make([]string, 0)
	for {
		d1 := rand.New(source).Intn(10)
		d2 := rand.New(source).Intn(10)
		d3 := rand.New(source).Intn(10)

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

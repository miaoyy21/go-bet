package xmd

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var latest = make(map[int]struct{})

func analysis(cache *Cache) error {
	if err := cache.Sync(200); err != nil {
		return err
	}

	issue := strconv.Itoa(cache.issue + 1)
	if !cache.user.isBetMode {
		time.Sleep(2 * time.Second)
	}

	// 设定是否进行投注
	if _, ok := latest[cache.result]; !ok && len(latest) > 0 && rand.Float32() <= 0.20 {
		latest = make(map[int]struct{})
		log.Printf("😤😤😤 第【%s】期：上一期开奖结果【%d】，由于投注失利，随机选择不进行投注 >>>>>>>>>> \n", issue, cache.result)

		return nil
	}

	// 当前账户可用余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	// 计算每个数字的间隔期数和当前赔率
	rts, exp, dev, err := RiddleDetail(cache.user, issue)
	if err != nil {
		return err
	}

	// 显示当前中奖情况
	latest = make(map[int]struct{})
	log.Printf("⭐️⭐️⭐️ 第【%d】期：开奖结果【%d】，下期「预估期望【%6.4f】，预估平均方差【%6.4f】」，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp, dev, surplus)

	// 仅投注当前赔率大于标准赔率的数字
	bets := make(map[int]float64)
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]

		var rx float64
		if r1/r0 >= 1.0 {
			rx = 1.0
		} else {
			rx = (r1/r0 - 0.98) * 100.0 / 2.0
		}

		if rx <= 0.01 {
			log.Printf("第【%s】期：竞猜数字【   %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
			continue
		}

		if rx >= 1.0 {
			latest[result] = struct{}{}
			log.Printf("第【%s】期：竞猜数字【 H %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
		} else {
			log.Printf("第【%s】期：竞猜数字【 L %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
		}

		bets[result] = rx
	}

	// 使用设定的投注模式
	if cache.user.isBetMode {
		if err := betMode(cache, issue, bets); err != nil {
			return err
		}
	} else {
		// 计算投注系数
		mrx := mrxFn(dev, cache.money)

		if err := betSingle(cache, issue, mrx, bets); err != nil {
			return err
		}
	}

	return nil
}

// 使用基于投注模式方式投注
func betMode(cache *Cache, issue string, bets map[int]float64) error {
	// 查询用户设定的投注模式
	m1Gold, err := hCustomModes(cache.user)
	if err != nil {
		return err
	}

	// 数字排序
	rs := make([]int, 0, len(bets))
	for result := range bets {
		rs = append(rs, result)
	}
	sort.Ints(rs)
	log.Printf("第【%s】期：预投注数字【%s】 >>>>>>>>>> \n", issue, fmtIntSlice(rs))

	// 确定投注模式ID
	modeId, modeName := modeFn(bets, 400)
	if modeId > 0 {
		if err := hModesBetting(issue, modeId, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：使用投注模式【%s】 >>>>>>>>>> \n", issue, modeName)
	}

	// 投注模式之外的数字
	extras := extraFn(modeId, m1Gold, bets)
	if len(extras) > 0 {
		log.Printf("第【%s】期：额外投注数字【%s】>>>>>>>>>> \n", issue, fmtIntSlice(m2sFn(extras)))
	}

	// 使用单数字投注模式，必须使用其提供的标准投注金额
	stdBets := []int{200000, 50000, 10000, 5000, 2000, 1000, 500}
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

	// 单数字投注
	for _, stdBet := range stdBets {
		if len(betMaps[stdBet]) > 0 {
			log.Printf("第【%s】期：押注金额【%-6d】，押注数字【%s】，投注成功 >>>>>>>>>> \n", issue, stdBet, fmtIntSlice(betMaps[stdBet]))
		}

		for _, result := range betMaps[stdBet] {
			if err := hBetting1(issue, stdBet, result, cache.user); err != nil {
				return err
			}
		}
	}

	return nil
}

func betSingle(cache *Cache, issue string, mrx float64, bets map[int]float64) error {
	// 查询用户设定的投注模式
	m1Gold, err := hCustomModes(cache.user)
	if err != nil {
		return err
	}

	for _, result := range SN28 {
		if _, ok := bets[result]; !ok || bets[result] <= 0.001 {
			continue
		}

		betGold := int(mrx * bets[result] * float64(2*m1Gold) * float64(stds[result]) / 1000)
		if err := hBetting1(issue, betGold, result, cache.user); err != nil {
			return err
		}
	}

	return nil
}

func mrxFn(dev float64, money int) float64 {
	mrx := 1.0

	// 投注金额 系数设定
	if money < 2<<24 {
		// 33,554,432
		mrx = 0.4
	} else if money < 2<<25 {
		// 67,108,864
		mrx = 0.7
	} else if money < 2<<26 {
		// 134,217,728
		mrx = 0.9
	} else {
		if money > 2<<28 {
			// 536,870,912
			mrx = 1.4
		} else if money > 2<<27 {
			// 268,435,456
			mrx = 1.2
		}
	}

	// 赔率标准方差 系数设定
	if dev > 1.1 {
		mrx = mrx * 1.30
	} else if dev > 1.05 {
		mrx = mrx * 1.20
	} else if dev > 1.00 {
		mrx = mrx * 1.10
	}

	return mrx
}

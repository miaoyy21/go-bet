package xmd

import (
	"log"
	"sort"
	"strconv"
	"time"
)

func analysis(cache *Cache) error {
	if err := cache.Sync(200); err != nil {
		return err
	}

	issue := strconv.Itoa(cache.issue + 1)

	// 当前账户可用余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	// 计算每个数字的间隔期数和当前赔率
	rts, exp, _, err := RiddleDetail(cache.user, issue)
	if err != nil {
		return err
	}

	// 显示当前中奖情况
	log.Printf("⭐️⭐️⭐️ 第【%d】期：开奖结果【%d】，下期预估期望返奖【%.2f%%】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, surplus)

	// 仅投注当前赔率大于标准赔率的数字
	bets := make(map[int]float64)
	x1s := make(map[int]struct{})
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]

		var rx float64
		if r1/r0 >= 1.0 {
			rx = 1.0
			x1s[result] = struct{}{}
		} else {
			rx = (r1/r0 - 0.98) * 100.0 / 2.0
		}

		if rx <= 0.01 {
			log.Printf("第【%s】期：竞猜数字【   %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
			continue
		}

		if rx >= 1.0 {
			log.Printf("第【%s】期：竞猜数字【 H %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
		} else {
			log.Printf("第【%s】期：竞猜数字【 L %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", issue, result, r0, r1, r1/r0)
		}

		bets[result] = rx
	}

	// 数字排序
	rs := make([]int, 0, len(bets))
	for result := range bets {
		rs = append(rs, result)
	}
	sort.Ints(rs)
	log.Printf("第【%s】期：预投注数字【%s】 >>>>>>>>>> \n", issue, fmtIntSlice(rs))

	// 确定投注模式ID
	modeId, modeName := modeFn(bets, 350)
	if modeId <= 0 {
		log.Printf("第【%s】期：无法确定投注模式【%s】 >>>>>>>>>> \n", issue, modeName)
		return nil
	}

	// 使用模式投注
	if err := hModesBetting(issue, modeId, cache.user); err != nil {
		return err
	}
	log.Printf("第【%s】期：使用投注模式【%s】 >>>>>>>>>> \n", issue, modeName)
	time.Sleep(500 * time.Millisecond)

	// 查询用户设定的投注模式
	mGold, err := hCustomModes(cache.user)
	if err != nil {
		return err
	}

	// 投注模式之外的数字
	extras := extraFn(modeId, mGold, x1s)
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

			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

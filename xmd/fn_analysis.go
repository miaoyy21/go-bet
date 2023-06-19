package xmd

import (
	"log"
	"sort"
	"strconv"
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
	modeId, modeName := modeFn(bets, 400)
	//modeId, modeName = 0, "暂时不启用模式" // TODO

	// 投注成功
	if modeId > 0 {
		if err := hModesBetting(issue, modeId, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：投注模式【%s】，投注成功 >>>>>>>>>> \n", issue, modeName)
	} else {
		log.Printf("第【%s】期：无法确定投注模式【%s】 >>>>>>>>>> \n", issue, modeName)
	}

	mGold, err := hCustomModes(cache.user)
	if err != nil {
		return err
	}

	// 其他的数字
	extras := extraFn(modeId, mGold, x1s)
	log.Printf("第【%s】期：额外投注【%s】，投注成功 >>>>>>>>>> \n", issue, fmtIntSlice(m2sFn(extras)))

	for result, betGold := range extras {
		if err := hBetting1(issue, betGold, result, cache.user); err != nil {
			return err
		}
	}

	return nil
}

func doBet() {
	//stdBets := []int{500, 1000, 2000, 5000, 10000, 50000}
}

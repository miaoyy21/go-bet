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

	nextIssue := strconv.Itoa(cache.issue + 1)

	// 当前账户可用余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	// 计算每个数字的间隔期数和当前赔率
	rts, exp, _, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}

	// 显示当前中奖情况
	log.Printf("⭐️⭐️⭐️ 第【%d】期：开奖结果【%d】，下期预估期望返奖【%.2f%%】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, surplus)

	// 仅投注当前赔率大于标准赔率的数字
	bets := make(map[int]float64)
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]

		var rx float64
		if r1/r0 >= 1.0 {
			rx = 1.0
		} else {
			rx = (r1/r0 - 0.99) * 100
		}

		if rx <= 0.01 {
			log.Printf("第【%s】期：竞猜数字【   %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", nextIssue, result, r0, r1, r1/r0)
			continue
		}

		log.Printf("第【%s】期：竞猜数字【 √ %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】\n", nextIssue, result, r0, r1, r1/r0)

		bets[result] = rx
	}

	// 数字排序
	rs := make([]int, 0, len(bets))
	for result := range bets {
		rs = append(rs, result)
	}
	sort.Ints(rs)
	log.Printf("第【%s】期：预投注数字【%s】 >>>>>>>>>> \n", nextIssue, fmtIntSlice(rs))

	// 确定投注模式ID
	modeId, modeName := parseModeId(bets)
	if modeId <= 0 {
		log.Printf("第【%s】期：无法确定投注模式ID【%d】 >>>>>>>>>> \n", nextIssue, modeId)
		return nil
	}

	// 投注成功
	if err := hModesBetting(nextIssue, modeId, cache.user); err != nil {
		return err
	}
	log.Printf("第【%s】期：投注模式【%s】，投注成功 >>>>>>>>>> \n", nextIssue, modeName)

	return nil
}

func parseModeId(bets map[int]float64) (int, string) {
	var m1, m2, m3, m4, m5, m6, m7, m8 int
	for result := range bets {
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

	m5 = int(float64(m5*44) / 56)
	m6 = int(float64(m5*56) / 44)
	log.Printf("模式权重：大数【%d】, 小数【%d】, 奇数【%d】, 偶数【%d】, 中数【%d】, 边数【%d】, 大尾数【%d】, 小尾数【%d】 \n", m1, m2, m3, m4, m5, m6, m7, m8)

	if m1 >= 350 && m1 >= m2 && m1 >= m3 && m1 >= m4 && m1 >= m5 && m1 >= m6 && m1 >= m7 && m1 >= m8 {
		return 1, "大数"
	}

	if m2 >= 350 && m2 >= m1 && m2 >= m3 && m2 >= m4 && m2 >= m5 && m2 >= m6 && m2 >= m7 && m2 >= m8 {
		return 2, "小数"
	}

	if m3 >= 350 && m3 >= m1 && m3 >= m2 && m3 >= m4 && m3 >= m5 && m3 >= m6 && m3 >= m7 && m3 >= m8 {
		return 3, "奇数"
	}

	if m4 >= 350 && m4 >= m1 && m4 >= m2 && m4 >= m3 && m4 >= m5 && m4 >= m6 && m4 >= m7 && m4 >= m8 {
		return 4, "偶数"
	}

	if m5 >= 350 && m5 >= m1 && m5 >= m2 && m5 >= m3 && m5 >= m4 && m5 >= m6 && m5 >= m7 && m5 >= m8 {
		return 5, "中数"
	}

	if m6 >= 350 && m6 >= m1 && m6 >= m2 && m6 >= m3 && m6 >= m4 && m6 >= m5 && m6 >= m7 && m6 >= m8 {
		return 6, "边数"
	}

	if m7 >= 350 && m7 >= m1 && m7 >= m2 && m7 >= m3 && m7 >= m4 && m7 >= m5 && m7 >= m6 && m7 >= m8 {
		return 7, "大尾数"
	}

	if m8 >= 350 && m8 >= m1 && m8 >= m2 && m8 >= m3 && m8 >= m4 && m8 >= m5 && m8 >= m6 && m8 >= m7 {
		return 8, "小尾数"
	}

	return 0, "未知"
}

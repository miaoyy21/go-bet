package xmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
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

	// 保存投注相关参数
	if xSurplus > 0 && cache.issue == issue {
		xRt := xRts[cache.result] / (1000.0 / float64(stds[cache.result]))
		query := fmt.Sprintf("INSERT INTO logs_%s(time, issue, result, money, member, user_gold,  exp, dev, rt, bet_gold, win_gold, gold) VALUES (?,?,?,?,?,?, ?,?,?,?,?,?)", cache.user.id)
		if _, err := cache.db.Exec(query,
			time.Now().Format("2006-01-02 15:04:05.999"), cache.issue, cache.result, cache.money, cache.member, xUserGold,
			xExp, xDev, xRt, xBetGold, surplus-xSurplus, surplus,
		); err != nil {
			return err
		}
	}
	issue = cache.issue + 1
	xSurplus = surplus
	xBetGold = 0
	xUserGold = cache.user.gold

	// 计算每个数字的间隔期数和当前赔率
	rts, exp, dev, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}
	xRts = rts
	xExp = exp
	xDev = dev

	// 显示当前中奖情况
	if len(latest) == 0 {
		log.Printf("⭐️⭐️⭐️ 第【✊ %d】期：开奖结果【%d】，下期预估期望返奖【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			log.Printf("⭐️⭐️⭐️ 第【👍 %d】期：开奖结果【%d】，下期预估期望返奖【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
		} else {
			log.Printf("⭐️⭐️⭐️ 第【👀 %d】期：开奖结果【%d】，下期预估期望返奖【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
		}
	}

	// 本期期望返奖大于设定的期望返奖时，才进行投注
	if dev <= cache.dev {
		latest = make(map[int]int)

		xBetGold = 0
		log.Printf("第【%s】期：赔率标准方差【%.2f】不足%.3f，放弃投注 >>>>>>>>>> \n", nextIssue, dev, cache.dev)
		return nil
	}

	// 投注金额 系数设定
	if cache.money < 2<<23 {
		// 16,777,216
		xUserGold = int(float64(xUserGold) * 0.2)
	} else if cache.money < 2<<24 {
		// 33,554,432
		xUserGold = int(float64(xUserGold) * 0.4)
	} else if cache.money < 2<<25 {
		// 67,108,864
		xUserGold = int(float64(xUserGold) * 0.7)
	} else if cache.money < 2<<26 {
		// 134,217,728
		xUserGold = int(float64(xUserGold) * 0.9)
	} else {
		// 268,435,456
		if cache.money > 2<<27 {
			xUserGold = int(float64(xUserGold) * 1.2)
		}
	}

	// 赔率标准方差 系数设定
	if dev > 1.1 {
		xUserGold = int(float64(xUserGold) * 1.30)
	} else if dev > 1.05 {
		xUserGold = int(float64(xUserGold) * 1.20)
	} else if dev > 1.00 {
		xUserGold = int(float64(xUserGold) * 1.10)
	}

	// 以万为单位进行投注
	if xUserGold > 100000 {
		xUserGold = xUserGold / 10000 * 10000
	}

	// 仅投注当前赔率大于标准赔率的数字
	latest = make(map[int]int)
	coverage := 0

	//spaces := SpaceFn(cache)
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]

		var rx float64
		if r1/r0 >= 1.0 {
			rx = 1.0
		} else {
			rx = (r1/r0 - 0.99) * 100
		}

		betGold := int(rx * float64(xUserGold) * float64(stds[result]) / 1000)
		if betGold <= 0 {
			//log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【     -】\n", nextIssue, result, r0, r1, r1/r0, spaces[result])
			continue
		}

		latest[result] = betGold
		coverage = coverage + int(float64(stds[result])*rx)
	}

	if float64(coverage) < 125 {
		latest = make(map[int]int)

		xBetGold = 0
		log.Printf("第【%s】期：覆盖率【%.2f%%】不足%.2f%%，放弃投注 >>>>>>>>>> \n", nextIssue, float64(coverage)/10, 12.5)
		return nil
	} else if float64(coverage) > 875 {
		latest = make(map[int]int)

		xBetGold = 0
		log.Printf("第【%s】期：覆盖率【%.2f%%】超过%.2f%%，放弃投注 >>>>>>>>>> \n", nextIssue, float64(coverage)/10, 87.5)
		return nil
	}

	total := 0
	rs := make([]int, 0, len(latest))
	for result, betGold := range latest {
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}

		//r0 := 1000.0 / float64(stds[result])
		//r1 := rts[result]
		//log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【% 6d】\n", nextIssue, result, r0, r1, r1/r0, spaces[result], betGold)

		rs = append(rs, result)
		total = total + betGold
		time.Sleep(25 * time.Millisecond)
	}
	sort.Ints(rs)

	// 显示投注的汇总结果
	surplus = surplus - total
	xBetGold = total
	log.Printf("第【%s】期：投注金额【%d】，投注数字【%s】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, fmtIntSlice(rs), surplus, float64(coverage)/10)

	return nil
}

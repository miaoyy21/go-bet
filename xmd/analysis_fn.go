package xmd

import (
	"fmt"
	"log"
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
	spaces := SpaceFn(cache)
	rts, exp, dev, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}
	xRts = rts
	xExp = exp
	xDev = dev

	// 显示当前中奖情况
	if len(latest) == 0 {
		log.Printf("⭐️⭐️⭐️ 第【✊ %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			log.Printf("⭐️⭐️⭐️ 第【👍 %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
		} else {
			log.Printf("⭐️⭐️⭐️ 第【👀 %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, exp*100, cache.user.gold, surplus)
		}
	}

	// 本期返奖率大于设定的返奖率时，才进行投注
	if exp <= cache.exp {
		latest = make(map[int]struct{})

		xBetGold = 0
		log.Printf("第【%s】期：预估返奖率【%.2f%%】不足%.2f%%，放弃投注 >>>>>>>>>> \n", nextIssue, exp*100, cache.exp*100)
		return nil
	}

	// 仅投注当前赔率大于标准赔率的数字
	latest = make(map[int]struct{})
	total, coverage := 0, 0
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]
		if r1 < r0 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【     -】\n", nextIssue, result, r0, r1, r1/r0, spaces[result])
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)

		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【% 6d】\n", nextIssue, result, r0, r1, r1/r0, spaces[result], betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	// 显示投注的汇总结果
	surplus = surplus - total
	xBetGold = total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	return nil
}
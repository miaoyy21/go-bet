package xmd

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func analysisA1(cache *Cache) error {
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
		query := fmt.Sprintf("INSERT INTO logs_%s(time, issue, result, user_gold,  rx, rt, bet_gold, win_gold, gold) VALUES (?,?,?,?, ?,?,?,?,?)", cache.user.id)
		if _, err := cache.db.Exec(query,
			time.Now().Format("2006-01-02 15:04"), cache.issue, cache.result, xUserGold,
			xRx, xRt, xBetGold, surplus-xSurplus, surplus,
		); err != nil {
			return err
		}
	}
	issue = cache.issue + 1
	xSurplus = surplus
	xUserGold = cache.user.gold

	// 计算每个数字的间隔期数和当前赔率
	spaces := SpaceFn(cache)
	rts, rx, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}
	xRts = rts
	xRx = rx

	// 显示当前中奖情况
	if len(latest) == 0 {
		log.Printf("⭐️⭐️⭐️ 第【✊ %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, cache.user.gold, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			log.Printf("⭐️⭐️⭐️ 第【👍 %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, cache.user.gold, surplus)
		} else {
			log.Printf("⭐️⭐️⭐️ 第【👀 %d】期：开奖结果【%d】，下期预估返奖率【%.2f%%】，下期基础投注【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, cache.user.gold, surplus)
		}
	}

	// 本期返奖率大于设定的返奖率时，才进行投注
	if rx <= cache.rx {
		latest = make(map[int]struct{})

		if cache.IsExtra() {
			log.Printf("️第【%s】期：预估返奖率【%.2f%%】不足%.2f%%，进行投注 20,000 >>>>>>>>>> \n", nextIssue, rx*100, cache.rx*100)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, rts, 20000); err != nil {
				return err
			}

			xBetGold = 20000
			return nil
		}

		xBetGold = 0
		log.Printf("第【%s】期：预估返奖率【%.2f%%】不足%.2f%%，放弃投注 >>>>>>>>>> \n", nextIssue, rx*100, cache.rx*100)
		return nil
	}

	// 本期是否存在当前赔率大于标准赔率的倍数的数字
	var c0 bool
	for _, result := range SN28 {
		if rts[result] >= 1000.0*cache.wx/float64(stds[result]) {
			c0 = true
			break
		}
	}

	// 当本期存在当前赔率大于标准赔率10%的数字时，才进行投注
	if !c0 {
		latest = make(map[int]struct{})

		if cache.IsExtra() {
			log.Printf("第【%s】期：不存在实际赔率超过%.2f%%的数字，仅投注 20,000 >>>>>>>>>> \n", nextIssue, cache.wx*100-100)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, rts, float64(20000)); err != nil {
				return err
			}

			xBetGold = 20000
			return nil
		}

		xBetGold = 0
		log.Printf("第【%s】期：不存在实际赔率超过%.2f%%的数字，放弃投注 >>>>>>>>>> \n", nextIssue, cache.wx*100-100)
		return nil
	}

	// 仅投注当前赔率大于标准赔率的数字
	latest = make(map[int]struct{})
	total, coverage := 0, 0
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]
		if r1 < r0 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【    -】\n", nextIssue, result, r0, r1, r1/r0, spaces[result])
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-6.4f】，间隔次数【%-4d】，投注金额【% 5d】\n", nextIssue, result, r0, r1, r1/r0, spaces[result], betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	// 显示投注的汇总结果
	surplus = surplus - total
	xBetGold = total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	// 如果处于活动奖励期间（每日投注金额超过2万达到指定的次数），按照活动要求不足2万投注金额
	if total < 20000 {
		if cache.IsExtra() {
			log.Printf("第【%s】期：投注金额不足，进行不足至 20,000  >>>>>>>>>> \n", nextIssue)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, rts, float64(20000-total)); err != nil {
				return err
			}

			xBetGold = 20000
		}
	}

	return nil
}

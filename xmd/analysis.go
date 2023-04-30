package xmd

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

var latest = make(map[int]struct{})
var xSurplus int
var xBetGold int

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

	if xSurplus > 0 {
		query := fmt.Sprintf("%s INTO logs(issue, time, bet_gold, win_gold, gold) VALUES (?,?,?,?,?)", "INSERT")
		if _, err := cache.db.Exec(query, cache.issue, time.Now().Format("2006-01-02 15:04"), xBetGold, surplus-xSurplus, surplus); err != nil {
			return err
		}
	}
	xSurplus = surplus

	spaces := SpaceFn(cache)
	rts, _, rx, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}

	// 输出
	if len(latest) == 0 {
		log.Printf("⭐️⭐️⭐️ 第【✊ %d】期：开奖结果【%d】，下一期预估返奖率【%.2f%%】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			log.Printf("⭐️⭐️⭐️ 第【👍 %d】期：开奖结果【%d】，下一期预估返奖率【%.2f%%】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, surplus)
		} else {
			log.Printf("⭐️⭐️⭐️ 第【👀 %d】期：开奖结果【%d】，下一期预估返奖率【%.2f%%】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, rx*100, surplus)
		}
	}

	// 返奖率小于0.95
	if rx < 0.95 {
		latest = make(map[int]struct{})
		if cache.IsExtra() {
			log.Printf("️第【%s】期：预估返奖率【%.2f%%】不足95%%，进行投注 20,000 >>>>>>>>>> \n", nextIssue, rx*100)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, 20000); err != nil {
				return err
			}

			xBetGold = 20000
			return nil
		}

		log.Printf("第【%s】期：预估返奖率【%.2f%%】不足95%%，仅投注 1,000 >>>>>>>>>> \n", nextIssue, rx*100)
		if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, 1000); err != nil {
			return err
		}

		xBetGold = 1000
		return nil
	}

	// 先初步看看赔率系数，是不是值得投注
	var c0 bool
	for _, result := range SN28 {
		if rts[result] > 1000.0*1.10/float64(stds[result]) {
			c0 = true
			break
		}
	}

	if !c0 {
		if cache.IsExtra() {
			log.Printf("第【%s】期：赔率超过5%%的覆盖率【0%%】，仅投注 20,000 >>>>>>>>>> \n", nextIssue)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, float64(20000)); err != nil {
				return err
			}

			xBetGold = 20000
			return nil
		}

		log.Printf("第【%s】期：赔率超过5%%的覆盖率【0%%】，仅投注 1,000 >>>>>>>>>> \n", nextIssue)
		if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, 1000); err != nil {
			return err
		}

		xBetGold = 1000
		return nil
	}

	latest = make(map[int]struct{})
	total, coverage := 0, 0
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]
		if r1 < r0 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-4.2f】，间隔次数【%-4d】，投注金额【    -】\n", nextIssue, result, r0, r1, r1/r0, spaces[result])
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-4.2f】，间隔次数【%-4d】，投注金额【% 5d】\n", nextIssue, result, r0, r1, r1/r0, spaces[result], betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	surplus = surplus - total
	xBetGold = total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	// 不足2万
	if total < 20000 {
		if cache.IsExtra() {
			log.Printf("第【%s】期：投注金额不足，进行不足至 20,000  ********** \n", nextIssue)
			if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, float64(20000-total)); err != nil {
				return err
			}

			xBetGold = 20000
		}
	}

	return nil
}

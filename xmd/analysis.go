package xmd

import (
	"log"
	"strconv"
)

var latest = make(map[int]struct{})

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

	rts, _, _, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}

	// 输出
	if len(latest) == 0 {
		log.Printf("第【✊ %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
	} else {
		if _, exists := latest[cache.result]; exists {
			log.Printf("第【👍 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
		} else {
			log.Printf("第【👀 %d】期：开奖结果【%d】，余额【%d】，开始执行分析 ...\n", cache.issue, cache.result, surplus)
		}
	}

	latest = make(map[int]struct{})
	total, coverage := 0, 0
	for _, result := range SN28 {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]
		if r1 < r0 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，投注金额【    -】\n", nextIssue, result, r0, r1)
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, r0, r1, betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	return nil
}

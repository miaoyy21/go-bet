package xmd

import (
	"log"
	"strconv"
	"time"
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

	// 添加
	if time.Now().Minute() == 0 {
		cache.hGolds = append(cache.hGolds, HGold{Time: time.Now().Format("2006-01-02 15:04"), Gold: surplus})

		log.Printf("输出每小时的金额情况如下：\n")
		for _, hGold := range cache.hGolds {
			log.Printf("【%s】：【% 9d】\n", hGold.Time, hGold.Gold)
		}
	}

	rts, _, _, err := RiddleDetail(cache.user, nextIssue)
	if err != nil {
		return err
	}

	spaces := SpaceFn(cache)

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
		if r1 < r0*1.05 {
			log.Printf("第【%s】期：竞猜数字【👀 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，间隔次数【%-4d】，投注金额【    -】\n", nextIssue, result, r0, r1, spaces[result])
			continue
		}

		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}
		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，间隔次数【%-4d】，投注金额【% 5d】\n", nextIssue, result, r0, r1, spaces[result], betGold)

		latest[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}

	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", nextIssue, total, surplus, float64(coverage)/10)

	// 不足2万
	if total < 20000 && time.Now().Hour() < 18 {
		log.Printf("第【%s】期：投注金额不足2万，进行不足至2.1万 ********** \n", nextIssue)

		if _, err := bet28(cache, nextIssue, surplus, SN28, spaces, float64(21000-total)); err != nil {
			return err
		}
	}

	return nil
}

package xmd

import (
	"log"
	"strconv"
)

func analysis(cache *Cache) error {
	if err := cache.Sync(2000); err != nil {
		return err
	}

	// 当前账户余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	target := []int{6, 7, 9, 12, 14, 16, 17, 18, 20}
	exists := make(map[int]struct{}, 0)
	for _, result := range target {
		exists[result] = struct{}{}
	}
	nextIssue := strconv.Itoa(cache.issue + 1)

	if _, ok := exists[cache.result]; !ok {
		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，不在投注范围内 ...\n", cache.issue, cache.result, surplus)
		return nil
	}

	betGold := int(float64(cache.user.gold) * float64(stds[cache.result]) / 1000)
	if err := hPostBet(nextIssue, betGold, cache.result, cache.user); err != nil {
		return err
	}
	log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】，余额【%d】 ...\n", nextIssue, cache.result, 1000.0/float64(stds[cache.result]), betGold, surplus-betGold)

	return nil
}

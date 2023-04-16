package xmd

import (
	"log"
	"strconv"
)

func analysis(cache *Cache) error {
	if err := cache.Sync(2000); err != nil {
		return err
	}

	nextIssue := strconv.Itoa(cache.issue + 1)

	// 当前账户余额
	surplus, err := hGetGold(cache.user)
	if err != nil {
		return err
	}

	//for i := len(cache.histories) - 1; i >= len(cache.histories)-8; i-- {
	//	result := cache.histories[i].result
	//	if result <= 5 || result >= 22 {
	//		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，暂时停止投注 ...\n", cache.issue, cache.result, surplus)
	//		return nil
	//	}
	//}

	size := len(cache.histories)
	r1 := cache.histories[size-1].result
	r2 := cache.histories[size-2].result

	if r1 < 10 || r1 > 17 {
		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，不符合投注条件A ...\n", cache.issue, cache.result, surplus)
		return nil
	}

	if r2 >= 10 && r2 <= 17 {
		log.Printf("第【%d】期：开奖结果【%d】，余额【%d】，不符合投注条件B ...\n", cache.issue, cache.result, surplus)
		return nil
	}

	var total int
	for _, result := range []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21} {
		betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
			return err
		}

		log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)
		total = total + betGold
	}

	surplus = surplus - total
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, total, surplus)

	return nil
}

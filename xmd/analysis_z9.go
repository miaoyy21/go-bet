package xmd

//import (
//	"log"
//	"strconv"
//)
//
//var lastSurplus int
//
//func analysis(cache *Cache) error {
//	if err := cache.Sync(200); err != nil {
//		return err
//	}
//
//	nextIssue := strconv.Itoa(cache.issue + 1)
//
//	// 当前账户余额
//	surplus, err := hGetGold(cache.user)
//	if err != nil {
//		return err
//	}
//
//	std, total := 1000, 0
//	if lastSurplus != 0 && surplus > lastSurplus {
//		std = cache.user.gold
//	}
//	lastSurplus = surplus
//
//	for _, result := range SN28 {
//		betGold := int(float64(std) * float64(stds[result]) / 1000)
//		if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
//			return err
//		}
//
//		total = total + betGold
//	}
//	log.Printf("第【%s】期：投注金额【%d】，余额【%d】 >>>>>>>>>> \n", nextIssue, total, surplus-total)
//
//	return nil
//}

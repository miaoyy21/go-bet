package xmd

//
//import (
//	"log"
//	"strconv"
//)
//
//func analysis(cache *Cache) error {
//	if err := cache.Sync(200); err != nil {
//		return err
//	}
//
//	nextIssue := strconv.Itoa(cache.issue + 1)
//
//	r1 := cache.histories[len(cache.histories)-1].result
//	r2 := cache.histories[len(cache.histories)-2].result
//
//	if r1 >= 10 && r1 <= 17 && (r2 <= 9 || r2 >= 18) {
//		for _, result := range []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21} {
//
//			betGold := int(float64(cache.user.gold) * float64(stds[result]) / 1000)
//			if err := hPostBet(nextIssue, betGold, result, cache.user); err != nil {
//				return err
//			}
//			log.Printf("第【%s】期：竞猜数字【👍 %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", nextIssue, result, 1000.0/float64(stds[result]), betGold)
//		}
//	}
//
//	return nil
//}

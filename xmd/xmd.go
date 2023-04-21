package xmd

import (
	"log"
	"time"
)

var SN10 = []int{0, 1, 2, 3, 4, 23, 24, 25, 26, 27}
var SN12 = []int{0, 1, 2, 3, 4, 5, 22, 23, 24, 25, 26, 27}
var SN28 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Run(cache *Cache) {
	log.Printf("当前设置投注基数为 %d ...\n", cache.user.gold)
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}
	calc()

	sec := 75.0
	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("%.2f秒后[%s]，将运行小鸡竞猜游戏 ...", sec-dua.Seconds(), time.Now().Add(time.Second*time.Duration(sec-dua.Seconds())).Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * time.Duration(sec-dua.Seconds()))
	if err := analysis(cache); err != nil {
		log.Println(err.Error())
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Println("游戏小鸡竞猜已启动 ...")
	for {
		select {
		case <-ticker.C:
			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func bet28(cache *Cache, issue string, surplus int, ns []int) (map[int]struct{}, error) {
	var total, coverage int

	bets := make(map[int]struct{})
	for _, result := range ns {
		betGold := int(rate * float64(cache.user.gold) * float64(stds[result]) / 1000)
		if err := hPostBet(issue, betGold, result, cache.user); err != nil {
			return nil, err
		}
		log.Printf("第【%s】期：竞猜数字【❤️ %02d】，标准赔率【%-7.2f】，投注金额【% 5d】\n", issue, result, 1000.0/float64(stds[result]), betGold)

		bets[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", issue, total, surplus-total, float64(coverage)/10)

	return bets, nil
}

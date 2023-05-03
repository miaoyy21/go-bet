package xmd

import (
	"log"
	"math"
	"time"
)

var SN10 = []int{0, 1, 2, 3, 4, 23, 24, 25, 26, 27}
var SN14 = []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
var SN28 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Run(cache *Cache) {
	log.Printf("当前设置活动状态（%t） ...\n", cache.isExtra)
	log.Printf("当前设置投注基数为 %d ...\n", cache.user.gold)
	log.Printf("当前设置当不存在超过实际赔率%.2f%%的数字时，仅进行全部投注 ...\n", cache.wx*100-100)
	log.Printf("当前设置当返奖率不超过%.2f%%时，仅进行全部投注 ...\n", cache.rx*100)
	log.Printf("当前设置本时内中奖比率变化量达到设定标准时，基础投注变为%.2f...\n", cache.ex)
	log.Printf("当前设置连续投注失败3次时，本期投注倍率变为%.2f倍...\n", cache.dx)
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}
	calc()

	sec := 53.5
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
			// 配置文件是否变化
			if ok, err := cache.Reload(); err != nil {
				log.Println(err.Error())
			} else {
				if ok {
					log.Println("配置文件变化，重新加载配置文件完成 ...")
				}
			}

			// 基础投注是否变化
			if ok, err := cache.Update(); err != nil {
				log.Println(err.Error())
			} else {
				if ok {
					log.Printf("由于本时内中奖比率变化量达到设定标准，基础投注变为【%d】 ... \n", cache.user.gold)
				}
			}

			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func bet28(cache *Cache, issue string, surplus int, ns []int, spaces map[int]int, rts map[int]float64, std float64) (map[int]struct{}, error) {
	var total, coverage int
	if std <= 0 {
		return nil, nil
	}

	bets := make(map[int]struct{})
	for _, result := range ns {
		r0 := 1000.0 / float64(stds[result])
		r1 := rts[result]

		betGold := int(math.Ceil(std * float64(stds[result]) / 1000))
		if err := hPostBet(issue, betGold, result, cache.user); err != nil {
			return nil, err
		}
		log.Printf("第【%s】期：竞猜数字【❤️ %02d】，标准赔率【%-7.2f】，实际赔率【%-7.2f】，赔率系数【%-4.2f】，间隔次数【%-4d】，投注金额【% 5d】\n", issue, result, r0, r1, r1/r0, spaces[result], betGold)

		bets[result] = struct{}{}
		total = total + betGold
		coverage = coverage + stds[result]
	}
	log.Printf("第【%s】期：投注金额【%d】，余额【%d】，覆盖率【%.2f%%】 >>>>>>>>>> \n", issue, total, surplus-total, float64(coverage)/10)

	return bets, nil
}

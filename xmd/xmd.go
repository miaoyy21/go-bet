package xmd

import (
	"log"
	"math"
	"time"
)

var SN28 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Run(cache *Cache) {
	log.Printf("当前设置仅投注实际赔率超过标准倍率%.2f%%的数字 ...\n", cache.wx*100-100)
	log.Printf("当前设置仅当返奖率超过%.2f%%时，才进行投注 ...\n", cache.exp*100)
	log.Printf("当前设置投注基数为 %d ...\n", cache.user.gold)
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}
	calc()

	// 加载前一次保存的全局变量
	if err := tempLoad(); err != nil {
		log.Println(err.Error())
	}

	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("%.2f秒后[%s]，将运行小鸡竞猜游戏 ...", cache.secs-dua.Seconds(), time.Now().Add(time.Second*time.Duration(cache.secs-dua.Seconds())).Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * time.Duration(cache.secs-dua.Seconds()))

	if err := analysis(cache); err != nil {
		log.Println(err.Error())
	}

	// 保存全局变量
	if err := tempSave(); err != nil {
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

			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}

			// 保存全局变量
			if err := tempSave(); err != nil {
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

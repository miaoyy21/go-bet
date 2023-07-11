package xmd

import (
	"log"
	"math/rand"
	"time"
)

var SN28 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Run(cache *Cache) {
	rand.Seed(time.Now().Unix())
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}

	log.Printf("当前是否启用设定投注模式【%t】 ... \n", cache.user.isBetMode)
	calc()

	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("%.2f秒后[%s]，将运行小鸡竞猜游戏 ...", cache.secs-dua.Seconds(), time.Now().Add(time.Second*time.Duration(cache.secs-dua.Seconds())).Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * time.Duration(cache.secs-dua.Seconds()))

	go func() {
		if err := cache.Sync(200); err != nil {
			log.Println(err.Error())
		}

		if _, err := cache.Reload(); err != nil {
			log.Println(err.Error())
		}

		//if isStop(cache) {
		//	return
		//}

		if err := analysis(cache); err != nil {
			log.Println(err.Error())
		}
	}()

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

			// 查询开奖历史
			if err := cache.Sync(200); err != nil {
				log.Println(err.Error())
			}

			//if isStop(cache) {
			//	continue
			//}

			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

//func isStop(cache *Cache) bool {
//	hm := time.Now().Format("15:04")
//
//	// 工作时间段
//	if (hm >= "04:00" && hm <= "05:00") || (hm >= "08:30" && hm <= "11:30") || (hm >= "14:30" && hm <= "17:00") {
//		if stops > 1 {
//			stops--
//
//			if rand.Float32() <= 0.50 {
//				latest = make(map[int]struct{})
//				log.Printf("😤😤😤 第【%d】期：上一期开奖结果【%d】，暂停投注，剩余期数【%d】 >>>>>>>>>> \n", cache.issue+1, cache.result, stops)
//				return true
//			}
//
//			return false
//		}
//
//		if len(latest) < 1 {
//			return false
//		}
//
//		if _, ok := latest[cache.result]; ok {
//			fails = 0
//			return false
//		}
//
//		fails++
//		if fails < 3 {
//			return false
//		}
//
//		fails, stops = 0, 10
//		latest = make(map[int]struct{})
//		log.Printf("😤😤😤 第【%d】期：上一期开奖结果【%d】，暂停投注，剩余期数【%d】 >>>>>>>>>> \n", cache.issue+1, cache.result, stops)
//		return true
//	}
//
//	fails, stops = 0, 0
//	if len(latest) < 1 {
//		return false
//	}
//
//	if _, ok := latest[cache.result]; !ok && rand.Float32() <= 0.50 {
//		latest = make(map[int]struct{})
//		log.Printf("😤😤😤 第【%d】期：上一期开奖结果【%d】，暂停投注，剩余期数【%d】 >>>>>>>>>> \n", cache.issue+1, cache.result, stops)
//		return true
//	}
//
//	return false
//}

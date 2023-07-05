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
		if _, err := cache.Reload(); err != nil {
			log.Println(err.Error())
		}

		if isStop(cache) {
			return
		}

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

			if isStop(cache) {
				continue
			}

			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func isStop(cache *Cache) bool {
	if len(latest) <= 0 {
		return false
	}

	if _, ok := latest[cache.result]; ok {
		return false
	}

	if rand.Float32() <= 0.50 {
		return false
	}

	latest = make(map[int]struct{})
	log.Printf("😤😤😤 第【%d】期：上一期开奖结果【%d】，由于投注失利，随机选择不进行投注 >>>>>>>>>> \n", cache.issue+1, cache.result)
	return true
}

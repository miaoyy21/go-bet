package xmd

import (
	"log"
	"time"
)

var SN28 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Run(cache *Cache) {
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}
	calc()

	dua := time.Now().Sub(time.Now().Truncate(time.Minute))
	log.Printf("%.2f秒后[%s]，将运行小鸡竞猜游戏 ...", cache.secs-dua.Seconds(), time.Now().Add(time.Second*time.Duration(cache.secs-dua.Seconds())).Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * time.Duration(cache.secs-dua.Seconds()))

	go func() {
		if _, err := cache.Reload(); err != nil {
			log.Println(err.Error())
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

			h := time.Now().Hour()
			if h == 3 || (h >= 9 && h <= 10) || (h >= 15 && h <= 16) {
				log.Printf("%s 属于投注暂停时间 ******** \n", time.Now().Format("2006-01-02 15:04:05"))
				continue
			}

			if err := analysis(cache); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

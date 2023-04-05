package xmd

import (
	"log"
	"time"
)

func Run(cache *Cache) {
	calc()
	if cache.user.isDebug {
		log.Println("当前设置为调试模式，不发送投注请求 ...")
	}

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

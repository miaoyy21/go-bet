package xmd

import (
	"fmt"
	"math"
	"time"
)

func (o *Cache) Update() (bool, error) {
	if time.Now().Minute() < 59 {
		return false, nil
	}

	var rate float64

	// 统计本时内中奖比率
	dt := time.Now().Format("2006-01-02 15")
	query := fmt.Sprintf("SELECT CONVERT(SUM(win_gold)/AVG(user_gold),DECIMAL(13,2)) AS rate FROM logs WHERE time LIKE '%s%%'", dt)
	if err := o.db.QueryRow(query).Scan(&rate); err != nil {
		return false, err
	}

	// 如果中奖比率在3倍以内，不改变基础投注
	if math.Abs(rate) < 3.0 {
		return false, nil
	}

	// 中奖比率超过3倍
	delta := 1.35
	if rate > 0 {
		o.user.gold = int(math.Ceil(float64(o.user.gold) * delta))
	} else {
		o.user.gold = int(math.Ceil(float64(o.user.gold) / delta))
	}

	return true, nil
}

package xmd

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
	"time"
)

func TestCache_Sync2(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/bet?charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true")
	if err != nil {
		log.Fatalf("sql.Open() fail : %s \n", err.Error())
	}

	var rows int
	if err := db.QueryRow("SELECT COUNT(1) AS XN FROM histories").Scan(&rows); err != nil {
		log.Fatalf("sql.QueryRow() fail : %s \n", err.Error())
	}

	log.Printf("样本容量： %d \n", rows)

	size := 1440
	target := make([]int, 0)

	log.Printf("样本分析周期：%d \n", size)
	for no := 0; no <= 27; no++ {
		rates := make([]float64, 0, int(math.Ceil(float64(rows)/float64(size))))
		for i := 0; i < int(math.Ceil(float64(rows)/float64(size))); i++ {
			var rate float64
			if err := db.QueryRow("CALL PROC(?,?,?)", no, size, i*size).Scan(&rate); err != nil {
				if err == sql.ErrNoRows {
					rates = append(rates, rate)
					continue
				}

				log.Fatalf("sql.Exec() fail : %s \n", err.Error())
			}

			rates = append(rates, rate)
		}

		var rate float64
		if err := db.QueryRow("CALL PROC(?,?,?)", no, size/2, rows-size/2).Scan(&rate); err != nil {
			if err != sql.ErrNoRows {
				log.Fatalf("sql.Exec() fail : %s \n", err.Error())
			}
		}

		rate1 := rates[len(rates)-1]
		rate2 := rates[len(rates)-2]

		if rate1 > 1 && rate2 > 1 && rate > 1 {
			target = append(target, no)
			log.Printf("[ %02d √ P1 ] %s [%.4f]\n", no, fmtF64(rates), rate)
		} else {
			if rate1 < 1 && rate < 1 {
				log.Printf("[ %02d x N1 ] %s [%.4f]\n", no, fmtF64(rates), rate)
			} else if rate1 > 1 && rate < 1 {
				log.Printf("[ %02d x N2 ] %s [%.4f]\n", no, fmtF64(rates), rate)
			} else if rate1 < 1 && rate > 1 {
				target = append(target, no)
				log.Printf("[ %02d √ P2 ] %s [%.4f]\n", no, fmtF64(rates), rate)
			} else if (rate1+rate2)/2 < rate {
				target = append(target, no)
				log.Printf("[ %02d √ P3 ] %s [%.4f]\n", no, fmtF64(rates), rate)
			} else {
				log.Printf("[ %02d _ __ ] %s [%.4f]\n", no, fmtF64(rates), rate)
			}
		}
	}
	log.Printf("所选取的结果：%s \n", fmtInt(target))

	log.Println(time.Now().Format("2006-01-02 15:04:05.999"))
}

func fmtF64(fs []float64) string {
	ss := make([]string, 0, len(fs))
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf("%.4f", f))
	}

	return strings.Join(ss, "  ")
}

func fmtInt(fs []int) string {
	ss := make([]string, 0, len(fs))
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf("%d", f))
	}

	return strings.Join(ss, ",")
}

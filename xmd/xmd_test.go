package xmd

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
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

	log.Printf("Rows Total Count is %d \n", rows)

	size := 1440
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

		log.Printf("[%02d] %s [%.4f]\n", no, fmtF64(rates), rate)
	}
}

func fmtF64(fs []float64) string {
	ss := make([]string, 0, len(fs))
	for _, f := range fs {
		ss = append(ss, fmt.Sprintf("%.4f", f))
	}

	return strings.Join(ss, "  ")
}

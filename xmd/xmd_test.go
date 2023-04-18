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

	for no := 0; no <= 27; no++ {
		rates := make([]string, 0, int(math.Ceil(float64(rows)/1440)))
		for i := 0; i < int(math.Ceil(float64(rows)/1440)); i++ {
			var rate float64
			if err := db.QueryRow("CALL PROC(?,?,?)", no, 1440, i*1440).Scan(&rate); err != nil {
				if err == sql.ErrNoRows {
					rates = append(rates, fmt.Sprintf("%.4f", rate))
					continue
				}

				log.Fatalf("sql.Exec() fail : %s \n", err.Error())
			}

			rates = append(rates, fmt.Sprintf("%.4f", rate))
		}

		var rate float64
		if err := db.QueryRow("CALL PROC(?,?,?)", no, 720, rows-720).Scan(&rate); err != nil {
			if err == sql.ErrNoRows {
				rates = append(rates, fmt.Sprintf("[%.4f]", rate))
			} else {
				log.Fatalf("sql.Exec() fail : %s \n", err.Error())
			}
		} else {
			rates = append(rates, fmt.Sprintf("[%.4f]", rate))
		}

		log.Printf("[%02d] %s \n", no, strings.Join(rates, "  "))
	}

}

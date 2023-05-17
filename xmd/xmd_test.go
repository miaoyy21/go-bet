package xmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"testing"
)

func TestCache_Sync2(t *testing.T) {

	dbHost, dbPass := "www.nsmei.com", "envl11a3ANsxZNCFQA0GO9HLObf8gi5"
	userId := "31591499"
	dt := "2023-05-16"

	// MySQL
	ds := fmt.Sprintf("root:%s@tcp(%s:3306)/bet?charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true", dbPass, dbHost)
	db, err := sql.Open("mysql", ds)
	if err != nil {
		t.Fatalf("sql.Open() fail :: %s \n", err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() fail :: %s \n", err.Error())
	}

	query := fmt.Sprintf("SELECT issue FROM logs_%s WHERE time > ? ORDER BY issue ASC", userId)
	res, err := sqlQuery(db, query, dt)
	if err != nil {
		t.Fatalf("sqlQuery() fail :: %s \n", err.Error())
	}

	for i, r := range res {
		ris := r["issue"]

		var wnRate int
		var wRate sql.NullFloat64
		wQuery := fmt.Sprintf("SELECT COUNT(1) AS nn,IFNULL(CONVERT(SUM(win_gold * IFNULL(rz,1.0))/AVG(user_gold * IFNULL(rz,1.0)),DECIMAL(13,2)),1.0) AS rate FROM logs_%s WHERE issue >= ? - 10 and issue < ?", userId)
		if err := db.QueryRow(wQuery, ris, ris).Scan(&wnRate, &wRate); err != nil {
			t.Fatalf("db.QueryRow(%d) 1 fail :: %s \n", ris, err.Error())
		}

		var fnRate int
		var fRate sql.NullFloat64
		fQuery := fmt.Sprintf("SELECT COUNT(1) AS nn,SUM(ABS(win_gold * rz)) AS rate FROM logs_%s WHERE issue >= ? - 10 and issue < ? and win_gold < 0", userId)
		if err := db.QueryRow(fQuery, ris, ris).Scan(&fnRate, &fRate); err != nil {
			t.Fatalf("db.QueryRow(%d) 2 fail :: %s \n", ris, err.Error())
		}

		rate := 1.0
		if i > 10 && wnRate+fnRate >= 8 {
			if !fRate.Valid {
				if wRate.Valid {
					rate = wRate.Float64
				}
			} else {
				if wRate.Valid {
					rate = wRate.Float64 / fRate.Float64
				}
			}
		}

		rate = math.Trunc(rate*100) / 100
		if rate < 0.1 {
			rate = 0.1
		} else if rate > 10 {
			rate = 10
		}

		uQuery := fmt.Sprintf("UPDATE logs_%s SET rz = ? WHERE issue = ?", userId)
		if _, err := db.Exec(uQuery, rate, ris); err != nil {
			t.Fatalf("db.Exec() fail :: %s \n", err.Error())
		}

		t.Logf("Update Issue of [%.2f%%] %d %.2f Successful ! \n", float64((i+1)*100)/float64(len(res)), ris, rate)
	}
}

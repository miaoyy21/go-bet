package xmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"testing"
)

func TestRun(t *testing.T) {
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

		var sn int
		var sr float64
		qs := fmt.Sprintf("SELECT COUNT(1) AS nn,IFNULL(CONVERT(SUM(win_gold * IFNULL(rz,1.0))/AVG(user_gold * IFNULL(rz,1.0))/AVG(IFNULL(rz,1.0)),DECIMAL(13,2)),1.0) AS rate FROM logs_%s WHERE issue >= ? - 20 and issue < ?", userId)
		if err := db.QueryRow(qs, ris, ris).Scan(&sn, &sr); err != nil {
			t.Fatalf("db.QueryRow(%d) fail :: %s \n", ris, err.Error())
		}

		rate := 1.0
		if sn >= 16 {
			rate = math.Trunc(math.Pow(1.25, sr)*100) / 100
			if rate < 0.25 {
				rate = 0.25
			} else if rate > 2 {
				rate = 2
			}
		}

		uQuery := fmt.Sprintf("UPDATE logs_%s SET rz = ? WHERE issue = ?", userId)
		if _, err := db.Exec(uQuery, rate, ris); err != nil {
			t.Fatalf("db.Exec() fail :: %s \n", err.Error())
		}

		t.Logf("Update Issue of %d/%d [%.2f%%] %d %.2f Successful ! \n", i+1, len(res), float64((i+1)*100)/float64(len(res)), ris, rate)
	}
}

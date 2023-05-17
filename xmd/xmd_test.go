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
	userId := "31783106"
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

		var sn int
		var sr float64
		qs := fmt.Sprintf("SELECT COUNT(1) AS nn,IFNULL(CONVERT(SUM(win_gold * IFNULL(rz,1.0))/AVG(user_gold * IFNULL(rz,1.0)),DECIMAL(13,2)),1.0) AS rate FROM logs_%s WHERE issue >= ? - 10 and issue < ?", userId)
		if err := db.QueryRow(qs, ris, ris).Scan(&sn, &sr); err != nil {
			t.Fatalf("db.QueryRow(%d) fail :: %s \n", ris, err.Error())
		}

		rate := 1.0
		if sn >= 8 {
			rate = math.Pow(1.25, sr)
		}

		uQuery := fmt.Sprintf("UPDATE logs_%s SET rz = ? WHERE issue = ?", userId)
		if _, err := db.Exec(uQuery, rate, ris); err != nil {
			t.Fatalf("db.Exec() fail :: %s \n", err.Error())
		}

		t.Logf("Update Issue of [%.2f%%] %d %.2f Successful ! \n", float64((i+1)*100)/float64(len(res)), ris, rate)
	}
}

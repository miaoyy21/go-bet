package xmd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/bet?charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true")
	if err != nil {
		log.Fatalf("sql.Open() fail : %s \n", err.Error())
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Getwd() fail : %s \n", err.Error())
	}

	dir0, _ := filepath.Split(dir)
	cache, err := NewCache(dir0)
	if err != nil {
		log.Fatalf("NewCache() fail : %s \n", err.Error())
	}

	histories, err := hGetHistories(1440*7, cache.user)
	if err != nil {
		log.Fatalf("hGetHistories() fail : %s \n", err.Error())
	}

	if _, err := db.Exec("TRUNCATE TABLE histories"); err != nil {
		log.Fatalf("hGetHistories() fail : %s \n", err.Error())
	}

	query := "INSERT INTO histories(issue,result) VALUES (?,?)"
	for i := 0; i <= len(histories)-1; i++ {
		if _, err := db.Exec(query, histories[i].Issue, histories[i].Result); err != nil {
			log.Printf("db.Exec(%s,%s) warning : %s \n", histories[i].Issue, histories[i].Result, err.Error())
			continue
		}

		log.Printf("Insert %q Successful ! \n", histories[i].Issue)
	}
}

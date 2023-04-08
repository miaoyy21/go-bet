package xmd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestCache_Sync(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/bet?charset=utf8mb4&collation=utf8mb4_general_ci&loc=Local&parseTime=true")
	if err != nil {
		log.Fatalf("sql.Open() fail : %s \n", err.Error())
	}

	unix := "1680814387"
	keyCode := "5e6e725a5abd96b9ec0b3bc929a4b4fd"
	deviceId := "0E6EE3CC-8184-4CD7-B163-50AE8AD4516F"
	userId := "31591499"
	token := "cbj7s576p3se6c87194kwqo1c1w2cq87sau8lc2s"

	userBase := NewUserBase(true, 0, unix, keyCode, deviceId, userId, token)
	histories, err := hGetHistories(200, userBase)
	if err != nil {
		log.Fatalf("hGetHistories() fail : %s \n", err.Error())
	}

	query := "INSERT INTO histories(issue,result) VALUES (?,?)"
	for i := 0; i <= len(histories)-1; i++ {
		if _, err := db.Exec(query, histories[i].Issue, histories[i].Result); err != nil {
			log.Printf("db.Exec(%s,%s) warning : %s \n", histories[i].Issue, histories[i].Result, err.Error())
			continue
		}

		log.Printf("Insert Issue %q Successful ! \n", histories[i].Issue)
	}
}

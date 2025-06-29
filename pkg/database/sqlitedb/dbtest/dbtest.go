package dbtest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/initial"
)

func WriteDBTest1() {
	db := sqldb.GetDB()
	var keyInt int = 1
	var value string = "db test1"
	for {
		key := strconv.Itoa(keyInt)
		// tx, err := db.Begin()
		// if err != nil {
		// 	fmt.Printf("begin. Exec error=%s", err)
		// 	continue
		// }

		insertClause := "insert into k_v_table(key,value) values($1,$2)"
		_, err := db.Exec(insertClause, key, value)
		if err != nil {
			fmt.Printf("dbtest write test1 failed: %s", err.Error())
		}
		//tx.Commit()
		keyInt++
	}
}

func WriteDBTest2() {
	db := sqldb.GetDB()
	var keyInt int = -1
	var value string = "db test2"
	for {
		key := strconv.Itoa(keyInt)

		insertClause := "insert into k_v_table(key,value) values($1,$2)"
		_, err := db.Exec(insertClause, key, value)
		if err != nil {
			fmt.Printf("dbtest write test2 failed: %s", err.Error())
		}
		keyInt--
	}
}

func WriteDBTest3() {
	db := sqldb.GetDB()
	var keyInt int = 1
	var value string = "db test3"
	for {
		key := strconv.Itoa(keyInt)
		insertClause := "update k_v_table set value =$1 where key=$2"
		_, err := db.Exec(insertClause, value, key)
		if err != nil {
			fmt.Printf("dbtest write test3 failed: %s", err.Error())
		}
		keyInt++
		time.Sleep(30 * time.Millisecond)
	}
}

func ReadDBTest1() {
	db := sqldb.GetDB()
	var offset = 0
	for {
		var key string
		var value string
		selectClause := fmt.Sprintf("select * from k_v_table limit 20 offset %d", offset)
		rows, _ := db.Query(selectClause)

		for rows.Next() {
			err := rows.Scan(&key, &value)
			if err != nil {
				fmt.Printf("Read test1 failed: %s\n", err.Error())
			} else {
				fmt.Printf("key %s, value %s\n", key, value)
			}
		}
		rows.Close()
		offset += 20
		time.Sleep(1 * time.Second)
	}
}

func Init() error {
	go WriteDBTest1()
	go WriteDBTest2()
	go WriteDBTest3()
	go ReadDBTest1()
	return nil
}

func init() {
	initial.Register("dbtest", initial.PhaseDefault, Init)
}

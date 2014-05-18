package sup

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import (
	"fmt"
	"log"
	"time"
)

var db *sql.DB

//
// STATUS
//

func SaveStatus(status *Status) error {
    createTimeStr := status.CreateTime.Format(time.RFC3339Nano)
	_, err := db.Exec("insert into status"+
		" (ip, user, create_time, tags, description)"+
		" values (?, ?, ?, ?, ?)",
        status.IP, status.User, createTimeStr, status.Tags, status.Description)
    return err
}

func GetStatuses(user string) ([]Status, error) {
    res, err := db.Query("select ip,user,create_time,tags,description"+
        " from status where user=?", user)
    ret := make([]Status, 0)
    if err != nil {
        return ret, err
    }
    for res.Next() {
        var s Status
        var createTimeStr string
        rowErr := res.Scan(
            &s.IP,
            &s.User,
            &createTimeStr,
            &s.Tags,
            &s.Description)
        if rowErr == nil {
            s.CreateTime, rowErr = time.Parse(time.RFC3339, createTimeStr)
        }
        if rowErr != nil {
            err = rowErr
        } else {
            ret = append(ret, s)
        }
    }
    return ret, err
}


// 
// INIT
//

func init() {
	conf := GetConfig()
	mysqlHost := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8",
		conf.DbUser,
		conf.DbPassword,
		conf.DbServer,
		conf.DbCatalog)

	// connect to the database, ping periodically to maintain the connection
	log.Printf("Connecting to %s\n", mysqlHost)
	var err error
	db, err = sql.Open("mysql", mysqlHost)
	if err != nil {
		panic(err)
	}
	go ping()

	// migrate the database
	migrateDb()
}

func ping() {
	ticker := time.Tick(time.Minute)
	for {
		<-ticker
		err := db.Ping()
		if err != nil {
			log.Printf("DB not ok: %v\n", err)
		} else {
			log.Printf("DB ok\n")
		}
	}
}

package dbm

// WIP
// Migrate mariadb to sqlite

//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/mattn/go-sqlite3"
//	"log"
//	"os"
//)
//
//func Create() error {
//	if _, err := os.Create("huetify.db"); err != nil {
//		return err
//	}
//
//	db, err := sql.Open("sqlite3", "huetify.db")
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	_, err = db.Exec("CREATE TABLE `customers` (`till_id` INTEGER PRIMARY KEY AUTOINCREMENT, `client_id` VARCHAR(64) NULL, `first_name` VARCHAR(255) NOT NULL, `last_name` VARCHAR(255) NOT NULL, `guid` VARCHAR(255) NULL, `dob` DATETIME NULL, `type` VARCHAR(1))")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	db.Close()
//	return nil
//}
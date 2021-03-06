package database

import (
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func ConnectGORM(user_id string, password string) *gorm.DB {

	DBMS := "mysql"
	USER := user_id
	PASS := password
	PROTOCOL := "tcp(0.0.0.0:3306)"
	DBNAME := "gateway"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	db.SingularTable(true)
	return db
}

func ConnectAdminGORM(user_id string, password string) *gorm.DB {

	DBMS := "mysql"
	USER := user_id
	PASS := password
	PROTOCOL := "tcp(0.0.0.0:3306)"
	DBNAME := "mysql"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		db.Close()
		panic(err.Error())
	}

	db.SingularTable(true)
	return db
}

func CheckJwt(tokenState *jwt.Token) (string, string) {
	claims := tokenState.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), claims["password"].(string)
}

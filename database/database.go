package database

import (
	"fmt"

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

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	db.SingularTable(true)
	return db
}

func CheckJwt(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), claims["password"].(string), nil
	} else {
		fmt.Println(err)
		return "", "", err
	}
}

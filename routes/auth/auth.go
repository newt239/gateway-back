package authRoute

import (
	"net/http"
	"os"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		loginPostParam := authPostParam{}
		if err := c.Bind(&loginPostParam); err != nil {
			return err
		}
		user_id, password := loginPostParam.UserId, loginPostParam.Password

		db := database.ConnectGORM(user_id, password)
		db.Close()

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user_id
		claims["password"] = password
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		godotenv.Load(".env")
		signature := os.Getenv("SIGNATURE")

		tokenString, err := token.SignedString([]byte(signature))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": tokenString,
		})
	}
}

func Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)
		var result user
		db.Where("user_id = ?", user_id).First(&user{}).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id":      result.UserId,
			"display_name": result.DisplayName,
			"user_type":    result.UserType,
			"available":    result.Available,
		})
	}
}

type authPostParam struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type user struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	UserType    string `json:"user_type"`
	Available   int    `json:"available"`
	CreatedBy   string `json:"created_by"`
}

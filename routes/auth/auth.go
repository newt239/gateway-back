package authRoute

import (
	"net/http"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := c.FormValue("user_id"), c.FormValue("password")
		db := database.ConnectGORM(user_id, password)

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user_id
		claims["password"] = password
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		db.Close()
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

type user struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	UserType    string `json:"user_type"`
	Role        string `json:"role"`
	Available   int    `json:"available"`
	Note        string `json:"note"`
	CreatedBy   string `json:"created_by"`
}

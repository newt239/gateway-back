package adminRoute

import (
	"fmt"
	"net/http"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		type createUserPostParam struct {
			UserId      string `json:"user_id"`
			Password    string `json:"password"`
			DisplayName string `json:"display_name"`
			UserType    string `json:"user_type"`
		}
		newUserData := createUserPostParam{}
		if err := c.Bind(&newUserData); err != nil {
			return err
		}
		db.Exec(fmt.Sprintf(`CREATE USER '%s'@'localhost' identified by '%s';`, newUserData.UserId, newUserData.Password))
		db.Exec("FLUSH PRIVILEGES;")

		var sql string
		switch newUserData.UserType {
		case "moderator":
			sql += fmt.Sprintf("GRANT ALL ON *.* TO '%s'@'localhost' WITH GRANT OPTION; ", newUserData.UserId)
		case "executive":
			sql += fmt.Sprintf("GRANT INSERT, UPDATE, SELECT ON gateway.* TO '%s'@'localhost'; ", newUserData.UserId)
		case "exhibit":
			sql += fmt.Sprintf("GRANT INSERT, UPDATE, SELECT ON gateway.* TO '%s'@'localhost'; ", newUserData.UserId)
		case "analysis":
			sql += fmt.Sprintf("GRANT SELECT ON gateway.* TO '%s'@'localhost'; ", newUserData.UserId)
		}
		db.Exec(sql)
		db.Exec(fmt.Sprintf("INSERT INTO gateway.user (user_id, display_name, user_type, created_by, available) VALUES ('%s', '%s', '%s', '%s', 1);", newUserData.UserId, newUserData.DisplayName, newUserData.UserType, user_id))
		db.Close()
		return c.NoContent(http.StatusOK) // status code 200で何も返さない
	}
}

func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		db.Exec(fmt.Sprintf("DROP USER '%s'@'localhost';", c.Param("user_id")))
		db.Exec(fmt.Sprintf("DELETE FROM gateway.user WHERE user_id='%s' AND created_by='%s';", c.Param("user_id"), user_id))
		return c.NoContent(http.StatusOK) // status code 200で何も返さない
	}
}

package adminRoute

import (
	"fmt"
	"net/http"
	"time"

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

func CreatedByMeUserList() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		var result []user
		db.Where("created_by = ?", user_id).Find(&user{}).Scan(&result)

		fmt.Println(result)
		return c.JSON(http.StatusOK, result)
	}
}

func CreateExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		type createExhibitPostParam struct {
			ExhibitId   string `json:"exhibit_id"`
			ExhibitName string `json:"exhibit_name"`
			RoomName    string `json:"room_name"`
			ExhibitType string `json:"exhibit_type"`
			Capacity    int    `json:"capacity"`
		}
		newExhibitData := createExhibitPostParam{}
		if err := c.Bind(&newExhibitData); err != nil {
			return err
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		exhibitEx := exhibit{
			ExhibitId:   newExhibitData.ExhibitId,
			ExhibitName: newExhibitData.ExhibitName,
			RoomName:    newExhibitData.RoomName,
			ExhibitType: newExhibitData.ExhibitType,
			Capacity:    newExhibitData.Capacity,
			LastUpdate:  now,
		}

		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)
		db.Table("exhibit").Omit("position", "status", "note").Create(&exhibitEx)
		return c.NoContent(http.StatusOK)
	}
}

func DeleteExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		db.Exec(fmt.Sprintf("DELETE FROM gateway.exhibit WHERE exhibit_id='%s';", c.Param("exhibit_id")))
		return c.NoContent(http.StatusOK) // status code 200で何も返さない
	}
}

type user struct {
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	UserType    string `json:"user_type"`
}

type exhibit struct {
	ExhibitId   string    `json:"exhibit_id"`
	ExhibitName string    `json:"exhibit_name"`
	RoomName    string    `json:"position_name"`
	ExhibitType string    `json:"exhibit_type"`
	Status      int       `json:"status"`
	Capacity    int       `json:"capacity"`
	Position    string    `json:"position"`
	LastUpdate  time.Time `json:"last_update"`
	Note        string    `json:"note"`
}

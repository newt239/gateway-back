package activityRoute

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Enter() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		enterPostParam := activityPostParam{}
		if err := c.Bind(&enterPostParam); err != nil {
			return err
		}

		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		session_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
		sessionEx := session{
			SessionId:      session_id,
			ExhibitId:      enterPostParam.ExhibitId,
			GuestId:        enterPostParam.GuestId,
			EnterAt:        now,
			EnterOperation: user_id,
			Available:      1,
		}

		db := database.ConnectGORM(user_id, password)
		db.Table("session").Omit("exit_at", "exit_operation", "note").Create(&sessionEx)
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"session_id": session_id,
		})
	}
}

func Exit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		exitPostParam := activityPostParam{}
		if err := c.Bind(&exitPostParam); err != nil {
			return err
		}

		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		sessionEx := session{
			ExitAt:        now,
			ExitOperation: user_id,
		}
		var result session

		db := database.ConnectGORM(user_id, password)
		db.Table("session").Where("guest_id = ?", exitPostParam.GuestId).Where("exhibit_id = ?", exitPostParam.ExhibitId).Where("exit_at is ?", gorm.Expr("NULL")).Updates(&sessionEx).Scan(&result)
		db.Close()
		return c.NoContent(http.StatusOK)
	}
}

type activityPostParam struct {
	ExhibitId string `json:"exhibit_id"`
	GuestId   string `json:"guest_id"`
}

type session struct {
	SessionId      string    `json:"session_id"`
	ExhibitId      string    `json:"exhibit_id"`
	GuestId        string    `json:"guest_id"`
	EnterAt        time.Time `json:"enter_at"`
	EnterOperation string    `json:"enter_operation"`
	ExitAt         time.Time `json:"exit_at"`
	ExitOperation  string    `json:"exit_operation"`
	Available      int       `json:"available"`
	Note           string    `json:"note"`
}

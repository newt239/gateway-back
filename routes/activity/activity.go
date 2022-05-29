package activityRoute

import (
	"net/http"
	"strconv"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Enter() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		session_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
		sessionEx := session{
			SessionId:      session_id,
			ExhibitId:      c.FormValue("exhibit_id"),
			GuestId:        c.FormValue("guest_id"),
			EnterAt:        now,
			EnterOperation: user_id,
			Available:      1,
		}
		var result session
		db.Omit("exit_at", "exit_operation", "note").Create(&sessionEx).Scan(&result)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"session_id": session_id,
		})
	}
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

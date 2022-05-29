package guestRoute

import (
	"net/http"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Info() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		guest_id := c.Param("guest_id")
		var result guest
		db.Where("guest_id = ?", guest_id).First(&guest{}).Scan(&result)
		db.Close()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guest_id":       guest_id,
			"guest_type":     result.GuestType,
			"reservation_id": result.ReservationId,
			"part":           result.Part,
			"available":      result.Available,
		})
	}
}

type guest struct {
	GuestId       string `json:"guest_id"`
	GuestType     string `json:"guest_type"`
	ReservationId string `json:"reservation_id"`
	ExhibitId     string `json:"exhibit_id"`
	Part          string `json:"part"`
	UserId        string `json:"user_id"`
	Available     string `json:"available"`
	RegisterAt    string `json:"register_at"`
	RevokeAt      string `json:"revoke_at"`
	Note          string `json:"note"`
}

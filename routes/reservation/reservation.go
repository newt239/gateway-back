package reservationRoute

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
		var result reservation
		db.Where("reservation_id = ?", c.Param("reservation_id")).First(&reservation{}).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"reservation_id": result.ReservationId,
			"guest_type":     result.GuestType,
			"part":           result.Part,
			"count":          result.Count,
			"registered":     result.Registered,
			"available":      result.Available,
		})
	}
}

type reservation struct {
	ReservationId string `json:"reservation_id"`
	GuestType     string `json:"guest_type"`
	Count         int    `json:"count"`
	Registered    int    `json:"registered"`
	Part          string `json:"part"`
	Available     int    `json:"available"`
	Note          string `json:"note"`
}

package activityRoute

import (
	"net/http"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func InfoEachExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		exhibit_id := c.Param("exhibit_id")
		var result exhibit
		db.Where("exhibit_id = ?", exhibit_id).First(&exhibit{}).Scan(&result)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"exhibit_id":   result.ExhibitId,
			"exhibit_name": result.ExhibitName,
			"exhibit_type": result.ExhibitType,
			"room_name":    result.RoomName,
			"Status":       result.Status,
		})
	}
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

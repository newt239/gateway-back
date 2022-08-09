package adminRoute

import (
	"net/http"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

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
		db.Close()

		return c.NoContent(http.StatusOK)
	}
}

func DeleteExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)
		db.Exec("DELETE FROM gateway.exhibit WHERE exhibit_id= ? ;", c.Param("exhibit_id"))
		db.Close()

		return c.NoContent(http.StatusOK) // status code 200で何も返さない
	}
}

type exhibit struct {
	ExhibitId   string    `json:"exhibit_id"`
	ExhibitName string    `json:"exhibit_name"`
	RoomName    string    `json:"room_name"`
	ExhibitType string    `json:"exhibit_type"`
	Status      int       `json:"status"`
	Capacity    int       `json:"capacity"`
	Position    string    `json:"position"`
	LastUpdate  time.Time `json:"last_update"`
	Note        string    `json:"note"`
}

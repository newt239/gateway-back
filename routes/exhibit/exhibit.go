package exhibitRoute

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func ExhibitList() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		type exhibitListParam struct {
			ExhibitId   string `json:"exhibit_id"`
			ExhibitName string `json:"exhibit_name"`
			Status      int    `json:"status"`
		}
		var result []exhibitListParam
		db.Table("exhibit").Find(&exhibit{}).Scan(&result)
		fmt.Println(result)
		return c.JSON(http.StatusOK, result)
	}
}

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

func CurrentEachExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		exhibit_id := c.Param("exhibit_id")
		result := db.Where("exhibit_id = ?", exhibit_id).Where("exit_at", gorm.Expr("NULL")).First(&session{})
		return c.JSON(http.StatusOK, map[string]interface{}{
			"count": result.RowsAffected,
		})
	}
}

func HistoryEachExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		type historyParam struct {
			Time  time.Time `json:"time"`
			Count int       `json:"count"`
		}
		var result []historyParam
		db.Raw("SELECT timestamp(DATE_FORMAT(enter_at, '%Y-%m-%d %H:00:00')) AS time, COUNT(*) AS count FROM gateway.session WHERE exhibit_id = ? AND DATE(enter_at) = ? GROUP BY DATE_FORMAT(enter_at, '%Y%m%d%H');", c.Param("exhibit_id"), c.Param("day")).Scan(&result)
		return c.JSON(http.StatusOK, result)
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

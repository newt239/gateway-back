package exhibitRoute

import (
	"fmt"
	"net/http"
	"time"

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
			GroupName   string `json:"group_name"`
			ExhibitType string `json:"exhibit_type"`
		}
		var result []exhibitListParam
		db.Table("exhibit").Where("status = 1").Find(&exhibit{}).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
	}
}

func InfoAllExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		type infoAllExhibitParam struct {
			GuestType string `json:"guest_type"`
			Count     int    `json:"count"`
		}
		var result []infoAllExhibitParam
		db.Raw("select count(*)as count, guest.guest_type from gateway.session inner join gateway.guest on session.guest_id = guest.guest_id  where exhibit_id = 'entrance' and exit_at is null group by guest.guest_type;").Scan(&result)
		db.Close()

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
		var countResult struct{ Current int }
		db.Raw("select count(*) as current from gateway.session where exhibit_id = ? and exit_at is null group by guest_id", exhibit_id).Scan(countResult)
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"exhibit_id":   result.ExhibitId,
			"exhibit_name": result.ExhibitName,
			"exhibit_type": result.ExhibitType,
			"room_name":    result.RoomName,
			"capacity":     result.Capacity,
			"current":      countResult.Current,
			"status":       result.Status,
		})
	}
}

func CurrentAllExhibitData() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		type currentEachExhibitParam struct {
			ID          string `json:"id"`
			ExhibitName string `json:"exhibit_name"`
			Count       int    `json:"count"`
			Capacity    int    `json:"capacity"`
		}
		var result []currentEachExhibitParam
		db := database.ConnectGORM(user_id, password)
		db.Raw("SELECT exhibit.exhibit_id AS id, exhibit_name, ifnull(count, 0) as count, capacity FROM exhibit LEFT JOIN current ON exhibit.exhibit_id = current.exhibit_id;").Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
	}
}

func CurrentEachExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		exhibit_id := c.Param("exhibit_id")

		type currentEachExhibitParam struct {
			ID        string `json:"id"`
			SessionId string `json:"session_id"`
			GuestType string `json:"guest_type"`
			EnterAt   string `json:"enter_at"`
		}
		var result []currentEachExhibitParam
		db.Raw(fmt.Sprintf("SELECT session.guest_id AS id, session.session_id, guest_type, enter_at FROM session INNER JOIN guest ON session.guest_id = guest.guest_id WHERE session.exhibit_id='%s' AND session.exit_at IS NULL;", exhibit_id)).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
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
		db.Close()

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

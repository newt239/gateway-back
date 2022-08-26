package exhibitRoute

import (
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
		db.Table("exhibit").Scan(&result)
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
		db.Raw(`
			select guest_type, count(*) as count 
			from gateway.guest 
			where guest_id in ( 
				select r.guest_id 
				from ( 
					select guest_id, count(*) as count 
					from gateway.activity 
					where exhibit_id = 'entrance' 
					group by guest_id 
				) as r 
				where mod(r.count, 2) = 1
			)
			group by guest_type;
			`).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
	}
}

func InfoEachExhibit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		exhibit_id := c.Param("exhibit_id")

		db := database.ConnectGORM(user_id, password)
		var result exhibit
		db.Table("exhibit").Where("exhibit_id = ?", exhibit_id).First(&exhibit{}).Scan(&result)
		type currentGuestListType struct {
			GuestId string
		}
		var currentGuestListResult []currentGuestListType
		db.Raw(`
			select r.guest_id 
			from ( 
				select guest_id, count(*) as count 
				from gateway.activity 
				where exhibit_id = ? 
				group by guest_id 
			) as r 
			where mod(r.count, 2) = 1;
		`, exhibit_id).Scan(&currentGuestListResult)
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"exhibit_id":   result.ExhibitId,
			"exhibit_name": result.ExhibitName,
			"exhibit_type": result.ExhibitType,
			"room_name":    result.RoomName,
			"capacity":     result.Capacity,
			"current":      len(currentGuestListResult),
		})
	}
}

func CurrentAllExhibitData() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		type currentEachExhibitParam struct {
			ExhibitID   string `json:"exhibit_id"`
			ExhibitName string `json:"exhibit_name"`
			GroupName   string `json:"group_name"`
			RoomName    string `json:"room_name"`
			ExhibitType string `json:"exhibit_type"`
			Count       int    `json:"count"`
			Capacity    int    `json:"capacity"`
		}
		var result []currentEachExhibitParam
		db := database.ConnectGORM(user_id, password)
		db.Raw(`
			select exhibit.exhibit_id, exhibit_name, group_name, room_name, exhibit_type, ifnull(count, 0) as count, capacity 
			from gateway.exhibit 
			left join (
				select r.exhibit_id, count(*) as count 
				from ( 
					select guest_id, count(*) as count, exhibit_id 
					from gateway.activity 
					group by guest_id, exhibit_id 
				) as r 
				where mod(r.count, 2) = 1 
				group by r.exhibit_id
			) as current
			on exhibit.exhibit_id = current.exhibit_id;
		`).Scan(&result)
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
			GuestType string `json:"guest_type"`
			EnterAt   string `json:"enter_at"`
		}
		var result []currentEachExhibitParam
		db.Raw(`
			select s.guest_id as id, guest_type, timestamp 
			FROM ( 
				select r.guest_id, r.timestamp 
				from ( 
					select guest_id, max(timestamp) as enter_at, count(*) as count 
					from gateway.activity 
					where exhibit_id = ? 
					group by guest_id 
				) as r 
				where mod(r.count, 2) = 1 
			) as s 
			inner join guest 
			on s.guest_id = guest.guest_id;
		`, exhibit_id).Scan(&result)
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
		db.Raw(`
			select timestamp(DATE_FORMAT(timestamp, '%Y-%m-%d %H:00:00')) as time, COUNT(timestamp) as count 
			from gateway.activity 
			where exhibit_id = ? and activity_type = 'enter' and DATE(timestamp) = ? 
			group by time;
		`, c.Param("exhibit_id"), c.Param("day")).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
	}
}

type exhibit struct {
	ExhibitId   string `json:"exhibit_id"`
	ExhibitName string `json:"exhibit_name"`
	RoomName    string `json:"room_name"`
	ExhibitType string `json:"exhibit_type"`
	Capacity    int    `json:"capacity"`
}

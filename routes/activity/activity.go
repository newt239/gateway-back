package activityRoute

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type activity struct {
	ActivityId   int64     `json:"activity_id"`
	GuestId      string    `json:"guest_id"`
	ExhibitId    string    `json:"exhibit_id"`
	ActivityType string    `json:"activity_type"`
	UserId       string    `json:"user_id"`
	Timestamp    time.Time `json:"timestamp"`
	Available    int       `json:"available"`
}

type activityPostParam struct {
	GuestId   string `json:"guest_id"`
	ExhibitId string `json:"exhibit_id"`
}

func Enter() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		enterPostParam := activityPostParam{}
		if err := c.Bind(&enterPostParam); err != nil {
			return err
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		activity_id := now.UnixMilli()
		activityEx := activity{
			ActivityId:   activity_id,
			ExhibitId:    enterPostParam.ExhibitId,
			GuestId:      enterPostParam.GuestId,
			ActivityType: "enter",
			Timestamp:    now,
			UserId:       user_id,
			Available:    1,
		}
		db := database.ConnectGORM(user_id, password)
		db.Table("activity").Create(&activityEx)
		db.Close()
		return c.NoContent(http.StatusOK)
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
		activity_id := now.UnixMilli()
		activityEx := activity{
			ActivityId:   activity_id,
			ExhibitId:    exitPostParam.ExhibitId,
			GuestId:      exitPostParam.GuestId,
			ActivityType: "exit",
			Timestamp:    now,
			UserId:       user_id,
			Available:    1,
		}
		db := database.ConnectGORM(user_id, password)
		db.Table("activity").Create(&activityEx)
		db.Close()
		return c.NoContent(http.StatusOK)
	}
}

func BatchExit() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		exitPostParams := []activityPostParam{}
		if err := c.Bind(&exitPostParams); err != nil {
			return err
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		str := "insert into gateway.activity (activity_id, exhibit_id, guest_id, activity_type, timestamp, user_id, available) values "
		var s []string
		for _, u := range exitPostParams {
			now := time.Now().In(jst)
			activity_id := now.UnixMilli()
			q := fmt.Sprintf("(%d, '%s', '%s', 'exit', now(), '%s', 1), ", activity_id, u.ExhibitId, u.GuestId, user_id)
			s = append(s, q)
		}
		query := strings.TrimRight(strings.Join(s, ""), ", ") + ";"
		db := database.ConnectGORM(user_id, password)
		db.Exec(str + query)
		db.Close()
		return c.NoContent(http.StatusOK)
	}
}

func History() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)
		t, _ := time.Parse("2006-01-02T15:04:05+09:00", c.Param("from"))
		type activityHistoryListType struct {
			ActivityId   string `json:"activity_id"`
			GuestId      string `json:"guest_id"`
			ExhibitId    string `json:"exhibit_id"`
			ActivityType string `json:"activity_type"`
			Timestamp    string `json:"timestamp"`
		}
		var activityList []activityHistoryListType
		db.Raw(`
			select activity_id, guest_id, exhibit_id, activity_type, timestamp 
			from gateway.activity 
			where timestamp > ? 
			limit 100
		`, t).Scan(&activityList)
		db.Close()
		return c.JSON(http.StatusOK, activityList)
	}
}

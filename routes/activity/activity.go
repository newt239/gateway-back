package activityRoute

import (
	"net/http"
	"strconv"
	"time"

	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type activity struct {
	ActivityId   string    `json:"activity_id"`
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
		activity_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
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
		return c.JSON(http.StatusOK, map[string]interface{}{
			"activity_id": activity_id,
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
		activity_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
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
		return c.JSON(http.StatusOK, map[string]interface{}{
			"activity_id": activity_id,
		})
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
		db.Table("activity").Where("timestamp > ?", t).Limit(100).Scan(&activityList)
		db.Close()
		return c.JSON(http.StatusOK, activityList)
	}
}

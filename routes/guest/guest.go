package guestRoute

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/newt239/gateway-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Info() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		db := database.ConnectGORM(user_id, password)

		guest_id := c.Param("guest_id")
		type guest struct {
			GuestId       string `json:"guest_id"`
			GuestType     string `json:"guest_type"`
			ReservationId string `json:"reservation_id"`
			Part          int    `json:"part"`
			Available     int    `json:"available"`
		}
		var guestInfoResult guest
		db.Table("guest").Where("guest_id = ?", guest_id).First(&guestInfoResult)

		type activityInfoResultParam struct {
			ExhibitId    string `json:"exhibit_id"`
			ActivityType string `json:"activity_type"`
		}
		var activityInfoResult activityInfoResultParam
		exhibit_id := ""
		err := db.Table("activity").Select("exhibit_id", "activity_type").Where("guest_id = ?", guest_id).Where("exhibit_id != 'entrance'").Order("timestamp desc").First(&activityInfoResult).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) && activityInfoResult.ActivityType == "enter" {
			exhibit_id = activityInfoResult.ExhibitId
		}

		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"guest_id":       guest_id,
			"guest_type":     guestInfoResult.GuestType,
			"reservation_id": guestInfoResult.ReservationId,
			"part":           guestInfoResult.Part,
			"available":      guestInfoResult.Available,
			"exhibit_id":     exhibit_id,
		})
	}
}

func Activity() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		guest_id := c.Param("guest_id")
		type resultParam struct {
			ExhibitId    string `json:"exhibit_id"`
			ActivityType string `json:"activity_type"`
			Timestamp    string `json:"timestamp"`
		}
		var result []resultParam
		db := database.ConnectGORM(user_id, password)
		db.Raw(`
			select exhibit_id, activity_type, timestamp 
			from gateway.activity 
			where guest_id = ?;
		`, guest_id).Scan(&result)
		db.Close()

		return c.JSON(http.StatusOK, result)
	}
}

func Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		type guestRegisterPostParam struct {
			ReservationId string   `json:"reservation_id"`
			GuestType     string   `json:"guest_type"`
			GuestIdList   []string `json:"guest_id"`
			Part          int      `json:"part"`
		}
		registerPostData := guestRegisterPostParam{}
		if err := c.Bind(&registerPostData); err != nil {
			return err
		}

		db := database.ConnectGORM(user_id, password)
		for _, guest_id := range registerPostData.GuestIdList {
			jst, _ := time.LoadLocation("Asia/Tokyo")
			now := time.Now().In(jst)
			activity_id := now.UnixMilli()
			db.Table("guest").Omit("exhibit_id", "revoke_at").Where("guest_id = ?", guest_id).Update(&guestParam{
				ReservationId: registerPostData.ReservationId,
				GuestType:     registerPostData.GuestType,
				Part:          registerPostData.Part,
				UserId:        user_id,
				RegisterAt:    now,
				Available:     1,
			})
			db.Table("activity").Create(&activity{
				ActivityId:   activity_id,
				GuestId:      guest_id,
				ExhibitId:    "entrance",
				ActivityType: "enter",
				UserId:       user_id,
				Timestamp:    now,
				Available:    1,
			})
		}
		db.Close()

		return c.NoContent(http.StatusOK)
	}
}

func Revoke() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		type guestRegisterPostParam struct {
			ReservationId string `json:"reservation_id"`
			GuestType     string `json:"guest_type"`
			NewGuestId    string `json:"new_guest_id"`
			OldGuestId    string `json:"old_guest_id"`
			Part          int    `json:"part"`
		}
		registerPostData := guestRegisterPostParam{}
		if err := c.Bind(&registerPostData); err != nil {
			return err
		}

		db := database.ConnectGORM(user_id, password)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		activity_id := now.UnixMilli()
		db.Table("guest").Omit("exhibit_id", "revoke_at").Create(&guestParam{
			ReservationId: registerPostData.ReservationId,
			GuestId:       registerPostData.NewGuestId,
			GuestType:     registerPostData.GuestType,
			Part:          registerPostData.Part,
			UserId:        user_id,
			RegisterAt:    now,
			Available:     1,
		})
		db.Table("activity").Create(&activity{
			ActivityId:   activity_id,
			GuestId:      registerPostData.NewGuestId,
			ExhibitId:    "entrance",
			ActivityType: "exit",
			UserId:       user_id,
			Timestamp:    now,
			Available:    1,
		})
		if registerPostData.OldGuestId != "" {
			db.Table("guest").Where("guest_id = ?", registerPostData.OldGuestId).Update(&guestParam{
				RevokeAt:  now,
				Available: 0,
			})
		}
		db.Close()

		return c.NoContent(http.StatusOK)
	}
}

type guestParam struct {
	ReservationId string    `json:"reservation_id"`
	GuestId       string    `json:"guest_id"`
	GuestType     string    `json:"guest_type"`
	Part          int       `json:"part"`
	UserId        string    `json:"user_id"`
	RegisterAt    time.Time `json:"register_at"`
	RevokeAt      time.Time `json:"revoke_at"`
	Available     int       `json:"available"`
	IsSpare       string    `json:"is_spare"`
}

type activity struct {
	ActivityId   int64     `json:"activity_id"`
	GuestId      string    `json:"guest_id"`
	ExhibitId    string    `json:"exhibit_id"`
	ActivityType string    `json:"activity_type"`
	UserId       string    `json:"user_id"`
	Timestamp    time.Time `json:"timestamp"`
	Available    int       `json:"available"`
}

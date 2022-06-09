package guestRoute

import (
	"errors"
	"fmt"
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
			Part          string `json:"part"`
			Available     int    `json:"available"`
		}
		var guestInfoResult guest
		db.Table("guest").Where("guest_id = ?", guest_id).First(&guestInfoResult)
		fmt.Println(guestInfoResult)

		type sessionInfoResultParam struct {
			ExhibitId string `json:"exhibit_id"`
		}
		var sessionInfoResult sessionInfoResultParam
		exhibit_id := ""
		err := db.Table("session").Select("exhibit_id").Where("guest_id = ?", guest_id).Where("exit_at IS NULL").First(&sessionInfoResult).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			exhibit_id = sessionInfoResult.ExhibitId
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

func Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)

		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		type guestRegisterPostParam struct {
			ReservationId string   `json:"reservation_id"`
			GuestType     string   `json:"guest_type"`
			GuestIdList   []string `json:"guest_id"`
			Part          string   `json:"part"`
		}
		registerPostData := guestRegisterPostParam{}
		if err := c.Bind(&registerPostData); err != nil {
			return err
		}

		type guestParam struct {
			ReservationId string    `json:"reservation_id"`
			GuestId       string    `json:"guest_id"`
			GuestType     string    `json:"guest_type"`
			Part          string    `json:"part"`
			UserId        string    `json:"user_id"`
			RegisterAt    time.Time `json:"register_at"`
			Available     int       `json:"available"`
		}

		db := database.ConnectGORM(user_id, password)
		for _, v := range registerPostData.GuestIdList {
			db.Omit("exhibit_id", "revoke_at", "note").Table("guest").Create(&guestParam{
				ReservationId: registerPostData.ReservationId,
				GuestId:       v,
				GuestType:     registerPostData.GuestType,
				Part:          registerPostData.Part,
				UserId:        user_id,
				RegisterAt:    now,
				Available:     1,
			})
		}
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{})
	}
}

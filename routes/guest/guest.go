package guestRoute

import (
	"errors"
	"net/http"
	"strconv"
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

		type sessionInfoResultParam struct {
			ExhibitId string `json:"exhibit_id"`
		}
		var sessionInfoResult sessionInfoResultParam
		exhibit_id := ""
		err := db.Table("session").Select("exhibit_id").Where("guest_id = ?", guest_id).Where("exit_at IS NULL").Where("exhibit_id != 'entrance'").First(&sessionInfoResult).Error
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

func Activity() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))
		guest_id := c.Param("guest_id")
		type resultParam struct {
			ExhibitId string `json:"exhibit_id"`
			EnterAt   string `json:"enter_at"`
			ExitAt    string `json:"exit_at"`
		}
		var result []resultParam
		db := database.ConnectGORM(user_id, password)
		db.Raw("select exhibit_id, enter_at, ifnull(exit_at, 'current') as exit_at from session where guest_id = ? order by enter_at;", guest_id).Scan(&result)
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
			session_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
			db.Table("guest").Omit("exhibit_id", "revoke_at", "note").Create(&guestParam{
				ReservationId: registerPostData.ReservationId,
				GuestId:       guest_id,
				GuestType:     registerPostData.GuestType,
				Part:          registerPostData.Part,
				UserId:        user_id,
				RegisterAt:    now,
				Available:     1,
			})
			db.Table("session").Omit("exit_at", "exit_operation", "note").Create(&sessionParam{
				SessionId:      session_id,
				GuestId:        guest_id,
				ExhibitId:      "entrance",
				EnterAt:        now,
				EnterOperation: user_id,
				Available:      1,
			})
		}
		db.Close()

		return c.JSON(http.StatusOK, map[string]interface{}{})
	}
}

func Revoke() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id, password := database.CheckJwt(c.Get("user").(*jwt.Token))

		type guestRegisterPostParam struct {
			ReservationId string `json:"reservation_id"`
			GuestType     string `json:"guest_type"`
			GuestId       string `json:"guest_id"`
			Part          int    `json:"part"`
		}
		registerPostData := guestRegisterPostParam{}
		if err := c.Bind(&registerPostData); err != nil {
			return err
		}

		db := database.ConnectGORM(user_id, password)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		now := time.Now().In(jst)
		session_id := "s" + strconv.FormatInt(now.UnixMilli(), 10)
		db.Table("guest").Omit("exhibit_id", "revoke_at").Create(&guestParam{
			ReservationId: registerPostData.ReservationId,
			GuestId:       registerPostData.GuestId,
			GuestType:     registerPostData.GuestType,
			Part:          registerPostData.Part,
			UserId:        user_id,
			RegisterAt:    now,
			Available:     1,
			Note:          "spare",
		})
		db.Table("session").Omit("exit_at", "exit_operation", "note").Create(&sessionParam{
			SessionId:      session_id,
			GuestId:        registerPostData.GuestId,
			ExhibitId:      "info_center",
			EnterAt:        now,
			EnterOperation: user_id,
			Available:      1,
		})
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
	Available     int       `json:"available"`
	Note          string    `json:"spare"`
}

type sessionParam struct {
	SessionId      string    `json:"session_id"`
	GuestId        string    `json:"guest_id"`
	ExhibitId      string    `json:"exhibit_id"`
	EnterAt        time.Time `json:"enter_at"`
	EnterOperation string    `json:"enter_operation"`
	ExitAt         time.Time `json:"exit_at"`
	ExitOperation  string    `json:"exit_operation"`
	Available      int       `json:"available"`
}

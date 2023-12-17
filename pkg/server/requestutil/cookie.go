package requestutil

import (
	"dokemon/pkg/crypto/ske"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthCookieContent struct {
	UserName string    `json:"userName"`
	Expiry   time.Time `json:"expiry"`
}

const AUTH_COOKIE_NAME = "dmauth"

func SetAuthCookie(c echo.Context, cc AuthCookieContent) {
	ccBytes, err := json.Marshal(cc)
	if err != nil {
		panic(err)
	}

	encryptedCC, err := ske.Encrypt(string(ccBytes))
	if err != nil {
		panic(err)
	}

	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Name = AUTH_COOKIE_NAME
	cookie.Value = encryptedCC
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func GetAuthCookie(c echo.Context) (*AuthCookieContent, error) {
	cookie, err := c.Cookie(AUTH_COOKIE_NAME)
	if err != nil {
		return nil, err
	}

	ccText, err := ske.Decrypt(cookie.Value)
	if err != nil {
		return nil, err
	}

	var cc AuthCookieContent
	err = json.Unmarshal([]byte(ccText), &cc)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

func DeleteAuthCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Name = AUTH_COOKIE_NAME
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Expires =  time.Unix(0, 0)
	c.SetCookie(cookie)
}
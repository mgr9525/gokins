package core

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const cookieName = "gokinsk"

func BindMapJSON(c *gin.Context) (*ruisUtil.Map, error) {
	pars := ruisUtil.NewMap()
	err := c.BindJSON(pars)
	return pars, err
}

func CreateToken(p *jwt.MapClaims, tmout time.Duration) (string, error) {
	claims := *p
	claims["times"] = time.Now().Format(time.RFC3339Nano)
	if tmout > 0 {
		claims["timeout"] = time.Now().Add(tmout).Format(time.RFC3339Nano)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokens, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return "", err
	}
	return tokens, nil
}
func SetToken(c *gin.Context, p *jwt.MapClaims, rem bool, doman ...string) (string, error) {
	tmout := time.Hour * 5
	if rem {
		tmout = time.Hour * 24 * 5
	}
	tokens, err := CreateToken(p, tmout)
	if err != nil {
		return "", err
	}
	cke := http.Cookie{Name: cookieName, Value: tokens, HttpOnly: false}
	if JwtCookiePath != "" {
		cke.Path = JwtCookiePath
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}

	cke.MaxAge = 60 * 60 * 5
	if rem {
		cke.MaxAge = 60 * 60 * 24 * 5
	}
	c.Header("Set-Cookie", cke.String())
	return tokens, nil
}

func ClearToken(c *gin.Context, doman ...string) {
	cke := http.Cookie{Name: cookieName, Value: "", HttpOnly: false}
	if JwtCookiePath != "" {
		cke.Path = JwtCookiePath
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}
	cke.MaxAge = -1
	c.Header("Set-Cookie", cke.String())
}

var secret = func(token *jwt.Token) (interface{}, error) {
	return []byte(JwtKey), nil
}

func GetToken(c *gin.Context) jwt.MapClaims {
	tks := ""
	ats := c.Request.Header.Get("Authorization")
	if ats != "" {
		aths, err := url.PathUnescape(ats)
		if err == nil && strings.HasPrefix(aths, "TOKEN ") {
			tks = strings.Replace(aths, "TOKEN ", "", 1)
		}
	}
	if tks == "" {
		tkc, err := c.Request.Cookie(cookieName)
		if err == nil {
			tks = tkc.Value
		}
	}
	if tks != "" {
		tk, err := jwt.Parse(tks, secret)
		if err == nil {
			claim, ok := tk.Claims.(jwt.MapClaims)
			if ok {
				return claim
			}
		}
	}
	return nil
}

func MidAccessAllow(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.String(200, "request ok!")
		c.Abort()
	}
}

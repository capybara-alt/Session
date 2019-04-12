package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

func secret() echo.HandlerFunc{
	return func(c echo.Context)error{
		//sessionを見る
		sess, err := session.Get("session", c)
		if err!=nil {
			return c.String(http.StatusInternalServerError, "Error")
		}
		//ログインしているか
		if b, _:=sess.Values["auth"];b!=true{
			return c.String(http.StatusUnauthorized, "401")
		}else {
			return c.String(http.StatusOK, sess.Values["foo"].(string))
		}
	}
}

func login(e *echo.Echo) echo.HandlerFunc{
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			//Path:でsessionの有効な範囲を指定｡指定無しで全て有効になる｡
			//有効な時間
			MaxAge:   86400 * 7,
			//trueでjsからのアクセス拒否
			HttpOnly: true,
		}
		//テキトウな値
		sess.Values["foo"] = "bar"
		//ログインしました
		sess.Values["auth"] = true
		//状態保存
		if err:=sess.Save(c.Request(), c.Response());err!=nil{
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func logout() echo.HandlerFunc{
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		//ログアウト
		sess.Values["auth"]=false
		//状態を保存
		if err:=sess.Save(c.Request(), c.Response());err!=nil{
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func main() {
	e := echo.New()
	e.GET("/login",login(e))
	e.GET("/logout",logout())
	e.GET("/secret",secret())
	e.Logger.Fatal(e.Start(":8080"))
}

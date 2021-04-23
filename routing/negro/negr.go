package negr

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	r "../router"
	"../router/Controller"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/urfave/negroni"
)

var instance *negroniWrapper
var once sync.Once

type negroniWrapper struct {
	negr      *negroni.Negroni
	webServer *http.Server

	socket_negr      *negroni.Negroni
	socket_webServer *http.Server
}

func GetInstance() *negroniWrapper {
	once.Do(func() {
		instance = &negroniWrapper{
			negr:        negroni.Classic(),
			socket_negr: negroni.Classic(),
		}
	})

	return instance
}

const (
	// 애플리케이션에서 사용할 세션의 키 정보
	sessionKey    = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
)

func InitNegroni() {

	// negroni 미들웨어 생성
	nw := GetInstance()

	// negroni에 session cookie 등록
	store := cookiestore.New([]byte(sessionSecret))
	nw.negr.Use(sessions.Sessions(sessionKey, store))
	nw.negr.Use(negroni.HandlerFunc(Controller.Auth))

	// negroni에 router를 핸들러로 등록
	router := r.GetInstance().Router
	nw.negr.UseHandler(router)

	nw.webServer = &http.Server{
		Addr:    ":8889",
		Handler: nw.negr,
	}

	go func() {
		nw.webServer.ListenAndServe()
		// nw.negr.Run(":8889")
	}()
}

func Fin() {
	// 일반 라우터 종료
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := GetInstance().webServer.Shutdown(ctx)
	if nil != err {
		fmt.Println("Error::Stop Router")
		fmt.Println(err.Error())
	}
}

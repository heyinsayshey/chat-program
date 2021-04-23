package router

import (
	"sync"

	"./Controller"
	"github.com/julienschmidt/httprouter"
)

var instance *RouterWrapper
var once sync.Once

type RouterWrapper struct {
	Router *httprouter.Router
	//WebServer *http.Server
}

func GetInstance() *RouterWrapper {
	once.Do(func() {
		instance = &RouterWrapper{}
	})

	return instance
}

func InitRouter() {
	//go GetInstance().startRoute()
	GetInstance().startRoute()
}

func (w *RouterWrapper) startRoute() {
	w.Router = httprouter.New()

	w.Router.GET("/test", Controller.Test)

	// login
	w.Router.GET("/loginview", Controller.LoginView)
	w.Router.POST("/login", Controller.Login)
	w.Router.GET("/logout", Controller.Logout)

	//main
	w.Router.GET("/main", Controller.MainView)

	//room
	w.Router.POST("/room/create", Controller.CreateRoom)
	w.Router.POST("/room/search", Controller.SearchRooms)

	//message
	w.Router.POST("/message/search", Controller.SearchMessages)

	w.Router.GET("/ws/:room_id", Controller.CreateClient)
}

func Fin() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// err := GetInstance().WebServer.Shutdown(ctx)
	// if nil != err {
	// 	fmt.Println("Error::Stop Router")
	// 	fmt.Println(err.Error())
	// }
}

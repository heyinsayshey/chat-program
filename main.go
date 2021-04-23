package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	//"github.com/julienschmidt/httprouter"
	//"github.com/unrolled/render"

	negr "./routing/negro"
	"./routing/renderer"
	"./routing/router"
)

var pidFile string

func main() {
	pidFile := "chat.pid"
	os.Remove(pidFile)

	if _, err := os.Stat(pidFile); os.IsExist(err) {
		fmt.Println("please check [", pidFile, "] file. exit.")
		fmt.Errorf(err.Error())
		return
	}

	f, e := os.OpenFile(pidFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if nil != e {
		fmt.Println(e)
		return
	}

	f.WriteString(strconv.FormatInt(int64(os.Getpid()), 10))
	f.Close()

	startService()

	os.Remove(pidFile)
}

func startService() {
	renderer.InitRenderer()
	fmt.Println("Render::Start - Render initialized")

	router.InitRouter()
	fmt.Println("Router:::Start - router initialized")

	negr.InitNegroni()
	fmt.Println("Negroni::Start - Negroni initialized")

	waitObject()

	stopService()
}

func waitObject() {
	chMainStop := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(chMainStop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)

	go func() {
		sig := <-chMainStop
		fmt.Printf("EXIT::Get OS signal[%v].\n", sig)
		done <- true
	}()

	<-done
}

func stopService() {
	negr.Fin()
	fmt.Println("Negroni::End - Stop Service.")

	router.Fin()
	fmt.Println("Router::End - Stop service.")

	renderer.Fin()
	fmt.Println("Render::End - Stop service.")
}

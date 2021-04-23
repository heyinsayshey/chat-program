package Controller

import (
	"fmt"
	"net/http"

	r "../../renderer"
	"github.com/julienschmidt/httprouter"
)

func MainView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ren := r.GetInstance()

	err := ren.HTML(w, http.StatusOK, "main", map[string]string{"host": "localhost:8889"})
	if nil != err {
		fmt.Println(err.Error())
	}
}

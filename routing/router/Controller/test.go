package Controller

import (
	"fmt"
	"net/http"

	r "../../renderer"
	"github.com/julienschmidt/httprouter"
)

func Test(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ren := r.GetInstance()

	err := ren.HTML(w, http.StatusOK, "example", map[string]string{"title": "Simple Haein chat!!!"})
	if nil != err {
		fmt.Println(err.Error())
	}
}

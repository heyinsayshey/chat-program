package Controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	parser "../../../common/JsonParser"
	"../../../common/RDBMS"
	r "../../renderer"
	"github.com/julienschmidt/httprouter"
)

func Login(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var reqBody map[string]string
	var err error

	reqBody, err = parser.HTMLRequestBodyParser(req)
	if nil != err {
		log.Fatal(err)
		return
	}

	if nil == reqBody {
		log.Fatal("empty reqBody.")
		return
	}

	id := reqBody["id"]
	pw := reqBody["password"]

	sql := fmt.Sprintf("call SP_LOGIN('%v', '%v')", id, pw)
	val, err := RDBMS.MySQLSelect(sql)
	if nil != err {
		log.Fatal(err)
		return
	}

	if strings.Compare(val, "200") == 0 {
		fmt.Println("login success.")

		sql := fmt.Sprintf("select thumbnail_path from tbl_member where id = '%v';", id)
		val, err := RDBMS.MySQLSelect(sql)
		if nil != err {
			log.Fatal(err)
			return
		}

		user := User{
			ID:             id,
			Thumbnail_path: val,
			Timestamp:      time.Now(),
		}

		SetCurrentUser(req, &user)

		http.Redirect(w, req, "/main", http.StatusFound)
		return

	} else {
		fmt.Println("wrong password.")
		ren := r.GetInstance()
		ren.HTML(w, http.StatusMovedPermanently, "login", map[string]string{"alert": "로그인 정보가 일치하지 않습니다."})
		//http.Redirect(w, req, "/loginview", http.StatusNotAcceptable)
		return
	}
}

func LoginView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ren := r.GetInstance()
	err := ren.HTML(w, http.StatusOK, "login", nil)
	if nil != err {
		fmt.Println(err.Error())
	}
}

func Logout(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	DeleteSession(req)

	http.Redirect(w, req, "/loginview", http.StatusFound)
}

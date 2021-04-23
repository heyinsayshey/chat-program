package Controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	r "../../renderer"
	sessions "github.com/goincremental/negroni-sessions"
)

type User struct {
	ID             string    `json:"id"`
	Thumbnail_path string    `json:"thumbnail_path"`
	Timestamp      time.Time `json:"timestamp"`
}

const (
	USERKEY          = "simple_chat_user_key"
	SESSION_DURATION = time.Hour * 1
)

func Auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var msg string
	var tmpl string

	for {
		if strings.HasPrefix(req.URL.RequestURI(), "/login") {
			next(rw, req)
			return
		}

		u := GetCurrentUser(req)
		if u == nil {
			tmpl = "login"
			msg = "로그인이 필요한 서비스 입니다."
			break
		}

		if !u.Valid() {
			tmpl = "login"
			msg = "로그인 세션이 만료되었습니다. 다시 로그인 해주세요."
			DeleteSession(req)
			break
		}

		SetCurrentUser(req, u)
		next(rw, req)
		return
	}

	ren := r.GetInstance()
	ren.HTML(rw, http.StatusBadRequest, tmpl, map[string]string{"alert": msg})
}

func (u *User) Valid() bool {
	return time.Now().Sub(u.Timestamp) < SESSION_DURATION
}

func (u *User) Refresh() {
	u.Timestamp = time.Now()
}

func GetCurrentUser(r *http.Request) *User {
	s := sessions.GetSession(r)

	if s.Get(USERKEY) == nil {
		return nil
	}

	data := s.Get(USERKEY).([]byte)
	var u User
	json.Unmarshal(data, &u)
	return &u
}

func SetCurrentUser(r *http.Request, u *User) {
	if u != nil {
		u.Refresh()
	}

	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(USERKEY, val)
}

func DeleteSession(r *http.Request) {
	s := sessions.GetSession(r)
	s.Delete(USERKEY)
	s.Clear()
}

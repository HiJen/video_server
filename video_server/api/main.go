package main

import (
	"go-learn1/video_server/api/session"
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "go-learn1/video_server/api/session"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHanders() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser) //CreateUser是一个func,not method

	router.POST("/user/:user_name", Login)

	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegisterHanders()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)
}

//访问:http://127.0.0.1:8002/user  没有s，not https
// listen-->RegisterHanders-->handlers，
// 每个handlers的goroutine 每个4k大小, 轻量级,类协程似的.

// handler--->validation{1.request, 2.user} -->business logic -->response.
// 1.data model
// 2.error handling.
// session

// request --->  main-->middleware-->defs(message, err)-->handlers-->dbops-->response-->

//middleware部分
// http.Handler a
// ab struct a
// duck type

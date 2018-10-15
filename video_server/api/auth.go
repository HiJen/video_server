package main

import (
	"go-learn1/video_server/api/defs"
	"go-learn1/video_server/api/session"
	"net/http"
)

var HEADER_FIELD_SESSION = "X-session-Id" //x开头是http自定义的header
var HEADER_FIELD_UNAME = "X-User_Name"

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}

	return true
}

// session部分
// 使用cache,减少db的访问量, 减少数据库压力.
// 1: session是否过期| 2: 存session的地方.
// 拉取session | 新用户新的session id(给他一个方法) | 过期或者不过期状态(判断用户是否合法).
package session

import (
	"go-learn1/video_server/api/dbops"
	"go-learn1/video_server/api/defs"
	"go-learn1/video_server/api/utils"
	"sync"
	"time"
)

// go1.9以后有的*sync.Map
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)   //del cache
	dbops.DeleteSession(sid) //del db
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := time.Now().UnixNano() / 1000000 //毫秒
	ttl := ct + 30*60*1000                //Serverside session valid time: 30min

	ss := &defs.SimpleSession{Username: un, TTL: ttl}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, un)

	return id
}

// sid=session id
//get cache by session
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		// ct := time.Now().UnixNano() / 100000
		// ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			//delete expired session
			deleteExpiredSession(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	} else {
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss == nil {
			return "", true
		}

		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		sessionMap.Store(sid, ss)
		return ss.Username, false
	}
	return "", true
}

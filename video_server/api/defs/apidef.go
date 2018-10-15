package defs

// requests
type UserCredential struct {
	Username string `json:"user_name"` //tag | go原生method处理json
	Pwd      string `json:"pwd"`
}

//Response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

// Data model
type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
	// CreateCtime  string
}

type Comment struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

type SimpleSession struct {
	Username string //LOGIN NAME
	TTL      int64  //过期
}

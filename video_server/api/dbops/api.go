package dbops

import (
	"database/sql"
	"log"
	"time"

	// "database/sql"

	"go-learn1/video_server/api/defs"
	"go-learn1/video_server/api/utils"

	_ "github.com/go-sql-driver/mysql"
)

// func openConn() *sql.DB {
// 	dbConn, err := sql.Open("mysql", "root:adminadmin@#@tcp(localhost:3306)/video_server?charset=utf8")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return dbConn
// }

//sql 返回值很标准的, 每一个语句后面都会跟一个err.
// driver里面的预编译--Prepare
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)") //不要用+拼接,很容易撞库和被攻击|Prepare预编译,更安全.
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close() //defer能不能就不用, 因为这里性能要求有余地. 目的怕每个if 错误跳出来就给它close掉.
	// stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()

	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=? ")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}

	_, err = stmDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmDel.Close()
	return nil
}

/*
// aid = auther id, name= video name
// standard, return all object, not one attribute!
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	//create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	// creatime-->db-->
	//go toutine 多个video写库的时候displayctime和createCtime跟数据库里面的顺序是永远保持一致的.
	t := time.Now()
	//you can add slash or adscore, but you can't modify any string!!!
	ctime := t.Format("Jan 02 2006, 15:04:05") //M D y, HH:MM:SS
	//换行时候,双引号换```，写库,预编译.
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info
	(id, author_id, name, display_ctime) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	//init object
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	defer stmtIns.Close()
	return res, nil
}
*/

//video 的增删改查!!
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info 
		(id, author_id, name, display_ctime) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}

	// res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}

	defer stmtIns.Close()
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}
	return res, nil
}

// func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
// 	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.display_ctime FROM video_info
// 	INNER JOIN users ON video_info.author_id=users.id
// 	WHERE users.login_name = ? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time <= FROM_UNIXTIME
// 	ORDER BY video_info.create_time DESC`)

// 	var res []*defs.VideoInfo

// 	if err != nil {
// 		return res, err
// 	}

// 	rows, err := stmtOut.Query(uname, from, to)
// 	if err != nil {
// 		log.Printf("%s", err)
// 		return res, err
// 	}

// 	for rows.Next() {
// 		var id, name, ctime string
// 		var aid int
// 		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
// 			return res, err
// 		}

// 		vi := &defs.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
// 		res = append(res, vi)
// 	}

// 	defer stmtOut.Close()

// 	return res, nil

// }

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

//评论
//aid=users id来输出Login_name | id,err==comments id
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

//<= 防止最后一条评论漏掉, 最后1秒时间和前一秒为同一秒的时候.
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	// users join comments (Out)comments -->author_id, video_id(in)
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.Login_name, comments.content FROM comments
	INNER JOIN users ON comments.author_id = users.id
	WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`)

	var res []*defs.Comment

	//返回rows指针.
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}
	//迭代器迭代每一行取值Scan出来
	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}

	defer stmtOut.Close()

	return res, nil

}

package db

import (
	"bytes"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"strings"
	. "tianchi/models"
	"time"
)

const (
	MaxChanLen        = 50000
	WriteIntervalNum  = 100
	WriteIntervalTime = 500 //毫秒
)

var WriteUserChan chan *User
var WriteRoomChan chan *Room
var WriteMsgChan chan *Message
var WriteUserTokenChan chan *Token
var UserEnterRoomChan chan string

func initWorkChan() {
	WriteUserChan = make(chan *User, MaxChanLen)
	WriteRoomChan = make(chan *Room, MaxChanLen)
	WriteMsgChan = make(chan *Message, MaxChanLen)
	WriteUserTokenChan = make(chan *Token, MaxChanLen)
	UserEnterRoomChan = make(chan string, MaxChanLen)

	timer := time.NewTicker(WriteIntervalTime * time.Millisecond)

	go func() {
		userList := make([]*User, 0)
		for {
			select {
			case user, ok := <-WriteUserChan:
				if ok {
					if len(userList) > WriteIntervalNum {
						userList = append(userList, user)
						writeList := userList
						userList = []*User{}
						go WriteUser(writeList)
					} else {
						userList = append(userList, user)
						continue
					}
				}

			case <-timer.C:
				if len(userList) > 0{
					writeList := userList
					userList = []*User{}
					go WriteUser(writeList)
				}
			}
		}
	}()

	go func() {
		roomList := make([]*Room, 0)
		for {
			select {
			case room, ok := <-WriteRoomChan:
				if ok {
					if len(roomList) > WriteIntervalNum {
						roomList = append(roomList, room)
						writeList := roomList
						roomList = []*Room{}
						go WriteRoom(writeList)
					} else {
						roomList = append(roomList, room)
						continue
					}
				}

			case <-timer.C:
				if len(roomList) > 0 {
					writeList := roomList
					roomList = []*Room{}
					go WriteRoom(writeList)
				}
			}
		}
	}()

	go func() {
		msgList := make([]*Message, 0)
		for {
			select {
			case msg, ok := <-WriteMsgChan:
				if ok {
					if len(msgList) > WriteIntervalNum {
						msgList = append(msgList, msg)
						writeList := msgList
						msgList = []*Message{}
						go WriteMsg(writeList)
					} else {
						msgList = append(msgList, msg)
						continue
					}
				}

			case <-timer.C:
				if len(msgList) > 0{
					writeList := msgList
					msgList = []*Message{}
					go WriteMsg(writeList)
				}
			}
		}
	}()

	go func() {
		tokenList := make([]*Token, 0)
		for {
			select {
			case userToken, ok := <-WriteUserTokenChan:
				if ok {
					if len(tokenList) > WriteIntervalNum {
						tokenList = append(tokenList, userToken)
						writeList := tokenList
						tokenList = []*Token{}
						go WriteUserToken(writeList)
					} else {
						tokenList = append(tokenList, userToken)
						continue
					}
				}

			case <-timer.C:
				if len(tokenList) > 0{
					writeList := tokenList
					tokenList = []*Token{}
					go WriteUserToken(writeList)
				}
			}
		}
	}()

	go func() {
		userRoomMap := make(map[string]string)
		for {
			select {
			case userRoom, ok := <-UserEnterRoomChan:
				if ok {
					ur := strings.Split(userRoom, "|")
					if len(ur) < 2 {
						continue
					}

					if len(userRoom) > WriteIntervalNum {
						userRoomMap[ur[0]] = ur[1]
						writeMap := userRoomMap
						userRoomMap = make(map[string]string)
						go UserEnterRoom(writeMap)
					} else {
						userRoomMap[ur[0]] = ur[1]
						continue
					}
				}

			case <-timer.C:
				if len(userRoomMap) > 0{
					writeMap := userRoomMap
					userRoomMap = make(map[string]string)
					go UserEnterRoom(writeMap)
				}
			}
		}
	}()
}

func WriteUser(userList []*User) {
	var buffer bytes.Buffer
	sql := "insert into `user` (`username`,`first_name`,`last_name`,`email`,`password`,`phone`) values"
	buffer.WriteString(sql)

	for i, user := range userList {
		if i == len(userList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s');", user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s'),", user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone))
		}
	}

	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("WriteUser err:%v", result.Error)
	}
	return
}

func WriteRoom(roomList []*Room) {
	var buffer bytes.Buffer
	sql := "insert into `room` (`id`,`name`) values"
	buffer.WriteString(sql)

	for i, room := range roomList {
		if i == len(roomList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s');", room.Id, room.Name))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s'),", room.Id, room.Name))
		}
	}
	l4g.Debug(buffer.String())
	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("WriteRoom err:%v", result.Error)
	}
	return
}

func WriteMsg(msgList []*Message) {
	var buffer bytes.Buffer
	sql := "insert into `message` (`id`,`text`,`time_stamp`,`room_id`) values"
	buffer.WriteString(sql)

	for i, msg := range msgList {
		if i == len(msgList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s');", msg.Id, msg.Text, msg.TimeStamp, msg.RoomId))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s'),", msg.Id, msg.Text, msg.TimeStamp, msg.RoomId))
		}
	}

	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("WriteMsg err:%v", result.Error)
	}
	return
}

func WriteUserToken(utList []*Token) {
	var buffer bytes.Buffer
	sql := "insert into `token` (`username`,`token`) values"
	buffer.WriteString(sql)

	for i, ut := range utList {
		if i == len(utList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s');", ut.Username, ut.Token))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s'),", ut.Username, ut.Token))
		}
	}

	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("WriteUserToken err:%v", result.Error)
	}
	return
}

func UserEnterRoom(urMap map[string]string) {
	var buffer bytes.Buffer
	sql := "update `user` set `room_id` = CASE `username` "
	buffer.WriteString(sql)

	for username, roomId := range urMap {
		buffer.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", username, roomId))
	}

	buffer.WriteString(`END;`)

	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("UserEnterRoom err:%v", result.Error)
	}

	return
}

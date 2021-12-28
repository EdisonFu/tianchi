package db

import (
	"bytes"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/panjf2000/ants/v2"
	"strings"
	. "tianchi/models"
	"time"
)

const (
	MaxChanLen        = 50000
	WriteIntervalNum  = 100
	WriteIntervalTime = 500 //毫秒
	PoolWorkerNum     = 30000

	User_Data_Type     = 1
	Room_Data_Type     = 2
	Msg_Data_Type      = 3
	Token_Data_Type    = 4
	UserRoom_Data_Type = 5
)

var WriteUserChan chan *User
var WriteRoomChan chan *Room
var WriteMsgChan chan *Message
var WriteUserTokenChan chan *Token
var UserEnterRoomChan chan string

type saveWorker struct {
	DataType int
	Data     interface{}
}

func initWorkChan() {
	pool, _ := ants.NewPoolWithFunc(PoolWorkerNum, func(i interface{}) {
		worker := i.(*saveWorker)
		switch worker.DataType {
		case User_Data_Type:
			user := worker.Data.(*User)
			WriteUser(user)

		case Room_Data_Type:
			rooms := worker.Data.([]*Room)
			WriteRoom(rooms)

		case Msg_Data_Type:
			msgs := worker.Data.([]*Message)
			WriteMsg(msgs)

		case Token_Data_Type:
			tokens := worker.Data.([]*Token)
			WriteUserToken(tokens)

		case UserRoom_Data_Type:
			userRoom := worker.Data.(map[string]string)
			UserEnterRoom(userRoom)
		}
	})

	WriteUserChan = make(chan *User, MaxChanLen)
	WriteRoomChan = make(chan *Room, MaxChanLen)
	WriteMsgChan = make(chan *Message, MaxChanLen)
	WriteUserTokenChan = make(chan *Token, MaxChanLen)
	UserEnterRoomChan = make(chan string, MaxChanLen)

	//使用一个定时器通知多个协程
	msgNotifyChan := make(chan bool, 1)
	tokenNotifyChan := make(chan bool, 1)
	userRoomNotifyChan := make(chan bool, 1)
	timer := time.NewTicker(WriteIntervalTime * time.Millisecond)

	//批量写数据库逻辑，不适用于用户，否则创建完用户立马操作进出房间会出现脏数据
	go func() {
		for {
			select {
			case user, ok := <-WriteUserChan:
				if ok {
					pool.Invoke(&saveWorker{
						DataType: User_Data_Type,
						Data:     user,
					})
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
						roomList = roomList[:0]
						pool.Invoke(&saveWorker{
							DataType: Room_Data_Type,
							Data:     writeList,
						})
					} else {
						roomList = append(roomList, room)
						continue
					}
				}

			case <-timer.C:
				msgNotifyChan <- true
				if len(roomList) > 0 {
					writeList := roomList
					roomList = roomList[:0]
					WriteRoom(writeList)
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
						msgList = msgList[:0]
						pool.Invoke(&saveWorker{
							DataType: Msg_Data_Type,
							Data:     writeList,
						})
					} else {
						msgList = append(msgList, msg)
						continue
					}
				}

			case <-msgNotifyChan:
				tokenNotifyChan <- true
				if len(msgList) > 0 {
					writeList := msgList
					msgList = msgList[:0]
					WriteMsg(writeList)
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
						tokenList = tokenList[:0]
						pool.Invoke(&saveWorker{
							DataType: Token_Data_Type,
							Data:     writeList,
						})
					} else {
						tokenList = append(tokenList, userToken)
						continue
					}
				}

			case <-tokenNotifyChan:
				userRoomNotifyChan <- true
				if len(tokenList) > 0 {
					writeList := tokenList
					tokenList = tokenList[:0]
					WriteUserToken(writeList)
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

					if len(userRoomMap) > WriteIntervalNum {
						userRoomMap[ur[0]] = ur[1]
						writeMap := userRoomMap
						userRoomMap = make(map[string]string)
						pool.Invoke(&saveWorker{
							DataType: UserRoom_Data_Type,
							Data:     writeMap,
						})
					} else {
						userRoomMap[ur[0]] = ur[1]
						continue
					}
				}

			case <-userRoomNotifyChan:
				if len(userRoomMap) > 0 {
					writeMap := userRoomMap
					userRoomMap = make(map[string]string)
					UserEnterRoom(writeMap)
				}
			}
		}
	}()
}

func WriteUser(user *User) {
	var buffer bytes.Buffer
	sql := "insert into `user` (`username`,`first_name`,`last_name`,`email`,`password`,`phone`) values"
	buffer.WriteString(sql)
	buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s');", user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone))

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

	var nameList string
	for username, roomId := range urMap {
		buffer.WriteString(fmt.Sprintf("WHEN '%s' THEN '%s' ", username, roomId))
		nameList = nameList + "'" + username + "',"
	}

	buffer.WriteString(fmt.Sprintf("END where username in (%s);", nameList[:len(nameList)-1]))
	result := db.Exec(buffer.String())
	if result.Error != nil {
		l4g.Error("UserEnterRoom err:%v", result.Error)
	}

	return
}

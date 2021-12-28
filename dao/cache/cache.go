package cache

import (
	l4g "github.com/alecthomas/log4go"
	"sync"
	db "tianchi/dao/mysql"
	"tianchi/models"
)

var UserRoomMap sync.Map    //userId:roomId
var UserInfoMap sync.Map    //userName:user
var RoomHelperMap sync.Map  //roomId:roomhelper
var RoomList []*models.Room //房间列表，按房间排序
var UserToken sync.Map

func ReloadCacheFromDB() {
	rooms := db.ReadAllRoom()
	users := db.ReadAllUser()
	msgs := db.ReadAllMessage()
	tokens := db.ReadAllToken()

	for _, value := range rooms {
		room := value
		helper := InitRoomHelper(&room)
		RoomHelperMap.Store(room.Id, helper)

		RoomList = append(RoomList, &room)
	}

	for _, token := range tokens {
		UserToken.Store(token.Token, token.Username)
	}

	for _, value := range users {
		user := value
		if user.RoomId != "null" {
			UserRoomMap.Store(user.Username, user.RoomId)
		}

		UserInfoMap.Store(user.Username, &user)

		helper, ok := RoomHelperMap.Load(user.RoomId)
		if ok && helper != nil {
			helper.(*RoomHelper).EnterRoom(user.Username)
		}
	}

	for _, value := range msgs {
		msg := value
		helper, ok := RoomHelperMap.Load(msg.RoomId)
		if ok && helper != nil {
			helper.(*RoomHelper).SendMsg(&msg)
		}
	}

	l4g.Info("ReloadCacheFromDB OK! room num:%d, user num:%d, msg num:%d", len(rooms), len(users), len(msgs))
}

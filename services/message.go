package services

import (
	"errors"
	"tianchi/dao/cache"
	db "tianchi/dao/mysql"
	"tianchi/models"
)

//发送消息
func SendMessage(username string, msg *models.Message) (err error) {
	if msg == nil {
		return errors.New("message is nil")
	}

	roomId, ok := cache.UserRoomMap.Load(username)
	if !ok {
		return errors.New("user not in room!")
	}

	roomHelper, ok := cache.RoomHelperMap.Load(roomId)
	if ok && roomHelper != nil {
		roomHelper.(*cache.RoomHelper).SendMsg(msg)
	}

	msg.RoomId = roomId.(string)
	db.WriteMsgChan <- msg
	return
}

//如果所取消息在cache中命中，则从cache中取
//否则从数据库取
func RetrieveMessage(username string, index, size int32) (messageList []*models.Message, err error) {
	var list []*models.Message

	roomId, ok := cache.UserRoomMap.Load(username)
	if !ok {
		return nil, errors.New("user not in room!")
	}

	//数据查询超出内存范围，查数据库
	if -(index * size) > cache.Max_Ring_Len {
		dbList := db.ReadMessage(roomId.(string), index, size)
		for i := 0; i < len(dbList); i++ {
			messageList = append(messageList, &dbList[i])
		}
		return
	}

	//查询内存数据
	roomHelper, ok := cache.RoomHelperMap.Load(roomId)
	if ok && roomHelper != nil {
		list = roomHelper.(*cache.RoomHelper).RetrieveMsg(cache.Max_Ring_Len)
	}

	l := int32(len(list))

	startIndex := index * size
	endIndex := (index + 1) * size

	if l <= -endIndex {
		return
	}

	s := l + startIndex
	if l+startIndex < 0 {
		s = 0
	}

	messageList = list[s : l+endIndex]
	return
}

package cache

import (
	"sync"
	"tianchi/models"
)

const Max_Ring_Len = 2000

//一个helper负责一个room,维护room内的用户和message
type RoomHelper struct {
	Id      string
	Name    string
	UserMap sync.Map
	MsgList *models.Ring
}

func InitRoomHelper(room *models.Room) *RoomHelper {
	helper := new(RoomHelper)
	helper.Id = room.Id
	helper.Name = room.Name
	helper.UserMap = sync.Map{}
	helper.MsgList = models.InitRing(Max_Ring_Len)

	return helper
}

func (helper *RoomHelper) EnterRoom(username string) {
	helper.UserMap.Store(username, true)
}

func (helper *RoomHelper) LeaveRoom(username string) {
	helper.UserMap.Store(username, false)
}

func (helper *RoomHelper) GetRoomUser() []string {
	userList := make([]string, 0)
	helper.UserMap.Range(func(key, value interface{}) bool {
		if value.(bool) {
			userList = append(userList, key.(string))
		}

		return true
	})

	return userList
}

func (helper *RoomHelper) SendMsg(msg *models.Message) {
	helper.MsgList.Insert(msg)
}

func (helper *RoomHelper) RetrieveMsg(num int) (msgList []*models.Message) {
	list := helper.MsgList.GetList(num)
	for _, data := range list {
		if data != nil {
			msg := data.(*models.Message)
			msgList = append(msgList, msg)
		}
	}

	return
}

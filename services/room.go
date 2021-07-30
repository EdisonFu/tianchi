package services

import (
	"errors"
	"strconv"
	"tianchi/dao/cache"
	db "tianchi/dao/mysql"
	"tianchi/models"
)

func CreateRoom(username string, room *models.Room) (id string, err error) {
	if room == nil {
		return "", errors.New("room is nil")
	}

	id = strconv.Itoa(len(cache.RoomList) + 1)
	room.Id = id
	helper := cache.InitRoomHelper(room)
	//helper.EnterRoom(username)
	cache.RoomHelperMap.Store(room.Id, helper)

	//cache.UserRoomMap.Store(username, room.Id)
	cache.RoomList = append(cache.RoomList, room)

	db.WriteRoomChan <- room
	//db.UserEnterRoomChan <- username + "|" + id

	return
}

func EnterRoom(username string, roomid string) (err error) {
	roomHelper, ok := cache.RoomHelperMap.Load(roomid)
	if ok && roomHelper != nil {
		roomHelper.(*cache.RoomHelper).EnterRoom(username)
	} else {
		return errors.New("room not exit!")
	}

	cache.UserRoomMap.Store(username, roomid)

	db.UserEnterRoomChan <- username + "|" + roomid

	return err
}

func LeaveRoom(username string) (err error) {
	roomId, ok := cache.UserRoomMap.Load(username)
	if !ok {
		return err
	}

	roomHelper, ok := cache.RoomHelperMap.Load(roomId)
	if ok && roomHelper != nil {
		roomHelper.(*cache.RoomHelper).LeaveRoom(username)
	} else {
		return errors.New("room not exit!")
	}

	db.UserEnterRoomChan <- username + "|" + "null"
	return err
}

func GetRoomInfo(roomId string) (name string, err error) {
	roomHelper, ok := cache.RoomHelperMap.Load(roomId)
	if ok && roomHelper != nil {
		return roomHelper.(*cache.RoomHelper).Name, nil
	} else {
		return "", errors.New("room not exit!")
	}
}

func GetUserList(roomId string) (userList []string, err error) {
	roomHelper, ok := cache.RoomHelperMap.Load(roomId)
	if ok && roomHelper != nil {
		return roomHelper.(*cache.RoomHelper).GetRoomUser(), nil
	} else {
		return nil, errors.New("room not exit!")
	}
}

func GetRoomList(index, size int32) (roomList []*models.Room, err error) {
	startIndex := index * size
	endIndex := (index + 1) * size
	l := int32(len(cache.RoomList))
	if l <= startIndex {
		return nil, errors.New("no request data!")
	}

	if l <= endIndex {
		roomList = cache.RoomList[startIndex:]
	} else {
		roomList = cache.RoomList[startIndex:endIndex]
	}

	return
}

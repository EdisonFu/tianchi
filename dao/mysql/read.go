package db

import (
	l4g "github.com/alecthomas/log4go"
	. "tianchi/models"
)

func ReadAllRoom() (rooms []Room) {
	result := db.Table("room").Find(&rooms)
	if result.Error != nil {
		l4g.Error("ReadAllRoom select err:%v", result.Error)
		return nil
	}

	return
}

func ReadAllToken() (tokens []Token) {
	result := db.Table("token").Find(&tokens)
	if result.Error != nil {
		l4g.Error("ReadAllToken select err:%v", result.Error)
		return nil
	}

	return
}

func ReadAllUser() (users []User) {
	result := db.Table("user").Find(&users)
	if result.Error != nil {
		l4g.Error("ReadAllUser select err:%v", result.Error)
		return nil
	}

	return
}

func ReadAllMessage() (msgs []Message) {
	result := db.Table("message").Find(&msgs)
	if result.Error != nil {
		l4g.Error("ReadAllMessage select err:%v", result.Error)
		return nil
	}

	return
}

func ReadMessage(roomid string, pageIndex, pageSize int32) (msgs []Message) {
	var count, limit int32
	db.Table("message").Where("room_id=?", roomid).Count(&count)
	if count <= 0 {
		return
	}

	if count+pageIndex*pageSize < -pageSize {
		return
	}

	if count+pageIndex*pageSize < 0 {
		limit = ((pageIndex + 1) * pageSize) + count
	} else {
		limit = pageSize
	}
	result := db.Table("message").Where("room_id=?", roomid).Offset(count + pageIndex*pageSize).Limit(limit).Find(&msgs)
	if result.Error != nil {
		l4g.Error("ReadMessage select err:%v", result.Error)
		return nil
	}

	return
}

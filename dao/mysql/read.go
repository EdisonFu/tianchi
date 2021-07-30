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
	result := db.Table("message").Where("roomId=?", roomid).Offset(pageIndex * pageSize).Limit(pageSize).Find(msgs)
	if result.Error != nil {
		l4g.Error("ReadMessage select err:%v", result.Error)
		return nil
	}

	return
}

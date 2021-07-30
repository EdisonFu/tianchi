package models

//user
type User struct {
	Username  string `json:"username" gorm:"primary_key"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	RoomId    string `json:"-" gorm:"index:room_idx"`
}

type UserResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Token struct {
	Username string `json:"username" gorm:"primary_key"`
	Token    string `json:"token"`
}

//room
type RoomControlData struct {
	PageIndex int32 `json:"pageIndex"`
	PageSize  int32 `json:"pageSize"`
}

type Room struct {
	Id   string `json:"id" gorm:"index:room_idx"`
	Name string `json:"name"`
}

//message
type MessageControlData struct {
	PageIndex int32 `json:"pageIndex"`
	PageSize  int32 `json:"pageSize"`
}

type Message struct {
	Id        string `json:"id" gorm:"index:msg_idx"`
	Text      string `json:"text"`
	TimeStamp string `json:"timestamp"`
	RoomId    string `json:"-" gorm:"index:msg_idx"`
}

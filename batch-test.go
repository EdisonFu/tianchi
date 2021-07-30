package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"tianchi/models"
	"time"
)

type Message struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RoomControlData struct {
	PageIndex int32 `json:"pageIndex"`
	PageSize  int32 `json:"pageSize"`
}

func main() {
	n := flag.Int("n", 100, "number")
	flag.Parse()

	////
	//createUser(*n)
	//token := loginUser(fmt.Sprintf("name%d", *n))
	//fmt.Println("token:", token)
	//fmt.Println(createRoom(*n, token))
	//time.Sleep(60 * time.Second)
	////

	var roomId string
	for i := 0; i < *n; i++ {
		go func(i int) {
			createUser(i)
			token := loginUser(fmt.Sprintf("name%d", i))
			if i%10 == 0 {
				roomId = createRoom(i, token)
			}
			getRoom(roomId)
			enterRoom(roomId, token)
			leaveRoom(token)
			enterRoom(roomId, token)
			getRoomUser(roomId)

			sendMsg(i, token)
			getMsg(token)
		}(i)

	}

	time.Sleep(20 * time.Second)

	for j := 0; j < *n; j++ {
		go func(j int) {
			getUser(fmt.Sprintf("name%d", j))
			getRoomList()
		}(j)

	}

	time.Sleep(100 * time.Second)
}

//---------------user---------------------
func createUser(id int) {
	Post("http://0.0.0.0:8080/user",
		"",
		User{
			Username:  fmt.Sprintf("name%d", id),
			FirstName: fmt.Sprintf("%d", id),
			LastName:  "",
			Email:     "",
			Password:  "123456",
			Phone:     "",
		})
}

func loginUser(username string) string {
	binData := Get(fmt.Sprintf("http://0.0.0.0:8080/userLogin?username=%s&password=123456", username),
		"")

	return string(binData)
}

func getUser(username string) string {
	binData := Get(fmt.Sprintf("http://0.0.0.0:8080/user/%s", username),
		"")

	return string(binData)
}

//---------------room---------------------
func createRoom(id int, token string) string {
	binData := Post("http://0.0.0.0:8080/room",
		token,
		Room{
			Name: fmt.Sprintf("room%d", id),
		})

	return string(binData)
}

func enterRoom(roomId string, token string) string {
	binData := Put(fmt.Sprintf("http://0.0.0.0:8080/room/%s/enter", roomId), token)

	return string(binData)
}

func leaveRoom(token string) string {
	binData := Put("http://0.0.0.0:8080/roomLeave",
		token)

	return string(binData)
}

func getRoom(roomId string) string {
	binData := Get(fmt.Sprintf("http://0.0.0.0:8080/room/%s", roomId),
		"")

	return string(binData)
}

func getRoomUser(roomId string) string {
	binData := Get(fmt.Sprintf("http://0.0.0.0:8080//room/%s/users", roomId),
		"")

	return string(binData)
}

func getRoomList() string {
	binData := Post("http://0.0.0.0:8080/roomList",
		"",
		RoomControlData{
			PageIndex: 0,
			PageSize:  100,
		})

	return string(binData)
}

//----------------------message--------------------
func sendMsg(id int, token string) string {
	binData := Post("http://0.0.0.0:8080/message/send",
		token,
		Message{
			Id:   fmt.Sprintf("%d", id),
			Text: fmt.Sprintf("hello,%d", id),
		})

	return string(binData)
}

func getMsg(token string) string {
	binData := Post("http://0.0.0.0:8080/message/retrieve",
		token,
		models.MessageControlData{
			PageIndex: -1,
			PageSize:  100,
		})

	return string(binData)
}

func Get(url string, auth string) []byte {
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", "Bearer "+auth)
	reqest.Header.Add("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil || response == nil {
		fmt.Printf("Do put url:%s err:%s\n", url, err.Error())
		return nil
	}
	defer response.Body.Close()

	buff := make([]byte, 200)
	n, _ := response.Body.Read(buff)
	return buff[:n]
}

func Put(url string, auth string) []byte {
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("PUT", url, nil)

	//增加header选项
	token := fmt.Sprintf("Bearer %s", auth)
	reqest.Header.Add("Authorization", token)
	reqest.Header.Add("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil || response == nil {
		fmt.Printf("Do put url:%s err:%s\n", url, err.Error())
		return nil
	}
	defer response.Body.Close()

	buff := make([]byte, 200)
	n, _ := response.Body.Read(buff)
	return buff[:n]
}

// 发送POST请求

func Post(url string, auth string, data interface{}) []byte {
	bin, err := json.Marshal(data)
	if err != nil {
		fmt.Println("post marshal req err:", err)
		return nil
	}

	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("POST", url, bytes.NewReader(bin))

	//增加header选项
	token := fmt.Sprintf("Bearer %s", auth)
	reqest.Header.Add("Authorization", token)
	reqest.Header.Add("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil || response == nil {
		fmt.Printf("Do POST url:%s err:%s\n", url, err.Error())
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("post response err:", response.StatusCode)
	}

	buff := make([]byte, 500)
	n, _ := response.Body.Read(buff)
	fmt.Println("post response data:", string(buff[:n]))
	return buff[:n]
}

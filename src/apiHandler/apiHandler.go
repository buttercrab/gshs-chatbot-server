package apiHandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ApiUserResponse struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Size    int        `json:"size"`
	List    []UserData `json:"list"`
}

type UserData struct {
	UserType string `json:"userType"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Etc      string `json:"etc"`
}

func IsLoggedIn(id string) (bool, *UserData) {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/getChatBotUser.do?user_key=" + id)
	defer res.Body.Close()

	var t ApiUserResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	if t.Code == "0000" {
		return true, &t.List[0]
	} else {
		return false, nil
	}
}

type ApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func LaptopApplyRequest(id string, user *UserData) error {
	goodsNo := ""
	switch user.Etc {
	case "1":
		goodsNo = "271"
	case "2":
		goodsNo = "272"
	case "3":
		goodsNo = "273"
	}
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/insertGoodsUse.do" +
		"?user_key=" + id +
		"&userType=" + user.UserType +
		"&userId=" + user.UserId +
		"&goodsNo=" + goodsNo +
		"&startDate=" + time.Now().Format("20060102") + "1900" +
		"&endDate=" + time.Now().Format("20060102") + "2100" +
		"&useTarget=" + "[%EC%B1%97%EB%B4%87]")

	var t ApiResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	if t.Code != "0000" {
		return errors.New(t.Code + " " + t.Message)
	}

	return nil
}

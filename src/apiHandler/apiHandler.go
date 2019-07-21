package apiHandler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type ApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApiSearchRequest struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Size    int           `json:"size"`
	List    []GoodsInform `json:"list"`
}

type GoodsInform struct {
	Site       string `json:"site"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	Accept     string `json:"accept"`
	Teacher2   string `json:"teacher2"`
	GoodsNo    int    `json:"goodsNo"`
	ManageNo   string `json:"manageNo"`
	Content    string `json:"content"`
	Teacher1   string `json:"teacher1"`
	Target     string `json:"target"`
	UserName   string `json:"userName"`
	GoodsUseNo int    `json:"goodsUseNo"`
	GoodsName  string `json:"goodsName"`
}

func GetLoginURL(id string) string {
	return "http://external.gs.hs.kr/external/regChatBot.do?user_key" + id
}

func GetUserData(id string) (*ApiResponse, *UserData) {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/getChatBotUser.do?user_key=" + id)
	defer res.Body.Close()

	var t ApiUserResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	return &ApiResponse{Code: t.Code, Message: t.Message}, &t.List[0]
}

func ExpireUser(id string) *ApiResponse {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/expireChatBotUser.do?user_key=" + id)
	defer res.Body.Close()

	var t ApiResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	return &t
}

func SearchGoodsUse(id string, goodsNo int) (*ApiResponse, *[]GoodsInform) {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/goodsUseList.do?user_key=" + id +
		"&goodsNo=" + strconv.Itoa(goodsNo))
	defer res.Body.Close()

	var t ApiSearchRequest
	_ = json.NewDecoder(res.Body).Decode(&t)

	return &ApiResponse{Code: t.Code, Message: t.Message}, &t.List
}

func RequestGoodsUse(id string, goodsNo int, startDate, endDate time.Time) *ApiResponse {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/insertGoodsUse.do?user_key=" + id +
		"&goodsNo=" + strconv.Itoa(goodsNo) +
		"&startDate=" + startDate.Format("200601021504") +
		"&endDate=" + endDate.Format("200601021504") +
		"&useTarget=" + "[%EC%B1%97%EB%B4%87]")
	defer res.Body.Close()

	var t ApiResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	return &t
}

func CancelGoodsUse(id string, goodsUseNo int) *ApiResponse {
	res, _ := http.Get("http://external.gs.hs.kr/external/chatbot/deleteGoodsUse.do?user_key=" + id +
		"&goodsUseNo=" + strconv.Itoa(goodsUseNo))
	defer res.Body.Close()

	var t ApiResponse
	_ = json.NewDecoder(res.Body).Decode(&t)

	return &t
}

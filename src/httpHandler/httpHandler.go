package httpHandler

import (
	"../apiHandler"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type chatBotJson struct {
	UserRequest userRequest `json:"userRequest"`
	Bot         bot         `json:"bot"`
	Action      action      `json:"action"`
}

type userRequest struct {
	Timezone  string `json:"timezone"`
	Utterance string `json:"utterance"`
	Lang      string `json:"lang"`
	User      user   `json:"user"`
}

type user struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Properties userProperties `json:"properties"`
}

type userProperties struct {
	PlusFriendUserKey string `json:"plusfriendUserKey"`
	AppUserId         string `json:"appUserId"`
}

type action struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	Params       map[string]string      `json:"params"`
	DetailParams map[string]detailParam `json:"detailParams"`
}

type detailParam struct {
	Origin    string `json:"origin"`
	Value     string `json:"value"`
	GroupName string `json:"groupName"`
}

type bot struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type chatBotResponse struct {
	Version  string        `json:"version"`
	Template skillTemplate `json:"template"`
}

type skillTemplate struct {
	Outputs []component `json:"outputs"`
}

type component struct {
	SimpleText simpleText `json:"simpleText"`
}

type simpleText struct {
	Text string `json:"text"`
}

func LaptopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", 400)
		return
	}

	var c chatBotJson
	if r.Body == nil {
		http.Error(w, "Please send a body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "JSON decoding error: "+err.Error(), 400)
		return
	}

	id := c.UserRequest.User.Properties.PlusFriendUserKey
	apiRes, user := apiHandler.GetUserData(id)

	var s []string
	goodsNo := 0
	switch user.Etc {
	case "1":
		goodsNo = 271
	case "2":
		goodsNo = 272
	case "3":
		goodsNo = 273
	}

	if apiRes.Code == "0000" {
		log.Println("/laptop name: " + user.UserName + ", id: " + user.UserId + ", key: " + id)
		apiRes, info := apiHandler.SearchGoodsUse(id, goodsNo)
		if apiRes.Code != "0000" {
			s = append(s, "오류가 발생했습니다. 다음 오류 메세제를 관리자에게 보여주세요.")
			s = append(s, apiRes.Message)
		} else {
			if len(*info) == 0 {
				start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 0, 0, 0, time.UTC)
				end := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 0, 0, 0, time.UTC)
				apiRes := apiHandler.RequestGoodsUse(id, goodsNo, start, end)
				if apiRes.Code != "0000" {
					s = append(s, "오류가 발생했습니다. 다음 오류 메세제를 관리자에게 보여주세요.")
					s = append(s, apiRes.Message)
				} else {
					s = append(s, "오늘 1차시에 노사실 신청을 완료하였습니다.")
					s = append(s, "승인이 나기 전까지는 취소를 할 수 있습니다.")
				}
			} else {
				s = append(s, "오늘 1차시에 이미 신청이 되어있습니다.")
			}
		}
	} else {
		s = append(s, "로그인이 필요합니다. 아래 링크를 눌러 로그인을 해주세요")
		s = append(s, apiHandler.GetLoginURL(id))
	}

	res := chatBotResponse{
		Version: "2.0",
		Template: skillTemplate{
			Outputs: ToComponent(s),
		},
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(res)
}

func ToComponent(s []string) []component {
	var comp []component
	for _, i := range s {
		comp = append(comp, component{
			SimpleText: simpleText{
				Text: i,
			},
		})
	}
	return comp
}

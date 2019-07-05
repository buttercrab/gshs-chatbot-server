package httpHandler

import (
	"../apiHandler"
	"encoding/json"
	"log"
	"net/http"
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
		http.Error(w, err.Error(), 400)
		return
	}

	id := c.UserRequest.User.Properties.PlusFriendUserKey

	log.Println(id)

	var s []string

	login, user := apiHandler.IsLoggedIn(id)

	if login {
		err := apiHandler.LaptopApplyRequest(id, user)
		if err != nil {
			s = append(s, "에러가 발생했습니다.")
			s = append(s, "관리자에게 아래 코드를 보여주세요.")
			s = append(s, err.Error())
		} else {
			s = append(s, "노사실 신청을 완료하였습니다.")
		}
	} else {
		s = append(s, "로그인이 필요합니다. 아래 링크를 눌러 로그인을 해주세요")
		s = append(s, "http://external.gs.hs.kr/external/regChatBot.do?user_key="+id)
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

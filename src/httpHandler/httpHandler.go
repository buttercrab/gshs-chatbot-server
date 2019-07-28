package httpHandler

import (
	"../apiHandler"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	SimpleText *simpleText `json:"simpleText,omitempty"`
	BasicCard  *basicCard  `json:"basicCard,omitempty"`
}

type simpleText struct {
	Text string `json:"text"`
}

type basicCard struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Buttons     []button `json:"buttons"`
}

type button struct {
	Label      string `json:"label"`
	Action     string `json:"action"`
	WebLinkUrl string `json:"webLinkUrl"`
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := c.Action.Params["place"]; !ok {
		http.Error(w, "No `place` param", 400)
		return
	}

	id := c.UserRequest.User.Properties.PlusFriendUserKey
	apiRes, user := apiHandler.GetUserData(id)
	place := c.Action.Params["place"]

	log.Println(c)

	var s []string

	if apiRes.Code == "0000" {
		log.Println("/request name: " + user.UserName + ", id: " + user.UserId + ", key: " + id)
		if place == "노사" {
			goodsNo := getLaptopNo(user.Etc)
			apiRes, info := apiHandler.SearchGoodsUse(id, user, goodsNo)
			if apiRes.Code != "0000" {
				s = append(s, "오류가 발생했습니다. 다음 오류 메세지를 관리자에게 보여주세요.")
				s = append(s, apiRes.Message)
			} else {
				if len(*info) == 0 {
					start, end := getTime("1차시")
					apiRes := apiHandler.RequestGoodsUse(id, user, goodsNo, start, end)
					if apiRes.Code != "0000" {
						s = append(s, "오류가 발생했습니다. 다음 오류 메세지를 관리자에게 보여주세요.")
						s = append(s, apiRes.Message)
					} else {
						s = append(s, "오늘 1차시에 노사실 신청을 완료하였습니다.")
						s = append(s, "승인이 나기 전까지는 취소를 할 수 있습니다.")
					}
				} else {
					s = append(s, "오늘 1차시에 이미 신청이 되어있습니다.")
				}
			}
		} else if place == "토학" {
			if no, ok := c.Action.Params["no"]; ok {
				no := getClubNo(no)
				var t string
				if val, ok := c.Action.Params["time"]; ok {
					t = val
				} else {
					t = "1차시"
				}
				start, end := getTime(t)
				apiRes := apiHandler.RequestGoodsUse(id, user, no, start, end)
				if apiRes.Code != "0000" {
					s = append(s, "오류가 발생했습니다. 다음 오류 메세지를 관리자에게 보여주세요.")
					s = append(s, apiRes.Message)
				} else {
					s = append(s, "오늘 "+t+" 토학실 신청을 완료하였습니다.")
					s = append(s, "승인이 나기 전까지는 취소를 할 수 있습니다.")
				}
			} else {

			}
		} else {
			http.Error(w, "Not existing place", 400)
			return
		}
	} else {
		res := chatBotResponse{
			Version: "2.0",
			Template: skillTemplate{
				Outputs: []component{
					{
						BasicCard: &basicCard{
							Title:       "로그인이 필요합니다.",
							Description: "아래 버튼을 눌러 로그인을 해주세요",
							Buttons: []button{
								{
									Label:      "로그인",
									Action:     "webLink",
									WebLinkUrl: apiHandler.GetLoginURL(id),
								},
							},
						},
					},
				},
			},
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(res)

		return
	}

	res := chatBotResponse{
		Version: "2.0",
		Template: skillTemplate{
			Outputs: toComponent(s),
		},
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(res)
}

func CancelHandler(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := c.Action.Params["place"]; !ok {
		http.Error(w, "No `place` param", 400)
		return
	}

	id := c.UserRequest.User.Properties.PlusFriendUserKey
	apiRes, user := apiHandler.GetUserData(id)
	place := c.Action.Params["place"]

	var s []string

	if apiRes.Code == "0000" {
		log.Println("/cancel name: " + user.UserName + ", id: " + user.UserId + ", key: " + id + ", place: " + place)

		if place == "노사" {
			apiRes, info := apiHandler.SearchGoodsUse(id, user, getLaptopNo(user.Etc))

			if apiRes.Code != "0000" {
				s = append(s, "오류가 발생했습니다. 다음 오류 메세지를 관리자에게 보여주세요.")
				s = append(s, apiRes.Message)
			} else {
				if len(*info) == 0 {
					s = append(s, "신청한 건이 없습니다.")
				} else {
					count := 0
					for _, v := range *info {
						start, _ := time.Parse("200601021504", v.StartDate)
						if start.Format("20060102") == time.Now().Format("20060102") && v.Accept == "N" {
							_ = apiHandler.CancelGoodsUse(id, user, v.GoodsUseNo)
							count++
						}
					}
					s = append(s, "오늘 신청된 승인되지 않은 "+strconv.Itoa(count)+"건을 취소하였습니다.")
				}
			}
		} else if place == "토학" {

		} else {
			http.Error(w, "Not existing place", 400)
			return
		}
	} else {
		res := chatBotResponse{
			Version: "2.0",
			Template: skillTemplate{
				Outputs: []component{
					{
						BasicCard: &basicCard{
							Title:       "로그인이 필요합니다.",
							Description: "아래 버튼을 눌러 로그인을 해주세요",
							Buttons: []button{
								{
									Label:      "로그인",
									Action:     "webLink",
									WebLinkUrl: apiHandler.GetLoginURL(id),
								},
							},
						},
					},
				},
			},
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(res)

		return
	}

	res := chatBotResponse{
		Version: "2.0",
		Template: skillTemplate{
			Outputs: toComponent(s),
		},
	}

	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(res)
}

func toComponent(s []string) []component {
	var comp []component
	for _, i := range s {
		comp = append(comp, component{
			SimpleText: &simpleText{
				Text: i,
			},
		})
	}
	return comp
}

func getLaptopNo(s string) int {
	if t, err := strconv.Atoi(s); 1 <= t && t <= 3 && err == nil {
		return t + 270
	}
	return -1
}

func getClubNo(s string) int {
	if t, err := strconv.Atoi(s); 101 <= t && t <= 113 && err == nil {
		return t + 25
	}
	return -1
}

func getTime(s string) (time.Time, time.Time) {
	switch s {
	case "1차시":
		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 19, 0, 0, 0, time.UTC),
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 0, 0, 0, time.UTC)
	case "2차시":
		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 21, 30, 0, 0, time.UTC),
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 24, 0, 0, 0, time.UTC)
	}
	return time.Now(), time.Now()
}

func getClubStatus(id string, user *apiHandler.UserData) (string, *[][]apiHandler.GoodsInform) {
	var res [][]apiHandler.GoodsInform
	for i := 126; i <= 138; i++ {
		apiRes, info := apiHandler.SearchGoodsUse(id, user, i)
		if apiRes.Code != "0000" {
			return apiRes.Message, nil
		}
		res = append(res, *info)
	}
	return "", &res
}

package xinge

import (
	"encoding/json"
	"time"
)

const (
	TYPE_NOTIFICATION        = 1
	TYPE_MESSAGE             = 2
	TYPE_APNS_NOTIFICATION   = 11
	TYPE_REMOTE_NOTIFICATION = 12
	DATETIMEFORMAT           = "2006-01-02 15:04:05"
)

type Message interface {
	IsValid() bool
	ToJSON() string
	GetType() int
	GetMultiPkg() int
	GetEnvironment() int
	GetLoopInterval() int
	GetLoopTimes() int
}

type MessageAndroid struct {
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	ExpireTime   int                    `json:"expire_time"`
	SendTime     string                 `json:"send_time"`
	AcceptTime   []TimeInterval         `json:"accept_time,omitempty"`
	Type         int                    `json:"message_type"`
	MultiPkg     int                    `json:"multi_pkg"`
	Style        *Style                 `json:"style,omitempty"`
	ClickAction  *ClickAction           `json:"action,omitempty"`
	Custom       map[string]interface{} `json:"custom_content,omitempty"`
	Raw          string                 `json:"raw"`
	LoopInterval int                    `json:"loopInterval"`
	LoopTimes    int                    `json:"loopTimes"`
}

func NewMessageAndroid() *MessageAndroid {
	return &MessageAndroid{
		Title:        "",
		Content:      "",
		SendTime:     time.Now().Format(DATETIMEFORMAT),
		AcceptTime:   nil,
		Type:         TYPE_NOTIFICATION,
		MultiPkg:     0,
		Raw:          "",
		LoopInterval: -1,
		LoopTimes:    -1,
		ClickAction:  NewClickAction(),
		Style:        NewStyle(0),
	}
}

func EasyMessageAndroid(title, content string) *MessageAndroid {
	msg := NewMessageAndroid()
	msg.Title = title
	msg.Content = content
	return msg
}

func (s *MessageAndroid) SetTitle(title string) {
	s.Title = title
}

func (s *MessageAndroid) SetContent(content string) {
	s.Content = content
}

func (s *MessageAndroid) SetCustom(custom map[string]interface{}) {
	s.Custom = custom
}

func (s *MessageAndroid) SetType(t int) {
	s.Type = t
}

func (s *MessageAndroid) SetAction(action *ClickAction) {
	s.ClickAction = action
}

func (s *MessageAndroid) SetStyle(style *Style) {
	s.Style = style
}

func (s *MessageAndroid) AddAcceptTime(acceptTime TimeInterval) {
	s.AcceptTime = append(s.AcceptTime, acceptTime)
}

func (s *MessageAndroid) SetMultiPkg(multiPkg int) {
	s.MultiPkg = multiPkg
}

func (s *MessageAndroid) GetType() int {
	return s.Type
}

func (s *MessageAndroid) GetMultiPkg() int {
	return s.MultiPkg
}

func (s *MessageAndroid) GetEnvironment() int {
	return 0
}

func (s *MessageAndroid) GetLoopInterval() int {
	return s.LoopInterval
}

func (s *MessageAndroid) GetLoopTimes() int {
	return s.LoopTimes
}

func (s *MessageAndroid) IsValid() bool {
	if s.Raw == "" {
		return true
	}

	if s.Type < TYPE_NOTIFICATION || s.Type > TYPE_MESSAGE {
		return false
	}

	if s.MultiPkg < 0 || s.MultiPkg > 1 {
		return false
	}

	if s.Type == TYPE_NOTIFICATION {
		if !s.Style.IsValid() {
			return false
		}

		if !s.ClickAction.IsValid() {
			return false
		}
	}

	if s.ExpireTime < 0 || s.ExpireTime > 3*24*60*60 {
		return false
	}

	_, err := time.Parse(DATETIMEFORMAT, s.SendTime)
	if err != nil {
		return false
	}

	if s.AcceptTime != nil {
		for _, v := range s.AcceptTime {
			if !v.IsValid() {
				return false
			}
		}
	}

	if s.LoopInterval > 0 && s.LoopTimes > 0 && ((s.LoopTimes-1)*s.LoopInterval+1) > 15 {
		return false
	}

	return true
}

func (s *MessageAndroid) ToJSON() string {
	if s.Raw != "" {
		return s.Raw
	}

	jsonObj := map[string]interface{}{}
	if s.Type == TYPE_NOTIFICATION {
		jsonObj["title"] = s.Title
		jsonObj["content"] = s.Content

		jsonObj["builder_id"] = s.Style.BuilderId
		jsonObj["ring"] = s.Style.Ring
		jsonObj["vibrate"] = s.Style.Vibrate
		jsonObj["clearable"] = s.Style.Clearable
		jsonObj["n_id"] = s.Style.NId
		jsonObj["ring_raw"] = s.Style.RingRaw
		jsonObj["lights"] = s.Style.Lights
		jsonObj["icon_type"] = s.Style.IconType
		jsonObj["icon_res"] = s.Style.IconRes
		jsonObj["style_id"] = s.Style.StyleId
		jsonObj["small_icon"] = s.Style.SmallIcon

		if s.ClickAction != nil {
			jsonObj["action"] = s.ClickAction
		}
	} else if s.Type == TYPE_MESSAGE {
		jsonObj["title"] = s.Title
		jsonObj["content"] = s.Content
	}

	if s.AcceptTime != nil {
		jsonObj["accept_time"] = s.AcceptTime
	}

	if s.Custom != nil {
		jsonObj["custom_content"] = s.Custom
	}

	byt, err := json.Marshal(jsonObj)
	if err != nil {
		return `{}`
	}

	return string(byt)
}

// iOS Message
type MessageIOS struct {
	ExpireTime   int                    `json:"expire_time"`
	SendTime     string                 `json:"send_time"`
	AcceptTime   []TimeInterval         `json:"accept_time,omitempty"`
	Type         int                    `json:"message_type"`
	Custom       map[string]interface{} `json:"custom,omitempty"`
	Raw          string                 `json:"raw,omitempty"`
	AlertStr     string                 `json:"alert,omitempty"`
	AlertJo      []string               `json:"alert,omitempty"`
	Badge        int                    `json:"badge"`
	Sound        string                 `json:"sound"`
	Category     string                 `json:"category"`
	LoopInterval int                    `json:"loop_interval"`
	LoopTimes    int                    `json:"loop_times"`
	Environment  int                    `json:"environment"`
}

func NewMessageIOS() *MessageIOS {
	return &MessageIOS{
		Type:         TYPE_APNS_NOTIFICATION,
		SendTime:     time.Now().Format(DATETIMEFORMAT),
		AcceptTime:   nil,
		Raw:          "",
		AlertStr:     "",
		AlertJo:      make([]string, 0),
		Badge:        1,
		Sound:        "beep.wav",
		Category:     "",
		LoopInterval: -1,
		LoopTimes:    -1,
		Environment:  IOSENV_DEV,
	}
}

func EasyMessageIOS(alert string, env int) *MessageIOS {
	msg := NewMessageIOS()
	msg.AlertStr = alert
	msg.Environment = env
	return msg
}

func (s *MessageIOS) SetAlert(alert string) {
	s.AlertStr = alert
}

func (s *MessageIOS) SetCustom(custom map[string]interface{}) {
	s.Custom = custom
}

func (s *MessageIOS) SetBadge(badge int) {
	s.Badge = badge
}

func (s *MessageIOS) SetType(t int) {
	s.Type = t
}

func (s *MessageIOS) SetEnvironment(env int) {
	s.Environment = env
}

func (s *MessageIOS) SetSound(sourd string) {
	s.Sound = sourd
}

func (s *MessageIOS) AddAcceptTime(acceptTime TimeInterval) {
	s.AcceptTime = append(s.AcceptTime, acceptTime)
}

func (s *MessageIOS) GetType() int {
	return s.Type
}

func (s *MessageIOS) GetMultiPkg() int {
	return 1
}

func (s *MessageIOS) GetEnvironment() int {
	return s.Environment
}

func (s *MessageIOS) GetLoopInterval() int {
	return s.LoopInterval
}

func (s *MessageIOS) GetLoopTimes() int {
	return s.LoopTimes
}

func (s *MessageIOS) IsValid() bool {
	if s.Raw == "" {
		return true
	}

	if s.Type < TYPE_APNS_NOTIFICATION || s.Type > TYPE_REMOTE_NOTIFICATION {
		return false
	}

	if s.ExpireTime < 0 || s.ExpireTime > 3*24*60*60 {
		return false
	}

	_, err := time.Parse(DATETIMEFORMAT, s.SendTime)
	if err != nil {
		return false
	}

	for _, v := range s.AcceptTime {
		if !v.IsValid() {
			return false
		}
	}

	if s.Type == TYPE_REMOTE_NOTIFICATION {
		return true
	}

	return s.AlertStr != "" || len(s.AlertJo) > 0
}

func (s *MessageIOS) ToJSON() string {
	if s.Raw != "" {
		return s.Raw
	}

	jsonObj := map[string]interface{}{}
	if s.Custom != nil {
		jsonObj["custom"] = s.Custom
	}

	if s.AcceptTime != nil {
		jsonObj["accept_time"] = s.AcceptTime
	}

	aps := map[string]interface{}{}
	if s.Type == TYPE_REMOTE_NOTIFICATION {
		aps["content-available"] = 1
	} else if s.Type == TYPE_APNS_NOTIFICATION {
		if len(s.AlertJo) > 0 {
			aps["alert"] = s.AlertJo
		} else {
			aps["alert"] = s.AlertStr
		}

		if s.Badge != 0 {
			aps["badge"] = s.Badge
		}

		if s.Sound != "" {
			aps["sound"] = s.Sound
		}

		if s.Category == "" {
			aps["category"] = s.Category
		}
	}
	jsonObj["aps"] = aps

	byt, err := json.Marshal(jsonObj)
	if err != nil {
		return `{}`
	}
	return string(byt)
}

package xinge

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	DEVICE_ALL                         int    = 0
	DEVICE_BROWSER                     int    = 1
	DEVICE_PC                          int    = 2
	DEVICE_WINPHONE                    int    = 5
	DEVICE_ANDROID                     int    = 3
	DEVICE_IOS                         int    = 4
	IOSENV_PROD                        int    = 1
	IOSENV_DEV                         int    = 2
	IOS_MIN_ID                         int64  = 2200000000
	RESTAPI_DOMAIN                     string = "http://openapi.xg.qq.com"
	HTTP_GET                           string = "GET"
	HTTP_POST                          string = "POST"
	CONTENT_TYPE_X_WWW_FORM_URLENCODED string = "application/x-www-form-urlencoded"
)

var (
	client                           *Client
	RESTAPI_PUSHSINGLEDEVICE         string = RESTAPI_DOMAIN + "/v2/push/single_device"
	RESTAPI_PUSHSINGLEACCOUNT        string = RESTAPI_DOMAIN + "/v2/push/single_account"
	RESTAPI_PUSHACCOUNTLIST          string = RESTAPI_DOMAIN + "/v2/push/account_list"
	RESTAPI_PUSHALLDEVICE            string = RESTAPI_DOMAIN + "/v2/push/all_device"
	RESTAPI_PUSHTAGS                 string = RESTAPI_DOMAIN + "/v2/push/tags_device"
	RESTAPI_QUERYPUSHSTATUS          string = RESTAPI_DOMAIN + "/v2/push/get_msg_status"
	RESTAPI_QUERYDEVICECOUNT         string = RESTAPI_DOMAIN + "/v2/application/get_app_device_num"
	RESTAPI_QUERYTAGS                string = RESTAPI_DOMAIN + "/v2/tags/query_app_tags"
	RESTAPI_CANCELTIMINGPUSH         string = RESTAPI_DOMAIN + "/v2/push/cancel_timing_task"
	RESTAPI_BATCHSETTAG              string = RESTAPI_DOMAIN + "/v2/tags/batch_set"
	RESTAPI_BATCHDELTAG              string = RESTAPI_DOMAIN + "/v2/tags/batch_del"
	RESTAPI_QUERYTOKENTAGS           string = RESTAPI_DOMAIN + "/v2/tags/query_token_tags"
	RESTAPI_QUERYTAGTOKENNUM         string = RESTAPI_DOMAIN + "/v2/tags/query_tag_token_num"
	RESTAPI_CREATEMULTIPUSH          string = RESTAPI_DOMAIN + "/v2/push/create_multipush"
	RESTAPI_PUSHACCOUNTLISTMULTIPLE  string = RESTAPI_DOMAIN + "/v2/push/account_list_multiple"
	RESTAPI_PUSHDEVICELISTMULTIPLE   string = RESTAPI_DOMAIN + "/v2/push/device_list_multiple"
	RESTAPI_QUERYINFOOFTOKEN         string = RESTAPI_DOMAIN + "/v2/application/get_app_token_info"
	RESTAPI_QUERYTOKENSOFACCOUNT     string = RESTAPI_DOMAIN + "/v2/application/get_app_account_tokens"
	RESTAPI_DELETETOKENOFACCOUNT     string = RESTAPI_DOMAIN + "/v2/application/del_app_account_tokens"
	RESTAPI_DELETEALLTOKENSOFACCOUNT string = RESTAPI_DOMAIN + "/v2/application/del_app_account_all_tokens"
)

// 信鸽 Client 结构体
type Client struct {
	accessId  int64
	secretKey string
}

// 实例化信鸽 Client 结构体，给 accessId, secretKey 赋值
func NewClient(accessId int64, secretKey string) *Client {
	if client == nil || accessId != client.accessId {
		client = &Client{accessId, secretKey}
	}
	return client
}

// 检验 Token 参数
func (c *Client) validateToken(token string) bool {
	if c.accessId >= IOS_MIN_ID {
		return len(token) == 64
	}
	return len(token) == 40 || len(token) == 64
}

// 检验设备类型
func (c *Client) validateMessageType(message Message) (deviceType int, err error) {
	if c == nil {
		return DEVICE_ALL, errors.New("xinge client nil!")
	}

	if c.accessId < IOS_MIN_ID {
		return DEVICE_ANDROID, nil
	} else if c.accessId >= IOS_MIN_ID && (message.GetEnvironment() == IOSENV_PROD || message.GetEnvironment() == IOSENV_DEV) {
		return DEVICE_IOS, nil
	}

	return DEVICE_ALL, errors.New("unknown message type!")
}

// 准备必要参数，调用信鸽的 Restful 接口， 正式发起 Push 推送（Push 推送专用函数）
func (c *Client) push(uri string, message Message, params map[string]interface{}) XgResponse {
	if _, err := c.validateMessageType(message); err != nil {
		return NewRespone(-1, err.Error())
	}

	// 消息类型：1：通知 2：透传消息。iOS平台请填0；默认1：通知
	params["message_type"] = message.GetType()
	//向iOS设备推送时必填，1表示推送生产环境；2表示推送开发环境。推送Android平台不填或填0
	params["environment"] = message.GetEnvironment()
	// 消息类型：1：通知 2：透传消息。iOS平台请填0；默认1：通知
	//0表示按注册时提供的包名分发消息；1表示按access id分发消息，所有以该access id成功注册推送的app均可收到消息。本字段对iOS平台无效
	params["multi_pkg"] = message.GetMultiPkg()

	params["expire_time"] = 600
	params["send_time"] = time.Now().Unix()

	return c.callRestful(uri, params)
}

//接收传入的必要参数， 调用信鸽的 Restful 接口，发起 POST 请求
func (c *Client) callRestful(uri string, params map[string]interface{}) XgResponse {
	params["access_id"] = c.accessId
	params["timestamp"] = time.Now().Unix()
	params["sign"] = generateSign(HTTP_POST, uri, c.secretKey, params)

	var buf bytes.Buffer
	for k, v := range params {
		buf.WriteString(fmt.Sprintf("%s=%v&", k, v))
	}

	// fmt.Println(strings.TrimRight(buf.String(), "&"))

	r, err := http.Post(uri, CONTENT_TYPE_X_WWW_FORM_URLENCODED, strings.NewReader(strings.TrimRight(buf.String(), "&")))
	if err != nil {
		return NewRespone(-1, "http post data err!")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewRespone(-1, "read response data err!")
	}

	// fmt.Println(string(body))

	var res XgResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return NewRespone(-1, fmt.Sprintf("json unmarshal fail:%s", err.Error()))
	}
	//fmt.Println(res)
	return res
}

// 信鸽响应结构体
type XgResponse struct {
	Code     int       `json:"ret_code"`
	Msg      string    `json:"err_msg,omitempty"`
	XgResult *XgResult `json:"result,omitempty"`
}

// 信鸽响应中的结果参数 结构体
type XgResult struct {
	PushId        int64          `json:"push_id,string,omitempty"`
	Tokens        []string       `json:"tokens,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	DeviceNum     int64          `json:"device_num,omitempty"`
	IsReg         int64          `json:"isReg,omitempty"`
	ConnTimestamp int64          `json:"connTimestamp,omitempty"`
	MsgsNum       int64          `json:"msgsNum,omitempty"`
	Total         int64          `json:"total,omitempty"`
	XgResultList  []XgResultList `json:"list,omitempty"`
}

// 信鸽响应中的结果列表参数 结构体
type XgResultList struct {
	PushId    string `json:"push_id,omitempty"`
	Status    int    `json:"status,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	Finished  int64  `json:"finished,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

func NewRespone(code int, msg string) XgResponse {
	return XgResponse{Code: code, Msg: msg}
}

func RespSuccess() XgResponse {
	return NewRespone(0, "")
}

// 生成签名
func generateSign(method, uri, secretKey string, params map[string]interface{}) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	buf.WriteString(strings.ToUpper(method) + u.Host + u.Path)

	keys := sortKey(params)
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf(`%s=%v`, k, params[k]))
	}
	buf.WriteString(secretKey)
	tmp := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(tmp[:])
}

//给参数的 key asc排序（升序）
func sortKey(p map[string]interface{}) []string {
	keys := make([]string, 0)
	for k, _ := range p {
		keys = append(keys, k)
	}
	list := sort.StringSlice(keys)
	sort.Sort(list)
	return []string(list)
}

// 初始化基本的参数
func initParams() map[string]interface{} {
	params := make(map[string]interface{})
	return params
}

// ==== 简易 api 接口 v1.1.4 引入 ====

// Android 简易 api 的接口

/**
 * Android 平台推送消息给单个设备
 */
func PushTokenAndroid(accessId int64, secretKey, title, content, deviceToken string) XgResponse {
	params := initParams()
	params["device_token"] = deviceToken
	message := EasyMessageAndroid(title, content)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHSINGLEDEVICE, message, params)
}

/**
 * Android 平台推送消息给单个账号
 */
func PushAccountAndroid(accessId int64, secretKey, title, content, account string) XgResponse {
	params := initParams()
	params["account"] = account
	message := EasyMessageAndroid(title, content)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHSINGLEACCOUNT, message, params)
}

/**
 * Android 平台推送消息给所有设备
 */
func PushAllAndroid(accessId int64, secretKey, title, content string) XgResponse {
	params := initParams()
	message := EasyMessageAndroid(title, content)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHALLDEVICE, message, params)
}

/**
 * Android 平台推送消息给标签选中设备
 */
func PushTagAndroid(accessId int64, secretKey, title, content, tag string) XgResponse {
	message := EasyMessageAndroid(title, content)
	tagList := []string{tag}
	c := NewClient(accessId, secretKey)
	return c.PushTags(tagList, "OR", message)
}

// iOS 简易 api 的接口

/**
 * iOS 平台推送消息给单个设备
 */
func PushTokenIOS(accessId int64, secretKey, content, deviceToken string, env int) XgResponse {
	params := initParams()
	params["device_token"] = deviceToken
	message := EasyMessageIOS(content, env)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHSINGLEDEVICE, message, params)
}

/**
 * iOS 平台推送消息给单个账号
 */
func PushAccountIOS(accessId int64, secretKey, content, account string, env int) XgResponse {
	params := initParams()
	params["account"] = account
	message := EasyMessageIOS(content, env)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHSINGLEACCOUNT, message, params)
}

/**
 * iOS 平台推送消息给所有设备
 */
func PushAllIOS(accessId int64, secretKey, content string, env int) XgResponse {
	params := initParams()
	message := EasyMessageIOS(content, env)
	params["message"] = message.ToJSON()
	c := NewClient(accessId, secretKey)
	return c.push(RESTAPI_PUSHALLDEVICE, message, params)
}

/**
 * iOS 平台推送消息给标签选中设备
 */
func PushTagIOS(accessId int64, secretKey, content, tag string, env int) XgResponse {
	message := EasyMessageIOS(content, env)
	c := NewClient(accessId, secretKey)
	tagList := []string{tag}
	return c.PushTags(tagList, "OR", message)
}

// ======================= 详细的api接口 =======================
/**
 * 推送给指定设备
 *
 * @param deviceToken 目标设备token
 * @param message 待推送的消息
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushSingleDevice(deviceToken string, message Message) XgResponse {
	params := initParams()
	params["device_token"] = deviceToken
	params["message"] = message.ToJSON()
	return c.push(RESTAPI_PUSHSINGLEDEVICE, message, params)
}

/**
 * 推送给指定账号
 *
 * @param deviceType 设备类型，请填0
 * @param account 目标账号
 * @param message 待推送的消息
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushSingleAccount(account string, message Message) XgResponse {
	params := initParams()
	params["account"] = account
	params["message"] = message.ToJSON()
	//`{"accept_time":[],"action":{"action_type":1,"browser":{},"aty_attr":{}},"builder_id":0,"clearable":1,"content":"测试信鸽推送 Android API","custom_content":null,"icon_res":"","icon_type":0,"lights":1,"n_id":0,"ring":0,"ring_raw":"","small_icon":"","style_id":1,"title":"哎菠菜","vibrate":1}`
	return c.push(RESTAPI_PUSHSINGLEACCOUNT, message, params)
}

/**
 * 推送给多个账号 <br/>
 * 如果目标账号数超过10000，建议改用{@link #pushAccountListMultiple}接口
 *
 * @param deviceType 设备类型，请填0
 * @param accountList 目标账号列表
 * @param message 待推送的消息
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushAccountList(accountList []string, message Message) XgResponse {
	params := initParams()
	account_list, err := json.Marshal(accountList)
	if err != nil {
		return NewRespone(-1, "json marshal fial!")
	}
	params["account_list"] = string(account_list)
	params["message"] = message.ToJSON()
	return c.push(RESTAPI_PUSHACCOUNTLIST, message, params)
}

/**
 * 推送给全量设备，限Android系统使用
 *
 * @param deviceType 请填0
 * @param message 待推送的消息
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushAllDevices(message Message) XgResponse {
	params := initParams()
	params["message"] = message.ToJSON()
	return c.push(RESTAPI_PUSHALLDEVICE, message, params)
}

/**
 * 推送给多个tags对应的设备
 *
 * @param deviceType 设备类型，请填0
 * @param tagList 指定推送的tag列表
 * @param tagOp 多个tag的运算关系，取值必须是下面之一： AND OR
 * @param message 待推送的消息
 * @param environment 推送的目标环境 必须是其中一种： {@link #IOSENV_PROD}生产环境 {@link #IOSENV_DEV}开发环境
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushTags(tagList []string, tagOp string, message Message) XgResponse {
	if !message.IsValid() || len(tagList) <= 0 || (tagOp != "AND" && tagOp != "OR") {
		return NewRespone(-1, "param invalid!")
	}

	params := initParams()
	tagListByt, err := json.Marshal(tagList)
	if err != nil {
		return NewRespone(-1, "json marshal fail!")
	}
	params["tags_list"] = string(tagListByt)
	params["tags_op"] = tagOp
	params["message"] = message.ToJSON()

	if message.GetLoopInterval() > 0 && message.GetLoopTimes() > 0 {
		params["loop_interval"] = message.GetLoopInterval()
		params["loop_times"] = message.GetLoopTimes()
	}

	return c.push(RESTAPI_PUSHTAGS, message, params)
}

/**
 * 创建大批量推送消息，后续可调用{@link #pushAccountListMultiple}或{@link #pushDeviceListMultiple}接口批量添加设备<br/>
 * 此接口创建的任务不支持定时推送
 *
 * @param message 待推送的消息
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) CreateMultipush(message Message) int64 {
	if !message.IsValid() {
		return 0
	}

	params := initParams()
	params["message"] = message.ToJSON()

	res := c.push(RESTAPI_CREATEMULTIPUSH, message, params)
	if res.XgResult == nil {
		return 0
	}

	if res.XgResult.PushId <= 0 {
		return 0
	}

	return res.XgResult.PushId
}

/**
 * 推送消息给大批量账号，可对同一个pushId多次调用此接口 <br/>
 * 建议用户采用此接口自行控制发送时间
 *
 * @param pushId {@link #createMultipush}返回的push_id
 * @param accountList 账号列表，数量最多为1000个
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushAccountListMultiple(pushId int64, accountList []string) XgResponse {
	if pushId <= 0 {
		return NewRespone(-1, "pushId invalid!")
	}

	if len(accountList) <= 0 {
		return NewRespone(-1, "param invalid!")
	}

	params := initParams()
	params["push_id"] = pushId
	accountListByt, err := json.Marshal(accountList)
	if err != nil {
		return NewRespone(-1, "json marshal fial!")
	}
	params["account_list"] = string(accountListByt)

	return c.callRestful(RESTAPI_PUSHACCOUNTLISTMULTIPLE, params)
}

/**
 * 推送消息给大批量设备，可对同一个pushId多次调用此接口 <br/>
 * 建议用户采用此接口自行控制发送时间
 *
 * @param pushId {@link #createMultipush}返回的push_id
 * @param deviceList 设备列表，数量最多为1000个
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) PushDeviceListMultiple(pushId int64, deviceList []string) XgResponse {
	if pushId <= 0 {
		return NewRespone(-1, "pushId invalid!")
	}

	if len(deviceList) <= 0 {
		return NewRespone(-1, "param invalid!")
	}

	params := initParams()
	params["push_id"] = pushId
	deviceListByt, err := json.Marshal(deviceList)
	if err != nil {
		return NewRespone(-1, "json marshal fial!")
	}
	params["device_list"] = string(deviceListByt)

	return c.callRestful(RESTAPI_PUSHDEVICELISTMULTIPLE, params)
}

/**
 * 查询群发消息的状态，可同时查询多个pushId状态
 *
 * @param pushIdList 各类推送任务返回的push_id，可以一次查询多个
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryPushStatus(pushIdList []string) XgResponse {
	l := len(pushIdList)
	if l == 0 {
		return NewRespone(-1, "pushId slice len eq zero!")
	}

	params := initParams()
	buf := bytes.NewBufferString("[")
	for i := 0; i < l; i++ {
		if i == l-1 {
			buf.WriteString(`{"push_id":"` + pushIdList[i] + `"}`)
		} else {
			buf.WriteString(`{"push_id":"` + pushIdList[i] + `"},`)
		}
	}
	buf.WriteString(`]`)
	params["push_ids"] = buf.String()

	return c.callRestful(RESTAPI_QUERYPUSHSTATUS, params)
}

/**
 * 查询应用覆盖的设备数
 *
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryDeviceCount() XgResponse {
	params := initParams()
	return c.callRestful(RESTAPI_QUERYDEVICECOUNT, params)
}

/**
 * 查询应用当前所有的tags
 *
 * @param start 从哪个index开始
 * @param limit 限制结果数量，最多取多少个tag
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryTags(start, limit int64) XgResponse {
	params := initParams()
	params["start"] = start
	params["limit"] = limit
	return c.callRestful(RESTAPI_QUERYTAGS, params)
}

/**
 * 查询应用所有的tags，如果超过100个，取前100个
 *
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryTagsBefore100() XgResponse {
	return c.QueryTags(0, 100)
}

/**
 * 查询带有指定tag的设备数量
 *
 * @param tag 指定的标签
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryTagTokenNum(tag string) XgResponse {
	params := initParams()
	params["tag"] = tag
	return c.callRestful(RESTAPI_QUERYTAGTOKENNUM, params)
}

/**
 * 查询设备下所有的tag
 *
 * @param deviceToken 目标设备token
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryTokenTags(device_token string) XgResponse {
	params := initParams()
	params["device_token"] = device_token
	return c.callRestful(RESTAPI_QUERYTOKENTAGS, params)
}

/**
 * 取消尚未推送的定时任务
 *
 * @param pushId 各类推送任务返回的push_id
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) CancelTimingPush(pushId string) XgResponse {
	params := initParams()
	params["push_id"] = pushId
	return c.callRestful(RESTAPI_CANCELTIMINGPUSH, params)
}

/**
 * 批量为token设备标签，每次调用最多输入20个pair
 *
 * @param tagTokenPairs 指定token对应的指定tag
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) BatchSetTag(tagTokenPairs []TagTokenPair) XgResponse {
	l := len(tagTokenPairs)
	if l == 0 {
		return NewRespone(-1, "invalid TagTokenPair length")
	}

	buf := bytes.NewBufferString(`[`)
	for i := 0; i < l; i++ {
		if !c.validateToken(tagTokenPairs[i].Token) {
			return NewRespone(-1, "invalid token "+tagTokenPairs[i].Token)
		}
		if i == l-1 {
			buf.WriteString(`["` + tagTokenPairs[i].Tag + `","` + tagTokenPairs[i].Token + `"]`)
		} else {
			buf.WriteString(`["` + tagTokenPairs[i].Tag + `","` + tagTokenPairs[i].Token + `"],`)
		}
	}
	buf.WriteString(`]`)

	params := initParams()
	params["tag_token_list"] = buf.String()
	return c.callRestful(RESTAPI_BATCHSETTAG, params)
}

/**
 * 批量为token删除标签，每次调用最多输入20个pair
 *
 * @param tagTokenPairs 指定token对应的指定tag
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) BatchDelTag(tagTokenPairs []TagTokenPair) XgResponse {
	l := len(tagTokenPairs)
	if l == 0 {
		return NewRespone(-1, "invalid TagTokenPair length")
	}

	buf := bytes.NewBufferString(`[`)
	for i := 0; i < l; i++ {
		if !c.validateToken(tagTokenPairs[i].Token) {
			return NewRespone(-1, "invalid token "+tagTokenPairs[i].Token)
		}
		if i == l-1 {
			buf.WriteString(`["` + tagTokenPairs[i].Tag + `","` + tagTokenPairs[i].Token + `"]`)
		} else {
			buf.WriteString(`["` + tagTokenPairs[i].Tag + `","` + tagTokenPairs[i].Token + `"],`)
		}
	}
	buf.WriteString(`]`)

	params := initParams()
	params["tag_token_list"] = buf.String()
	return c.callRestful(RESTAPI_BATCHDELTAG, params)
}

/**
 * 查询token相关的信息，包括最近一次活跃时间，离线消息数等
 *
 * @param deviceToken 目标设备token
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryInfoOfToken(deviceToken string) XgResponse {
	params := initParams()
	params["device_token"] = deviceToken
	return c.callRestful(RESTAPI_QUERYINFOOFTOKEN, params)
}

/**
 * 查询账号绑定的token
 *
 * @param account 目标账号
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) QueryTokensOfAccount(account string) XgResponse {
	params := initParams()
	params["account"] = account
	return c.callRestful(RESTAPI_QUERYTOKENSOFACCOUNT, params)
}

/**
 * 删除指定账号和token的绑定关系（token仍然有效）
 *
 * @param account 目标账号
 * @param deviceToken 目标设备token
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) DeleteTokenOfAccount(account, deviceToken string) XgResponse {
	params := initParams()
	params["account"] = account
	params["device_token"] = deviceToken
	return c.callRestful(RESTAPI_DELETETOKENOFACCOUNT, params)
}

/**
 * 删除指定账号绑定的所有token（token仍然有效）
 *
 * @param account 目标账号
 * @return 服务器执行结果， XgResponse 实体
 */
func (c *Client) DeleteAllTokensOfAccount(account string) XgResponse {
	params := initParams()
	params["account"] = account
	return c.callRestful(RESTAPI_DELETEALLTOKENSOFACCOUNT, params)
}

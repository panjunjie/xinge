xinge 是根据官方腾讯信鸽推送 Reset API 接口开发出来的 SDK Go 版本的实现，基本和官方的各大 SDK 接口同步。后期如果官方有改动，本 SDK 也将会尽快跟进！

### 使用指南

安装：
go get github.com/panjunjie/xinge

使用例子：
```go
import(
	"github.com/panjunjie/xinge"
)

func main(){
    accessId := 2100259827 //自己的信鸽 accessId
    secretKey := "c2bb3a21f49715748a3a240dd17e6bd4" //自己的信鸽 secretKey
    
    //快捷接口给苹果用户发推送消息
    env := 2  // 1：正式环境；2：测试环境
    xinge.PushTokenIOS(accessId, secretKey, "给用户发信鸽 iOS 消息", "token", env) //给指定 Token 用户发消息
    xinge.PushAccountIOS(accessId, secretKey, "给用户发信鸽 iOS 消息", "accountId", env) //给指定用户账号发消息
    xinge.PushAllIOS(accessId, secretKey, "给用户发信鸽 iOS 消息", env) //给全部用户发消息
    xinge.PushTagIOS(accessId, secretKey, "给用户发信鸽 iOS 消息", "tag", env) //给指定 Tag 用户发消息
    
    //快捷接口给安卓用户发推送消息
    xinge.PushTokenAndroid(accessId, secretKey, "给用户发信鸽 Android 消息", "token") //给指定 Token 用户发消息
    xinge.PushAccountAndroid(accessId, secretKey, "给用户发信鸽 Android 消息", "accountId") //给指定用户账号发消息
    xinge.PushAllAndroid(accessId, secretKey, "给用户发信鸽 Android 消息") //给全部用户发消息
    xinge.PushTagAndroid(accessId, secretKey, "给用户发信鸽 Android 消息", "tag") //给指定 Tag 用户发消息
    
    
    //高级接口使用
    clientXG := xinge.NewClient(accessId, secretKey)
    
    messageIOS := xinge.EasyMessageIOS("给用户发信鸽 iOS 消息", env) // iOS 简单消息体实例化
    //定义自定义参数
    custom := map[string]interface{}{}
    custom["customTxt"] = 1
    //给消息体设置更多参数
    messageIOS.SetCustom(custom)
    messageIOS.SetBadge(100)
    messageIOS.SetSound("ring.ogg")
    messageIOS.AddAcceptTime(...)
    
    messageAndroid := xinge.EasyMessageAndroid("推送的标题","给用户发信鸽 Android 消息") // Android 简单消息体实例化
    //给消息体设置更多参数
    messageAndroid.SetCustom(custom)
    
    clientXG.PushSingleAccount("accountId",messageIOS) //给指定 iOS 用户账号发消息
    clientXG.PushSingleAccount("accountId",messageAndroid) //给指定 Android 用户账号发消息
    ...
}
```

### 已实现的接口
本 SDK 接口大致，分快捷接口和高级接口。快捷接口使用一行代码完成推送操作，注：快捷方式只支持推送通知，不支持透传消息。 快捷方式不支持定时推送。
高级接口可以自定义更细的消息参数，实现自定义消息，定时发送，透传等等高级场景

#### 快捷接口：
1. Android 平台推送消息给单个设备 PushTokenAndroid
2. Android 平台推送消息给单个账号 PushAccountAndroid
3. Android 平台推送消息给所有设备 PushAllAndroid
4. Android 平台推送消息给标签选中设备 PushTagAndroid
5. IOS 平台推送消息给单个设备 PushTokenIOS
6. IOS 平台推送消息给单个账号 PushAccountIOS
7. IOS 平台推送消息给所有设备 PushAllIOS
8. IOS 平台推送消息给标签选中设备 PushTagIOS

#### 高级接口：
1. PushSingleDevice 推送消息给单个设备
2. PushSingleAccount 推送消息给单个账号
3. PushAccountList 推送消息给多个账号
4. PushAllDevice 推送消息给单个 app 的所有设备
5. PushTags 推送消息给 tags 指定的设备
6. createMultipush 创建大批量推送消息
7. PushAccountListMultiple 推送消息给大批量账号(可多次) 
8. PushDeviceListMultiple 推送消息给大批量设备(可多次)
9. QueryPushStatus 查询群发消息发送状态
10. QueryDeviceCount 查询应用覆盖的设备数
11. QueryTags 查询应用的 tags.
12. CancelTimingPush 取消尚未推送的定时消息
13. BatchSetTag 批量为 token 设置标签
14. BatchDelTag 批量为 token 删除标签 
15. QueryTokenTags 查询 token 的 tags 
16. QueryTagTokenNum 查询 tag 下 token 的数目
17. QueryInfoOfToken 查询 token 的相关信息
18. QueryTokensOfAccount 查询 account 绑定的 token
19. DeleteTokenOfAccount 删除 account 绑定的 token
20. DeleteAllTokensOfAccount 删除 account 绑定的所有 token

### 消息体定义

消息体接口、Android 消息体、 iOS 消息体，
```go
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

    // iOS Message
    type MessageIOS struct {
        ExpireTime   int                    `json:"expire_time"`
        SendTime     string                 `json:"send_time"`
        AcceptTime   []TimeInterval         `json:"accept_time"`
        Type         int                    `json:"message_type"`
        Custom       map[string]interface{} `json:"custom,omitempty"`
        Raw          string                 `json:"raw"`
        AlertStr     string                 `json:"alert"`
        AlertJo      []string               `json:"alert"`
        Badge        int                    `json:"badge"`
        Sound        string                 `json:"sound"`
        Category     string                 `json:"category"`
        LoopInterval int                    `json:"loop_interval"`
        LoopTimes    int                    `json:"loop_times"`
        Environment  int                    `json:"environment"`
    }
```

我们提供简易的消息体实例化
EasyMessageIOS(alert)
EasyMessageAndroid(title,content)

需要了解消息体结构，才能配置更细的参数，调用高级接口很有帮助。


### 需要你的帮助
如果你在使用的过程中，发现任何可疑的 Bug，请不吝反馈，我会尽快检查修复，谢谢。
"# xinge" 

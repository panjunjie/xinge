/**
 * xinge 是根据官方腾讯信鸽推送 Reset API 接口开发出来的 SDK Go 版本的实现，基本和官方的各大 SDK 接口同步。后期如果官方有改动，本 SDK 也将会尽快跟进！
 *
 * 使用指南
 *
 *
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

 *
*/

package xinge

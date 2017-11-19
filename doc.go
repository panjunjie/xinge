/**
 * xinge 是根据官方腾讯信鸽推送的 Reset API 接口开发出来的 Go 版本实现，基本和官方的各大 SDK 接口同步。
 *
 * 使用指南
 *
 *
 func main() {
	var accessId int64 = 2100259000
	var secretKey string = "c2cb3a21f49715748a3a240dd17e6bd4"
	client * Client = NewClient(accessId, secretKey)

	//实例化 iOS 消息体
	var messageIOS *MessageIOS = EasyMessageIOS("潘军杰测试信鸽推送duo iOS API", xinge.IOSENV_PROD)

	//自定义参数
	custom := map[string]interface{}{}
	custom["business"] = 1
	custom["time"] = "2017-05-30 10:19:00"
	messageIOS.SetCustom(custom)

	//iOS 单个账号推送
	client.PushSingleAccount(xinge.DEVICE_IOS, "100048", messageIOS)
}

 *
*/

package xinge

// func main() {
// 	var accessId int64 = 2100259000
// 	var secretKey string = "c2cb3a21f49715748a3a240dd17e6bd4"
// 	client * Client = NewClient(accessId, secretKey)

// 	//实例化 iOS 消息体
// 	var messageIOS *MessageIOS = EasyMessageIOS("潘军杰测试信鸽推送duo iOS API", xinge.IOSENV_PROD)

// 	//自定义参数
// 	custom := map[string]interface{}{}
// 	custom["business"] = 1
// 	custom["time"] = "2017-05-30 10:19:00"
// 	messageIOS.SetCustom(custom)

// 	//iOS 单个账号推送
// 	client.PushSingleAccount(xinge.DEVICE_IOS, "100048", messageIOS)
// }

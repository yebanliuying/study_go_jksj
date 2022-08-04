package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"xshop_api/user-web/forms"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"xshop_api/user-web/global"
)

//GenerateSmsCode 生成短信验证码
func GenerateSmsCode(witdh int) string {
	num := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(num)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", num[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsCode := forms.SendSmsForm{}
	//绑定json数据
	if err := ctx.ShouldBind(&sendSmsCode); err != nil {
		HandleValidateError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	if err != nil {
		panic(err)
	}
	//mobile := "17600669474"
	smsCode := GenerateSmsCode(4)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsCode.Mobile                            //手机号
	request.QueryParams["SignName"] = "码脑科技"                                       //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_164509045"                          //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	//  fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
	//json数据解析

	//保存验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsCode.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":"发送成功",
	})

}

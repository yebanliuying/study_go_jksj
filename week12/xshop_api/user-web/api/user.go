package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"strings"
	"time"
	"xshop_api/user-web/middlewares"
	"xshop_api/user-web/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"xshop_api/user-web/forms"
	"xshop_api/user-web/global"
	"xshop_api/user-web/global/response"
	"xshop_api/user-web/proto"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

//HandleGrpcErrorToHttp grpc状态码转换http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": fmt.Sprintf("未知错误:%s", e.Message()), //todo 上线前关闭暴露错误
				})
			}
			return
		}
	}
}

//HandleValidateError 处理验证错误
func HandleValidateError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(), //todo 没有校验类型，返回错误没有翻译
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {


	//接收参数
	pnStr := ctx.DefaultQuery("pn", "1")
	pn, _ := strconv.Atoi(pnStr)
	pSizeStr := ctx.DefaultQuery("psize", "10")
	pSize, _ := strconv.Atoi(pSizeStr)

	rRsp, rErr := global.UserSrvCln.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(pSize),
	})
	if rErr != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】 失败")
		HandleGrpcErrorToHttp(rErr, ctx)
		return
	}

	result := make([]interface{}, 0)

	for _, value := range rRsp.Data {

		user := response.UserResponse{
			Id:       value.Id,
			Nickname: value.Nickname,
			Mobile:   value.Mobile,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)

}

func PasswordLogin(ctx *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}

	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidateError(ctx, err)
		return
	}

	//验证验证码
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//登录
	rsp, err := global.UserSrvCln.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		//检验密码
		rRsp, pRrr := global.UserSrvCln.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		})
		if pRrr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if rRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := model.CustomClaims{
					ID:          uint(rsp.Id),
					Nickname:    rsp.Nickname,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天国企
						//Issuer: global.ServerConfig.Name,
						Issuer: "test", //todo 后期设置成项目名称
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nickname":   rsp.Nickname,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //过期时间
				})

				//ctx.JSON(http.StatusOK, map[string]string{
				//	"msg" : "登录成功",
				//})
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败",
				})
			}

		}
	}
}

//Register 用户注册
func Register(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidateError(ctx, err)
		return
	}


	//验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	user, err := global.UserSrvCln.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户】失败：%s", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//ctx.JSON(http.StatusBadRequest)
	j := middlewares.NewJWT()
	claims := model.CustomClaims{
		ID:          uint(user.Id),
		Nickname:    user.Nickname,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天国企
			//Issuer: global.ServerConfig.Name,
			Issuer: "test", //todo 后期设置成项目名称
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nickname":   user.Nickname,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //过期时间
	})

}

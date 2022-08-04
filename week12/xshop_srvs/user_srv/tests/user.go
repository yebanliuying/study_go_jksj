package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"xshop_srvs/user_srv/proto"
)

var conn *grpc.ClientConn
var userClient proto.UserClient


func Init()  {
	var err error
	conn, err = grpc.Dial("0.0.0.0:9527", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	userClient = proto.NewUserClient(conn)
}

//TestGetUserList 测试-用户列表 && 密码验证
func TestGetUserList()  {
	rsp,err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn: 1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _,user := range rsp.Data{
		fmt.Println(user.Mobile,user.Nickname,user.Password)
		cRsp,cErr := userClient.CheckPassword(context.Background(),&proto.PasswordCheckInfo{
			Password: "123456",
			EncryptedPassword: user.Password,
		})
		if cErr != nil {
			panic(cErr)
		}
		fmt.Println(cRsp.Success)
	}

}

//TestCreateUser 测试-创建用户
func TestCreateUser()  {

	for i := 0; i < 10; i++ {
		rsp,err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Nickname: fmt.Sprintf("老王_%d", i),
			Password: "123456",
			Mobile: fmt.Sprintf("1522222222%d", i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func main()  {
	Init()
	defer conn.Close()

	//TestCreateUser()
	TestGetUserList()


}
package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"
	"xshop_srvs/user_srv/utils"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"xshop_srvs/user_srv/global"
	"xshop_srvs/user_srv/model"
	"xshop_srvs/user_srv/proto"
)

type UserServer struct {
	//向前兼容 或者使用 protoc --go-grpc_opt=require_unimplemented_servers=false
	proto.UnimplementedUserServer
}

//CheckPassword 校验密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {

	options := &utils.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := utils.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}

//ModelToResponse 用户信息返回结构 model 取值 to 返回类型
func ModelToResponse(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id: user.ID,
		Password: user.Password,
		Nickname:user.Nikename,
		Gender: uint32(user.Gender),
		Role: uint32(user.Role),
	}
	//处理默认值为nil的字段 避免直接放入grpc容易出错
	if user.Birthday != nil{
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

//ModelToResponseOfId 用户id返回结构
func ModelToResponseOfId(user model.User) proto.UserInfoResponse {
	userIdRsp := proto.UserInfoResponse{
		Id: user.ID,
	}
	return userIdRsp
}

//GetUserList 服务端-获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//取用户全部数据
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("用户列表")
	//声明rpc
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	//取用户分页数据
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	//赋值rpc
	for _,user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

//GetUserByMobile 服务端-获取用户信息by手机号
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

//GetUserById 服务端-获取用户信息byId
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

//CreateUser 服务端-新增用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists ,"用户已存在")
	}

	user.Mobile = req.Mobile
	user.Nikename = req.Nickname

	//密码加密
	options := &utils.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := utils.Encode(req.Password, options)
	//todo 后期优化支持更换密码的加密算法
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userIdRsp := ModelToResponseOfId(user)
	return &userIdRsp, nil
}

//UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result :=  global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthday := time.Unix(int64(req.Birthday), 0)
	user.Nikename = req.Nickname
	user.Birthday = &birthday
	user.Gender = int(req.Gender)

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}


//Paginate 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func (s *UserServer) mustEmbedUnimplementedUserServer() {}
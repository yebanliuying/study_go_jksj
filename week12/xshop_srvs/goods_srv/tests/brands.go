package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"xshop_srvs/goods_srv/proto"
)

var conn *grpc.ClientConn
var brandClient proto.GoodsClient


func TestGetBrandList()  {
	rsp,err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _,brand := range rsp.Data{
		fmt.Println(brand.Name)
	}

}


func Init()  {
	var err error
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	brandClient = proto.NewGoodsClient(conn)
}

func main()  {
	Init()
	defer conn.Close()

	//TestCreateUser()
	//TestGetBrandList()
	TestGetBrandList()


}
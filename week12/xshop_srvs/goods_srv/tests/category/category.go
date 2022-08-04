package main


import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"xshop_srvs/goods_srv/proto"
)


var conn *grpc.ClientConn
var brandClient proto.GoodsClient

func TestGetCategoryList()  {
	rsp,err := brandClient.GetAllCategorysList(context.Background(), &empty.Empty{
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _,category := range rsp.Data{
		fmt.Println(category.Name)
	}

}

func TestGetSubCategoryList()  {
	rsp,err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.SubCategorys)
	//Info         *CategoryInfoResponse   `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	//SubCategorys []*CategoryInfoResponse `protobuf:"bytes,3,rep,name=subCategorys,proto3" json:"subCategorys,omitempty"`

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
	TestGetSubCategoryList()


}

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"xshop_srvs/inventory_srv/proto"
)


var conn *grpc.ClientConn
var invClient proto.InventoryClient

func TestSetInv(goodsId, num int32)  {
	_,err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:    num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")

}

func TestInvDetail(goodsId int32)  {
	rsp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)

}

func TestSell (wg *sync.WaitGroup) {
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId:421, Num: 1},
			//{GoodsId:422, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestReback()  {
	_, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId:421, Num: 2},
			//{GoodsId:423, Num: 1},
			{GoodsId:422, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存归还成功")
}

func Init()  {
	var err error
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	invClient = proto.NewInventoryClient(conn)
}

func main()  {
	Init()
	defer conn.Close()

	//并发请求
	var wg sync.WaitGroup
	wg.Add(30)
	for i := 0; i<30;i++ {
		go TestSell(&wg)
	}

	wg.Wait()

	//for i := 421; i < 840; i++ {
	//	TestSetInv(int32(i),100)
	//}
	//TestSetInv(422,20)

	//TestInvDetail(421)
	//TestReback()

}

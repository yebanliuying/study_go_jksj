package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/bwmarrin/snowflake"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
	"xshop_srvs/order_srv/global"
	"xshop_srvs/order_srv/model"
	"xshop_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn() {
	//年月日时分秒+用户id+2位随机数

}


func (s OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rps proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rps.Total = int32(result.RowsAffected)
	}

	for _, value := range shopCarts {
		item := proto.ShopCartInfoResponse{
			Id:      value.ID,
			UserId:  value.User,
			GoodsId: value.Goods,
			Nums:    value.Nums,
			Checked: value.Checked,
		}
		rps.Data = append(rps.Data, &item)
	}

	return &rps, nil
}

func (s OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		//如果记录已经存在，则合并购物车记录, 更新操作
		shopCart.Nums += req.Nums
	} else {
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (s *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	//更新购物车记录，更新数量和选中状态
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (s *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	//分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (s *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	//web层应该先查询一下订单id是否是当前用户的
	//web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

type OrderListener struct {
	Code codes.Code
	Detail string
	ID int32
	OrderAmount float32
	Ctx context.Context

}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState{
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	parentSpan := opentracing.SpanFromContext(o.Ctx)


	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)

	//获取购物车商品
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	for _, value := range shopCarts {
		goodsIds = append(goodsIds, value.Goods)
		goodsNumsMap[value.Goods] = value.Nums
	}

	//跨服务调用商品微服务 批量查询商品
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()


	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	var orderAmount float32
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		//拼接订单商品数据 缺少订单号Order
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		//拼接扣减库存数据
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	//先扣库存 再通过事务控制是否需要归还。失败就归还，成功rollback
	//跨服务调用库存微服务 查询库存扣减
	queryInvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{OrderSn:orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	queryInvSpan.Finish()

	// 开始事务
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount

	//生成订单主表信息
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	//重新给订单商品数据加上订单号
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	//批量插入订单商品表信息
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_ordergoods", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	//删除购物车记录
	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		// 遇到错误时回滚事务
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	deleteShopCartSpan.Finish()

	//发送mq延时消息 订单超时失效
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"139.198.21.42:9876"}))

	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动Start失败")
	}

	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(5)//16
	_, err = p.SendSync(context.Background(), msg)
	//res, err := p.SendSync(context.Background(), primitive.NewMessage("topic1", []byte("this is message1")))
	if err != nil {
		zap.S().Errorf("发送延时消息失败：%v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	//if err = p.Shutdown(); err != nil {
	//	panic("关闭producer失败")
	//}



	// 否则，提交事务
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
	//return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState{
	//return primitive.CommitMessageState
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	//没查到创建的订单号 退还库存
	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	//否则撤掉归还库存的事务
	return primitive.RollbackMessageState
}


func (s *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	orderListener := OrderListener{Ctx:ctx}

	//连接mq
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{"139.198.21.42:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败:%s", err.Error())
		return nil, err
	}

	if err = p.Start(); err != nil {
		zap.S().Errorf("启动Start失败:%s", err.Error())
		return nil, err
	}

	//生成订单
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "生成订单id失败")
	}
	id := node.Generate().String()

	order := model.OrderInfo{
		OrderSn:      id,
		//OrderMount:   orderAmount,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}

	jsonString,_ := json.Marshal(order)

	//调用mq事务消息 实现两个方法
	_, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败:%s\n",err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}

	if orderListener.Code == codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: order.OrderSn, Total: orderListener.OrderAmount}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeout (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){

	for _, msg := range msgs {

		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msg.Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn:orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		
		//找到订单判断状态
		if order.Status != "TRADE_SUCCESS" {

			//不是成功的状态 就归还库存
			tx := global.DB.Begin()


			//修改订单状态为已关闭
			order.Status = "TRADE_CLOSE"
			tx.Save(&order)

			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"139.198.21.42:9876"}))

			if err != nil {
				panic("生成producer失败")
			}

			if err = p.Start(); err != nil {
				panic("启动Start失败")
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msg.Body))
			//res, err := p.SendSync(context.Background(), primitive.NewMessage("topic1", []byte("this is message1")))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n",err)
				return consumer.ConsumeRetryLater, nil

			}

			//if err = p.Shutdown(); err != nil {
			//	panic("关闭producer失败")
			//}

		}
	}
	return consumer.ConsumeSuccess, nil
}
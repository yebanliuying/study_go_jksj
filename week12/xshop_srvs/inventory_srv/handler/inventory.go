package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"xshop_srvs/inventory_srv/global"
	"xshop_srvs/inventory_srv/model"
	"xshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (s *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (s InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

//var m sync.Mutex //互斥锁
//Sell 扣减库存
func (s InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	// 开始事务
	tx := global.DB.Begin()
	sellDetail := model.StockSellDetail{
		OrderSn:req.OrderSn,
		Status: 1,
	}
	var details []model.GoodsDetail
	//m.Lock() //获取锁
	//遍历订单商品
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods:goodInfo.GoodsId,
			Num:goodInfo.Num,
		})
		var inv model.Inventory
		//查询库存 悲观锁
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() // 遇到错误时回滚事务
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}
		mutex := global.RS.NewMutex(fmt.Sprintf("redsync_goods_%d", goodInfo.GoodsId))
		//mutex := global.RS.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 遇到错误时回滚事务
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//当库存小于扣减数量
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 遇到错误时回滚事务
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//乐观锁
		//tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1})

	}
	sellDetail.Detail = details
	//下入sellDetail表
	if result := tx.Create(&sellDetail); result.RowsAffected == 0{
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}

	// 否则，提交事务
	tx.Commit()
	//m.Unlock() //释放锁
	return &emptypb.Empty{}, nil
}

//Reback 归还库存
func (s InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开始事务
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//查询库存
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 遇到错误时回滚事务
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		//加库存
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	// 否则，提交事务
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func AutoReback (ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
	type OrderInfo struct {
		OrderSn string
	}

	for i := range msgs {
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败: %v\n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		//将inv库存加回去，将sellDetail的status设置为2，要在事务中进行
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		//没有查询到该订单库存记录，直接扔掉消息
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn:orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		//查询到，逐个归还
		for _,orderGood := range sellDetail.Detail {
			//加库存失败 过段时间重新消费
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods:orderGood.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		//更新状态失败 过段时间重新消费
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn:orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
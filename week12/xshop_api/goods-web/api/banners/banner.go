package banners

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"xshop_api/goods-web/api"
	"xshop_api/goods-web/forms"
	"xshop_api/goods-web/global"
	"xshop_api/goods-web/proto"
)

func List(ctx *gin.Context) {
	rsp, err := global.GoodsSrvCln.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["index"] = value.Index
		reMap["image"] = value.Image
		reMap["url"] = value.Url

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func New(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvCln.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
		Image: bannerForm.Image,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id
	response["index"] = rsp.Index
	response["url"] = rsp.Url
	response["image"] = rsp.Image

	ctx.JSON(http.StatusOK, response)
}

func Update(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvCln.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvCln.DeleteBanner(context.Background(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "")
}

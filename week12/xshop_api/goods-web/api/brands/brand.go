package brands

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"xshop_api/goods-web/api"
	"xshop_api/goods-web/forms"
	"xshop_api/goods-web/global"
	"xshop_api/goods-web/proto"
)

func List(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.GoodsSrvCln.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pn:    int32(pnInt),
		PSize: int32(pSizeInt),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	for _, value := range rsp.Data[pnInt : pnInt*pSizeInt+pSizeInt] {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvCln.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	ctx.JSON(http.StatusOK, request)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvCln.DeleteBrand(context.Background(), &proto.BrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvCln.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvCln.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	ctx.JSON(http.StatusOK, result)
}

func CategoryBrandList(ctx *gin.Context) {
	rsp, err := global.GoodsSrvCln.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["category"] = map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		}
		reMap["brand"] = map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, reMap)
	}

	reMap["data"] = result
	ctx.JSON(http.StatusOK, reMap)
}

func NewCategoryBrand(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvCln.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	ctx.JSON(http.StatusOK, response)
}

func UpdateCategoryBrand(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvCln.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvCln.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "")
}

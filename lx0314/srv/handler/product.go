package handler

import (
	"context"
	"errors"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/model"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	//pb.UnimplementedGreeterServer
	__.UnimplementedServiceServer
}

// AddProducts implements helloworld.GreeterServer
func (s *Server) AddProducts(_ context.Context, in *__.AddProductsReq) (*__.AddProductsResp, error) {
	var product model.Product
	product = model.Product{
		Name:     in.Name,
		Subtitle: in.Subtitle,
		Desc:     in.Desc,
		Images:   in.Images,
		Price:    in.Price,
		Stock:    int(in.Stock),
		Status:   int(in.Status),
	}
	err := product.CreateProduct(nacos.DB)
	if err != nil {
		return nil, errors.New("创建商品失败" + err.Error())
	}
	return &__.AddProductsResp{
		Id: int64(product.ID),
	}, nil
}

func (s *Server) UpdateProducts(_ context.Context, in *__.UpdateProductsReq) (*__.UpdateProductsResp, error) {
	var product model.Product
	err := product.FindProductById(nacos.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	product = model.Product{
		Name:     in.Name,
		Subtitle: in.Subtitle,
		Desc:     in.Images,
		Images:   in.Images,
		Price:    in.Price,
		Stock:    int(in.Stock),
		Status:   int(in.Status),
	}
	err = product.UpdateProduct(nacos.DB, in.Id)
	if err != nil {
		return nil, errors.New("更新商品失败" + err.Error())
	}
	return &__.UpdateProductsResp{
		Id: int64(product.ID),
	}, nil
}

func (s *Server) DelProducts(_ context.Context, in *__.DelProductsReq) (*__.DelProductsResp, error) {
	var product model.Product
	err := product.FindProductById(nacos.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	err = product.DeleteProduct(nacos.DB)
	if err != nil {
		return nil, errors.New("删除商品失败" + err.Error())
	}
	return &__.DelProductsResp{
		Id: in.Id,
	}, nil
}

func (s *Server) GetProductsById(_ context.Context, in *__.GetProductsByIdReq) (*__.GetProductsByIdResp, error) {
	var product model.Product
	err := product.FindProductById(nacos.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	return &__.GetProductsByIdResp{
		Products: &__.Products{
			Id:       int64(product.ID),
			Name:     product.Name,
			Subtitle: product.Subtitle,
			Images:   product.Images,
			Price:    product.Price,
			Stock:    int64(product.Stock),
			Status:   int64(product.Status),
		},
	}, nil
}

func (s *Server) SearchProducts(_ context.Context, in *__.SearchProductsReq) (*__.SearchProductsResp, error) {
	var product model.Product
	list, err := product.SearchProduct(nacos.DB, in)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	return &__.SearchProductsResp{
		Products: list,
	}, nil
}

package handler

import (
	"context"
	"errors"
	"lx0318/config"
	__ "lx0318/proto"
	"lx0318/rpc/model"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	//pb.UnimplementedGreeterServer
	__.UnimplementedServiceServer
}

//func (s *Server) Register(_ context.Context, in *__.RegisterReq) (*__.RegisterResp, error) {
//	var user model.MemberRegister
//	err := user.FindUserByUsername(config.DB, in.Username)
//	if err != nil {
//		return nil, errors.New("用户查询失败")
//	}
//	if user.ID != 0 {
//		return nil, errors.New("用户已存在")
//	}
//	user = model.User{
//		Username: in.Username,
//		Password: in.Password,
//	}
//	err = user.CreateUser(config.DB)
//	if err != nil {
//		return nil, errors.New("用户注册失败")
//	}
//	return &__.RegisterResp{
//		Id: int64(user.ID),
//	}, nil
//}

// AddProducts implements helloworld.GreeterServer
func (s *Server) AddProducts(_ context.Context, in *__.AddProductsReq) (*__.AddProductsResp, error) {
	var product model.Product
	product = model.Product{
		Name:   in.Name,
		Price:  in.Price,
		Stock:  int(in.Stock),
		Status: int(in.Status),
	}
	err := product.CreateProduct(config.DB)
	if err != nil {
		return nil, errors.New("创建商品失败" + err.Error())
	}
	return &__.AddProductsResp{
		Id: int64(product.ID),
	}, nil
}

func (s *Server) UpdateProducts(_ context.Context, in *__.UpdateProductsReq) (*__.UpdateProductsResp, error) {
	var product model.Product
	err := product.FindProductById(config.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	product = model.Product{
		Name:   in.Name,
		Price:  in.Price,
		Stock:  int(in.Stock),
		Status: int(in.Status),
	}
	err = product.UpdateProduct(config.DB, in.Id)
	if err != nil {
		return nil, errors.New("更新商品失败" + err.Error())
	}
	return &__.UpdateProductsResp{
		Id: int64(product.ID),
	}, nil
}

func (s *Server) DelProducts(_ context.Context, in *__.DelProductsReq) (*__.DelProductsResp, error) {
	var product model.Product
	err := product.FindProductById(config.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	err = product.DeleteProduct(config.DB)
	if err != nil {
		return nil, errors.New("删除商品失败" + err.Error())
	}
	return &__.DelProductsResp{
		Id: in.Id,
	}, nil
}

func (s *Server) GetProductsById(_ context.Context, in *__.GetProductsByIdReq) (*__.GetProductsByIdResp, error) {
	var product model.Product
	err := product.FindProductById(config.DB, in.Id)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	if product.ID == 0 {
		return nil, errors.New("商品不存在")
	}
	return &__.GetProductsByIdResp{
		Products: &__.Products{
			Id:     int64(product.ID),
			Name:   product.Name,
			Price:  product.Price,
			Stock:  int64(product.Stock),
			Status: int64(product.Status),
		},
	}, nil
}

func (s *Server) SearchProducts(_ context.Context, in *__.SearchProductsReq) (*__.SearchProductsResp, error) {
	var product model.Product
	list, err := product.SearchProduct(config.DB, in)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	return &__.SearchProductsResp{
		Products: list,
	}, nil
}

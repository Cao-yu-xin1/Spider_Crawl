package handler

import (
	"context"
	"errors"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/nacos"
	__ "github.com/Cao-yu-xin1/Spider_Crawl/lx0314/proto"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/srv/model"
)

func (s *Server) Register(_ context.Context, in *__.RegisterReq) (*__.RegisterResp, error) {
	var user model.User
	err := user.FindUserByUsername(nacos.DB, in.Username)
	if err != nil {
		return nil, errors.New("用户查询失败")
	}
	if user.ID != 0 {
		return nil, errors.New("用户已存在")
	}
	user = model.User{
		Username: in.Username,
		Password: in.Password,
	}
	err = user.CreateUser(nacos.DB)
	if err != nil {
		return nil, errors.New("用户注册失败")
	}
	return &__.RegisterResp{
		Id: int64(user.ID),
	}, nil
}

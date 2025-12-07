// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/Nozomi9967/wmss-user-api/common"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户详情
func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.Response, err error) {
	userId := l.ctx.Value("user_id").(string)
	var user *types.UserInfo
	user, _ = l.svcCtx.SysUserModel.SelectOneDetail(l.ctx, userId)

	if user == nil {
		return &types.Response{
			Code: 200,
			Msg:  "查询失败，当前登录用户不存在",
			Data: nil,
		}, nil
	}

	if user.RoleID != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("查询用户失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户失败，权限不足",
		}, nil
	}

	user, _ = l.svcCtx.SysUserModel.SelectOneDetail(l.ctx, req.UserID)

	if user == nil {
		return &types.Response{
			Code: 200,
			Msg:  "查询失败，暂无数据",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "查询成功",
		Data: user,
	}, nil
}

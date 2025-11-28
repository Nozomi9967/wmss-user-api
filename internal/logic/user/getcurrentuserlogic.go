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

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.Response, err error) {
	userId := l.ctx.Value("user_id").(string)
	var user *types.UserInfo
	user, err = l.svcCtx.SysUserModel.SelectOneDetail(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户失败",
		}, nil
	}
	if user.RoleID != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("查询用户失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户失败，权限不足",
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "查询成功",
		Data: user,
	}, nil
	return
}

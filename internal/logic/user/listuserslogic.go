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

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户列表
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
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

	var userList *[]types.UserInfo
	userList, err = l.svcCtx.SysUserModel.SelectBatchDetail(l.ctx, req)
	if err != nil {
		l.Logger.Errorf("查询用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户失败",
		}, nil
	}
	l.Logger.Errorf("查询用户成功")
	return &types.Response{
		Code: 200,
		Msg:  "查询用户成功",
		Data: userList,
	}, nil
}

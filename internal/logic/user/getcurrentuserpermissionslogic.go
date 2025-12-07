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

type GetCurrentUserPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户权限权限
func NewGetCurrentUserPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserPermissionsLogic {
	return &GetCurrentUserPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserPermissionsLogic) GetCurrentUserPermissions() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("user_id").(string)
	var user *types.UserInfo
	user, err = l.svcCtx.SysUserModel.SelectOneDetail(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户权限失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户权限失败",
		}, nil
	}
	if user == nil {
		return &types.Response{
			Code: 200,
			Msg:  "查询失败，暂无数据",
			Data: nil,
		}, nil
	}
	if user.RoleID != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("查询用户权限失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户权限失败，权限不足",
		}, nil
	}

	var permissionsInfo *[]common.RawPermissionInfo
	permissionsInfo, err = l.svcCtx.SysPermissionModel.SelectPermissionsInfoByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户权限失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查询用户权限失败",
		}, nil
	}
	return &types.Response{
		Code: 200,
		Msg:  "查询用户权限成功",
		Data: permissionsInfo,
	}, nil
}

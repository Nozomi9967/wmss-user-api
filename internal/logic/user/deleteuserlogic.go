// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/Nozomi9967/wmss-user-api/common"
	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.Response, err error) {
	userId := l.ctx.Value("user_id").(string)
	var user *model.SysUser
	user, err = l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("删除用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "删除用户失败",
		}, nil
	}
	if user.RoleId != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("删除用户失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "删除用户失败，权限不足",
		}, nil
	}
	err = l.svcCtx.SysUserModel.SoftDelete(l.ctx, req.UserID)
	if err != nil {
		l.Logger.Errorf("删除用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "删除用户失败",
		}, nil
	}
	return &types.Response{
		Code: 200,
		Msg:  "删除用户成功",
	}, nil
}

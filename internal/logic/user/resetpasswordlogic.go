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

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置用户密码
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("user_id").(string)
	var user *model.SysUser
	user, err = l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("修改密码失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "修改密码失败",
		}, nil
	}
	if user.RoleId != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("修改密码失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "修改密码失败，权限不足",
		}, nil
	}

	user, err = l.svcCtx.SysUserModel.FindOne(l.ctx, req.UserID)
	if err != nil {
		l.Logger.Errorf("修改密码失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "修改密码失败",
		}, nil
	}
	newPassword := common.EncryptPassword(req.NewPassword)
	user.Password = newPassword
	err = l.svcCtx.SysUserModel.Update(l.ctx, user)
	if err != nil {
		l.Logger.Errorf("修改密码失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "修改密码失败",
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "修改密码成功",
	}, nil
}

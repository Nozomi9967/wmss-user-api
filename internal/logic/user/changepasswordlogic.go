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

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户密码
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.Response, err error) {
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

	// 判断旧密码是否正确
	isRight := common.VerifyPassword(req.OldPassword, user.Password)
	if !isRight {
		l.Logger.Errorf("密码错误：%s", req.OldPassword)
		return &types.Response{
			Code: 400,
			Msg:  "修改密码失败，密码错误",
		}, nil
	}

	// 加密新密码
	hashedNewPassword := common.EncryptPassword(req.NewPassword)
	user.Password = hashedNewPassword
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

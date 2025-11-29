// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户信息
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.Response, err error) {
	// 1. 从 Context 中获取要操作的用户ID（假设用户只能更新自己的信息）
	userId := l.ctx.Value("user_id").(string)

	// 2. 查找现有用户记录
	var user *model.SysUser
	user, err = l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		// 统一处理查找失败，可能是 ErrNotFound 或数据库错误
		l.Logger.Errorf("查找用户失败，无法更新: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "查找用户失败",
		}, nil
	}

	// 3. 映射并更新字段（只更新 req 中提供的值）

	// 基础字符串字段
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.RoleID != "" {
		user.RoleId = req.RoleID
	}
	if req.UserStatus != "" {
		user.UserStatus = req.UserStatus
	}

	// sql.NullString 字段
	if req.Department != "" {
		user.Department = sql.NullString{String: req.Department, Valid: true}
	}
	// 如果 req.Department 为空字符串，则保持 user.Department 的原有状态（不做更新）

	if req.Position != "" {
		user.Position = sql.NullString{String: req.Position, Valid: true}
	}

	if req.ContactPhone != "" {
		user.ContactPhone = sql.NullString{String: req.ContactPhone, Valid: true}
	}

	// 日期字段：PasswordExpireTime (需要从 string 转换为 time.Time)
	if req.PasswordExpireTime != "" {
		// 假设日期格式为 "YYYY-MM-DD"
		const layout = "2006-01-02"

		// 使用 time.Local 解析，以确保时区设置正确
		parsedTime, parseErr := time.ParseInLocation(layout, req.PasswordExpireTime, time.Local)
		if parseErr != nil {
			l.Logger.Errorf("更新用户失败: 密码过期时间格式错误: %v", parseErr)
			return &types.Response{
				Code: 400, // 400 Bad Request
				Msg:  "密码过期时间格式错误，请使用 YYYY-MM-DD 格式",
			}, nil
		}
		user.PasswordExpireTime = parsedTime
	}

	// 4. 调用 Model 层进行更新
	err = l.svcCtx.SysUserModel.Update(l.ctx, user)
	if err != nil {
		l.Logger.Errorf("执行数据库更新用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "更新用户数据失败",
		}, nil
	}

	// 5. 返回成功响应
	return &types.Response{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}

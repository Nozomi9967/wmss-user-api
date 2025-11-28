// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Nozomi9967/wmss-user-api/common"
	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建用户
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.Response, err error) {
	userId := l.ctx.Value("user_id").(string)
	var user *model.SysUser
	user, err = l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("新增用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "新增用户失败",
		}, nil
	}
	if user.RoleId != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("新增用户失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "新增用户失败，权限不足",
		}, nil
	}

	var sysUser *model.SysUser
	rawId := uuid.New().String()
	sysUserId := strings.ReplaceAll(rawId, "-", "")
	hashedPasword := common.EncryptPassword(req.Password)
	now := time.Now()
	passwordExpireTime := now.AddDate(0, 0, 90)
	sysUser = &model.SysUser{
		UserId:   sysUserId,
		UserName: req.UserName,
		RealName: req.RealName,
		Password: hashedPasword,
		RoleId:   req.RoleID,
		Department: sql.NullString{
			String: req.Department,
			Valid:  true,
		},
		Position: sql.NullString{
			String: req.Position,
			Valid:  true,
		},
		ContactPhone: sql.NullString{
			String: req.ContactPhone,
			Valid:  true,
		},
		UserStatus:         req.UserStatus,
		LastLoginTime:      sql.NullTime{Valid: false},
		PasswordExpireTime: passwordExpireTime,
		CreateTime:         now,
		UpdateTime:         now,
	}
	_, err = l.svcCtx.SysUserModel.Insert(l.ctx, sysUser)
	if err != nil {
		l.Logger.Errorf("新增用户失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "新增用户失败",
		}, nil
	}
	return &types.Response{
		Code: 200,
		Msg:  "新增用户成功",
	}, nil
}

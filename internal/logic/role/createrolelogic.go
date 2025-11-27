// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Nozomi9967/wmss-user-api/common"
	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建角色
func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) (resp *types.Response, err error) {
	userName := l.ctx.Value("username").(string)
	sysUser, err := l.svcCtx.SysUserModel.FindOneByUserName(l.ctx, userName)
	if err != nil || sysUser.RoleId != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("插入角色失败，权限不足: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "创建角色失败，权限不足",
		}, err
	}
	// 生成唯一ID，可用 UUID
	roleId := strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err = l.svcCtx.SysRoleModel.Insert(l.ctx, &model.SysRole{
		RoleId:   roleId,
		RoleName: req.RoleName,
		RoleDesc: sql.NullString{String: req.RoleDesc, Valid: true},
	})
	if err != nil {
		l.Logger.Errorf("插入角色失败: %v", err)
		return &types.Response{
			Code: 500,
			Msg:  "创建角色失败",
		}, err
	}
	return &types.Response{
		Code: 200,
		Msg:  "创建角色成功",
	}, nil
}

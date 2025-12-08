// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"
	"fmt"

	"github.com/Nozomi9967/wmss-user-api/common"
	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/svc"
	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分配角色权限
func NewAssignRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolePermissionsLogic {
	return &AssignRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRolePermissionsLogic) AssignRolePermissions(req *types.AssignRolePermissionsReq) (resp *types.Response, err error) {
	userId, ok := l.ctx.Value("user_id").(string)
	if !ok || userId == "" {
		l.Logger.Errorf("登录过期：用户ID为空或类型错误")
		return &types.Response{
			Code: 401,
			Msg:  "登录过期，请重新登录",
		}, nil
	}

	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if user == nil || err != nil {
		l.Logger.Errorf("用户[%s]不存在", userId)
		return &types.Response{
			Code: 404,
			Msg:  "用户不存在",
		}, nil
	}

	if user.RoleId != common.SUPER_ADMIN_ROLE_ID {
		l.Logger.Errorf("用户[%s]权限不足，尝试更新角色[%s]", userId, req.RoleID)
		return &types.Response{
			Code: 403,
			Msg:  "权限不足，仅超级管理员可更新角色",
		}, nil
	}

	// 1. 查询该角色现有的所有权限
	existingPermissionIds, err := l.svcCtx.SysRolePermissionModel.FindPermissionsByRoleId(l.ctx, req.RoleID)
	if err != nil {
		l.Logger.Errorf("查询角色[%s]现有权限失败: %v", req.RoleID, err)
		return &types.Response{
			Code: 500,
			Msg:  "查询现有权限失败",
		}, nil
	}

	// 2. 将请求的权限ID转为 map，方便查找
	newPermissionMap := make(map[string]bool)
	for _, permissionId := range req.PermissionIDs {
		newPermissionMap[permissionId] = true
	}

	// 3. 找出需要删除的权限（存在于数据库但不在新列表中）
	var toDelete []string
	for _, existingId := range existingPermissionIds {
		if !newPermissionMap[existingId] {
			toDelete = append(toDelete, existingId)
		}
	}

	// 4. 找出需要新增的权限（在新列表中但不在数据库中）
	existingPermissionMap := make(map[string]bool)
	for _, existingId := range existingPermissionIds {
		existingPermissionMap[existingId] = true
	}

	var toInsert []*model.SysRolePermission
	for _, permissionId := range req.PermissionIDs {
		if !existingPermissionMap[permissionId] {
			toInsert = append(toInsert, &model.SysRolePermission{
				RoleId:       req.RoleID,
				PermissionId: permissionId,
			})
		}
	}

	// 5. 执行删除操作
	if len(toDelete) > 0 {
		err = l.svcCtx.SysRolePermissionModel.DeleteByRoleIdAndPermissionIds(l.ctx, req.RoleID, toDelete)
		if err != nil {
			l.Logger.Errorf("删除角色[%s]的权限失败: %v", req.RoleID, err)
			return &types.Response{
				Code: 500,
				Msg:  "删除权限失败",
			}, nil
		}
		l.Logger.Infof("删除角色[%s]的 %d 个权限", req.RoleID, len(toDelete))
	}

	// 6. 执行新增操作
	if len(toInsert) > 0 {
		_, err = l.svcCtx.SysRolePermissionModel.InsertBatch(l.ctx, toInsert)
		if err != nil {
			l.Logger.Errorf("新增角色[%s]的权限失败: %v", req.RoleID, err)
			return &types.Response{
				Code: 500,
				Msg:  "新增权限失败",
			}, nil
		}
		l.Logger.Infof("新增角色[%s]的 %d 个权限", req.RoleID, len(toInsert))
	}

	return &types.Response{
		Code: 200,
		Msg:  fmt.Sprintf("分配权限成功，删除 %d 个，新增 %d 个", len(toDelete), len(toInsert)),
	}, nil
}

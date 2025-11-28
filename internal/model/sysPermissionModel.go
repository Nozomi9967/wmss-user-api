package model

import (
	"context"

	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysPermissionModel = (*customSysPermissionModel)(nil)

type (
	// SysPermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysPermissionModel.
	SysPermissionModel interface {
		sysPermissionModel
		withSession(session sqlx.Session) SysPermissionModel
		SelectPermissionsInfoByUserId(ctx context.Context, userId string) (*[]types.PermissionInfo, error)
	}

	customSysPermissionModel struct {
		*defaultSysPermissionModel
	}
)

// NewSysPermissionModel returns a model for the database table.
func NewSysPermissionModel(conn sqlx.SqlConn) SysPermissionModel {
	return &customSysPermissionModel{
		defaultSysPermissionModel: newSysPermissionModel(conn),
	}
}

func (m *customSysPermissionModel) withSession(session sqlx.Session) SysPermissionModel {
	return NewSysPermissionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSysPermissionModel) SelectPermissionsInfoByUserId(ctx context.Context, userId string) (*[]types.PermissionInfo, error) {
	var dbpermissionInfo []additionTypes

	// SQL 逻辑解析：
	// 1. 从权限表 (p) 出发，获取权限详情。
	// 2. 内连接 角色-权限关联表 (rp)，找到该权限属于哪些角色。
	// 3. 内连接 用户表 (u)，通过 rp.role_id = u.role_id 找到拥有该角色的用户。
	// 4. 条件筛选：指定 userId，且 权限、用户 均未被逻辑删除。
	// 5. IFNULL/COALESCE：处理数据库 NULL 值，防止映射到 Go string 字段时报错。
	query := `
		SELECT
			p.permission_id,
			p.permission_name,
			p.permission_code,
			p.permission_type,
			COALESCE(p.parent_permission_id, '') AS parent_permission_id
		FROM
			sys_permission p
		INNER JOIN
			sys_role_permission rp ON p.permission_id = rp.permission_id
		INNER JOIN
			sys_user u ON rp.role_id = u.role_id
		WHERE
			u.user_id = ?
			AND p.deleted_at IS NULL
			AND u.deleted_at IS NULL
	`

	// 执行查询
	err := m.conn.QueryRowsCtx(ctx, &dbpermissionInfo, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		// 如果用户没有分配角色，或者角色没有权限，返回空切片而不是错误，视业务需求而定
		return &[]types.PermissionInfo{}, nil
	default:
		return nil, err
	}
}

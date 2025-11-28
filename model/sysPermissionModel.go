package model

import (
	"context"

	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
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
	sql := `
        SELECT DISTINCT
            sp.permission_id,
            sp.permission_name,
            sp.permission_code,
            sp.permission_type,
            sp.parent_permission_id,
            sp.create_time,
            sp.update_time
        FROM sys_permission sp
        INNER JOIN sys_role_permission srp ON sp.permission_id = srp.permission_id
        INNER JOIN sys_user_role sur ON srp.role_id = sur.role_id
        WHERE sur.user_id = ? 
            AND sp.deleted_at IS NULL
            AND sur.deleted_at IS NULL
        ORDER BY sp.create_time ASC
    `

	var permissions []types.PermissionInfo
	err := m.conn.QueryRowsCtx(ctx, &permissions, sql, userId)
	if err != nil {
		logx.WithContext(ctx).Errorf("SelectPermissionsInfoByUserId failed: userId=%s, err=%v", userId, err)
		return nil, err
	}

	return &permissions, nil
}

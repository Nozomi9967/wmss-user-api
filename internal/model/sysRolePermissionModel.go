package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRolePermissionModel = (*customSysRolePermissionModel)(nil)

type (
	// SysRolePermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRolePermissionModel.
	SysRolePermissionModel interface {
		sysRolePermissionModel
		withSession(session sqlx.Session) SysRolePermissionModel
		DeleteByRoleIdAndPermissionIds(ctx context.Context, roleId string, permissionIds []string) error
	}

	customSysRolePermissionModel struct {
		*defaultSysRolePermissionModel
	}
)

// NewSysRolePermissionModel returns a model for the database table.
func NewSysRolePermissionModel(conn sqlx.SqlConn) SysRolePermissionModel {
	return &customSysRolePermissionModel{
		defaultSysRolePermissionModel: newSysRolePermissionModel(conn),
	}
}

func (m *customSysRolePermissionModel) withSession(session sqlx.Session) SysRolePermissionModel {
	return NewSysRolePermissionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultSysRolePermissionModel) DeleteByRoleIdAndPermissionIds(ctx context.Context, roleId string, permissionIds []string) error {
	if len(permissionIds) == 0 {
		return nil
	}

	placeholders := make([]string, len(permissionIds))
	args := make([]interface{}, len(permissionIds)+1)
	args[0] = roleId

	for i, id := range permissionIds {
		placeholders[i] = "?"
		args[i+1] = id
	}
	placeholderStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf(
		"delete from %s where `role_id` = ? and `permission_id` in (%s)",
		m.table, placeholderStr,
	)

	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

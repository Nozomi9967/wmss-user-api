package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleModel = (*customSysRoleModel)(nil)

type (
	// SysRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleModel.
	SysRoleModel interface {
		sysRoleModel
		withSession(session sqlx.Session) SysRoleModel
		FindOneLogical(ctx context.Context, roleId string) (*SysRole, error)
	}

	customSysRoleModel struct {
		*defaultSysRoleModel
	}
)

// NewSysRoleModel returns a model for the database table.
func NewSysRoleModel(conn sqlx.SqlConn) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: newSysRoleModel(conn),
	}
}

func (m *customSysRoleModel) withSession(session sqlx.Session) SysRoleModel {
	return NewSysRoleModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSysRoleModel) FindOneLogical(ctx context.Context, roleId string) (*SysRole, error) {
	query := fmt.Sprintf("select %s from %s where `role_id` = ? and `deleted_at` is null limit 1", sysRoleRows, m.table)
	var resp *SysRole
	_ = m.conn.QueryRowCtx(ctx, &resp, query, roleId)
	return resp, nil
}

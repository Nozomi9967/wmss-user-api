package svc

import (
	"github.com/Nozomi9967/wmss-user-api/internal/config"
	"github.com/Nozomi9967/wmss-user-api/internal/model"
	"github.com/Nozomi9967/wmss-user-api/internal/repository"
	"github.com/Nozomi9967/wmss-user-api/middleware"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config   `json:"config"`
	DB             *gorm.DB        `json:"db,omitempty"`
	AuthMiddleware rest.Middleware `json:"auth_middleware,omitempty"`

	SysUserModel           model.SysUserModel              `json:"sys_user_model,omitempty"`
	SysRoleModel           model.SysRoleModel              `json:"sys_role_model,omitempty"`
	SysPermissionModel     model.SysPermissionModel        `json:"sys_permission_model,omitempty"`
	SysRolePermissionModel model.SysRolePermissionModel    `json:"sys_role_permission_model,omitempty"`
	RoleRepository         repository.RoleRepository       `json:"role_repository"`
	PermissionRepository   repository.PermissionRepository `json:"permission_repository"`
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	//zrpc.WithUnaryClientInterceptor(middleware.RpcMetaInterceptor())
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware().Handle,

		// 初始化 Model（无缓存版本）
		SysUserModel:           model.NewSysUserModel(conn),
		SysRoleModel:           model.NewSysRoleModel(conn),
		SysPermissionModel:     model.NewSysPermissionModel(conn),
		SysRolePermissionModel: model.NewSysRolePermissionModel(conn),
		RoleRepository:         *repository.NewRoleRepository(conn, model.NewSysRoleModel(conn)),
		PermissionRepository:   *repository.NewPermissionRepository(conn, model.NewSysPermissionModel(conn)),
	}
}

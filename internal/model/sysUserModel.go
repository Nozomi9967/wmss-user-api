package model

import (
	"context" // 新增：用于上下文
	"fmt"
	"strings"
	"time" // 新增：处理时间字段

	"github.com/Nozomi9967/wmss-user-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserModel = (*customSysUserModel)(nil)

type UserDetail struct {
	UserId             string     `db:"user_id"`
	UserName           string     `db:"user_name"`
	RealName           string     `db:"real_name"`
	Password           string     `db:"password"`
	RoleId             string     `db:"role_id"`
	RoleName           string     `db:"role_name"`
	Department         *string    `db:"department"`    // 替换 sql.NullString
	Position           *string    `db:"position"`      // 替换 sql.NullString
	ContactPhone       *string    `db:"contact_phone"` // 替换 sql.NullString
	UserStatus         string     `db:"user_status"`
	LastLoginTime      *time.Time `db:"last_login_time"` // 替换 sql.NullTime
	PasswordExpireTime time.Time  `db:"password_expire_time"`
	CreateTime         time.Time  `db:"create_time"`
	UpdateTime         time.Time  `db:"update_time"`
	DeletedAt          *time.Time `db:"deleted_at"`
}

// 2. 扩展 SysUserModel 接口：添加 SoftDelete 逻辑删除方法
type (
	SysUserModel interface {
		sysUserModel                                   // 继承原有方法（FindOne、Insert、Update等）
		withSession(session sqlx.Session) SysUserModel // 原有：事务会话支持
		// 新增：逻辑删除方法（入参：用户ID；返回：错误）
		SoftDelete(ctx context.Context, userId string) error
		SelectOneDetail(ctx context.Context, userId string) (sysUser *types.UserInfo, err error)
		SelectBatchDetail(ctx context.Context, req *types.ListUsersReq) (*[]types.UserInfo, error)
	}

	customSysUserModel struct {
		*defaultSysUserModel // 组合默认模型，复用原有逻辑
	}
)

// 原有：模型初始化方法（不变）
func NewSysUserModel(conn sqlx.SqlConn) SysUserModel {
	return &customSysUserModel{
		defaultSysUserModel: newSysUserModel(conn),
	}
}

// 原有：事务会话方法（不变）
func (m *customSysUserModel) withSession(session sqlx.Session) SysUserModel {
	return NewSysUserModel(sqlx.NewSqlConnFromSession(session))
}

// SoftDelete 逻辑删除系统用户：更新 deleted_at 为当前时间，标记为已删除
func (m *customSysUserModel) SoftDelete(ctx context.Context, userId string) error {
	// 1. 定义SQL：仅更新「未删除」的用户（deleted_at IS NULL），避免重复删除
	sqlStmt := `
		UPDATE sys_user 
		SET 
			deleted_at = ?,    -- 标记删除时间
			update_time = ?,   -- 同步更新时间
			user_status = '已删除' -- 可选：同步更新用户状态（符合业务语义）
		WHERE 
			user_id = ?        -- 目标用户ID
			AND deleted_at IS NULL -- 仅删除未删除的记录
	`

	// 2. 获取当前时间（推荐UTC时间，避免跨时区偏差）
	now := time.Now().UTC()

	// 3. 执行SQL：传入3个参数（deleted_at、update_time、user_id）
	_, err := m.conn.ExecCtx(
		ctx,     // 上下文（用于超时控制、事务传递）
		sqlStmt, // SQL语句
		now,     // 参数1：deleted_at = 当前时间
		now,     // 参数2：update_time = 当前时间
		userId,  // 参数3：目标用户ID
	)

	// 4. 错误处理：打印日志便于排查
	if err != nil {
		logx.WithContext(ctx).Errorf(
			"SoftDelete sys_user failed: userId=%s, err=%v",
			userId, err,
		)
		return err // 返回错误，让业务层处理
	}

	// 5. 可选：校验是否实际修改了数据（避免用户已被删除的情况）
	// （如果需要严格校验，可执行查询判断影响行数，示例如下）
	// result, err := m.conn.ExecCtx(...)
	// if err != nil { ... }
	// rowsAffected, _ := result.RowsAffected()
	// if rowsAffected == 0 {
	//     return fmt.Errorf("user %s already deleted", userId)
	// }

	return nil
}

func (m *customSysUserModel) SelectOneDetail(ctx context.Context, userId string) (sysUser *types.UserInfo, err error) {
	sql := `
    SELECT sys_user.*, sys_role.role_name 
    FROM sys_user
    LEFT JOIN sys_role ON sys_user.role_id = sys_role.role_id
    WHERE sys_user.deleted_at IS NULL AND user_id = ?
    `

	// 方案1：直接声明结构体变量（推荐）
	var userDetail UserDetail
	err = m.conn.QueryRowCtx(ctx, &userDetail, sql, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			fmt.Println("userDetail not found, userId:", userId)
			return nil, nil
		}
		return nil, err
	}

	userInfo := &types.UserInfo{
		UserID:             userDetail.UserId,
		UserName:           userDetail.UserName,
		RealName:           userDetail.RealName,
		RoleID:             userDetail.RoleId,
		RoleName:           userDetail.RoleName,
		UserStatus:         userDetail.UserStatus,
		PasswordExpireTime: userDetail.PasswordExpireTime.Format("2006-01-02 15:04:05"),
		CreateTime:         userDetail.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:         userDetail.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	// 处理可能为 nil 的指针字段
	if userDetail.Department != nil {
		userInfo.Department = *userDetail.Department
	}
	if userDetail.Position != nil {
		userInfo.Position = *userDetail.Position
	}
	if userDetail.ContactPhone != nil {
		userInfo.ContactPhone = *userDetail.ContactPhone
	}
	if userDetail.LastLoginTime != nil {
		userInfo.LastLoginTime = userDetail.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return userInfo, nil
}
func (m *customSysUserModel) SelectBatchDetail(ctx context.Context, req *types.ListUsersReq) (*[]types.UserInfo, error) {
	// 构建基础SQL
	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString(`
        SELECT 
            sys_user.*,
            sys_role.role_name 
        FROM sys_user
        LEFT JOIN sys_role ON sys_user.role_id = sys_role.role_id
        WHERE sys_user.deleted_at IS NULL
    `)

	// 构建查询条件和参数
	var args []interface{}

	// 添加条件过滤
	if req.UserName != "" {
		sqlBuilder.WriteString(" AND sys_user.user_name LIKE ?")
		args = append(args, "%"+req.UserName+"%")
	}
	if req.RealName != "" {
		sqlBuilder.WriteString(" AND sys_user.real_name LIKE ?")
		args = append(args, "%"+req.RealName+"%")
	}
	if req.RoleID != "" {
		sqlBuilder.WriteString(" AND sys_user.role_id = ?")
		args = append(args, req.RoleID)
	}
	if req.UserStatus != "" {
		sqlBuilder.WriteString(" AND sys_user.user_status = ?")
		args = append(args, req.UserStatus)
	}
	if req.Department != "" {
		sqlBuilder.WriteString(" AND sys_user.department LIKE ?")
		args = append(args, "%"+req.Department+"%")
	}

	// 添加排序
	sqlBuilder.WriteString(" ORDER BY sys_user.create_time DESC")

	// 添加分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		sqlBuilder.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, req.PageSize, offset)
	}

	sql := sqlBuilder.String()

	// 执行查询
	var userDetails []UserDetail
	err := m.conn.QueryRowsCtx(ctx, &userDetails, sql, args...)
	if err != nil {
		logx.WithContext(ctx).Errorf("SelectBatchDetail failed: %v, sql: %s, args: %v", err, sql, args)
		return nil, err
	}

	// 转换为目标类型
	userInfos := make([]types.UserInfo, 0, len(userDetails))
	for _, detail := range userDetails {
		userInfo := m.convertUserDetailToInfo(&detail)
		userInfos = append(userInfos, *userInfo)
	}

	return &userInfos, nil
}

// 提取的转换方法（可以在多个地方复用）
func (m *customSysUserModel) convertUserDetailToInfo(detail *UserDetail) *types.UserInfo {
	userInfo := &types.UserInfo{
		UserID:             detail.UserId,
		UserName:           detail.UserName,
		RealName:           detail.RealName,
		RoleID:             detail.RoleId,
		RoleName:           detail.RoleName,
		UserStatus:         detail.UserStatus,
		PasswordExpireTime: detail.PasswordExpireTime.Format("2006-01-02 15:04:05"),
		CreateTime:         detail.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:         detail.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	// 安全处理可能为nil的指针字段
	if detail.Department != nil {
		userInfo.Department = *detail.Department
	}
	if detail.Position != nil {
		userInfo.Position = *detail.Position
	}
	if detail.ContactPhone != nil {
		userInfo.ContactPhone = *detail.ContactPhone
	}
	if detail.LastLoginTime != nil {
		userInfo.LastLoginTime = detail.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return userInfo
}

package common

type RawPermissionInfo struct {
	PermissionId   string `db:"permission_id"`   // 权限唯一标识
	PermissionName string `db:"permission_name"` // 权限名称
	PermissionCode string `db:"permission_code"` // 权限编码
	PermissionType string `db:"permission_type"` // 权限类型
	// 注意：因为 SQL 中使用了 COALESCE(..., '') AS parent_permission_id
	// 将 NULL 转换成了空字符串，所以这里可以直接使用 string 类型。
	ParentPermissionId string `db:"parent_permission_id"` // 父权限ID
}

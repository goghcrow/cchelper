package module

import "database/sql"

type ArchRow struct {
	Id    sql.NullInt64
	Pid   sql.NullInt64
	Depth sql.NullInt64
	Order sql.NullInt64
	Path  sql.NullString
	Name  sql.NullString
}

type EmployeeRow struct {
	Id       sql.NullInt64  `id`
	ArchId   sql.NullInt64  `部门id`
	Erp      sql.NullString `erp`
	User     sql.NullString `域名`
	Name     sql.NullString `姓名`
	Identity sql.NullString `职级`
	To       sql.NullString `汇报对象`
}

package arch

import (
	"database/sql"
	_ "lib/mysql"
	"xf/module"
)

func (self *ArchModule) ArchFetchUser() (result map[sql.NullString]*module.EmployeeRow, err error) {
	db, err := sql.Open("mysql", module.Opt.Mysql())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT 
	id,
	部门id AS archid,
	erp,
	a.域名 AS 'user',
	a.姓名 AS 'name',
	职级 AS identity,
	b.域名 AS 'to'
FROM Employee a LEFT JOIN (
	SELECT 域名,姓名 FROM Employee
) b ON a.汇报对象 = b.姓名
`)
	if err != nil {
		return
	}

	id := &sql.NullInt64{}
	archId := &sql.NullInt64{}
	erp := &sql.NullString{}
	user := &sql.NullString{}
	name := &sql.NullString{}
	identity := &sql.NullString{}
	to := &sql.NullString{}

	result = make(map[sql.NullString]*module.EmployeeRow)

	for rows.Next() {
		err = rows.Scan(id, archId, erp, user, name, identity, to)
		if err != nil {
			return
		}

		//if user.Valid {
		result[*user] = &module.EmployeeRow{
			Id:       *id,
			ArchId:   *archId,
			Erp:      *erp,
			User:     *user,
			Name:     *name,
			Identity: *identity,
			To:       *to,
		}
		//}
	}
	return
}

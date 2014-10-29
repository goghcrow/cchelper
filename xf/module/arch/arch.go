package arch

import (
	"database/sql"
	_ "lib/mysql"
	"xf/module"
)

type ArchModule struct{}

type UserEntity struct {
}

// Todo: 内置组织架构...读取mysql???
//var x link.Channel

func init() {
	module.Arch = &ArchModule{}
}

func (self *ArchModule) ChannelInit() {

}

func (self *ArchModule) ArchFetchAll() (result map[sql.NullInt64]*module.ArchRow, err error) {
	db, err := sql.Open("mysql", module.Opt.Mysql())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Architecture")
	if err != nil {
		return
	}

	id := &sql.NullInt64{}
	pid := &sql.NullInt64{}
	depth := &sql.NullInt64{}
	order := &sql.NullInt64{}
	path := &sql.NullString{}
	name := &sql.NullString{}

	result = make(map[sql.NullInt64]*module.ArchRow)

	for rows.Next() {
		err = rows.Scan(id, pid, depth, order, path, name)
		if err != nil {
			return
		}

		result[*id] = &module.ArchRow{
			Id:    *id,
			Pid:   *pid,
			Depth: *depth,
			Order: *order,
			Path:  *path,
			Name:  *name,
		}
	}
	return
}

package sqlx

type Table struct {
	Name    string // 表名, snake name
	Sql     string // sql语句
	Comment string // 注释
}

type Schema struct {
	Name     string
	Entities []*Table
}

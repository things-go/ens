package schema

type Index struct {
	Table       string
	KeyName     string
	PrimaryKey  bool
	Unique      bool
	IsComposite bool
	Priority    int
	IndexType   string
	Columns     []string
}

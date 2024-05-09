package ens

// EntityDescriptor Each table corresponds to an EntityDescriptor
type EntityDescriptor struct {
	Name        string                  // entity name
	Comment     string                  // entity comment
	Table       TableDef                // entity table define
	Fields      []*FieldDescriptor      // field information
	Indexes     []*IndexDescriptor      // index information
	ForeignKeys []*ForeignKeyDescriptor // foreign key information
}

type EntityDescriptorSlice []*EntityDescriptor

func (t EntityDescriptorSlice) Len() int           { return len(t) }
func (t EntityDescriptorSlice) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t EntityDescriptorSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

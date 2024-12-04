package migration_types

type Migratable interface {
	Save(filename string) error
	ToNext() (Migratable, error)
}

package migration_types

type Migratable interface {
	Save(filename string) error
	ToNext() (Migratable, error)
}

type NullConfig struct{}

func (NullConfig) Save(_ string) error {
	return nil
}

func (NullConfig) ToNext() (Migratable, error) {
	return NullConfig{}, nil
}

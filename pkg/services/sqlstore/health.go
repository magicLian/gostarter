package sqlstore

type HealthStore interface {
	DBPing() error
}

func (cs *SqlStore) DBPing() error {
	return cs.GetDB().Exec("SELECT 1").Error
}

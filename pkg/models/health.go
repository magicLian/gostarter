package models

const (
	DATABASE_OK      = "ok"
	DATABASE_FAILING = "failing"
)

type Health struct {
	Database   string `json:"database"` //ok or failing
	ApiVersion string `json:"apiVersion"`
}

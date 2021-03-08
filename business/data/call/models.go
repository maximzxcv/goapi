package call

type Call struct {
	ID     string `json:"id" db:"id"`
	Method string `json:"method" db:"method"`
	Path   string `json:"path" db:"path"`
	UserID string `json:"userId" db:"userId"`
}

type CreateCall struct {
	Method string
	Path   string
}

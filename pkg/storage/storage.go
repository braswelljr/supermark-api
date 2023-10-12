package storage

type File struct {
	ID          string `validate:"omitempty"`
	Name        string `validate:"omitempty"`
	Size        string `validate:"omitempty"`
	ContentType string `validate:"omitempty"`
	Data        []byte `validate:"required"`
}

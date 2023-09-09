package globals

const (
	// Error messages
	ErrFileNotFound     = "FILE_NOT_FOUND"
	ErrPermissionDenied = "PERMISSION_DENIED"
	ErrNotCustomized    = "NOT_CUSTOMIZED"
)

type ValueError struct {
	Value string
	Err   error
}

func NewValueError(value string, err error) *ValueError {
	return &ValueError{
		Value: value,
		Err:   err,
	}
}

func (e *ValueError) Error() string {
	return e.Err.Error()
}

package database

type ObjectInterface interface {
	ToString() string
	GetOid() []byte
	SetOid([]byte)
	GetType() string
	GetMode() string
}

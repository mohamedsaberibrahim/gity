package database

type Entry struct {
	name string
	oid  []byte
}

func (e *Entry) New(name string, oid []byte) {
	e.name = name
	e.oid = oid
}

func (e *Entry) GetName() string {
	return e.name
}

func (e *Entry) GetOid() []byte {
	return e.oid
}

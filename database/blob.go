package database

type Blob struct {
	data []byte
	oid  []byte
}

func (b *Blob) New(data []byte) {
	b.data = data
}

func (b *Blob) ToString() string {
	return string(b.data[:])
}

func (b *Blob) GetOid() []byte {
	return b.oid
}

func (b *Blob) SetOid(oid []byte) {
	b.oid = oid
}

func (b *Blob) GetType() string {
	return "blob"
}

package transaction

type Database interface {
	Get(hash string) (interface{}, error)
	Put(val interface{}) (string, error)
	GetData(hash string) (uint64, []byte, error)
	PutData(typ uint64, data []byte) (string, error)
	Resolve(addr string) (string, error)
	Publish(addr, path string) error
	KeyGen() (string, error)
	KeyList() ([]string, error)
}

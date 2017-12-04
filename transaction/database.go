package transaction

type Database interface {
	Get(hash string) (interface{}, error)
	Put(val interface{}) (string, error)
	Resolve(addr string) (string, error)
	Publish(addr, path string) error
	KeyGen() (string, error)
	KeyList() ([]string, error)
}

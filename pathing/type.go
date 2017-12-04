package pathing

type Type int

const (
	Type_Map Type = iota
	Type_Array
	Type_String
	Type_Bytes
	Type_Uint64
	Type_Int64
	Type_Bool
	Type_Nil
	Type_Unknown
)

func (t Type) String() string {
	switch t {
	case Type_Map:
		return "map"
	case Type_Array:
		return "array"
	case Type_String:
		return "string"
	case Type_Bytes:
		return "bytes"
	case Type_Uint64:
		return "uint64"
	case Type_Int64:
		return "int64"
	case Type_Bool:
		return "bool"
	case Type_Nil:
		return "nil"
	case Type_Unknown:
		return "unknown"
	}
	return "invalid"
}

package dataformat

type DataFormat interface {
	Unmarshal(data []byte, targetType any) (any, error)
	Marshal(data any) (string, error)
}

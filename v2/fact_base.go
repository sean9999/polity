package polity

type FactTable struct {
	Rows map[string]FactRow
}
type FactBase struct {
	Tables map[string]FactTable
	Events chan FactResult[any]
}

type FactRow map[string]any

func (fr FactRow) Update(key string, newVal any) (any, any, error) {
	oldVal := fr[key]
	fr[key] = newVal
	return oldVal, newVal, nil
}

type FactResult[T any] struct {
	OldVal T
	NewVal T
	Err    error
	Msg    string
}

func (fr FactResult[_]) Error() string {
	return fr.Err.Error()
}

func (fr FactResult[_]) String() string {
	return fr.Msg
}

type Alives map[Peer[Addresser]]bool

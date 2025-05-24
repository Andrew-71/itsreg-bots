package bots

type EntryPoint struct {
	Key   string
	State int
}

func (e EntryPoint) IsZero() bool {
	return e == EntryPoint{}
}

func NewEntryPoint(key string, state int) (EntryPoint, error) {
	if key == "" {
		return EntryPoint{}, NewInvalidInputError(
			"invalid-entry-key",
			"expected non-empty entrypoint key",
		)
	}

	if state == 0 {
		return EntryPoint{}, NewInvalidInputError(
			"invalid-entry-key",
			"expected non-empty entrypoint state",
		)
	}

	return EntryPoint{
		Key:   key,
		State: state,
	}, nil
}

func MustNewEntryPoint(key string, state int) EntryPoint {
	e, err := NewEntryPoint(key, state)
	if err != nil {
		panic(err)
	}
	return e
}

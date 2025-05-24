package bots

type Mailing struct {
	Name         string
	EntryKey     string
	RequireState int
}

func NewMailing(
	name string,
	entryKey string,
	requireState int,
) (Mailing, error) {
	if name == "" {
		return Mailing{}, NewInvalidInputError(
			"invalid-mailing",
			"expected non-empty name",
		)
	}

	if entryKey == "" {
		return Mailing{}, errEntryPointIsEmpty
	}

	return Mailing{
		Name:         name,
		EntryKey:     entryKey,
		RequireState: requireState,
	}, nil
}

func MustNewMailing(
	name string,
	entryKey string,
	requireState int,
) Mailing {
	m, err := NewMailing(name, entryKey, requireState)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Mailing) IsZero() bool {
	return m == Mailing{}
}

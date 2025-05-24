package bots

// InvalidInputError возникает тогда и только тогда, когда ошибка
// возникла по причине невалидного ввода пользователя.
type InvalidInputError struct {
	slug string
	msg  string
}

func NewInvalidInputError(slug string, msg string) InvalidInputError {
	return InvalidInputError{
		slug: slug,
		msg:  msg,
	}
}

func (e InvalidInputError) Error() string {
	return e.msg
}

func (e InvalidInputError) Slug() string {
	return e.slug
}

var (
	errBlockIsEmpty = NewInvalidInputError(
		"invalid-block",
		"expected non-empty block",
	)
	errEntryPointIsEmpty = NewInvalidInputError(
		"invalid-entry-key",
		"expected non-empty entry point",
	)
	errMailingIsEmpty = NewInvalidInputError(
		"invalid-mailing",
		"expected non-empty mailing",
	)
)

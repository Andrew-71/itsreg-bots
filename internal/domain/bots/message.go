package bots

import (
	"slices"
)

type Message struct {
	Text    string
	Buttons []string
}

func (m Message) IsZero() bool {
	return m.Text == ""
}

func NewPlainMessage(text string) (Message, error) {
	if text == "" {
		return Message{}, NewInvalidInputError(
			"invalid-message",
			"expected non-empty message text",
		)
	}

	return Message{
		Text:    text,
		Buttons: make([]string, 0),
	}, nil
}

func MustNewPlainMessage(text string) Message {
	m, err := NewPlainMessage(text)
	if err != nil {
		panic(err)
	}
	return m
}

func NewMessageWithButtons(
	text string,
	options []Option,
) (Message, error) {
	if text == "" {
		return Message{}, NewInvalidInputError(
			"invalid-message",
			"expected non-empty message text",
		)
	}

	if len(options) == 0 {
		return Message{}, NewInvalidInputError(
			"invalid-message",
			"expected non-empty message options",
		)
	}

	buttons := make([]string, len(options))
	for i, option := range options {
		buttons[i] = option.Text
	}

	return Message{
		Text:    text,
		Buttons: buttons,
	}, nil
}

func MustNewMessageWithButtons(
	text string,
	options []Option,
) Message {
	m, err := NewMessageWithButtons(text, options)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Message) Equal(o Message) bool {
	return m.Text == o.Text && buttonsEqual(m.Buttons, o.Buttons)
}

func buttonsEqual(a, b []string) bool {
	return slices.Equal(b, a)
}

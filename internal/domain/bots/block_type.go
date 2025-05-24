package bots

import "fmt"

type BlockType struct {
	s string
}

var (
	MessageBlock   = BlockType{s: "message"}
	QuestionBlock  = BlockType{s: "question"}
	SelectionBlock = BlockType{s: "selection"}
)

func (b BlockType) String() string {
	return b.s
}

func (b BlockType) IsZero() bool {
	return b == BlockType{}
}

func NewBlockTypeFromString(s string) (BlockType, error) {
	switch s {
	case "message":
		return MessageBlock, nil
	case "question":
		return QuestionBlock, nil
	case "selection":
		return SelectionBlock, nil
	}
	return BlockType{}, NewInvalidInputError(
		"invalid-block",
		fmt.Sprintf("expected one of ['message', 'question', 'selection'], got %s", s),
	)
}

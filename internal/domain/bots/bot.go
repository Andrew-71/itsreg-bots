package bots

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/bmstu-itstech/itsreg-bots/pkg/funcs/maps"
)

const startEntryKey = "start"

type Bot struct {
	UUID string

	OwnerUUID string

	entryPoints map[string]EntryPoint
	blocks      map[int]Block
	mailings    map[string]Mailing

	Name   string
	Token  string
	Status Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

func commonBotValidate(
	uuid string,
	ownerUUID string,
	entries []EntryPoint,
	blocks []Block,
	name string,
	token string,
) error {
	if uuid == "" {
		return NewInvalidInputError(
			"invalid-bot",
			"expected non-empty bot uuid",
		)
	}

	if ownerUUID == "" {
		return NewInvalidInputError(
			"invalid-bot",
			"expected non-empty owner uuid",
		)
	}

	if len(entries) == 0 {
		return NewInvalidInputError(
			"invalid-bot",
			"expected at least one entry",
		)
	}

	if len(blocks) == 0 {
		return NewInvalidInputError(
			"invalid-bot",
			"expected at least one block",
		)
	}

	if name == "" {
		return NewInvalidInputError(
			"invalid-bot",
			"expected non-empty name",
		)
	}

	if token == "" {
		return NewInvalidInputError(
			"invalid-token",
			"expected non-empty token",
		)
	}

	if !regexp.MustCompile("[0-9]{8,10}:[a-zA-Z0-9_-]{35}").MatchString(token) {
		return NewInvalidInputError(
			"invalid-token",
			"token doesn't satisfy regex",
		)
	}

	return nil
}

func NewBot(
	uuid string,
	ownerUUID string,
	entries []EntryPoint,
	mailings []Mailing,
	blocks []Block,
	name string,
	token string,
) (*Bot, error) {
	if err := commonBotValidate(uuid, ownerUUID, entries, blocks, name, token); err != nil {
		return nil, err
	}

	if mailings == nil {
		mailings = make([]Mailing, 0)
	}

	es, err := mapEntries(entries)
	if err != nil {
		return nil, err
	}

	if _, ok := es[startEntryKey]; !ok {
		return nil, errEntryPointIsEmpty
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return nil, err
	}

	if err = hasMessageLoops(bs); err != nil {
		return nil, err
	}

	vs := vertices(bs)
	for _, entry := range entries {
		err := colorizeVertices(vs, entry.State)
		if err != nil {
			return nil, err
		}
	}

	if whiteVertexState := findWhiteVertex(vs); whiteVertexState > 0 {
		return nil, newUnusedBlockFoundError(whiteVertexState)
	}

	ms, err := mapMailings(mailings)
	if err != nil {
		return nil, err
	}

	for entryKey, mailing := range ms {
		if _, ok := es[entryKey]; !ok {
			return nil, NewInvalidInputError(
				"invalid-mailing",
				fmt.Sprintf("mailing %q has non-existent entry key %q",
					mailing.Name, entryKey),
			)
		}
	}

	return &Bot{
		UUID:        uuid,
		OwnerUUID:   ownerUUID,
		entryPoints: es,
		blocks:      bs,
		mailings:    ms,
		Name:        name,
		Token:       token,
		Status:      Stopped,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func MustNewBot(
	uuid string,
	ownerUUID string,
	entryPoints []EntryPoint,
	mailings []Mailing,
	blocks []Block,
	name string,
	token string,
) *Bot {
	b, err := NewBot(uuid, ownerUUID, entryPoints, mailings, blocks, name, token)
	if err != nil {
		panic(err)
	}
	return b
}

// note: мы теоретически можем объединить NewBot и UnmarshallBotFromDB,
// они по логике очень похожи

func UnmarshallBotFromDB(
	uuid string,
	ownerUUID string,
	entries []EntryPoint,
	mailings []Mailing,
	blocks []Block,
	name string,
	token string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Bot, error) {
	if err := commonBotValidate(uuid, ownerUUID, entries, blocks, name, token); err != nil {
		return nil, err
	}

	if mailings == nil {
		mailings = make([]Mailing, 0)
	}

	if status == "" {
		return nil, NewInvalidInputError(
			"invalid-status",
			"expected non-empty status",
		)
	}

	if createdAt.IsZero() {
		return nil, NewInvalidInputError(
			"invalid-timestamp",
			"expected non-empty created-at timestamp",
		)
	}

	if updatedAt.IsZero() {
		return nil, NewInvalidInputError(
			"invalid-timestamp",
			"expected non-empty updated-at timestamp",
		)
	}

	es, err := mapEntries(entries)
	if err != nil {
		return nil, err
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return nil, err
	}

	if err = hasMessageLoops(bs); err != nil {
		return nil, err
	}

	ms, err := mapMailings(mailings)
	if err != nil {
		return nil, err
	}

	st, err := NewStatusFromString(status)
	if err != nil {
		return nil, err
	}

	return &Bot{
		UUID:        uuid,
		OwnerUUID:   ownerUUID,
		entryPoints: es,
		blocks:      bs,
		mailings:    ms,
		Name:        name,
		Token:       token,
		Status:      st,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (b *Bot) Children(block Block) []Block {
	children := make([]Block, 0)

	if block.NextState != 0 {
		next := b.blocks[block.NextState]
		children = append(children, next)
		children = append(children, b.Children(next)...)
	}

	for _, opt := range block.Options {
		next := b.blocks[opt.Next]
		children = append(children, next)
		children = append(children, b.Children(next)...)
	}

	return children
}

func (b *Bot) Blocks() []Block {
	blocks := make([]Block, 0, len(b.blocks))
	for _, block := range b.blocks {
		blocks = append(blocks, block)
	}
	return blocks
}

func (b *Bot) Entries() []EntryPoint {
	entries := make([]EntryPoint, 0, len(b.entryPoints))
	for _, entry := range b.entryPoints {
		entries = append(entries, entry)
	}
	return entries
}

func (b *Bot) Mailings() []Mailing {
	mailings := make([]Mailing, 0, len(b.mailings))
	for _, m := range b.mailings {
		mailings = append(mailings, m)
	}
	return mailings
}

type MailingNotFoundError struct {
	EntryKey string
}

func (e MailingNotFoundError) Error() string {
	return fmt.Sprintf("mailing with entry key '%s' not found", e.EntryKey)
}

func (b *Bot) Mailing(entryKey string) (Mailing, error) {
	m, ok := b.mailings[entryKey]
	if !ok {
		return Mailing{}, MailingNotFoundError{entryKey}
	}
	return m, nil
}

func (b *Bot) AddMailing(name string, requireState int, entry EntryPoint, blocks []Block) error {
	for _, block := range blocks {
		if block.IsZero() {
			return errBlockIsEmpty
		}
		if _, ok := b.blocks[block.State]; ok {
			return newBlockIsDuplicatedError(block.State)
		}
	}

	if _, ok := b.entryPoints[entry.Key]; ok {
		return newEntryIsDuplicatedError(entry.Key)
	}

	if _, ok := b.mailings[entry.Key]; ok {
		return newMailingIsDuplicatedError(entry.Key)
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return err
	}

	bs = maps.Join(bs, b.blocks)
	vs := vertices(bs)
	err = colorizeVertices(vs, entry.State)
	if err != nil {
		return err
	}
	for _, entry := range b.entryPoints {
		err = colorizeVertices(vs, entry.State)
		if err != nil {
			return err
		}
	}

	if whiteVertexState := findWhiteVertex(vs); whiteVertexState > 0 {
		return newUnusedBlockFoundError(whiteVertexState)
	}

	if err = hasMessageLoops(bs); err != nil {
		return err
	}

	mailing, err := NewMailing(name, entry.Key, requireState)
	if err != nil {
		return err
	}

	// Apply
	for _, block := range blocks {
		b.blocks[block.State] = block
	}

	b.entryPoints[entry.Key] = entry
	b.mailings[entry.Key] = mailing

	return nil
}

func (b *Bot) SetStatus(status Status) {
	b.Status = status
}

var ErrPermissionDenied = errors.New("permission denied")

func (b *Bot) CanSeeBot(userUUID string) error {
	if b.OwnerUUID != userUUID {
		return ErrPermissionDenied
	}
	return nil
}

type vertex struct {
	Block Block
	Color color
}

type color int

const (
	white color = iota
	grey
	black
)

func newBlockNotFoundError(state int) error {
	return NewInvalidInputError(
		"invalid-block-not-found",
		fmt.Sprintf("block with state %d not found", state),
	)
}

func newBlockIsDuplicatedError(state int) error {
	return NewInvalidInputError(
		"invalid-block-duplicate",
		fmt.Sprintf("block with state %d is duplicated", state),
	)
}

func newBlockLoopError(state int) error {
	return NewInvalidInputError(
		"invalid-block-loop",
		fmt.Sprintf("block with state %d causes an infinite loop", state),
	)
}

func newUnusedBlockFoundError(state int) error {
	return NewInvalidInputError(
		"invalid-block-unused",
		fmt.Sprintf("block with state %d is unused", state),
	)
}

func newEntryIsDuplicatedError(key string) error {
	return NewInvalidInputError(
		"invalid-entry-duplicate",
		fmt.Sprintf("entry with key '%s' is duplicated", key),
	)
}

func newMailingIsDuplicatedError(key string) error {
	return NewInvalidInputError(
		"invalid-mailing-duplicate",
		fmt.Sprintf("mailing with key '%s' is duplicated", key),
	)
}

func mapBlocks(blocks []Block) (map[int]Block, error) {
	mapped := make(map[int]Block)
	for _, block := range blocks {
		if block.IsZero() {
			return nil, errBlockIsEmpty
		}
		if _, ok := mapped[block.State]; ok {
			return nil, newBlockIsDuplicatedError(block.State)
		}
		mapped[block.State] = block
	}
	return mapped, nil
}

func mapEntries(entries []EntryPoint) (map[string]EntryPoint, error) {
	mapped := make(map[string]EntryPoint)
	for _, entry := range entries {
		if entry.IsZero() {
			return nil, errEntryPointIsEmpty
		}
		if _, ok := mapped[entry.Key]; ok {
			return nil, newEntryIsDuplicatedError(entry.Key)
		}
		mapped[entry.Key] = entry
	}
	return mapped, nil
}

func mapMailings(mailings []Mailing) (map[string]Mailing, error) {
	mapped := make(map[string]Mailing)
	for _, mailing := range mailings {
		if mailing.IsZero() {
			return nil, errMailingIsEmpty
		}
		if _, ok := mapped[mailing.EntryKey]; ok {
			return nil, newMailingIsDuplicatedError(mailing.EntryKey)
		}
		mapped[mailing.EntryKey] = mailing
	}
	return mapped, nil
}

func vertices(blocks map[int]Block) map[int]*vertex {
	v := make(map[int]*vertex)
	for state, block := range blocks {
		v[state] = &vertex{
			Block: block,
			Color: white,
		}
	}
	return v
}

func colorizeVertices(vertices map[int]*vertex, currentState int) error {
	current, ok := vertices[currentState]
	if !ok {
		return newBlockNotFoundError(currentState)
	}

	current.Color = grey

	childrenStates := current.Block.ChildrenStates()

	for _, nextState := range childrenStates {
		next, ok := vertices[nextState]
		if !ok {
			return newBlockNotFoundError(nextState)
		}

		if next.Color == white {
			err := colorizeVertices(vertices, nextState)
			if err != nil {
				return err
			}
			next.Color = black
		}
	}

	return nil
}

func findWhiteVertex(vertices map[int]*vertex) int {
	for _, v := range vertices {
		if v.Color == white {
			return v.Block.State
		}
	}
	return 0
}

func (b *Bot) Traverse(startState int) []Block {
	vs := vertices(b.blocks)
	b.traverseRecursive(vs, startState)

	notWhitePredicate := func(s int, v *vertex) bool {
		return v.Color != white
	}

	travelled := make([]Block, 0, len(vs))
	for _, v := range maps.Filter(vs, notWhitePredicate) {
		travelled = append(travelled, v.Block)
	}
	return travelled
}

func (b *Bot) traverseRecursive(vertices map[int]*vertex, currentState int) {
	current := vertices[currentState]
	current.Color = grey

	childrenStates := current.Block.ChildrenStates()

	for _, nextState := range childrenStates {
		next := vertices[nextState]

		if next.Color == white {
			b.traverseRecursive(vertices, nextState)
			next.Color = black
		}
	}
}

// hasMessageLoops saves us from some ill-formatted bots crashing entire server
func hasMessageLoops(blocks map[int]Block) error {
	visited := make(map[int]bool)
	for i := range blocks {
		visited[i] = false
	}
	for i := range blocks {
		if err := hasMessageLoopsRecursive(i, blocks, visited); err != nil {
			fmt.Println(blocks)
			return err
		}
	}
	return nil
}

func hasMessageLoopsRecursive(curr int, blocks map[int]Block, visited map[int]bool) error {
	if blocks[curr].Type == MessageBlock {
		if visited[curr] {
			return newBlockLoopError(curr)
		}
		visited[curr] = true
		i := blocks[curr].NextState
		return hasMessageLoopsRecursive(i, blocks, visited)
	}
	return nil
}

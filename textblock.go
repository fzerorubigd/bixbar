package bixbar

import "fmt"

type textBlock struct {
	text      string
	name, ins string
}

func (tb textBlock) Name() string {
	return tb.name
}

func (tb textBlock) Instance() string {
	return tb.ins
}

func (tb *textBlock) Click(x int, y int, b Button) {
	tb.text = fmt.Sprintf("Clicked at %d,  %d, %s", x, y, b.String())
}

func (tb textBlock) FullText() string {
	return tb.text
}

func (tb textBlock) ShortText() string {
	return tb.text
}

func (tb textBlock) MinWidth() StringInt {
	return StringInt{String: tb.text}
}

func (tb textBlock) Align() Align {
	return Align("left")
}

func (tb textBlock) Color() (*Color, bool) {
	return nil, false
}

func (tb textBlock) Background() (*Color, bool) {
	return nil, false
}

func (tb textBlock) Border() (*Color, bool) {
	return nil, false
}

func (tb textBlock) Separator() bool {
	return true
}

func (tb textBlock) SeparatorBlockWidth() int {
	return 15
}

func (tb textBlock) Urgent() bool {
	return false
}

func (tb textBlock) Markup() Markup {
	return Markup("none")
}

func NewTextBlock(text, name, ins string) SimpleBlock {
	return &textBlock{
		text: text,
		name: name,
		ins:  ins,
	}
}

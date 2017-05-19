package main

import (
	"math/rand"

	"time"

	"github.com/fzerorubigd/bixbar"
)

type staticBlock struct {
	fullText string
	color    *bixbar.Color
}

func (sb staticBlock) FullText() string {
	return sb.fullText
}

func (sb staticBlock) ShortText() string {
	return sb.fullText
}

func (sb staticBlock) MinWidth() bixbar.StringInt {
	return bixbar.StringInt{String: sb.fullText}
}

func (sb staticBlock) Align() bixbar.Align {
	return bixbar.Align("left")
}

func (sb staticBlock) Color() (*bixbar.Color, bool) {
	return sb.color, sb.color != nil
}

func (sb staticBlock) Background() (*bixbar.Color, bool) {
	return nil, false
}

func (sb staticBlock) Border() (*bixbar.Color, bool) {
	return nil, false
}

func (sb staticBlock) Separator() bool {
	return true
}

func (sb staticBlock) SeparatorBlockWidth() int {
	return 15
}

func (sb staticBlock) Urgent() bool {
	return false
}

func (sb staticBlock) Markup() bixbar.Markup {
	return bixbar.Markup("none")
}

func (sb *staticBlock) Click(x int, y int, b bixbar.Button) {
	sb.color = &bixbar.Color{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
	}
	sb.fullText = time.Now().String()
}

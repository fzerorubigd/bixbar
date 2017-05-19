package main

import "github.com/fzerorubigd/bixbar"

// header is the header for i3 protocol
type header struct {
	Version     int    `json:"version"`
	StopSignal  string `json:"stop_signal"`
	ContSignal  string `json:"cont_signal"`
	ClickEvents bool   `json:"click_events"`
}

type block struct {
	FullText            string           `json:"full_text"`
	ShortText           string           `json:"short_text,omitempty"`
	Color               string           `json:"color,omitempty"`
	Background          string           `json:"background,omitempty"`
	Border              string           `json:"border,omitempty"`
	MinWidth            bixbar.StringInt `json:"min_width,omitempty"`
	Align               bixbar.Align     `json:"align,omitempty"`
	Urgent              bool             `json:"urgent,omitempty"`
	Separator           bool             `json:"separator,omitempty"`
	SeparatorBlockWidth int              `json:"separator_block_width,omitempty"`
	Markup              bixbar.Markup    `json:"markup,omitempty"`
	Name                string           `json:"name,omitempty"`
	Instance            string           `json:"instance,omitempty"`
}

// click is the click event
type click struct {
	Name     string        `json:"name"`
	Instance string        `json:"instance"`
	X        int           `json:"x"`
	Y        int           `json:"y"`
	Button   bixbar.Button `json:"button"`
}

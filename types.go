package bixbar

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Align is the align string
type Align string

// Markup is pango or none
type Markup string

// Button is the button on click
type Button int

const (
	// LeftButton means left button
	LeftButton Button = 1
	// MiddleButton means the middle button
	MiddleButton = 2
	// RightButton means the right button
	RightButton = 3
	ScrollDown  = 4
	ScrollUp    = 5
)

// StringInt is the custom string or int version
type StringInt struct {
	String string
	Int    int
}

// Color is the color
type Color struct {
	R, G, B uint8
}

// String is the stringer version of button
func (b Button) String() string {
	switch b {
	case LeftButton:
		return "LeftButton"
	case MiddleButton:
		return "MiddleButton"
	case RightButton:
		return "RightButton"
	case ScrollDown:
		return "ScrollDown"
	case ScrollUp:
		return "ScrollUp"
	default:
		return fmt.Sprint(b)
	}
}

// MarshalJSON is the custom marshaller
func (m Markup) MarshalJSON() ([]byte, error) {
	switch m {
	case "pango", "none":
		return []byte(fmt.Sprintf(`"%s"`, m)), nil
	}

	return nil, errors.New("invalid value, valids are pango,none")
}

// MarshalJSON is the custom marshaller
func (a Align) MarshalJSON() ([]byte, error) {
	switch a {
	case "center", "right", "left":
		return []byte(fmt.Sprintf(`"%s"`, a)), nil
	}

	return nil, errors.New("invalid value, valids are center,right,left")
}

// MarshalJSON is the custom marshaller
func (c Color) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"#%02x%02x%02x"`, c.R, c.G, c.B)), nil
}

// String return the html color version of this
func (c Color) String() string {
	return fmt.Sprintf(`"#%02x%02x%02x"`, c.R, c.G, c.B)
}

// MarshalJSON is the custom marshaller
func (si StringInt) MarshalJSON() ([]byte, error) {
	if si.String != "" {
		return json.Marshal(si.String)
	}
	return []byte(fmt.Sprint(si.Int)), nil
}

// header is the header for i3 protocol
type header struct {
	Version     int    `json:"version"`
	StopSignal  string `json:"stop_signal"`
	ContSignal  string `json:"cont_signal"`
	ClickEvents bool   `json:"click_events"`
}

type block struct {
	FullText            string    `json:"full_text"`
	ShortText           string    `json:"short_text,omitempty"`
	Color               string    `json:"color,omitempty"`
	Background          string    `json:"background,omitempty"`
	Border              string    `json:"border,omitempty"`
	MinWidth            StringInt `json:"min_width,omitempty"`
	Align               Align     `json:"align,omitempty"`
	Urgent              bool      `json:"urgent,omitempty"`
	Separator           bool      `json:"separator,omitempty"`
	SeparatorBlockWidth int       `json:"separator_block_width,omitempty"`
	Markup              Markup    `json:"markup,omitempty"`
	Name                string    `json:"name,omitempty"`
	Instance            string    `json:"instance,omitempty"`
}

// click is the click event
type click struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Button   Button `json:"button"`
}

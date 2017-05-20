package bixbar

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	var re string = "none"
	switch m {
	case "pango", "none":
		re = string(m)
	}

	return []byte(fmt.Sprintf(`"%s"`, re)), nil
}

// MarshalJSON is the custom marshaller
func (a Align) MarshalJSON() ([]byte, error) {
	var re string = "left"
	switch a {
	case "center", "right", "left":
		re = string(a)
	}
	return []byte(fmt.Sprintf(`"%s"`, re)), nil
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
		// for handling special case, Using the original marshal is very better
		return json.Marshal(si.String)
	}
	return []byte(fmt.Sprint(si.Int)), nil
}

var colorStr = regexp.MustCompile("^#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$")

// NewColor from string
func NewColor(s string) (*Color, error) {
	parts := colorStr.FindStringSubmatch(s)
	if len(parts) != 4 {
		return nil, fmt.Errorf("the %s is not a valid color. use the #AABBCC format", s)
	}
	toByte := func(in string) uint8 {
		i, _ := strconv.ParseInt(in, 16, 8)
		return byte(i)
	}
	res := Color{
		R: toByte(parts[1]),
		G: toByte(parts[2]),
		B: toByte(parts[3]),
	}

	return &res, nil
}

package main

import (
	"bufio"
	"bytes"
	"strings"
)

const (
	fullText int = iota
	shortText
	color
	minWidth
	align
	name
	instance
	urgent
	separator
	separatorBlockWidth
	markup
)

type blockOutput [markup + 1]string

func newBlockOutput(o []byte) blockOutput {
	reader := bufio.NewReader(bytes.NewBuffer(o))

	bl := blockOutput{}
	for i := 0; i <= markup; i++ {
		line, err := reader.ReadString('\n')
		bl[i] = strings.Trim(line, "\n\r\t")
		if err != nil {
			// We are done reading
			break
		}
	}

	return bl
}

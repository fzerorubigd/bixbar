package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/fzerorubigd/bixbar"
)

type shellBlock struct {
	command  string
	interval time.Duration
	label    string
	format   string

	name, ins string

	lock   sync.RWMutex
	output blockOutput
	// TODO
	// signal
}

func (sb *shellBlock) execute(click bool, x, y int, button bixbar.Button) error {
	sb.lock.Lock()
	defer sb.lock.Unlock()

	cmd := exec.Command("bash", "-c", sb.command)
	env := os.Environ()
	env = append(
		env,
		fmt.Sprintf("BLOCK_NAME=%s", sb.name),
		fmt.Sprintf("BLOCK_INSTANCE=%s", sb.ins),
	)
	if click {
		env = append(
			env,
			fmt.Sprintf("BLOCK_BUTTON=%d", button),
			fmt.Sprintf("BLOCK_X=%d", x),
			fmt.Sprintf("BLOCK_Y=%d", y),
		)
	}
	// TODO : log system?
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	sb.output = newBlockOutput(out)
	// its a bit tricky here. name and ins are from the above, but in i3blocks are used as
	// env passed to script. so I change them here, but not in the output json.
	if sb.output[name] != "" {
		sb.name = sb.output[name]
	}
	if sb.output[instance] != "" {
		sb.ins = sb.output[instance]
	}

	return nil
}

func (sb *shellBlock) FullText() string {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	if sb.label != "" {
		return sb.label + sb.output[fullText]
	}

	return sb.output[fullText]
}

func (sb *shellBlock) ShortText() string {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	return sb.output[shortText]
}

func (sb *shellBlock) MinWidth() bixbar.StringInt {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	min := sb.output[minWidth]

	if i, err := strconv.ParseInt(min, 10, 32); err == nil {
		return bixbar.StringInt{Int: int(i)}
	}

	return bixbar.StringInt{String: min}
}

func (sb *shellBlock) Align() bixbar.Align {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	return bixbar.Align(sb.output[align])
}

func (sb *shellBlock) Color() (*bixbar.Color, bool) {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	c, err := bixbar.NewColor(sb.output[color])
	return c, err == nil
}

func (sb *shellBlock) Background() (*bixbar.Color, bool) {
	return nil, false
}

func (sb *shellBlock) Border() (*bixbar.Color, bool) {
	return nil, false
}

func (sb *shellBlock) Separator() bool {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	b, _ := strconv.ParseBool(sb.output[separator])
	return b
}

func (sb *shellBlock) SeparatorBlockWidth() int {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	i, err := strconv.ParseInt(sb.output[separatorBlockWidth], 10, 32)
	if err != nil {
		return 0
	}

	return int(i)
}

func (sb *shellBlock) Urgent() bool {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	b, _ := strconv.ParseBool(sb.output[urgent])
	return b
}

func (sb *shellBlock) Markup() bixbar.Markup {
	sb.lock.RLock()
	defer sb.lock.RUnlock()

	return bixbar.Markup(sb.output[markup])
}

func (sb *shellBlock) Click(x int, y int, b bixbar.Button) {
	err := sb.execute(true, x, y, b)
	if err != nil {
		// TODO : log it
	}
}

func newShellBlock(name, ins, command string, interval int, label, format string) bixbar.InteractiveBlock {
	sb := &shellBlock{
		name:     name,
		ins:      ins,
		command:  command,
		interval: time.Duration(interval) * time.Second,
		label:    label,
		format:   format,
	}

	go func() {
		sb.execute(false, 0, 0, bixbar.LeftButton)
		if sb.interval == 0 {
			return
		}
		ticker := time.NewTicker(sb.interval)
		for range ticker.C {
			sb.execute(false, 0, 0, bixbar.LeftButton)
		}
	}()

	return sb
}

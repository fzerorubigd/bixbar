package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fzerorubigd/bixbar"
)

// Bar is a single bar in system
type Bar interface {
	AddBlock(string, string, bixbar.SimpleBlock)

	Start()

	Stop()
}

type singleBlock struct {
	name, ins string
	bl        bixbar.SimpleBlock
}

type bar struct {
	lock   sync.RWMutex
	blocks []singleBlock
	stop   bool
	tick   time.Duration
	reader io.Reader
	writer io.Writer

	iBlocks map[string]bixbar.InteractiveBlock

	close chan struct{}
}

func (b *bar) AddBlock(name, ins string, l bixbar.SimpleBlock) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.blocks = append(b.blocks, singleBlock{
		name: name,
		ins:  ins,
		bl:   l,
	})
	if ib, ok := l.(bixbar.InteractiveBlock); ok {
		b.iBlocks[fmt.Sprintf("%s-%s", name, ins)] = ib
	}
}

func (b *bar) writeLoop() {
	// First write the header
	v := header{
		Version:     1,
		ClickEvents: b.reader != nil,
		ContSignal:  "SIGCONT",
		StopSignal:  "SIGSTOP",
	}
	err := json.NewEncoder(b.writer).Encode(v)
	if err != nil {
		panic(err) // TODO : real error handling :)
	}
	var comma bool
theWriteLoop:
	for {
		select {
		case <-time.After(b.tick):
		case <-b.close:
			break theWriteLoop
		}
		b.write(comma)
		comma = true
	}
}

func (b *bar) callClick(c click) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if ib, ok := b.iBlocks[fmt.Sprintf("%s-%s", c.Name, c.Instance)]; ok {
		ib.Click(c.X, c.Y, c.Button)
	}
}

func (b *bar) readLoop() {
	// TODO : write into channel to use select
	buf := []byte{' '}
	for {
		n, err := b.reader.Read(buf)
		if err != nil {
			break
		}

		if n < 1 {
			continue
		}
		if buf[0] != '[' && buf[0] != ',' {
			continue
		}
		c := click{}
		err = json.NewDecoder(b.reader).Decode(&c)
		if err != nil {
			break
		}
		b.callClick(c)
		// Update the bar, since there is a chance that the bar needs to change
		b.write(true)
	}
}

func (b *bar) write(comma bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.stop {
		return
	}
	var all = make([]block, len(b.blocks))
	for i := range b.blocks {
		bl := block{
			FullText:            b.blocks[i].bl.FullText(),
			ShortText:           b.blocks[i].bl.ShortText(),
			Separator:           b.blocks[i].bl.Separator(),
			Markup:              b.blocks[i].bl.Markup(),
			Align:               b.blocks[i].bl.Align(),
			MinWidth:            b.blocks[i].bl.MinWidth(),
			SeparatorBlockWidth: b.blocks[i].bl.SeparatorBlockWidth(),
			Urgent:              b.blocks[i].bl.Urgent(),
		}
		if c, ok := b.blocks[i].bl.Color(); ok {
			bl.Color = c.String()
		}
		if c, ok := b.blocks[i].bl.Background(); ok {
			bl.Background = c.String()
		}
		if c, ok := b.blocks[i].bl.Border(); ok {
			bl.Border = c.String()
		}

		if _, ok := b.blocks[i].bl.(bixbar.InteractiveBlock); ok {
			bl.Name = b.blocks[i].name
			bl.Instance = b.blocks[i].ins
		}
		all[i] = bl
	}
	sep := []byte(",")
	if !comma {
		sep = []byte("[")
	}
	_, err := b.writer.Write(sep)
	if err != nil {
		panic(err) // TODO : :)
	}

	err = json.NewEncoder(b.writer).Encode(all)
	if err != nil {
		panic(err) // TODO : :)
	}
}

func (b *bar) sigWatch() {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGCONT, syscall.SIGSTOP)

bigLoop:
	for {
		select {

		case s := <-sigs:
			b.lock.Lock()
			if s == syscall.SIGSTOP {
				b.stop = true
			}
			if s == syscall.SIGCONT {
				b.stop = false
			}
			b.lock.Unlock()
		case <-b.close:
			break bigLoop
		}
	}
}

func (b *bar) Start() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.close = make(chan struct{})

	go b.writeLoop()
	go b.sigWatch()
	if b.reader != nil {
		go b.readLoop()
	}
}

func (b *bar) Stop() {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.close == nil {
		return
	}

	close(b.close)
}

// NewBar create a bar and start it
func NewBar(d time.Duration, out io.Writer, in io.Reader) Bar {
	res := &bar{
		iBlocks: make(map[string]bixbar.InteractiveBlock),
		tick:    d,
		writer:  out,
	}
	if in != nil {
		res.reader = in
	}
	return res
}

package gost

import (
	"bytes"
	"strings"
)

type spliterFunc func(*Spliter) spliterFunc

type Spliter struct {
	fullString []byte
	current    int
	c          chan []byte
	errChan    chan error // won't recv any errors in current version
	patterns   string
	closed     bool

	// step back
	stepBack bool
	lastEmit []byte
}

func NewSpliter(fullString []byte, splitPatterns []rune) *Spliter {
	p := &Spliter{fullString: fullString}
	p.c = make(chan []byte)
	p.errChan = make(chan error)
	p.patterns = string(splitPatterns)
	return p
}

// in the path spliting case, when emitting a section, the char at pos (should be "/" or "\") would be igored
// so after the section is emitted, set start over it to get next char
func (self *Spliter) emit(end int) {
	self.c <- self.fullString[self.current:end]
	self.current = end + 1
}

func (self *Spliter) emitAll() {
	self.c <- self.fullString[self.current:]
}

func (self *Spliter) mainTask() {
	f := findSection
	for f != nil && !self.closed {
		f = f(self)
	}
	close(self.c)
}

func (self *Spliter) Split() {
	go self.mainTask()
}

func (self *Spliter) Next() ([]byte, error) {
	if self.stepBack {
		self.stepBack = false
		return self.lastEmit, nil
	}

	select {
	case section, ok := <-self.c:
		if !ok {
			close(self.errChan)
		}
		self.lastEmit = section
		return section, nil
	case err := <-self.errChan:
		close(self.errChan)
		return nil, err
	}
}

func (self *Spliter) StepBack() {
	self.stepBack = true
}

func (self *Spliter) Close() {
	self.closed = true
}

func findSection(spliter *Spliter) spliterFunc {
	if spliter.current == len(spliter.fullString) {
		return nil
	}
	if strings.IndexByte(spliter.patterns, spliter.fullString[spliter.current]) >= 0 {
		spliter.current++
		return findSection
	}

	return getSection
}

func getSection(spliter *Spliter) spliterFunc {
	if i := bytes.IndexAny(spliter.fullString[spliter.current:], spliter.patterns); i != -1 {
		spliter.emit(spliter.current + i)
		return findSection
	}
	spliter.emitAll()
	return nil
}

package sappress

import (
	"fmt"
	"strings"
	"sync"
)

type Progress struct {
	total       int
	current     int
	barChar     string
	emptyChar   string
	barLength   int
	description string
	mu          sync.Mutex
}

func NewProgressBar(total, barLength int, description string) *Progress {
	p := &Progress{
		total:       total,
		barLength:   barLength,
		description: description,
		mu:          sync.Mutex{},
		barChar:     "█",
		emptyChar:   "░",
	}
	return p
}

func (p *Progress) increment() {
	p.current++
}

func (p *Progress) Add() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.increment()
	if p.current > p.total {
		p.current = p.total
	}
	p.Print()
}

func (p *Progress) Print() {
	precentage := float64(p.current) / float64(p.total)
	filled := int(precentage * float64(p.barLength))
	barAsString := strings.Repeat(p.barChar, filled) + strings.Repeat(p.emptyChar, p.barLength-filled)
	fmt.Printf("\r%s [%s] %d%%", p.description, barAsString, int(precentage*100))
	if precentage == 1 {
		fmt.Println()
	}
}

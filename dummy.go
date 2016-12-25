package main

import (
	"bufio"
	"io"
	"strings"
)

type displayLine []rune

type dummyDisplay struct {
	out      *bufio.Writer
	row      int
	col      int
	width    int
	height   int
	display  []displayLine
	inverted bool
}

func NewDummyDisplay(rows, cols int, out io.Writer) *dummyDisplay {
	lines := make([]displayLine, rows)
	for r := 0; r < rows; r++ {
		lines[r] = make(displayLine, cols)
	}
	display := &dummyDisplay{
		out:     bufio.NewWriter(out),
		row:     0,
		col:     0,
		width:   cols,
		height:  rows,
		display: lines,
	}
	display.Clear()

	return display
}

func (d *dummyDisplay) Clear() error {
	for r := 0; r < d.height; r++ {
		for c := 0; c < d.width; c++ {
			d.display[r][c] = ' '
		}
	}
	d.row = 0
	d.col = 0
	return d.showDisplay()
}

func (d *dummyDisplay) Write(text string) error {
	for _, c := range text {
		if c == '\n' {
			d.nextLine()

			continue
		}

		d.display[d.row][d.col] = c
		d.col++
		if d.col >= d.width {
			d.nextLine()
		}
	}

	return d.showDisplay()
}

func (d *dummyDisplay) SetDisplayInverted(inverted bool) error {
	d.inverted = inverted
	return d.showDisplay()
}

func (d *dummyDisplay) nextLine() {
	d.col = 0
	d.row++
	if d.row >= d.height {
		d.row = 0
	}
}

func (d *dummyDisplay) showDisplay() error {
	header := strings.Repeat("#", d.width+2) + "\n"
	// Header
	if _, err := d.out.WriteString(header); err != nil {
		return err
	}

	// Lines
	for _, l := range d.display {
		d.out.WriteRune('#')
		for _, r := range l {
			if r == ' ' && d.inverted {
				d.out.WriteRune('#')
			} else {
				d.out.WriteRune(r)
			}
		}
		d.out.WriteString("#\n")
	}

	// Footer
	if _, err := d.out.WriteString(header); err != nil {
		return err
	}

	return d.out.Flush()
}

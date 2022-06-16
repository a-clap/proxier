package modifier

import (
	"fmt"
	"strings"
)

type Modifier struct {
	lines []string
}

func New(buf []byte) *Modifier {
	m := &Modifier{lines: strings.Split(string(buf), "\n")}
	return m
}

func (m *Modifier) RemoveLines() error {
	return fmt.Errorf("not implemented yet")
}

func (m *Modifier) Get() []byte {
	return []byte(strings.Join(m.lines, "\n"))
}

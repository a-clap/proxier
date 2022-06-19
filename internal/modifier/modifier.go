package modifier

import (
	"proxier/pkg/logger"
	"regexp"
	"strings"
)

type Modifier struct {
	lines []string
	log   logger.Logger
}

type lineMatcher interface {
	Match(s string) bool
}

func New(buf []byte, logger logger.Logger) *Modifier {
	m := &Modifier{
		lines: strings.Split(string(buf), "\n"),
		log:   logger,
	}
	return m
}

func (m *Modifier) RemoveLines(pattern string) (linesRemoved int, err error) {
	r, err := regexp.Compile(pattern)
	var lm lineMatcher
	if err != nil {
		m.log.Infof("Pattern not compiled, trying with exact string")
		lm = &regexMatcher{r}
	} else {
		lm = &stringsMatcher{p: pattern}
	}

	return m.removeLines(lm)
}

func (m *Modifier) removeLines(lm lineMatcher) (linesRemoved int, err error) {
	var removedLines []int

	for i, line := range m.lines {
		if ok := lm.Match(line); ok {
			removedLines = append(removedLines, i)
		}
	}
	for i, lineNumber := range removedLines {
		m.lines = append(m.lines[:lineNumber-i], m.lines[lineNumber-i+1:]...)
	}
	linesRemoved = len(removedLines)
	return
}

func (m *Modifier) Get() []byte {
	return []byte(strings.Join(m.lines, "\n"))
}

type regexMatcher struct {
	*regexp.Regexp
}

type stringsMatcher struct {
	p string
}

func (r *regexMatcher) Match(s string) bool {
	return r.Match(s)
}
func (sm *stringsMatcher) Match(s string) bool {
	return strings.Contains(s, sm.p)
}

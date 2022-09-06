package modifier

import (
	"github.com/a-clap/logger"
	"regexp"
	"strings"
)

var Logger logger.Logger = logger.NewNop()

type Modifier struct {
	lines []string
}

type lineMatcher interface {
	Match(s string) bool
}

func New(buf []byte) *Modifier {
	m := &Modifier{
		lines: strings.Split(string(buf), "\n"),
	}
	return m
}

func (m *Modifier) Get() []byte {
	return []byte(strings.Join(m.lines, "\n"))
}

func (m *Modifier) RemoveLines(pattern string) (linesRemoved int, err error) {
	r, err := regexp.Compile(pattern)
	var lm lineMatcher
	if err != nil {
		Logger.Debugf("Pattern not compiled, trying with exact string %s", pattern)
		lm = &regexMatcher{r}
	} else {
		lm = &stringsMatcher{p: pattern}
	}

	return m.removeLines(lm)
}

func (m *Modifier) AppendLines(lines []string) (linesAppended int) {
	Logger.Infof("Adding lines %s", lines)
	m.lines = append(m.lines, lines...)
	return len(lines)
}

func (m *Modifier) removeLines(lm lineMatcher) (linesRemoved int, err error) {
	var removedLines []int
	for i, line := range m.lines {
		if ok := lm.Match(line); ok {
			Logger.Infof("Removing line \"%s\"", line)
			removedLines = append(removedLines, i)
		}
	}
	for i, lineNumber := range removedLines {
		m.lines = append(m.lines[:lineNumber-i], m.lines[lineNumber-i+1:]...)
	}
	linesRemoved = len(removedLines)
	return
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

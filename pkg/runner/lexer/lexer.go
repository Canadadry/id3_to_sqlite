package lexer

func Lex(in string) []string {
	out := []string{}
	l := newLexer(in)
	for next, ok := l.getNextToken(); ok; next, ok = l.getNextToken() {
		if next != "" {
			out = append(out, next)
		}
	}
	return out
}

type lexer struct {
	source  string
	current int
	read    int
	ch      byte
	eof     bool
}

func newLexer(source string) *lexer {
	l := &lexer{
		source:  source,
		current: 0,
		read:    0,
		ch:      0,
	}
	l.readChar()
	return l
}

func (l *lexer) getNextToken() (string, bool) {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}

	if l.eof {
		return "", false
	}

	if l.ch == '"' {
		return l.readQuoted(), true
	}

	return l.readUnQuoted(), true
}

func (l *lexer) readChar() {
	l.ch = 0
	l.eof = true
	if l.read < len(l.source) {
		l.ch = l.source[l.read]
		l.eof = false
	}

	l.current = l.read
	l.read += 1
}

func (l *lexer) readUnQuoted() string {
	start := l.current
	for !isWhiteSpace(l.ch) && l.eof == false {
		l.readChar()
	}
	return l.source[start:l.current]
}

func (l *lexer) readQuoted() string {
	l.readChar()
	start := l.current
	for l.ch != '"' && l.eof == false {
		l.readChar()
	}
	return l.source[start:l.current]

}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

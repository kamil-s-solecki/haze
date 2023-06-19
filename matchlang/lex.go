package matchlang

type TokenType int

const (
	CodeToken TokenType = iota
	SizeToken
	TextToken
	EqualsToken
	NotEqualsToken
	AndToken
	OrToken
	LiteralToken
)

type LexToken struct {
	Type  TokenType
	Value string
}

type SplitterState int

const (
	SplitterConsumingState SplitterState = iota
	SplitterConsumedState
	SplitterConsumingStringLiteralState
	SplitterDoneState
)

type Splitter struct {
	state          SplitterState
	start, current int
	chunk          string
	str            string
}

func (s *Splitter) consume() bool {
	if s.state == SplitterDoneState {
		return false
	}

	if s.current == len(s.str) {
		s.state = SplitterDoneState
		s.chunk = s.str[s.start:s.current]
		return true
	}

	switch s.str[s.current] {
	case ' ':
		s.consumeSpace()
	case '\'':
		s.consumeQuote()
	default:
		s.consumeOther()
	}
	s.current++
	return true
}

func (s *Splitter) consumeSpace() {
	switch s.state {
	case SplitterConsumingState:
		s.chunk = s.str[s.start:s.current]
		s.state = SplitterConsumedState
		s.start = s.current + 1
	}
}

func (s *Splitter) consumeQuote() {
	switch s.state {
	case SplitterConsumedState:
		s.state = SplitterConsumingStringLiteralState
		s.start = s.current + 1
	case SplitterConsumingStringLiteralState:
		s.state = SplitterConsumedState
		s.chunk = s.str[s.start:s.current]
	}
}

func (s *Splitter) consumeOther() {
	switch s.state {
	case SplitterConsumingStringLiteralState:
	default:
		s.state = SplitterConsumingState
	}
}

func (s *Splitter) emit() (string, bool) {
	if s.chunk == "" {
		return "", false
	}
	res := s.chunk
	s.chunk = ""
	return res, true
}

func split(s string) []string {
	result := []string{}
	splitter := Splitter{state: SplitterConsumedState, str: s}
	for splitter.consume() {
		if chunk, ok := splitter.emit(); ok {
			result = append(result, chunk)
		}
	}
	return result
}

func lex(s string) []LexToken {
	result := []LexToken{}
	for _, word := range split(s) {
		var token LexToken
		switch word {
		case "code":
			token = LexToken{Type: CodeToken}
		case "size":
			token = LexToken{Type: SizeToken}
		case "text":
			token = LexToken{Type: TextToken}
		case "=":
			token = LexToken{Type: EqualsToken}
		case "!=":
			token = LexToken{Type: NotEqualsToken}
		case "and":
			token = LexToken{Type: AndToken}
		case "or":
			token = LexToken{Type: OrToken}
		default:
			token = LexToken{Type: LiteralToken, Value: word}
		}
		result = append(result, token)
	}
	return result
}

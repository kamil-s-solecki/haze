package matchlang

type TokenType int

const (
	CodeToken TokenType = iota
	SizeToken
	TextToken
	EqualsToken
	NotEqualsToken
	LiteralToken
)

type LexToken struct {
	Type  TokenType
	Value string
}

type SplitterState int

const (
	ConsumingState SplitterState = iota
	ConsumedState
	ConsumingStringLiteralState
	DoneState
)

type Splitter struct {
	state          SplitterState
	start, current int
	chunk          string
	str            string
}

func (s *Splitter) consume() bool {
	if s.state == DoneState {
		return false
	}

	if s.current == len(s.str) {
		s.state = DoneState
		s.chunk = s.str[s.start:s.current]
		return true
	}

	switch s.str[s.current] {
	case ' ':
		s.consumeSpace()
	case '\'':
		s.consumeQuote()
	default:
		s.state = ConsumingState
	}
	s.current++
	return true
}

func (s *Splitter) consumeSpace() {
	if s.state == ConsumingState {
		s.chunk = s.str[s.start:s.current]
		s.state = ConsumedState
	}
	s.start = s.current + 1
}

func (s *Splitter) consumeQuote() {
	switch s.state {
	case ConsumedState:
		s.state = ConsumingStringLiteralState
		s.start = s.current + 1
	case ConsumingStringLiteralState:
		s.state = ConsumedState
		s.chunk = s.str[s.start : s.current-1]
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
	splitter := Splitter{state: ConsumedState, str: s}
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
		default:
			token = LexToken{Type: LiteralToken, Value: word}
		}
		result = append(result, token)
	}
	return result
}

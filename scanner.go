package main

import "strconv"

type Scanner struct {
	source  string
	tokens  []token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token{Type: EOF, lexeme: "", literal: nil, line: s.line})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addTokens(LEFT_PAREN)
		break
	case ')':
		s.addTokens(RIGHT_PAREN)
		break
	case '{':
		s.addTokens(LEFT_BRACE)
		break
	case '}':
		s.addTokens(RIGHT_BRACE)
		break
	case ',':
		s.addTokens(COMMA)
		break
	case '.':
		s.addTokens(DOT)
		break
	case '-':
		s.addTokens(MINUS)
		break
	case '+':
		s.addTokens(PLUS)
		break
	case ';':
		s.addTokens(SEMICOLON)
		break
	case '*':
		s.addTokens(STAR)
		break
	case '!':
		if s.match('=') {
			s.addTokens(BANG_EQUAL)
		} else {
			s.addTokens(BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.addTokens(EQUAL_EQUAL)
		} else {
			s.addTokens(EQUAL)
		}
		break
	case '<':
		if s.match('=') {
			s.addTokens(LESS_EQUAL)
		} else {
			s.addTokens(LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.addTokens(GREATER_EQUAL)
		} else {
			s.addTokens(GREATER)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break

	case '\n':
		s.line++
		break
	case '"':
		s.isString()
	case 'o':
		if s.match('r') {
			s.addTokens(OR)
		}
		break
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			lox.error(s.line, "Unexpected character.")
		}
		break
	}

}

func (s *Scanner) identifier() {
	keywords := make(map[string]int)
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE

	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string([]rune(s.source)[s.start:s.current])
	Type, ok := keywords[text]
	if !ok {
		s.addTokens(IDENTIFIER)
	}
	s.addTokens(TokenType(Type))
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	float, err := strconv.ParseFloat(string([]rune(s.source)[s.start:s.current]), 64)
	if err != nil {
		lox.error()
	}
	s.addToken(NUMBER, float)
}

func (s *Scanner) isDigit(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}
func (s *Scanner) isString() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		lox.error("Unterminated string value")
		return
	}

	s.advance()
	// the closing "

	stringVal := []rune(s.source)[s.start+1 : s.current-1]
	s.addToken(STRING, stringVal)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}

	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if len(s.source) <= s.current+1 {
		return '\000'
	}

	return rune(s.source[s.current+1])
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if []rune(s.source)[s.current] == expected {
		s.current++
		return true
	}
	return false
}

func (s *Scanner) advance() rune {
	src := []rune(s.source)
	curr := s.current
	s.current++
	return src[curr]
}

func (s *Scanner) addTokens(t TokenType) {
	s.addToken(t, nil)
}

func (s *Scanner) addToken(t TokenType, literal any) {
	src := []rune(s.source)
	text := src[s.start : s.current+1]
	s.tokens = append(s.tokens, token{Type: t, lexeme: string(text), literal: literal, line: s.line})
}

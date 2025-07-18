package scanner

import (
	"fmt"
)

type Tokentype int
type Token struct {
	//Represents the token type of the word.
	Ttype Tokentype
	// line represents the line number in which the token was emitted.
	Line  int
	Lexem string
}

type sprop struct {
	// line represents the current line number in the source code.
	line int
	// start is the index where the current token starts.
	start int
	// current is the index of the character currently being scanned.
	current int
}

// enum tokentypes
const (
	LEFT_PAREN Tokentype = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	QUESTION
	COLON
	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	TENARY
	// Literals.
	IDENTIFIER
	STRING
	NUMBER
	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

func isAlphabet(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumber(c byte) bool {
	return isAlphabet(c) || isDigit(c)
}

func addToken(tokens []*Token, line int, lexem string, ttype Tokentype) []*Token {
	token := Token{}
	token.Line = line
	token.Lexem = lexem
	token.Ttype = ttype
	return append(tokens, &token)
}

func peekNext(source []byte, current int) (byte, bool) {
	isEOF := current+1 >= len(source)
	if isEOF {
		//null character
		return 0, true
	}
	return source[current+1], false
}
func peek(source []byte, current int) (byte, bool) {
	isEOF := current >= len(source)
	if isEOF {
		//null character
		return 0, true
	}
	return source[current], false
}

func ScanTokens(source []byte) ([]*Token, error) {
	tokens := make([]*Token, 0, 5)

	sp := sprop{line: 1}

	for sp.current < len(source) {
		sp.start = sp.current
		c := source[sp.current]
		//update current to hold the array index of the next character
		sp.current++
		switch c {
		case '(':
			tokens = addToken(tokens, sp.line, "", LEFT_PAREN)
		case ')':
			tokens = addToken(tokens, sp.line, "", RIGHT_PAREN)
		case '{':
			tokens = addToken(tokens, sp.line, "", LEFT_BRACE)
		case '}':
			tokens = addToken(tokens, sp.line, "", RIGHT_BRACE)
		case '+':
			tokens = addToken(tokens, sp.line, "", PLUS)
		case '-':
			tokens = addToken(tokens, sp.line, "", MINUS)
		case '*':
			tokens = addToken(tokens, sp.line, "", STAR)
		case ';':
			tokens = addToken(tokens, sp.line, "", SEMICOLON)
		case ',':
			tokens = addToken(tokens, sp.line, "", COMMA)
		case '?':
			tokens = addToken(tokens, sp.line, "", QUESTION)
		case ':':
			tokens = addToken(tokens, sp.line, "", COLON)
		case '.':
			tokens = addToken(tokens, sp.line, "", DOT)
		case '!':
			tokenType := BANG
			c, _ := peek(source, sp.current)
			if c == '=' {
				tokenType = BANG_EQUAL
				sp.current++
			}
			tokens = addToken(tokens, sp.line, "", tokenType)
		case '=':
			tokenType := EQUAL
			c, _ := peek(source, sp.current)
			if c == '=' {
				tokenType = EQUAL_EQUAL
				//consume '='
				sp.current++
			}
			tokens = addToken(tokens, sp.line, "", tokenType)
		case '<':
			tokenType := LESS
			c, _ := peek(source, sp.current)
			if c == '=' {
				tokenType = LESS_EQUAL
				//consume '='
				sp.current++
			}
			tokens = addToken(tokens, sp.line, "", tokenType)
		case '>':
			tokenType := GREATER
			c, _ := peek(source, sp.current)
			if c == '=' {
				tokenType = GREATER_EQUAL
				//consume '='
				sp.current++
			}
			tokens = addToken(tokens, sp.line, "", tokenType)
		case '/':
			c, _ := peek(source, sp.current)
			if c == '/' || c == '*' {
				//handle line comment
				if c == '/' {
					//consume slash
					sp.current++
					for {
						c, isEOF := peek(source, sp.current)
						if c == '\n' || isEOF {
							if c == '\n' {
								sp.line++
								//consume new line
								sp.current++
							}
							break
						}
						sp.current++
					}
				} else {
					//handle c-style block comment
					//consume *
					sp.current++
					slashStarCount := 1
					for {
						c, isEOF := peek(source, sp.current)
						cn, _ := peekNext(source, sp.current)
						if c == '*' && cn == '/' {
							slashStarCount--
							//consumes */
							sp.current = sp.current + 2
						} else if c == '/' && cn == '*' {
							slashStarCount++
							//consumes /*
							sp.current = sp.current + 2
						} else {
							if c == '\n' {
								sp.line++
							}
							//skip characters
							sp.current++
						}
						if slashStarCount == 0 {
							break
						}
						if isEOF {
							return nil, fmt.Errorf("Expected proper comment statements")
						}
					}
				}
			} else {
				//regular slash
				tokens = addToken(tokens, sp.line, "", SLASH)
			}
		case '"':
			//string
			for {
				c, isEOF := peek(source, sp.current)
				//support for multi-line strings
				if c == '\n' {
					sp.line++
				}
				if c == '"' || isEOF {
					if isEOF {
						//report error
						return nil, fmt.Errorf("Unterminated string at line %d", sp.line)
					}
					//consume "
					sp.current++
					break
				}
				sp.current++
			}
			//what kind of allocation does the string memory array have???
			tokens = addToken(tokens, sp.line, string(source[sp.start+1:sp.current-1]), STRING)
		case ' ':
		case '\r':
		case '\t':
		// Ignore whitespace.
		case '\n':
			sp.line++
		default:
			if isDigit(c) {
				//number
				seenDot := false
				for {
					c, _ := peek(source, sp.current)
					if isDigit(c) {
						sp.current++
						continue
					} else if c == '.' && !seenDot {
						seenDot = true
						//consume '.
						sp.current++
						c, _ := peek(source, sp.current)
						if isDigit(c) {
							sp.current++
							continue
						}
						//c on this line is a non-digit
						return nil, fmt.Errorf("Invalid Float")
					} else {
						break
					}
				}
				tokens = addToken(tokens, sp.line, string(source[sp.start:sp.current]), NUMBER)
			} else if isAlphaNumber(c) || c == '_' {
				for {
					c, _ := peek(source, sp.current)
					if isAlphaNumber(c) || c == '_' {
						sp.current++
						continue
					}
					break
				}
				//handle identifiers/keywords here
				tt := IDENTIFIER
				id := string(source[sp.start:sp.current])
				keywords := map[string]Tokentype{
					"and":    AND,
					"class":  CLASS,
					"else":   ELSE,
					"false":  FALSE,
					"for":    FOR,
					"fun":    FUN,
					"if":     IF,
					"nil":    NIL,
					"or":     OR,
					"print":  PRINT,
					"return": RETURN,
					"super":  SUPER,
					"this":   THIS,
					"true":   TRUE,
					"var":    VAR,
					"while":  WHILE,
				}
				ttV, ok := keywords[id]
				if ok {
					tt = ttV
				}
				tokens = addToken(tokens, sp.line, id, tt)
			} else {
				return nil, fmt.Errorf("Unexpected Character '%c'", c)
			}
		}
	}

	tokens = addToken(tokens, sp.line, "", EOF)
	return tokens, nil
}

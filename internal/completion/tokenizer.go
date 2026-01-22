package completion

import (
	"strings"
	"unicode"
)

// Token represents a parsed SQL token.
type Token struct {
	Value string
	Type  TokenType
}

// TokenType represents the type of SQL token.
type TokenType int

const (
	TokenWord TokenType = iota
	TokenOperator
	TokenString
	TokenNumber
	TokenPunctuation
	TokenWhitespace
)

// Tokenize splits SQL input into tokens for completion analysis.
func Tokenize(input string) []Token {
	var tokens []Token
	runes := []rune(input)
	i := 0

	for i < len(runes) {
		ch := runes[i]

		switch {
		case unicode.IsSpace(ch):
			start := i
			for i < len(runes) && unicode.IsSpace(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenWhitespace,
			})

		case ch == '\'' || ch == '"':
			quote := ch
			start := i
			i++
			for i < len(runes) && runes[i] != quote {
				if runes[i] == '\\' && i+1 < len(runes) {
					i++
				}
				i++
			}
			if i < len(runes) {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenString,
			})

		case ch == '`':
			start := i
			i++
			for i < len(runes) && runes[i] != '`' {
				i++
			}
			if i < len(runes) {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenWord,
			})

		case unicode.IsLetter(ch) || ch == '_':
			start := i
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenWord,
			})

		case unicode.IsDigit(ch):
			start := i
			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenNumber,
			})

		case strings.ContainsRune("(),;*", ch):
			tokens = append(tokens, Token{
				Value: string(ch),
				Type:  TokenPunctuation,
			})
			i++

		case strings.ContainsRune("=<>!+-/%", ch):
			start := i
			for i < len(runes) && strings.ContainsRune("=<>!", runes[i]) {
				i++
			}
			if i == start {
				i++
			}
			tokens = append(tokens, Token{
				Value: string(runes[start:i]),
				Type:  TokenOperator,
			})

		default:
			tokens = append(tokens, Token{
				Value: string(ch),
				Type:  TokenPunctuation,
			})
			i++
		}
	}

	return tokens
}

// GetLastWord returns the last incomplete word from input for completion.
func GetLastWord(input string) string {
	if len(input) == 0 {
		return ""
	}

	runes := []rune(input)
	end := len(runes)

	// Find the start of the last word
	start := end
	for start > 0 {
		ch := runes[start-1]
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '`' {
			start--
		} else {
			break
		}
	}

	return string(runes[start:end])
}

// GetNonWhitespaceTokens returns only non-whitespace tokens.
func GetNonWhitespaceTokens(tokens []Token) []Token {
	var result []Token
	for _, t := range tokens {
		if t.Type != TokenWhitespace {
			result = append(result, t)
		}
	}
	return result
}

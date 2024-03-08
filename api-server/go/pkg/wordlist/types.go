package wordlist

type NumberToWord map[int]string
type WordToNumber map[string]int

type WordList struct {
	NumberToWord  *NumberToWord
	WordToNumber  *WordToNumber
	Count         int
	MaxWordLength int
	MinWordLength int
}

package wordlist

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
)

/*
 * The format for the wordlist is super simple:	it's a JSON array of strings.
 * The index of the strings is also the number id for looking up the string.
 */

func LoadWordList(WordListConfig *config.WordListConfig) *WordList {
	// check wordlist file hash
	fileHash := getFileSha256(WordListConfig.FilePath)
	if fileHash != WordListConfig.FileHash {
		fmt.Println("Wordlist file hash does not match config hash. " + WordListConfig.FilePath + "\tgot: " + fileHash + "\texpected: " + WordListConfig.FileHash)
		fmt.Println("Changing the wordlist has consequences for existing shorturls. Please update the wordlist hash in the config file to match the new wordlist file.")
		os.Exit(1)
	}

	wordList := new(WordList)
	// Load the word list from the file
	file, err := os.Open(WordListConfig.FilePath)
	if err != nil {
		fmt.Println("Error opening config ("+WordListConfig.FilePath+"):", err)
		return wordList
	}
	defer file.Close()

	// decode to an array of []string
	loadedWordList := make([]string, 8192)

	// Decode the JSON from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&loadedWordList)
	if err != nil {
		fmt.Println("Error decoding config.json:", err)
		return wordList
	}

	// Create the word list
	numberToWord := make(NumberToWord, 0)
	wordToNumber := make(WordToNumber, 0)
	maxWordLength := 0
	minWordLength := 20
	// loop through the word list and create the number to word and word to number maps
	for i, word := range loadedWordList {
		// make word TitleCase by making first letter uppercase
		word = toTitleCase(word)
		numberToWord[i] = word
		wordToNumber[word] = i

		// update the min and max word lengths
		if len(word) > maxWordLength {
			maxWordLength = len(word)
		}
		if len(word) < minWordLength {
			minWordLength = len(word)
		}
	}

	// Create the word list
	wordList.NumberToWord = &numberToWord
	wordList.WordToNumber = &wordToNumber
	wordList.Count = len(numberToWord)
	wordList.MaxWordLength = maxWordLength
	wordList.MinWordLength = minWordLength

	return wordList
}

func toTitleCase(s string) string {
	if len(s) < 1 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getFileSha256(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return ""
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// Converts a word to a number id
func GetIdFromWord(WordList *WordList, word string) int {
	// get the id from the word
	id, ok := (*WordList.WordToNumber)[word]
	if !ok {
		return -1
	}
	return id
}

// Converts a list of words to a number id
func GetIdFromWords(WordList *WordList, words []string) string {
	// get the ids from the words
	ids := []string{}
	for _, word := range words {
		// convert the word to an id
		id := GetIdFromWord(WordList, word)
		if id == -1 {
			fmt.Println("error getting id from word: ", word)
			return ""
		}
		idStr := strconv.Itoa(id)
		// pad the id with 0s to 4 characters
		for len(idStr) < 4 {
			idStr = "0" + idStr
		}
		ids = append(ids, idStr)
	}
	return strings.Join(ids, "-")
}

// Converts a number id to a word id
func GetWordFromId(WordList *WordList, id int) string {
	// get the word from the id
	id = id % WordList.Count
	return (*WordList.NumberToWord)[id]
}

// Converts number ids to a word id
func GetWordsFromId(WordList *WordList, id string) string {
	// get the words from the ids
	// split the id into the individual ids
	ids := strings.Split(id, "-")

	words := []string{}
	for _, id := range ids {
		// convert the id to an int
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("error converting id to int:", err)
			return ""
		}

		// get the word from the id
		word := GetWordFromId(WordList, id)
		if word == "" {
			fmt.Println("error getting word from id:", err)
			return ""
		}
		words = append(words, word)
	}
	return strings.Join(words, "")
}

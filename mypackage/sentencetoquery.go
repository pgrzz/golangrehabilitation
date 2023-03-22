package mypackage

import (
	"fmt"
	"strings"
)

var wordEmbeddings = map[string][]float64{
	"hello": {0.1, 0.2, 0.3},
	"world": {0.4, 0.5, 0.6},
	// add more word embeddings here
}

func SentenceToQuery(sentence string) []float64 {
	words := strings.Split(sentence, " ")
	var query []float64
	for i := 0; i < len(wordEmbeddings["hello"]); i++ {
		sum := 0.0
		for _, word := range words {
			wordEmbedding, ok := wordEmbeddings[word]
			if ok {
				sum += wordEmbedding[i]
			}
		}
		query = append(query, sum/float64(len(words)))
	}
	return query
}

func main() {
	sentence := "hello world"
	query := SentenceToQuery(sentence)
	fmt.Println(query)
}

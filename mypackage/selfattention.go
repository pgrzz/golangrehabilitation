package mypackage

import (
	"fmt"
	"math"
)

type SelfAttention struct {
	Query    []float64
	Key      []float64
	Value    []float64
	NumHeads int
}

func (sa *SelfAttention) Compute() []float64 {
	dK := float64(len(sa.Key)) / float64(sa.NumHeads)

	queries := splitIntoHeads(sa.Query, dK)
	keys := splitIntoHeads(sa.Key, dK)
	values := splitIntoHeads(sa.Value, dK)

	results := make([]float64, len(sa.Query))

	for i, q := range queries {
		for j, k := range keys {
			attention := dotProduct(q, k) / math.Sqrt(dK)
			expAttention := math.Exp(attention)
			softmax := math.Exp(attention) / sumOfExps([]float64{expAttention})
			v := values[j]
			for _, val := range v {
				results[i] += softmax * val
			}
		}
	}

	return results
}

func splitIntoHeads(array []float64, dK float64) [][]float64 {
	var result [][]float64
	start := 0
	for start < len(array) {
		end := int(math.Min(float64(len(array)), float64(start)+dK))
		result = append(result, array[start:end])
		start = end
	}
	return result
}

func dotProduct(a, b []float64) float64 {
	result := 0.0
	for i := range a {
		result += a[i] * b[i]
	}
	return result
}

func sumOfExps(array []float64) float64 {
	result := 0.0
	for _, val := range array {
		result += math.Exp(val)
	}
	return result
}

func tt() {
	sa := SelfAttention{
		Query:    []float64{1, 2, 3},
		Key:      []float64{4, 5, 6, 7, 8, 9},
		Value:    []float64{10, 11, 12, 13, 14, 15},
		NumHeads: 2,
	}

	result := sa.Compute()
	fmt.Println(result)
}

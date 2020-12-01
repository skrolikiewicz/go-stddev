package main

import "math"

func calculateStdDev(values []int) float64 {
	itemsCount := len(values)
	sum := sum(values)
	mean := calculateMean(sum, itemsCount)

	var result float64
	for i := 0; i < itemsCount; i++ {
		result += math.Pow(float64(values[i])-mean, 2)
	}

	return math.Sqrt(result / float64(itemsCount))
}

func sum(values []int) int {
	result := 0
	for _, value := range values {
		result += value
	}

	return result
}

func calculateMean(sum int, itemsCount int) float64 {
	return float64(sum) / float64(itemsCount)
}

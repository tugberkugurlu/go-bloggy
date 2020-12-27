package readingtime

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const wordsPerMinute = 200;

// http://www.craigabbott.co.uk/blog/how-to-calculate-reading-time-like-medium
func Calculate(texts []string) (time.Duration, error) {
	var numberOfWords int
	for _, text := range texts {
		numberOfWords += len(strings.Fields(text))
	}
	minutes := int(math.Ceil(float64(numberOfWords) / float64(wordsPerMinute)))
	return time.ParseDuration(fmt.Sprintf("%dm", minutes))
}
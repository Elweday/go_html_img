package helpers

import (
	"encoding/json"
	"math"
	"os"
)

// Item represents a single item in the JSON array
type Item struct {
	Time  float64 `json:"time"`
	Word  string  `json:"word"`
	Frames int     `json:"frames"`
}

// ReadAndConvertToFrames reads the JSON file and converts time to frames
func ReadAndConvertToFrames(filePath string, frameRate int) ([]Item, error) {
	// Read JSON file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data
	var items []Item
	err = json.Unmarshal(fileData, &items)
	if err != nil {
		return nil, err
	}

	// Convert time to frames
	for i := range items {
		items[i].Frames = int(math.Round(items[i].Time * float64(frameRate)))
	}

	return items, nil
}
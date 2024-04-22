package util

import "settings-loader/internal/service/json"

func SplitArrayIntoBatches(arr []json.SettingsDTO, batchSize int) [][]json.SettingsDTO {
	var batches [][]json.SettingsDTO
	totalElements := len(arr)

	for i := 0; i < totalElements; i += batchSize {
		end := i + batchSize
		if end > totalElements {
			end = totalElements
		}
		batches = append(batches, arr[i:end])
	}

	return batches
}

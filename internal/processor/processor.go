package processor

import "okta-ingestor/internal/models"

func Deduplicate(events []models.LogEvent) []models.LogEvent {
	seen := make(map[string]struct{})
	var result []models.LogEvent

	for _, e := range events {
		if _, ok := seen[e.UUID]; ok {
			continue
		}
		seen[e.UUID] = struct{}{}
		result = append(result, e)
	}

	return result
}
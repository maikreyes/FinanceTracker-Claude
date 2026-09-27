package services

import "finance-tracker/pkg/config"

func newTestServices() *Services {
	return NewServices(&config.Config{
		NotionAPIKey:               "fake-key",
		NotionDatabaseID:           "fake-data-source",
		NotionCategoryDataSourceID: "fake-data-source",
	})
}

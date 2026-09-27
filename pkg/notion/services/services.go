package services

import (
	"net/http"
	"sync"
	"time"

	"finance-tracker/pkg/config"
	"finance-tracker/pkg/ports"
)

// apiBaseURL es var, no const, para poder apuntarla a un httptest.Server
// en los tests de este paquete.
var apiBaseURL = "https://api.notion.com/v1"

// "Expenses" es una multi-source database: parent debe ser
// data_source_id, no database_id. Validado end-to-end con la REST API
// real (token de integración, no solo MCP) el 2026-09-20 (ver
// .claude/history/2026-09-20-notion-client-real-end-to-end.md).
const apiVersion = "2025-09-03"

// Services es el cliente HTTP de la API de Notion, sin dependencias
// externas. Cubre las dos data sources: "Expenses" (movimientos, recibos,
// resumen, último movimiento) y "Category".
type Services struct {
	apiKey               string
	expensesDataSourceID string // NOTION_DATABASE_ID
	categoryDataSourceID string // NOTION_CATEGORY_DATA_SOURCE_ID
	httpClient           *http.Client

	categoriesMu       sync.Mutex
	categoriesCache    map[string]categoryEntry // nombre en minúsculas -> entrada
	categoriesCachedAt time.Time
}

func NewServices(cfg *config.Config) *Services {
	return &Services{
		apiKey:               cfg.NotionAPIKey,
		expensesDataSourceID: cfg.NotionDatabaseID,
		categoryDataSourceID: cfg.NotionCategoryDataSourceID,
		httpClient:           &http.Client{Timeout: 10 * time.Second},
	}
}

var (
	_ ports.ExpenseWriter           = (*Services)(nil)
	_ ports.ReceiptUploader         = (*Services)(nil)
	_ ports.MonthlySummarizer       = (*Services)(nil)
	_ ports.LatestTransactionFinder = (*Services)(nil)
	_ ports.CategoryResolver        = (*Services)(nil)
)

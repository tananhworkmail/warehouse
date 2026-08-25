package main

import (
	"fmt"
	"log"

	"web-api/internal/api/services"
	"web-api/internal/pkg/config"
)

func main() {
	if err := config.Setup("data/config.yml"); err != nil {
		log.Fatal(err)
	}
	result, err := services.MigrateLegacyGoreTexWaterproofForms()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("migrated reports=%d revisions=%d skipped=%d\n",
		result.MigratedReports, result.MigratedRevisions, result.SkippedReports)
}

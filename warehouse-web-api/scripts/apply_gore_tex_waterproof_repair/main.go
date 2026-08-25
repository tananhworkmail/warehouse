package main

import (
	"log"
	"os"

	"web-api/internal/pkg/config"
	"web-api/internal/pkg/database"
)

func main() {
	if err := config.Setup("data/config.yml"); err != nil {
		log.Fatal(err)
	}
	script, err := os.ReadFile("scripts/gore_tex_waterproof_unicode_business_key.sql")
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.LYS_WEB_Connection()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Exec(string(script)).Error; err != nil {
		log.Fatal(err)
	}
	log.Print("waterproof Unicode business-key repair completed")
}

package main

import (
	"log"

	"github.com/manuhdez/films-api-test/database/seeders"
)

// main - Seeds the database with auto-generated content for local testing
func main() {
	log.Println("Seeding database...")
	defer log.Println("Database seeding done")

	seeder := seeders.NewSeeder()
	seeder.Seed()
}

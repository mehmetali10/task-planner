package migrate

import (
	"log"

	"github.com/mehmetali10/task-planner/internal/pkg/database/postgres"
	"github.com/mehmetali10/task-planner/internal/pkg/database/postgres/tables"
)

func MigrateAndSeed() {
	log.Print("Initializing database connection...")
	postgres.ConnectToDB()
	defer postgres.CloseDB()
	defer func() {
		log.Print("Closing database connection...")
		sqlDB, err := postgres.DB.DB()
		if err != nil {
			log.Fatalf("Failed to get database connection: %v", err)
		}
		sqlDB.Close()
	}()

	log.Print("Starting database migration...")
	if err := postgres.DB.AutoMigrate(
		&tables.Task{},
		&tables.Developer{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Print("Database migration completed successfully.")

	log.Print("Seeding developers...")
	developers := []tables.Developer{
		{ID: 1, FirstName: "DEV1", LastName: "One", Email: "dev1@example.com", Capacity: 1},
		{ID: 2, FirstName: "DEV2", LastName: "Two", Email: "dev2@example.com", Capacity: 2},
		{ID: 3, FirstName: "DEV3", LastName: "Three", Email: "dev3@example.com", Capacity: 3},
		{ID: 4, FirstName: "DEV4", LastName: "Four", Email: "dev4@example.com", Capacity: 4},
		{ID: 5, FirstName: "DEV5", LastName: "Five", Email: "dev5@example.com", Capacity: 5},
	}

	for _, dev := range developers {
		if err := postgres.DB.Create(&dev).Error; err != nil {
			log.Printf("Failed to seed developer: %v", err)
		}
	}
	log.Print("Developers seeded successfully.")
}

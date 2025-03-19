package migrate

import (
	"github.com/mehmetali10/task-planner/pkg/database/postgres"
	"github.com/mehmetali10/task-planner/pkg/database/postgres/tables"
	"github.com/mehmetali10/task-planner/pkg/log"
)

// func Migrate(db *gorm.DB, logger log.Logger) {
// 	logger.Info("Starting database migration...")
// 	if err := db.AutoMigrate(
// 		&tables.Task{},
// 		&tables.Developer{},
// 		// ...add other tables here...
// 	); err != nil {
// 		logger.Fatal("Migration failed: %v", err)
// 	}
// 	logger.Info("Database migration completed successfully.")
// }

// func SeedDevelopers(db *gorm.DB, logger log.Logger) {
// 	logger.Info("Seeding developers...")
// 	developers := []tables.Developer{
// 		{FirstName: "DEV1", LastName: "One", Email: "dev1@example.com", Capacity: 1},
// 		{FirstName: "DEV2", LastName: "Two", Email: "dev2@example.com", Capacity: 2},
// 		{FirstName: "DEV3", LastName: "Three", Email: "dev3@example.com", Capacity: 3},
// 		{FirstName: "DEV4", LastName: "Four", Email: "dev4@example.com", Capacity: 4},
// 		{FirstName: "DEV5", LastName: "Five", Email: "dev5@example.com", Capacity: 5},
// 	}

// 	for _, dev := range developers {
// 		if err := db.Create(&dev).Error; err != nil {
// 			logger.Error("Failed to seed developer: %v", err)
// 		}
// 	}
// 	logger.Info("Developers seeded successfully.")
// }

func MigrateAndSeed(logger log.Logger) {
	logger.Info("Initializing database connection...")
	postgres.ConnectToDB()
	defer postgres.CloseDB()
	defer func() {
		logger.Info("Closing database connection...")
		sqlDB, _ := postgres.DB.DB()
		sqlDB.Close()
	}()

	logger.Info("Starting database migration...")
	if err := postgres.DB.AutoMigrate(
		&tables.Task{},
		&tables.Developer{},
	); err != nil {
		logger.Fatal("Migration failed: %v", err)
	}
	logger.Info("Database migration completed successfully.")

	logger.Info("Seeding developers...")
	developers := []tables.Developer{
		{FirstName: "DEV1", LastName: "One", Email: "dev1@example.com", Capacity: 1},
		{FirstName: "DEV2", LastName: "Two", Email: "dev2@example.com", Capacity: 2},
		{FirstName: "DEV3", LastName: "Three", Email: "dev3@example.com", Capacity: 3},
		{FirstName: "DEV4", LastName: "Four", Email: "dev4@example.com", Capacity: 4},
		{FirstName: "DEV5", LastName: "Five", Email: "dev5@example.com", Capacity: 5},
	}

	for _, dev := range developers {
		if err := postgres.DB.Create(&dev).Error; err != nil {
			logger.Error("Failed to seed developer: %v", err)
		}
	}
	logger.Info("Developers seeded successfully.")
}

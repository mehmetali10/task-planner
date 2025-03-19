package postgres

import (
	"fmt"

	"github.com/mehmetali10/task-planner/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetApp().DBHost,
		config.GetApp().DBPort,
		config.GetApp().DBUser,
		config.GetApp().DBPassword,
		config.GetApp().DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return DB, nil
}

func CloseDB() error {
	if DB == nil {
		return nil
	}

	db, err := DB.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	DB = nil
	return nil
}

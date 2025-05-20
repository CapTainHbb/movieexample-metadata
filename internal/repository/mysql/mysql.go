package mysql

import (
	"context"
	"fmt"
	"log"
	"os"

	"movieexample-metadata/pkg/model"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New() (*Repository, error) {
	user := os.Getenv("MYSQL_USER")
    password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUrl := os.Getenv("DATABASE_URL")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, dbUrl, dbName)
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&model.Metadata{}); err != nil {
		log.Fatal("Migration failed:", err)
		return nil, nil
	}

	return &Repository{db}, 
	nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var metadata model.Metadata
	result := r.db.Where("id = ?", id, &metadata)
	if result.Error != nil {
		return nil, result.Error
	}

	return &metadata, nil
}

func (r *Repository) Put(ctx context.Context, id string, m *model.Metadata) error {
	newMetadata := model.Metadata { 
		Title: m.Title, 
		Description: m.Description, 
		Director: m.Director,
	}
	
	result := r.db.Create(&newMetadata)
	return result.Error
}

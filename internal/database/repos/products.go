package repos

import (
	"gorm.io/gorm"
	"update-microservice/internal/database/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := repo.db.First(&product, id)
	return &product, result.Error
}

func (repo *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := repo.db.Find(&products)
	return products, result.Error
}

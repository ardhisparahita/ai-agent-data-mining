package repository

import (
	"context"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *domain.Product) error {
	return r.DB.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryImpl) FindAll(ctx context.Context, req request.ProductQueryRequest) ([]domain.Product, int64, error) {
	var products []domain.Product
	var totalRows int64

	query := r.DB.WithContext(ctx).Model(&domain.Product{}).Preload("Category")

	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.MinPrice > 0 {
		query = query.Where("price >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		query = query.Where("price <= ?", req.MaxPrice)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
	}

	if req.SortBy == "" || allowedSort[req.SortBy] {
		req.SortBy = "created_at"
	}

	if req.Order != "asc" && req.Order != "desc" {
		req.Order = "desc"
	}

	offset := (req.Page - 1) * req.Limit

	err := query.
		Order(req.SortBy + " " + req.Order).
		Limit(req.Limit).
		Offset(offset).
		Find(&products).
		Error

	return products, totalRows, err
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.Product, error) {
	var product domain.Product

	err := r.DB.WithContext(ctx).Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, err
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *domain.Product) error {
	return r.DB.WithContext(ctx).Save(product).Error
}

func (r *ProductRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, product *domain.Product) error {
	return tx.WithContext(ctx).Save(product).Error
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

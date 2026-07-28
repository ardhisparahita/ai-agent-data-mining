package service

import (
	"context"
	"errors"
	"math"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/mapper"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		Repo: repo,
	}
}

func (s *ProductServiceImpl) Create(ctx context.Context, req request.CreateProductRequest) (*response.ProductResponse, error) {
	product := domain.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.Repo.Create(ctx, &product); err != nil {
		return nil, err
	}

	createdProduct, err := s.Repo.FindByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return mapper.ToProductResponse(createdProduct), nil
}

func (s *ProductServiceImpl) FindAll(ctx context.Context, req request.ProductQueryRequest) (*response.ProductListResponse, error) {
	req.Page, req.Limit = utils.NormalizePagination(
		req.Page, req.Limit,
	)

	products, totalRows, err := s.Repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	return &response.ProductListResponse{
		Items:      mapper.ToProductResponses(products),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalRows:  totalRows,
		TotalPages: totalPage,
	}, nil

}

func (s *ProductServiceImpl) FindByID(ctx context.Context, id uint64) (*response.ProductResponse, error) {
	product, err := s.Repo.FindByID(ctx, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("product not found")
		}

		return nil, err
	}

	return mapper.ToProductResponse(product), nil
}

func (s *ProductServiceImpl) Update(ctx context.Context, id uint64, req request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.Repo.FindByID(ctx, id)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("product not found")
		}

		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.Repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return mapper.ToProductResponse(product), nil
}

func (s *ProductServiceImpl) Delete(ctx context.Context, id uint64) error {
	_, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("product not found")
		}

		return err
	}

	return s.Repo.Delete(ctx, id)
}

func (s *ProductServiceImpl) UploadImage(ctx context.Context, id uint64, imageURL string) (*response.ProductResponse, error) {
	product, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("product not found")
		}
		return nil, err
	}

	product.ImageURL = imageURL

	if err := s.Repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return mapper.ToProductResponse(product), nil
}

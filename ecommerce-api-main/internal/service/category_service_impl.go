package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/request"
	"github.com/ardhisparahita/ecommerce-api/internal/dto/response"
	"github.com/ardhisparahita/ecommerce-api/internal/mapper"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"gorm.io/gorm"
)

type CategoryServiceImpl struct {
	Repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		Repo: repo,
	}
}

func (s *CategoryServiceImpl) Create(ctx context.Context, req request.CreateCategoryRequest) (*response.CategoryResponse, error) {
	category := domain.Category{
		Name: req.Name,
	}
	err := s.Repo.Create(ctx, &category)
	if err != nil {
		return nil, err
	}

	return mapper.ToCategoryResponse(&category), nil
}

func (s *CategoryServiceImpl) FindAll(ctx context.Context) ([]response.CategoryResponse, error) {
	categories, err := s.Repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToCategoryResponses(categories), nil
}

func (s *CategoryServiceImpl) Update(ctx context.Context, id uint64, req request.UpdateCategoryRequest) (*response.CategoryResponse, error) {
	category, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("category not found")
		}
		return nil, err
	}

	category.Name = req.Name

	if err := s.Repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return &response.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil

}

func (s *CategoryServiceImpl) FindByID(ctx context.Context, id uint64) (*response.CategoryResponse, error) {
	category, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("category not found")
		}
		return nil, err
	}

	return &response.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

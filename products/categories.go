package products

import (
	"context"
	"encore.app/products/cs"
	"fmt"
	"github.com/go-playground/validator/v10"

	"encore.app/pkg/middleware"
)

// =====================================================================================================================
// CATEGORY
// =====================================================================================================================

// CreateCategory - Create a new category
//
//		@param ctx - context.Context
//		@param payload - *CategoryRequest
//	 @return error
//
// encore:api auth method=POST path=/categories/create
func CreateCategory(ctx context.Context, payload *cs.CategoryRequest) error {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) {
		return fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// validate payload
	if err := validator.New().Struct(payload); err != nil {
		return err
	}

	// create category
	if err := cs.Create(ctx, payload); err != nil {
		return err
	}

	return nil
}

// GetCategory - Get a product
//
//	@param ctx - context.Context
//	@param id
//	@return product
//	@return error
//
// encore:api auth method=GET path=/categories/get/:id
func GetCategory(ctx context.Context, id string) (*cs.Category, error) {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &cs.Category{}, err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) {
		return &cs.Category{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// get category
	category, err := cs.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory - Update a category
//
//	@param ctx - context.Context
//	@param id - string
//	@param payload
//	@return error
//
// encore:api auth method=PATCH path=/categories/update/:id
func UpdateCategory(ctx context.Context, id string, payload *cs.UpdateCategoryRequest) error {
	// validate payload
	if err := validator.New().Struct(payload); err != nil {
		return err
	}

	// update category
	if err := cs.Update(ctx, id, payload); err != nil {
		return err
	}

	// return nil if no error
	return nil
}

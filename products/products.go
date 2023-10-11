package products

import (
	"context"
	"fmt"

	"encore.app/pkg/middleware"
	"encore.app/products/store"
)

// Get - Get a product
//
//	@param ctx - context.Context
//	@param id
//	@return product
//	@return error
//
// encore:api auth method=GET path=/products/:id
func Get(ctx context.Context, id string) (*store.Product, error) {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &store.Product{}, err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) {
		return &store.Product{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// return user
	return &store.Product{}, nil
}

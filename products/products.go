package products

import (
	"context"
	"fmt"

	"encore.app/pkg/middleware"
	"encore.app/products/ps"
)

// Get - Get a product
//
//	@param ctx - context.Context
//	@param id
//	@return product
//	@return error
//
// encore:api auth method=GET path=/products/:id
func Get(ctx context.Context, id string) (*ps.Product, error) {
	// check for claims
	claims, err := middleware.GetVerifiedClaims(ctx, "")
	if err != nil {
		return &ps.Product{}, err
	}

	// check for the roles
	if !claims.HasRole(middleware.RoleSuperAdmin, middleware.RoleAdmin) {
		return &ps.Product{}, fmt.Errorf("unauthorized: you are not authorized to perform this action")
	}

	// return user
	return &ps.Product{}, nil
}

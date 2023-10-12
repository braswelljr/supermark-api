package ps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.dev/storage/sqldb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"encore.app/pkg/database"
)

// get the service name
var productsDatabase = sqlx.NewDb(sqldb.Named("products").Stdlib(), "postgres")

// FindOneByField - get product by field
//
//	@param ctx - context.Context
//	@param field - string
//	@param ops - string
//	@param value - interface{}
//	@return product
//	@return error
func FindOneByField(ctx context.Context, field, ops string, value interface{}) (Product, error) {
	// set the data fields for the query
	data := map[string]interface{}{
		field: value,
	}

	// query statement to be executed
	q := "SELECT * FROM products WHERE %v %v :%v LIMIT 1"
	// format query parameters
	q = fmt.Sprintf(q, field, ops, field)

	// declare product
	var product Product
	// execute query
	if err := database.NamedStructQuery(ctx, productsDatabase, q, data, &product); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return Product{}, ErrNotFound
		}
		return Product{}, fmt.Errorf("selecting products by ID[%v]: %w", value, err)
	}

	return product, nil
}

// Create - Create is a function that creates a new product.
//
// @param ctx - context.Context
// @param payload
// @return product
// @return error
func Create(ctx context.Context, payload *ProductRequest) (Product, error) {
	// create a new product
	product := Product{
		Id:          uuid.New().String(),
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		CategoryId:  payload.CategoryId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// insert product into database
	if _, err := FindOneByField(ctx, "name", "=", product.Name); err == nil {
		return Product{}, fmt.Errorf("product with name[%v] already exists", product.Name)
	}

	query := `
    INSERT INTO products (id, name, description, price, category_id, created_at, updated_at)
    VALUES (:id, :name, :description, :price, :category_id, :created_at, :updated_at)
`

	// execute query
	if err := database.NamedExecQuery(ctx, productsDatabase, query, product); err != nil {
		return Product{}, fmt.Errorf("inserting product: %w", err)
	}

	// query data from database
	p, err := FindOneByField(ctx, "id", "=", product.Id)
	if err != nil {
		return Product{}, fmt.Errorf("selecting product: %w", err)
	}

	return p, nil
}

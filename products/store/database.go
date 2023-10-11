package store

import (
  "context"
  "encore.dev/storage/sqldb"
  "errors"
  "fmt"
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

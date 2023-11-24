package cs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"encore.app/pkg/pagination"

	"encore.dev/storage/sqldb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"encore.app/pkg/database"
)

// get the service name
var categoriesDatabase = sqlx.NewDb(sqldb.Named("products").Stdlib(), "postgres")

// FindOneByField - get user by field
//
//	@param ctx - context.Context
//	@param field - string
//	@param ops - string
//	@param value - interface{}
//	@return user
//	@return error
func FindOneByField(ctx context.Context, field, ops string, value interface{}) (Category, error) {
	// set the data fields for the query
	data := map[string]interface{}{
		field: value,
	}

	// query statement to be executed
	q := "SELECT * FROM categories WHERE %v %v :%v LIMIT 1"
	// format query parameters
	q = fmt.Sprintf(q, field, ops, field)

	// declare category
	var category Category
	// execute query
	if err := database.NamedStructQuery(ctx, categoriesDatabase, q, data, &category); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return Category{}, ErrNotFound
		}
		return Category{}, fmt.Errorf("selecting categories by ID[%v]: %w", value, err)
	}

	return category, nil
}

// FindManyByField - get categories by field
//
//	@param ctx - context.Context
//	@param field - string
//	@param ops - string
//	@param value - interface{}
//	@return categories
//	@return error
func FindManyByField(ctx context.Context, field, ops string, value interface{}) ([]Category, error) {
	// set the data fields for the query
	data := map[string]interface{}{
		field: value,
	}

	// query statement to be executed
	q := "SELECT * FROM categories WHERE %v %v :%v"
	// format query parameters
	q = fmt.Sprintf(q, field, ops, field)

	// declare categories
	var categories []Category
	// execute query
	if err := database.NamedStructQuery(ctx, categoriesDatabase, q, data, &categories); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return []Category{}, ErrNotFound
		}
		return []Category{}, fmt.Errorf("selecting categorys by ID[%v]: %w", value, err)
	}

	return categories, nil
}

// Create - Create a new category
//
//	@param ctx - context.Context
//	@param payload - *CreateCategoryPayload
//	@return error
func Create(ctx context.Context, payload *CategoryRequest) error {
	// check if category already exists
	cat, err := FindOneByField(ctx, "name", "=", strings.ToLower(payload.Name))
	if err == nil {
		return err
	}
	// check if category ID is empty (if not, category already exists)
	if len(strings.TrimSpace(cat.Id)) > 0 {
		return ErrAlreadyExists
	}

	// create category
	category := Category{
		Id:          uuid.New().String(),
		Name:        strings.ToLower(payload.Name),
		Description: payload.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// query statement to be executed
	query := `
    INSERT INTO categories (id, name, description, created_at, updated_at)
    VALUES (:id, :name, :description, :created_at, :updated_at)
  `

	// create category
	if err := database.NamedExecQuery(ctx, categoriesDatabase, query, category); err != nil {
		return fmt.Errorf("creating category: %w", err)
	}

	return nil
}

// Get - Get is a function that gets a category.
//
// @param ctx - context.Context
// @param id - string
// @return category
// @return error
func Get(ctx context.Context, id string) (*Category, error) {
	// check if category exists
	category, err := FindOneByField(ctx, "id", "=", id)
	if err != nil {
		return nil, fmt.Errorf("selecting category: %w", err)
	}

	// return category
	return &category, nil
}

// GetMany - GetMany is a function that gets many categories.
//
// @param ctx - context.Context
// @param ids - []string
// @return categories
// @return error
func GetMany(ctx context.Context, ids []string) ([]Category, error) {
	// check if category exists
	categories, err := FindManyByField(ctx, "id", "=", ids)
	if err != nil {
		return nil, fmt.Errorf("selecting category: %w", err)
	}

	// return category
	return categories, nil
}

// Delete - Delete is a function that deletes a category.
//
// @param ctx - context.Context
// @param id - string
// @return error
func Delete(ctx context.Context, id string) error {
	// query statement to be executed
	q := "DELETE FROM categories WHERE id = :id"

	// execute query
	if err := database.NamedExecQuery(ctx, categoriesDatabase, q, map[string]interface{}{"id": id}); err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}

	// Delete was successful
	return nil
}

// DeleteMany - DeleteMany is a function that deletes many categories.
//
// @param ctx - context.Context
// @param ids - []string
// @return error
func DeleteMany(ctx context.Context, ids []string) error {
	// query statement to be executed
	q := `
    DELETE FROM categories
    WHERE id IN (:ids)
  `

	// execute query
	if err := database.NamedExecQuery(ctx, categoriesDatabase, q, map[string]interface{}{
		"ids": ids,
	}); err != nil {
		return fmt.Errorf("deleting categories: %w", err)
	}

	// Delete was successful
	return nil
}

// Update - Update is a function that updates a category.
//
// @param ctx - context.Context
// @param id - string
// @param payload
// @return category
// @return error
func Update(ctx context.Context, id string, payload *UpdateCategoryRequest) error {
	// check if category exists
	category, err := FindOneByField(ctx, "id", "=", id)
	if err != nil {
		return fmt.Errorf("selecting category: %w", err)
	}

	// map for query fields
	fields := map[string]interface{}{}

	// if not empty, update category field
	vp := reflect.ValueOf(payload)

	// loop through payload fields and check for empty values
	for i := 0; i < vp.NumField(); i++ {
		// get the db tag name of the field
		field := vp.Type().Field(i).Tag.Get("db")
		// get the value of the field
		value := vp.Field(i).Interface()

		// if the value is not empty, add it to the fields map
		if len(strings.TrimSpace(value.(string))) > 0 {
			fields[field] = value
		}
	}

	// create query fields
	var ks []string

	fields["updated_at"] = time.Now().UTC()

	// loop through fields and create query fields
	for k := range fields {
		ks = append(ks, fmt.Sprintf("%v = :%v", k, k))
	}

	// query statement to be executed
	q := fmt.Sprintf("UPDATE categories SET %v WHERE id = :%v RETURNING *", strings.Join(ks, ", "), category.Id)

	// execute query
	if err := database.NamedExecQuery(ctx, categoriesDatabase, q, fields); err != nil {
		return fmt.Errorf("updating category: %w", err)
	}

	// return category
	return nil
}

// GetAll - GetAll is a function that gets all users.
//
//	@param ctx - context.Context
//	@return users
//	@return error
func GetAll(ctx context.Context, pag *pagination.Options) (*PaginatedCategoriesResponse, error) {
	var categories []Category

	// create query
	countQuery := `SELECT COUNT(*) FROM categories`

	// get total count of users
	// get count of categories
	count, err := database.NamedCountQuery(ctx, categoriesDatabase, countQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("getting count of categories: %w", err)
	}

	// set limit to 20 if it is less than 0 or greater than count
	if pag.Limit < 1 || pag.Limit > count {
		pag.Limit = 50
	}

	// calculate for pagination
	// var paginate pagination.Paginate
	// initialize pagination
	paging := pagination.New(pag.Page, pag.Limit, count)

	// if page is greater than total pages, set page to total pages
	if pag.Page > paging.Pages() {
		paging.SetPage(paging.Pages())
	}

	// query to set offset and limit
	const query = `SELECT * FROM categories LIMIT :limit OFFSET :offset`
	// data to be passed to the query
	p := struct {
		Limit  int `db:"limit" json:"limit" validate:"omitempty" url:"limit"`
		Offset int `db:"offset" json:"offset" validate:"omitempty" url:"offset"`
	}{
		Limit:  paging.PerPage(),
		Offset: paging.Offset(),
	}

	// execute query
	if err := database.NamedSliceQuery(ctx, categoriesDatabase, query, p, &categories); err != nil {
		return nil, fmt.Errorf("getting categories: %w", err)
	}

	return &PaginatedCategoriesResponse{
		TotalPages:  paging.Pages(),
		Total:       paging.Total(),
		CurrentPage: paging.Page(),
		Categories:  categories,
	}, nil
}

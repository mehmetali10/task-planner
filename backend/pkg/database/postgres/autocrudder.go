package postgres

import (
	"context"

	"github.com/mehmetali10/task-planner/pkg/automapper"

	"gorm.io/gorm"
)

// Create performs a database record creation based on the given request.
//
// It establishes a database connection and defers its closure.
// The function maps the request data to a new item using the automapper package.
//
// It then creates a new record in the database with the mapped item.
//
// If any error occurs during the creation, the function returns the error.
// Otherwise, it maps the created item to the response type and returns it.
func Create[Dest any, Source any](ctx context.Context, req any) (Dest, error) {
	ConnectToDB()
	defer CloseDB()

	var resp Dest
	var newItem Source
	automapper.MapLoose(req, &newItem)

	db := DB.WithContext(ctx).Model(&newItem).Create(&newItem)

	if db.Error != nil {
		return resp, db.Error
	}

	automapper.MapLoose(newItem, &resp)

	return resp, nil
}

// Read fetches database record based on the provided rule and conditions.
// It queries the database and populates a destination response.
// If successful, it returns the response; otherwise, an error is returned.
func Read[Dest any, Source any](ctx context.Context, rule any, args ...any) (Dest, error) {
	ConnectToDB()
	defer CloseDB()

	var resp Dest
	var existingItem Source

	db := DB.WithContext(ctx).Model(&existingItem).Where(map[string]interface{}{"IsDeleted": false}).Where(rule, args...).Find(&resp)

	if err := db.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, nil
		}
		return resp, err
	}

	return resp, nil
}

func ReadWithOrCondition[Dest any, Source any](ctx context.Context, rule1, rule2 any, args ...any) (Dest, error) {
	ConnectToDB()
	defer CloseDB()

	var resp Dest
	var existingItem Source

	db := DB.WithContext(ctx).Model(&existingItem).Where(args).Where(map[string]interface{}{"IsDeleted": false}).Where(rule1).Or(rule2).Find(&resp)

	if err := db.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return resp, nil
		}
		return resp, err
	}

	return resp, nil
}

// Update updates a database record based on the provided rule and request.
// It maps the request data, updates the record, and returns the updated response.
// If an error occurs during the update, it is returned.
func Update[Dest any, Source any](ctx context.Context, rule any, req any) (Dest, error) {
	ConnectToDB()
	defer CloseDB()

	var resp Dest
	var existingItem Source
	automapper.MapLoose(req, &existingItem)

	db := DB.WithContext(ctx).Model(&existingItem).Where(map[string]interface{}{"IsDeleted": false}).Where(rule).Updates(&existingItem)

	if db.Error != nil {
		return resp, db.Error
	}

	automapper.MapLoose(existingItem, &resp)

	return resp, nil
}

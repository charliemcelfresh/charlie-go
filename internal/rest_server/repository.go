package rest_server

import (
	"context"

	"github.com/charliemcelfresh/charlie-go/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	pool *sqlx.DB
}

func NewRepository() repository {
	return repository{pool: config.GetDB()}
}

func (r repository) GetItems(ctx context.Context, page int) ([]Item, error) {
	itemsToReturn := []Item{}
	userID := getUserIdFromContext(ctx)
	statement := `
		SELECT
			i.id, i.name, i.created_at, i.updated_at
		FROM
			items i
		JOIN
			user_items ui ON ui.item_id = i.id
		WHERE
		    ui.user_id = $1
		LIMIT 10
		OFFSET $2

	`
	err := r.pool.SelectContext(ctx, &itemsToReturn, statement, userID, page)
	return itemsToReturn, err
}

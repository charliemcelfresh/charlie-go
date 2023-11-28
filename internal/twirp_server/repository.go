package twirp_server

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Item struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type repository struct {
	db database
}

func NewRepository(db *sqlx.DB) repository {
	return repository{
		db,
	}
}

func (r repository) CreateItem(ctx context.Context, name string) error {
	userID := getUserIdFromContext(ctx)
	statement := `
		WITH inserted_item AS (
		    INSERT INTO items (name) values ($1) RETURNING id
	    )
		INSERT INTO
		    user_items (user_id, item_id)
	    VALUES ($2, (SELECT id FROM inserted_item));
	`
	_, err := r.db.ExecContext(ctx, statement, name, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) GetItem(ctx context.Context, itemID string) (Item, error) {
	itemToReturn := Item{}
	userID := getUserIdFromContext(ctx)
	statement := `
		SELECT
			i.id, i.name, i.created_at, i.updated_at
		FROM
			items i
		JOIN
			user_items ui ON ui.item_id = i.id
		WHERE
			i.id = $1 AND ui.user_id = $2;
	`
	err := r.db.GetContext(ctx, &itemToReturn, statement, itemID, userID)
	return itemToReturn, err
}

package rest_server

type GetItemRequest struct {
	ItemId string
}

type Item struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

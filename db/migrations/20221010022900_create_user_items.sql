-- migrate:up

CREATE TABLE user_items (
    user_id uuid NOT NULL REFERENCES users(id) DEFERRABLE,
    item_id uuid NOT NULL REFERENCES items(id) DEFERRABLE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_user_items__user_id_item_id_uniq ON user_items (user_id, item_id);

-- migrate:down

DROP TABLE user_items;
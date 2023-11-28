-- migrate:up

CREATE TABLE items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_name_uniq ON items (name);

-- migrate:down

DROP TABLE items;
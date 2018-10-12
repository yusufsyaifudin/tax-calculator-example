-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS taxes (
  "id" BIGSERIAL NOT NULL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "name" VARCHAR NOT NULL,
  "tax_code" INTEGER NOT NULL,
  "price" INTEGER NOT NULL CHECK (price >= 0), -- since pg don't have unsigned int, we need to CHECK it
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- add foreign key check
ALTER TABLE taxes ADD CONSTRAINT taxes_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- add index for faster query (where user_id = ?)
CREATE INDEX IF NOT EXISTS unique_idx_taxes_on_user_id ON taxes(user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS taxes;

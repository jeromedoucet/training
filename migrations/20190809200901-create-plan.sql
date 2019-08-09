
-- +migrate Up
CREATE TABLE "training"."plan" (
  id uuid PRIMARY KEY,
  creator_id uuid REFERENCES "training"."user"("id") NOT NULL,
  trainee_id uuid REFERENCES "training"."user"("id") NOT NULL,
  name text NOT NULL
);

-- +migrate Down
DROP TABLE "training"."plan" CASCADE;


-- +migrate Up
CREATE TABLE "training"."plan_session" (
  id uuid PRIMARY KEY,
  plan_id uuid REFERENCES "training"."plan"("id") NOT NULL,
  day date NOT NULL,
  description text NOT NULL,
  comments text
);

-- +migrate Down

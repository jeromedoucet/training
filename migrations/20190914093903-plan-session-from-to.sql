
-- +migrate Up

ALTER TABLE "training"."plan_session" ADD COLUMN "from" timestamp with time zone NOT NULL,
  ADD COLUMN "to" timestamp with time zone NOT NULL,
  DROP COLUMN "day";

-- +migrate Down
ALTER TABLE "training"."plan_session" DROP COLUMN "from",
  DROP COLUMN "to",
  ADD COLUMN "day" date NOT NULL;

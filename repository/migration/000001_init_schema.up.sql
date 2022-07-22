CREATE TABLE "user_policies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "nym_id" VARCHAR NOT NULL,
  "recipient" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE "routine_transaction_policies" (
  "schedule_start_date" DATE NOT NULL,
  "schedule_end_date" DATE NOT NULL,
  "frequency" VARCHAR NOT NULL,
  "amount" JSONB NOT NULL
)
INHERITS (
  "user_policies"
);

CREATE TABLE "transaction_trigger_policies" (
  "targeted_balance" JSONB NOT NULL,
  "amount" JSONB NOT NULL
)
INHERITS (
  "user_policies"
);

CREATE TABLE "policy_events" (
  "nym_id" VARCHAR NOT NULL,
  "transfer_sequence" BIGINT NOT NULL,
  "policy_id" BIGINT NOT NULL,
  PRIMARY KEY ("nym_id", "transfer_sequence")
);

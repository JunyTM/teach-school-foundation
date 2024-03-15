CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,   -- Tên chủ tài khoản
    "balance" bigint NOT NULL,  -- Số dư tài khoản
    "currency" varchar NOT NULL, -- Đơn vị tiền tệ
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,   -- Số lượng nhận vào
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,           -- Số lượng giao dịch
    "currency" varchar NOT NULL,        -- Dơn vị tiền tệ
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

-- Create Indexes
CREATE INDEX "idx_accounts_owner" ON "accounts" ("owner");

CREATE INDEX "idx_entries_account_id" ON "entries" ("account_id");

CREATE INDEX "idx_transfers_from_account_id" ON "transfers" ("from_account_id");

CREATE INDEX "idx_transfers_to_account_id" ON "transfers" ("to_account_id");

CREATE INDEX "idx_transfers_from_to_account_id" ON "transfers" ("from_account_id", "to_account_id");

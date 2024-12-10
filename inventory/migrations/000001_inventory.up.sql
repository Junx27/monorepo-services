CREATE TABLE IF NOT EXISTS inventories (
    uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id text NOT NULL,
    product_id text NOT NULL,
    transaction_id text NOT NULL,
    quantity int NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

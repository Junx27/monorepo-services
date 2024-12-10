CREATE TABLE IF NOT EXISTS transactions (
    uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    status varchar(50) NOT NULL,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
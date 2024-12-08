CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    price int NOT NULL,
    stock int NOT NULL,
    created_at timestamptz DEFAULT NOW()
);

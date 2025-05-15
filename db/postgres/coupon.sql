-- Recreate usage_type ENUM
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'usage_type') THEN
        ALTER TYPE usage_type RENAME TO usage_type_old;
        CREATE TYPE usage_type AS ENUM ('single_use', 'multi_use');
        ALTER TABLE coupons
            ALTER COLUMN usage_type DROP DEFAULT,
            ALTER COLUMN usage_type TYPE usage_type USING usage_type::text::usage_type,
            ALTER COLUMN usage_type SET DEFAULT 'single_use'::usage_type;
        DROP TYPE usage_type_old;
    ELSE
        CREATE TYPE usage_type AS ENUM ('single_use', 'multi_use');
    END IF;
END $$;

-- Recreate discount_type ENUM
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'discount_type') THEN
        ALTER TYPE discount_type RENAME TO discount_type_old;
        CREATE TYPE discount_type AS ENUM ('flat', 'percentage');
        ALTER TABLE coupons
            ALTER COLUMN discount_type DROP DEFAULT,
            ALTER COLUMN discount_type TYPE discount_type USING discount_type::text::discount_type,
            ALTER COLUMN discount_type SET DEFAULT 'flat'::discount_type;
        DROP TYPE discount_type_old;
    ELSE
        CREATE TYPE discount_type AS ENUM ('flat', 'percentage');
    END IF;
END $$;

-- Recreate discount_target ENUM
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'discount_target') THEN
        ALTER TYPE discount_target RENAME TO discount_target_old;
        CREATE TYPE discount_target AS ENUM ('delivery', 'total_order_value');
        ALTER TABLE coupons
            ALTER COLUMN discount_target DROP DEFAULT,
            ALTER COLUMN discount_target TYPE discount_target USING discount_target::text::discount_target,
            ALTER COLUMN discount_target SET DEFAULT 'total_order_value'::discount_target;
        DROP TYPE discount_target_old;
    ELSE
        CREATE TYPE discount_target AS ENUM ('delivery', 'total_order_value');
    END IF;
END $$;

-- Create coupons table
CREATE TABLE IF NOT EXISTS coupons (
    coupon_code TEXT PRIMARY KEY,
    expiry_date TIMESTAMPTZ NOT NULL,
    usage_type usage_type NOT NULL DEFAULT 'single_use'::usage_type,
    applicable_medicine_ids JSONB NOT NULL DEFAULT '[]',
    applicable_categories JSONB NOT NULL DEFAULT '[]',
    min_order_value DOUBLE PRECISION NOT NULL DEFAULT 0,
    valid_start TIMESTAMPTZ NOT NULL,
    valid_end TIMESTAMPTZ NOT NULL,
    terms_and_conditions TEXT NOT NULL DEFAULT '',
    discount_type discount_type NOT NULL DEFAULT 'flat'::discount_type,
    discount_value DOUBLE PRECISION NOT NULL DEFAULT 0,
    max_usage_per_user INTEGER NOT NULL DEFAULT 1,
    discount_target discount_target NOT NULL DEFAULT 'total_order_value'::discount_target,
    max_discount_amount DOUBLE PRECISION NOT NULL DEFAULT 0
);

-- Create coupon_usages table
CREATE TABLE IF NOT EXISTS coupon_usages (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    coupon_code TEXT NOT NULL REFERENCES coupons(coupon_code) ON DELETE CASCADE,
    used_at TIMESTAMPTZ DEFAULT NOW()
);

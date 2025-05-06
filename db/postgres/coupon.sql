CREATE TABLE coupons (
    coupon_code TEXT PRIMARY KEY,
    expiry_date TIMESTAMPTZ,
    usage_type TEXT,
    applicable_medicine_ids JSONB,
    applicable_categories JSONB,
    min_order_value DOUBLE PRECISION,
    valid_start TIMESTAMPTZ,
    valid_end TIMESTAMPTZ,
    terms_and_conditions TEXT,
    discount_type TEXT,
    discount_value DOUBLE PRECISION,
    max_usage_per_user INTEGER,
    discount_target TEXT,
    max_discount_amount DOUBLE PRECISION
);

CREATE TABLE coupon_usages (
    id SERIAL PRIMARY KEY,
    user_id TEXT,
    coupon_code TEXT REFERENCES coupons(coupon_code) ON DELETE CASCADE,
    used_at TIMESTAMPTZ DEFAULT NOW()
);

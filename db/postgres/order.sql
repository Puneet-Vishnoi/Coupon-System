CREATE TABLE coupons (
    coupon_code TEXT PRIMARY KEY,
    expiry_date TIMESTAMPTZ,
    usage_type TEXT,
    applicable_medicine_ids JSONB,
    applicable_categories JSONB,
    min_order_value FLOAT,
    valid_start TIMESTAMPTZ,
    valid_end TIMESTAMPTZ,
    terms_and_conditions TEXT,
    discount_type TEXT,
    discount_value FLOAT,
    max_usage_per_user INT,
    discount_target TEXT
);

CREATE TABLE coupon_usages (
    id SERIAL PRIMARY KEY,
    user_id TEXT,
    coupon_code TEXT REFERENCES coupons(coupon_code),
    used_at TIMESTAMPTZ
);

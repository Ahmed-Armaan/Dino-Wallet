-- Assets
INSERT INTO assets (id, code, created_at) VALUES
    (gen_random_uuid(), 'gold',  now()),
    (gen_random_uuid(), 'gem',   now()),
    (gen_random_uuid(), 'coins', now());

-- Users
INSERT INTO users (id, user_name, created_at) VALUES
    (gen_random_uuid(), 'alice', now()),
    (gen_random_uuid(), 'bob',   now());

-- System Accounts
-- treasury (allow_negative = true)
INSERT INTO accounts (id, user_id, asset_id, type, balance, allow_negative, created_at)
SELECT
    gen_random_uuid(),
    NULL,
    a.id,
    'treasury',
    0,
    true,
    now()
FROM assets a;

-- revenue  (allow_negative = false)
INSERT INTO accounts (id, user_id, asset_id, type, balance, allow_negative, created_at)
SELECT
    gen_random_uuid(),
    NULL,
    a.id,
    'revenue',
    0,
    false,
    now()
FROM assets a;

-- user wallets
INSERT INTO accounts (id, user_id, asset_id, type, balance, allow_negative, created_at)
SELECT
    gen_random_uuid(),
    u.id,
    a.id,
    'normal',
    CASE
        -- Alice balances
        WHEN u.user_name = 'alice' AND a.code = 'gold'  THEN 100
        WHEN u.user_name = 'alice' AND a.code = 'gem'   THEN 25
        WHEN u.user_name = 'alice' AND a.code = 'coins' THEN 500

        -- Bob balances
        WHEN u.user_name = 'bob'   AND a.code = 'gold'  THEN 50
        WHEN u.user_name = 'bob'   AND a.code = 'gem'   THEN 10
        WHEN u.user_name = 'bob'   AND a.code = 'coins' THEN 200

        ELSE 0
    END,
    false,
    now()
FROM users u
CROSS JOIN assets a;

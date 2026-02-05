TRUNCATE TABLE ledger_entries CASCADE;
TRUNCATE TABLE transactions CASCADE;
TRUNCATE TABLE accounts CASCADE;
TRUNCATE TABLE idempotency_keys CASCADE;

INSERT INTO accounts (owner_id, asset_type, balance, created_at, updated_at)
VALUES ('SYSTEM_TREASURY', 'GOLD_COIN', 1000000000, NOW(), NOW()),
       ('SYSTEM_TREASURY', 'DIAMOND', 1000000000, NOW(), NOW()),
       ('SYSTEM_TREASURY', 'LOYALTY_POINT', 1000000000, NOW(), NOW());

INSERT INTO accounts (owner_id, asset_type, balance, created_at, updated_at)
VALUES ('USER_1', 'GOLD_COIN', 500, NOW(), NOW()), -- Can spend immediately
       ('USER_1', 'DIAMOND', 10, NOW(), NOW()),    -- Has some premium currency
       ('USER_1', 'LOYALTY_POINT', 100, NOW(), NOW());

INSERT INTO accounts (owner_id, asset_type, balance, created_at, updated_at)
VALUES ('USER_2', 'GOLD_COIN', 100, NOW(), NOW()), -- Just enough to test a small spend
       ('USER_2', 'DIAMOND', 0, NOW(), NOW());

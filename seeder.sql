TRUNCATE status;
INSERT INTO
  status (id, name)
VALUES
  ('0', 'Inactive'), ('1', 'Active');

TRUNCATE status;
INSERT INTO
  product_status (id, name)
VALUES
  ('0', 'Inactive'), ('1', 'Active'), ('2', 'Sold');

INSERT INTO
  code_sequences (prefix, sequence, name, year)
VALUES
  ("BAG", 0, "Bags", 2024),;

-- DAYS
TRUNCATE days;
INSERT INTO
  days (id, name, name_en)
VALUES
  ('1', 'Minggu', 'Sunday'),
  ('2', 'Senin', 'Monday'),
  ('3', 'Selasa', 'Tuesday'),
  ('4', 'Rabu', 'Wednesday'),
  ('5', 'Kamis', 'Thursday'),
  ('6', 'Jumat', 'Friday'),
  ('7', 'Sabtu', 'Saturday');

-- PAYMENT TYPE
INSERT INTO
  payment_type (id, name, status_id, created_at, created_by, updated_at, updated_by)
VALUES
  ("1", 'Cash', 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1),
  ("2", 'Card', 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1),
  ("3", 'Transfer', 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1, UTC_TIMESTAMP + INTERVAL 7 HOUR, 1);

-- CARD PROVIDERS
INSERT INTO
  card_providers (id, name, created_at, created_by, updated_at, updated_by)
VALUES
  ('1', 'Visa', UTC_TIMESTAMP + INTERVAL 7 HOUR, 0, UTC_TIMESTAMP + INTERVAL 7 HOUR, 0),
  ('2', 'MasterCard', UTC_TIMESTAMP + INTERVAL 7 HOUR, 0, UTC_TIMESTAMP + INTERVAL 7 HOUR, 0),
  ('3', 'American Express', UTC_TIMESTAMP + INTERVAL 7 HOUR, 0, UTC_TIMESTAMP + INTERVAL 7 HOUR, 0),
  ('4', 'JCB', UTC_TIMESTAMP + INTERVAL 7 HOUR, 0, UTC_TIMESTAMP + INTERVAL 7 HOUR, 0),
  ('5', 'UnionPay', UTC_TIMESTAMP + INTERVAL 7 HOUR, 0, UTC_TIMESTAMP + INTERVAL 7 HOUR, 0);

TRUNCATE card_types;
INSERT INTO
  card_type (id, name)
VALUES
  ('1', 'Debit'),
  ('2', 'Credit');

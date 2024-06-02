CREATE TABLE categories (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  status_id VARCHAR(255) DEFAULT "1",

  created_at DATETIME NULL,
  created_by INT NULL,
  updated_at DATETIME NULL,
  updated_by INT NULL,
  PRIMARY KEY (id)
);
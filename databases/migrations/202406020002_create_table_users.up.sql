CREATE TABLE banks (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  status_id VARCHAR(255) DEFAULT "1",
  created_at DATETIME NULL,
  created_by INT NULL,
  updated_at DATETIME NULL,
  updated_by INT NULL,
  PRIMARY KEY (id),
  INDEX index_username (username)
);
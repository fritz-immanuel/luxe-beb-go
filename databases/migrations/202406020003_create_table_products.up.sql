CREATE TABLE products (
  id VARCHAR(255) NOT NULL,
  code VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  price DECIMAL(25,2) DEFAULT 0,
  brand_id VARCHAR(255) NOT NULL,
  category_id VARCHAR(255) NOT NULL,
  description LONGTEXT NOT NULL,
  status_id VARCHAR(255) DEFAULT "1",

  created_at DATETIME NULL,
  created_by INT NULL,
  updated_at DATETIME NULL,
  updated_by INT NULL,
  PRIMARY KEY (id),
  INDEX index_code (code),
  INDEX index_brand_id (brand_id),
  INDEX index_category_id (category_id)
);
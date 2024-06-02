// Code generated by rice embed-go; DO NOT EDIT.
package databases

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "202406020000_create_table_status.up.sql",
		FileModTime: time.Unix(1717341152, 0),

		Content: string("CREATE TABLE status (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  PRIMARY KEY (id),\r\n);"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "202406020001_create_table_banks.up.sql",
		FileModTime: time.Unix(1717341152, 0),

		Content: string("CREATE TABLE banks (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  status_id VARCHAR(255) DEFAULT \"1\",\r\n  created_at DATETIME NULL,\r\n  created_by INT NULL,\r\n  updated_at DATETIME NULL,\r\n  updated_by INT NULL,\r\n  PRIMARY KEY (id),\r\n  INDEX index_buffet_id (buffet_id),\r\n  INDEX index_business_id (business_id),\r\n  INDEX index_shift_id (shift_id)\r\n);"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "202406020002_create_table_users.up.sql",
		FileModTime: time.Unix(1717341152, 0),

		Content: string("CREATE TABLE banks (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  email VARCHAR(255) NOT NULL,\r\n  username VARCHAR(255) NOT NULL,\r\n  password VARCHAR(255) NOT NULL,\r\n  status_id VARCHAR(255) DEFAULT \"1\",\r\n  created_at DATETIME NULL,\r\n  created_by INT NULL,\r\n  updated_at DATETIME NULL,\r\n  updated_by INT NULL,\r\n  PRIMARY KEY (id),\r\n  INDEX index_username (username)\r\n);"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "202406020003_create_table_products.up.sql",
		FileModTime: time.Unix(1717345076, 0),

		Content: string("CREATE TABLE products (\r\n  id VARCHAR(255) NOT NULL,\r\n  code VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  price DECIMAL(25,2) DEFAULT 0,\r\n  brand_id VARCHAR(255) NOT NULL,\r\n  category_id VARCHAR(255) NOT NULL,\r\n  description LONGTEXT NOT NULL,\r\n  status_id VARCHAR(255) DEFAULT \"1\",\r\n\r\n  created_at DATETIME NULL,\r\n  created_by INT NULL,\r\n  updated_at DATETIME NULL,\r\n  updated_by INT NULL,\r\n  PRIMARY KEY (id),\r\n  INDEX index_code (code),\r\n  INDEX index_brand_id (brand_id),\r\n  INDEX index_category_id (category_id)\r\n);"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "202406020004_create_table_brands.up.sql",
		FileModTime: time.Unix(1717345263, 0),

		Content: string("CREATE TABLE brands (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  status_id VARCHAR(255) DEFAULT \"1\",\r\n\r\n  created_at DATETIME NULL,\r\n  created_by INT NULL,\r\n  updated_at DATETIME NULL,\r\n  updated_by INT NULL,\r\n  PRIMARY KEY (id)\r\n);"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "202406020005_create_table_categories.up.sql",
		FileModTime: time.Unix(1717345539, 0),

		Content: string("CREATE TABLE categories (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  status_id VARCHAR(255) DEFAULT \"1\",\r\n\r\n  created_at DATETIME NULL,\r\n  created_by INT NULL,\r\n  updated_at DATETIME NULL,\r\n  updated_by INT NULL,\r\n  PRIMARY KEY (id)\r\n);"),
	}
	file8 := &embedded.EmbeddedFile{
		Filename:    "202406020006_create_table_product_status.up.sql",
		FileModTime: time.Unix(1717346132, 0),

		Content: string("CREATE TABLE product_status (\r\n  id VARCHAR(255) NOT NULL,\r\n  name VARCHAR(255) NOT NULL,\r\n  PRIMARY KEY (id),\r\n);"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1717346128, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "202406020000_create_table_status.up.sql"
			file3, // "202406020001_create_table_banks.up.sql"
			file4, // "202406020002_create_table_users.up.sql"
			file5, // "202406020003_create_table_products.up.sql"
			file6, // "202406020004_create_table_brands.up.sql"
			file7, // "202406020005_create_table_categories.up.sql"
			file8, // "202406020006_create_table_product_status.up.sql"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./migrations`, &embedded.EmbeddedBox{
		Name: `./migrations`,
		Time: time.Unix(1717346128, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"202406020000_create_table_status.up.sql":         file2,
			"202406020001_create_table_banks.up.sql":          file3,
			"202406020002_create_table_users.up.sql":          file4,
			"202406020003_create_table_products.up.sql":       file5,
			"202406020004_create_table_brands.up.sql":         file6,
			"202406020005_create_table_categories.up.sql":     file7,
			"202406020006_create_table_product_status.up.sql": file8,
		},
	})
}

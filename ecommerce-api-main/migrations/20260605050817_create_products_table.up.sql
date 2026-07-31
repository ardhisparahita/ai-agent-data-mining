CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    category_id BIGINT NOT NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    image_url VARCHAR(500),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_products_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
);
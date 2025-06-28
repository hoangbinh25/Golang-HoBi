CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    idUser INT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (idUser) REFERENCES users(idUser)
);

CREATE TABLE order_items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price BIGINT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

CREATE TABLE cart_items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    idUser INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (idUser) REFERENCES users(idUser),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE reviews (
    id INT PRIMARY KEY AUTO_INCREMENT,
    idUser INT NOT NULL,
    product_id INT NOT NULL,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (idUser) REFERENCES users(idUser),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE payments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT NOT NULL,
    payment_method VARCHAR(50), -- cash, credit_card, paypal, vnpay...
    status VARCHAR(50) DEFAULT 'pending',
    paid_at DATETIME,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- Bảng lưu token xác nhận email
CREATE TABLE email_verifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    type ENUM('register', 'forgot_password') NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    used BOOLEAN DEFAULT FALSE
);

-- Thêm cột email_verified vào bảng users
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

ALTER TABLE payments
MODIFY payment_method ENUM('cash', 'credit_card', 'paypal', 'vnpay') NOT NULL;

ALTER TABLE payments
MODIFY status ENUM('pending', 'paid', 'failed') NOT NULL DEFAULT 'pending';

ALTER TABLE orders
MODIFY total_amount BIGINT;

ALTER TABLE orders
ADD column order_date DATETIME;

ALTER TABLE order_items
MODIFY price BIGINT;

select * from orders;
select * from order_items;

select * from users;
ALTER TABLE users ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT 0;




CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'rootpassword';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS warehouse_auth;
USE warehouse_auth;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role ENUM('admin', 'staff') DEFAULT 'staff',
    position VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Data Dummy (Password default: 'password123' - harus di-hash di aplikasi)
INSERT INTO users (username, password_hash, email, role, position) VALUES 
('admin_gudang', '$2a$10$ExK06zL1.7abcde...', 'admin@warehouse.com', 'admin', 'Head of Warehouse'),
('staff_joko', '$2a$10$ExK06zL1.7abcde...', 'joko@warehouse.com', 'staff', 'Inventory Clerk');
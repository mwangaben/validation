-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS gotest_two;

USE gotest;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    age INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_name (name)
    );

-- Create products table
CREATE TABLE IF NOT EXISTS products (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_code (code)
    );

-- Insert test data
INSERT IGNORE INTO users (name, email, age) VALUES
('John Doe', 'john@example.com', 30),
('Jane Smith', 'jane@example.com', 25),
('Bob Johnson', 'bob@example.com', 35);

INSERT IGNORE INTO products (code, name, price) VALUES
('P001', 'Laptop', 999.99),
('P002', 'Mouse', 29.99);
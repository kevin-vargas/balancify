-- Create database if not exists
CREATE DATABASE IF NOT EXISTS balancify;

-- Use the database
USE balancify;

-- Create Transactions table

CREATE TABLE IF NOT EXISTS Users (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Email VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS Transactions (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Date DATE,
    Amount DECIMAL(10, 2),
    UserID INT,
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);



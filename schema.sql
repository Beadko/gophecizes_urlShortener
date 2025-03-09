-- Create the URL Mappings table
CREATE DATABASE IF NOT EXISTS defaultdb;

-- Use the database
USE paths;

-- Create the url_mappings table
CREATE TABLE IF NOT EXISTS paths (
    path STRING UNIQUE NOT NULL, 
    url STRING NOT NULL 
);

INSERT INTO url_mappings (path, url) VALUES
('/urlshort', 'https://github.com/gophercises/urlshort'),
('/urlshort-final', 'https://github.com/gophercises/urlshort/tree/final');
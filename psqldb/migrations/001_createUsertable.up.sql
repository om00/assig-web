CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100) UNIQUE,
    status INT DEFAULT 1,        -- 1 = active, 0 = inactive
    blockReason TEXT DEFAULT NULL, -- Optional text, default NULL
    blockReasonCode INT DEFAULT 0, -- 0 means no block reason, otherwise code
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);



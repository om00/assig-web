CREATE TABLE Admin (
    id SERIAL PRIMARY
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phoneNumber varchar(15)  NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
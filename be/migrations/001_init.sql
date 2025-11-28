-- ENUM Definitions
CREATE TYPE donation_status AS ENUM (
    'pending',
    'verified_for_auction',
    'verified_for_donation'
);

CREATE TYPE verification_decision AS ENUM ('auction', 'donation');

CREATE TYPE auction_item_status AS ENUM ('scheduled', 'ongoing', 'finished');

CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');

-- Tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE donations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    title VARCHAR(255),
    description TEXT,
    category VARCHAR(255),
    condition VARCHAR(255),
    status donation_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE donation_photos (
    id SERIAL PRIMARY KEY,
    donation_id INT REFERENCES donations(id),
    url VARCHAR(255)
);

CREATE TABLE verifications (
    id SERIAL PRIMARY KEY,
    donation_id INT REFERENCES donations(id),
    verifier_id INT REFERENCES users(id),
    condition VARCHAR(255),
    category VARCHAR(255),
    decision verification_decision NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auction_sessions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    start_time TIMESTAMP,
    end_time TIMESTAMP
);

CREATE TABLE auction_items (
    id SERIAL PRIMARY KEY,
    donation_id INT REFERENCES donations(id),
    title VARCHAR(255),
    description TEXT,
    category VARCHAR(255),
    starting_price INT,
    status auction_item_status NOT NULL,
    session_id INT REFERENCES auction_sessions(id)
);

CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    auction_item_id INT REFERENCES auction_items(id),
    user_id INT REFERENCES users(id),
    amount INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    auction_item_id INT REFERENCES auction_items(id),
    amount INT,
    status payment_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE final_donations (
    id SERIAL PRIMARY KEY,
    donation_id INT REFERENCES donations(id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    content TEXT,
    week INT,
    created_at TIMESTAMP DEFAULT NOW()
);

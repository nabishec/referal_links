CREATE TABLE users (
    id  SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE referrals (
    id SERIAL PRIMARY KEY,
    referrer_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referral_name TEXT NOT NULL,
    referral_email TEXT UNIQUE NOT NULL,
    date DATE NOT NULL
);
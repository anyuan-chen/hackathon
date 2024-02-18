CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT,
    company TEXT,
    email TEXT,
    phone TEXT,
    role TEXT,
    salt TEXT,
    hashed_secret TEXT
);

CREATE TABLE skills (
    user_id INTEGER,
    skill TEXT,
    rating INTEGER,
    PRIMARY KEY (user_id, skill),
    FOREIGN KEY (user_id) REFERENCES "users"(id) 
);

CREATE INDEX idx_user_id ON skills(user_id);

CREATE TABLE tokens (
    id INTEGER PRIMARY KEY,
    bearer_token TEXT UNIQUE,
    expiry_time BIGINT,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES "users"(id)
);

CREATE INDEX idx_user_id_tokens ON tokens(user_id);

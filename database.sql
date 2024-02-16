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

ALTER TABLE users ADD UNIQUE (name, company);


CREATE TABLE skills (
    id INTEGER PRIMARY KEY,
    user_id TEXT,
    skill TEXT,
    rating INTEGER,
    FOREIGN KEY (user_id) REFERENCES "users"(id)
);

CREATE TABLE user_histories (
    id INTEGER PRIMARY KEY,
    date_change DATE,
    who_changed_id TEXT,
    name TEXT,
    company TEXT,
    email TEXT,
    phone TEXT,
    role TEXT,
    hashed_secret TEXT,
    FOREIGN KEY (who_changed_id) REFERENCES "users"(id)
);

CREATE TABLE skill_histories (
    id INTEGER PRIMARY KEY,
    date_change DATE,
    who_changed_id TEXT,
    user_id TEXT,
    skill TEXT,
    rating INTEGER,
    FOREIGN KEY (who_changed_id) REFERENCES "users"(id),
    FOREIGN KEY (user_id) REFERENCES "users"(id)
);

CREATE TABLE tokens (
    id INTEGER PRIMARY KEY,
    bearer_token TEXT UNIQUE,
    expiry_time BIGINT,
    FOREIGN KEY (id) REFERENCES "users"(id)
);


CREATE TABLE test (
    id INTEGER PRIMARY KEY
);
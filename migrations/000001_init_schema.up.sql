-- +migrate Up
CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    subject VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE volunteers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    location VARCHAR(255),
    skills TEXT[] DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE registrations (
    id BIGSERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    dob VARCHAR(50),
    gender VARCHAR(50),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    state_of_origin VARCHAR(255),
    state_of_residence VARCHAR(255),
    education VARCHAR(255),
    previous_office VARCHAR(255),
    interested_office VARCHAR(255),
    previous_contest VARCHAR(255),
    card_carrying_member BOOLEAN NOT NULL DEFAULT FALSE,
    party_membership_doc_link TEXT,
    motivation TEXT,
    political_understanding TEXT,
    assistance_needed TEXT[], -- array of strings
    other_support TEXT,
    preferred_communication VARCHAR(100),
    consent BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

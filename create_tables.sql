
DROP TABLE IF EXISTS registrations CASCADE;
DROP TABLE IF EXISTS volunteers CASCADE;
DROP TABLE IF EXISTS contacts CASCADE;

CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE volunteers (
    id BIGSERIAL PRIMARY KEY,
    fullname VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    skillsets TEXT[],    
    state_of_residence TEXT,
    consent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE registrations (
    id BIGSERIAL PRIMARY KEY,
    fullname VARCHAR(150) NOT NULL,
    dob TEXT,
    gender TEXT,
    email VARCHAR(150) NOT NULL,
    phone TEXT,
    state_of_origin TEXT,
    state_of_residence TEXT,
    education TEXT,
    previous_office TEXT,
    interested_office TEXT,
    previous_contest TEXT,
    card_carrying_member BOOLEAN DEFAULT FALSE,
    party_membership_doc_link TEXT,
    motivation TEXT,
    political_understanding TEXT,
    assistance_needed TEXT[], 
    other_support TEXT,
    availability TEXT[], 
    preferred_communication TEXT,
    consent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

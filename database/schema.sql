CREATE TABLE IF NOT EXISTS headquarters (
    swift_code VARCHAR(11) PRIMARY KEY CHECK (swift_code LIKE '%XXX'),
    bank_name TEXT NOT NULL,
    country_iso2 CHAR(2) NOT NULL,
    country_name TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS branches (
    swift_code VARCHAR(11) PRIMARY KEY CHECK (swift_code NOT LIKE '%XXX'),
    hq_swift_code VARCHAR(11) NOT NULL REFERENCES headquarters(swift_code) ON DELETE CASCADE,
    bank_name TEXT NOT NULL,
    country_iso2 CHAR(2) NOT NULL,
    country_name TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_branches_country ON branches(country_iso2);
CREATE INDEX IF NOT EXISTS idx_headquarters_country ON headquarters(country_iso2);

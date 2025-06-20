-- Canadian Data Sovereignty Database Extensions
-- Educational Note: These extensions support Canadian privacy law compliance
-- and data sovereignty requirements for Gander Social

-- Create databases that align with Bluesky's existing structure
-- but add Canadian-specific compliance features
CREATE DATABASE plc_dev;
CREATE DATABASE pds_dev;
CREATE DATABASE bgs_dev;
CREATE DATABASE bsky_dev;

-- Create Canadian-specific database for sovereignty features
CREATE DATABASE canada_dev;

-- Grant permissions (matches existing Bluesky development patterns)
GRANT ALL PRIVILEGES ON DATABASE plc_dev TO bsky;
GRANT ALL PRIVILEGES ON DATABASE pds_dev TO bsky;
GRANT ALL PRIVILEGES ON DATABASE bgs_dev TO bsky;
GRANT ALL PRIVILEGES ON DATABASE bsky_dev TO bsky;
GRANT ALL PRIVILEGES ON DATABASE canada_dev TO bsky;

-- Switch to Canadian database for sovereignty extensions
\c canada_dev;

-- Enable encryption extension for Canadian compliance
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create compliance audit table for PIPEDA requirements
-- Educational Note: This tracks all data access for Canadian privacy law compliance
CREATE TABLE compliance_audit (
    id SERIAL PRIMARY KEY,
    user_did TEXT,
    operation VARCHAR(50) NOT NULL,
    table_name VARCHAR(100),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    jurisdiction VARCHAR(5) DEFAULT 'CA',
    compliance_notes TEXT,
    INDEX(timestamp),
    INDEX(user_did),
    INDEX(operation)
);

COMMENT ON TABLE compliance_audit IS 'PIPEDA compliance audit trail for Canadian data operations';

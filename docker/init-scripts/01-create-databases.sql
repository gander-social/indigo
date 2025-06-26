-- AT Protocol Development Database Setup
-- This script creates all the necessary databases for the AT Protocol stack
-- Educational Note: We separate databases by service to maintain clear boundaries
-- and enable service-specific optimizations and backup strategies

-- Create the main PLC (Public Ledger of Credentials) database
-- This stores DID documents and identity resolution data
CREATE DATABASE plc_dev;
COMMENT ON DATABASE plc_dev IS 'PLC directory for decentralized identity management';

-- Create the Personal Data Server database  
-- This stores user repositories, records, and authentication data
CREATE DATABASE pds_dev;
COMMENT ON DATABASE pds_dev IS 'Personal Data Server storage for user repositories';

-- Create the Big Graph Service database
-- This stores relay metadata, federation info, and firehose state
CREATE DATABASE bgs_dev;
COMMENT ON DATABASE bgs_dev IS 'Big Graph Service for relay and firehose operations';

-- Create the Gander AppView database
-- This stores social graph data, feeds, and application-layer information
CREATE DATABASE gndr_dev;
COMMENT ON DATABASE gndr_dev IS 'Gander social application data and indexed content';

-- Create the Canadian sovereignty database
-- This stores Canadian-specific data with enhanced encryption and compliance
CREATE DATABASE canada_dev;
COMMENT ON DATABASE canada_dev IS 'Canadian data sovereign storage with compliance tracking';

-- Create a dedicated database for session and cache overflow
-- When Redis needs persistent backup or large object storage
CREATE DATABASE sessions_dev;
COMMENT ON DATABASE sessions_dev IS 'Session and cache persistence storage';

-- Grant appropriate permissions
-- Educational Note: In production, each service should have its own user
-- with minimal required permissions following the principle of least privilege
GRANT ALL PRIVILEGES ON DATABASE plc_dev TO gndr;
GRANT ALL PRIVILEGES ON DATABASE pds_dev TO gndr;
GRANT ALL PRIVILEGES ON DATABASE bgs_dev TO gndr;
GRANT ALL PRIVILEGES ON DATABASE gndr_dev TO gndr;
GRANT ALL PRIVILEGES ON DATABASE canada_dev TO gndr;
GRANT ALL PRIVILEGES ON DATABASE sessions_dev TO gndr;

-- Set up Canadian database with specific compliance settings
\c canada_dev;

-- Enable row-level security for fine-grained access control
-- This is crucial for data sovereignty compliance
ALTER DATABASE canada_dev SET row_security = on;

-- Create extension for encryption capabilities
-- Educational Note: This enables transparent data encryption at the database level
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create a compliance tracking table
-- This tracks data access for audit purposes as required by Canadian privacy laws
CREATE TABLE IF NOT EXISTS compliance_audit (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    operation VARCHAR(20) NOT NULL,
    user_did TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    jurisdiction VARCHAR(10) DEFAULT 'CA',
    data_classification VARCHAR(50) DEFAULT 'PERSONAL',
    access_reason TEXT,
    compliance_flags JSONB DEFAULT '{}'::jsonb
);

COMMENT ON TABLE compliance_audit IS 'Audit trail for Canadian data sovereignty compliance';

-- Create indexes for efficient compliance reporting
CREATE INDEX idx_compliance_timestamp ON compliance_audit(timestamp);
CREATE INDEX idx_compliance_user ON compliance_audit(user_did);
CREATE INDEX idx_compliance_operation ON compliance_audit(operation);

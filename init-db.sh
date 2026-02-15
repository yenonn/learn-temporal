#!/bin/bash
# =============================================================================
# Initialize additional databases required by Temporal
# This script runs automatically on first PostgreSQL startup
# =============================================================================

set -e

echo "Creating Temporal visibility database..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Visibility database (for workflow search & filtering)
    CREATE DATABASE temporal_visibility;
    GRANT ALL PRIVILEGES ON DATABASE temporal_visibility TO temporal;

    -- Optional: separate database for advanced visibility
    -- CREATE DATABASE temporal_visibility_advanced;
    -- GRANT ALL PRIVILEGES ON DATABASE temporal_visibility_advanced TO temporal;
EOSQL

echo "Temporal databases initialized successfully."

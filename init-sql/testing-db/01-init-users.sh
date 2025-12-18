#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$TESTING_DB_USER" --dbname "$TESTING_DB_NAME" <<-EOSQL
    -- CREATE SCHEMA IF NOT EXISTS testing;
    
    -- DO
    -- \$do\$
    -- BEGIN
    --    IF NOT EXISTS (
    --       SELECT FROM pg_catalog.pg_roles
    --       WHERE  rolname = '${TESTING_ADMIN_USER}') THEN
    --
    --       CREATE ROLE ${TESTING_ADMIN_USER} WITH LOGIN PASSWORD '${TESTING_ADMIN_PASSWORD}';
    --    END IF;
    -- END
    -- \$do\$;
    
    -- GRANT ALL PRIVILEGES ON DATABASE ${TESTING_DB_NAME} TO ${TESTING_DB_USER};
    -- GRANT ALL PRIVILEGES ON SCHEMA public TO ${TESTING_DB_USER};
EOSQL
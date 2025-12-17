#!/bin/bash
# Exit immediately if a command exits with a non-zero status.
set -e

# This script creates multiple databases and users for different services
# within a single PostgreSQL instance. It reads database names, usernames,
# and passwords from environment variables passed by docker-compose.

# Function to create a database and a dedicated user with full privileges.
# It checks for existence before creating to ensure idempotency.
# Arguments: db_name, user_name, user_password
echo "****** MULTIPLE DATABASE INITIALIZATION STARTED ******"
create_database_and_user() {
    local db_name=$1
    local user_name=$2
    local user_password=$3

    # Use psql to execute the SQL commands
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        -- Create the database if it does not exist
        SELECT 'CREATE DATABASE $db_name'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$db_name')\gexec

        -- Create the user role if it does not exist
        DO \$\$
        BEGIN
           IF NOT EXISTS (
              SELECT FROM pg_catalog.pg_roles
              WHERE  rolname = '$user_name') THEN

              CREATE ROLE $user_name WITH LOGIN PASSWORD '$user_password' CREATEROLE;
           END IF;
        END
        \$\$;

        -- Grant all privileges on the new database to the new user
        GRANT ALL PRIVILEGES ON DATABASE $db_name TO $user_name;
EOSQL
}

# --- Database Creation ---
# The environment variables below are expected to be defined in the .env file
# and passed to this container via the 'env_file' property in docker-compose.

# Create database for the KNOWLEDGE-BASE service
if [ -n "$KNOWLEDGE_BASE_DB_NAME" ] && [ -n "$KNOWLEDGE_BASE_DB_USER" ] && [ -n "$KNOWLEDGE_BASE_DB_PASSWORD" ]; then
    echo "Creating database '$KNOWLEDGE_BASE_DB_NAME' for user '$KNOWLEDGE_BASE_DB_USER'..."
    create_database_and_user "$KNOWLEDGE_BASE_DB_NAME" "$KNOWLEDGE_BASE_DB_USER" "$KNOWLEDGE_BASE_DB_PASSWORD"
else
    echo "Skipping Public Side DB creation: environment variables not set."
fi

# Create database for the TESTING service
if [ -n "$TESTING_DB_NAME" ] && [ -n "$TESTING_DB_USER" ] && [ -n "$TESTING_DB_PASSWORD" ]; then
    echo "Creating database '$TESTING_DB_NAME' for user '$TESTING_DB_USER'..."
    create_database_and_user "$TESTING_DB_NAME" "$TESTING_DB_USER" "$TESTING_DB_PASSWORD"
else
    echo "Skipping Admin Panel DB creation: environment variables not set."
fi

# Create database for the PERSONAL-ACCOUNT service
if [ -n "$PERSONAL_ACCOUNT_DB_NAME" ] && [ -n "$PERSONAL_ACCOUNT_DB_USER" ] && [ -n "$PERSONAL_ACCOUNT_DB_PASSWORD" ]; then
    echo "Creating database '$PERSONAL_ACCOUNT_DB_NAME' for user '$PERSONAL_ACCOUNT_DB_USER'..."
    create_database_and_user "$PERSONAL_ACCOUNT_DB_NAME" "$PERSONAL_ACCOUNT_DB_USER" "$PERSONAL_ACCOUNT_DB_PASSWORD"
else
    echo "Skipping Personal Account DB creation: environment variables not set."
fi

echo "****** MULTIPLE DATABASE INITIALIZATION COMPLETE ******"
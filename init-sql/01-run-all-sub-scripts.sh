#!/bin/bash
set -e

BASE_DIR="/docker-entrypoint-initdb.d/init-sql"

# Helper function to process scripts in a given directory
# Takes the directory name and the names of the environment variables to use
process_scripts_for_user() {
    local dir_name=$1
    local db_name_var=$2
    # The application ADMIN user who will own the tables
    local app_admin_user_var=$3
    local app_admin_password_var=$4

    local FULL_DIR="$BASE_DIR/$dir_name"
    
    # --- Dereference variable names to get actual values ---
    local DB_NAME="${!db_name_var}"
    local APP_ADMIN_USER="${!app_admin_user_var}"
    local APP_ADMIN_PASSWORD="${!app_admin_password_var}"

    if [ ! -d "$FULL_DIR" ]; then
        echo "Directory $FULL_DIR not found, skipping."
        return
    fi
    
    # --- Phase 1: Run shell scripts as the default postgres user ---
    # These scripts are expected to handle their own psql connections,
    # typically to create application users/roles.
    echo "---- [Phase 1] Processing .sh scripts in: $FULL_DIR ----"
    find "$FULL_DIR" -maxdepth 1 -name "*.sh" | sort | while read -r SCRIPT; do
        if [ -z "$DB_NAME" ]; then
            echo "Database name for shell scripts in $dir_name is not set, skipping $SCRIPT."
            continue
        fi
        echo "--- Executing shell script: $SCRIPT (as default user) ---"
        # chmod +x "$SCRIPT"
        bash "$SCRIPT"
    done

    # --- Phase 2: Run SQL scripts as the specific application admin user ---
    echo "---- [Phase 2] Processing .sql scripts in: $FULL_DIR for user $APP_ADMIN_USER ----"
    if [ -z "$DB_NAME" ] || [ -z "$APP_ADMIN_USER" ] || [ -z "$APP_ADMIN_PASSWORD" ]; then
        echo "Application admin credentials for SQL scripts in $dir_name are not fully set, skipping."
        return
    fi
    
    find "$FULL_DIR" -maxdepth 1 -name "*.sql" | sort | while read -r SCRIPT; do
        echo "--- Executing SQL script: $SCRIPT ---"
        PGPASSWORD="$APP_ADMIN_PASSWORD" psql -v ON_ERROR_STOP=1 --username "$APP_ADMIN_USER" --dbname "$DB_NAME" -f "$SCRIPT"
    done
}

echo "======== Running Custom Init Scripts for Specific Users ========"

# The script 00-init-dbs.sh is assumed to have run first, creating the DBs and SUPERUSERS.
# This script then runs .sh scripts in each folder (to create app users)
# and then runs .sql scripts as those newly created app users.

# For 'knowledge-base-db' folder, run .sh scripts, then run .sql scripts as KNOWLEDGE_BASE_ADMIN_USER.
process_scripts_for_user "knowledge-base-db" "KNOWLEDGE_BASE_DB_NAME" "KNOWLEDGE_BASE_DB_USER" "KNOWLEDGE_BASE_DB_PASSWORD"

# For 'testing-db' folder, we assume the main user is TESTING_DB_USER, who will own the objects.
process_scripts_for_user "testing-db" "TESTING_DB_NAME" "TESTING_DB_USER" "TESTING_DB_PASSWORD"

# For 'personal-account-db' folder, we assume the main user is PERSONAL_ACCOUNT_DB_USER.
process_scripts_for_user "personal-account-db" "PERSONAL_ACCOUNT_DB_NAME" "PERSONAL_ACCOUNT_DB_USER" "PERSONAL_ACCOUNT_DB_PASSWORD"

echo "======== Custom Init Scripts Finished ========"

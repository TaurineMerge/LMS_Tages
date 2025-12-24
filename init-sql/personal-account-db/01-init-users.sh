#!/bin/bash
set -e

# Создание ролей для базы personal_account_prod
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$PERSONAL_ACCOUNT_DB_NAME" <<-EOSQL
    -- Логирование
    DO \$\$
    BEGIN
        RAISE NOTICE 'Creating roles for personal_account_prod';
    END \$\$;

    -- Создать схему personal_account, если её нет
    CREATE SCHEMA IF NOT EXISTS personal_account;

    -- Создать основной пользователь (PERSONAL_ACCOUNT_DB_USER)
    DO \$\$
    BEGIN
       IF NOT EXISTS (
          SELECT FROM pg_catalog.pg_roles
          WHERE  rolname = '${PERSONAL_ACCOUNT_DB_USER}') THEN

          CREATE ROLE ${PERSONAL_ACCOUNT_DB_USER} WITH LOGIN PASSWORD '${PERSONAL_ACCOUNT_DB_PASSWORD}';
       END IF;
    END
    \$\$;

    -- Предоставить права
    GRANT ALL PRIVILEGES ON DATABASE ${PERSONAL_ACCOUNT_DB_NAME} TO ${PERSONAL_ACCOUNT_DB_USER};
    GRANT ALL PRIVILEGES ON SCHEMA personal_account TO ${PERSONAL_ACCOUNT_DB_USER};
    GRANT USAGE ON SCHEMA personal_account TO ${PERSONAL_ACCOUNT_DB_USER};
    -- Добавить права на public для создания функций
    GRANT ALL PRIVILEGES ON SCHEMA public TO ${PERSONAL_ACCOUNT_DB_USER};
    GRANT USAGE ON SCHEMA public TO ${PERSONAL_ACCOUNT_DB_USER};

    -- Логирование завершения
    DO \$\$
    BEGIN
        RAISE NOTICE 'Roles for personal_account_prod created successfully';
    END \$\$;
EOSQL
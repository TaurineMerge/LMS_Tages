#!/bin/sh
set -e

# Запуск MinIO в фоне
minio server /data --console-address ':9001' &
MINIO_PID=$!

# Ожидание готовности MinIO
echo "Waiting for MinIO to start..."
sleep 5

until curl -sf http://localhost:9000/minio/health/live > /dev/null 2>&1; do
    echo "Waiting for MinIO..."
    sleep 2
done

echo "MinIO started, initializing buckets..."

# Настройка mc alias
mc alias set myminio http://localhost:9000 "${MINIO_ROOT_USER}" "${MINIO_ROOT_PASSWORD}"

# Создание бакетов
mc mb myminio/"${MINIO_BUCKET_CERTIFICATES}" --ignore-existing
mc mb myminio/"${MINIO_BUCKET_SNAPSHOTS}" --ignore-existing
mc mb myminio/"${MINIO_BUCKET_IMAGES}" --ignore-existing

# Включение версионирования
mc version enable myminio/"${MINIO_BUCKET_CERTIFICATES}"

echo "Buckets created:"
mc ls myminio/

echo "MinIO initialization complete!"

# Держим MinIO запущенным
wait $MINIO_PID
#!/bin/bash

# MinIO initialization script
# Uses environment variables from .env

echo "Initializing MinIO..."

# Set MinIO alias (server will be started soon)
# Actually, since we are replacing entrypoint, we need to start server first.

# Start MinIO server in background
minio server /data &

# Wait for server to be ready
echo "Waiting for MinIO server to start..."
sleep 5

# Set MinIO alias
mc alias set myminio http://localhost:9000 "$MINIO_ACCESS_KEY" "$MINIO_SECRET_KEY"

# Create buckets if they don't exist
mc mb myminio/"$MINIO_BUCKET_CERTIFICATES" --ignore-existing
mc mb myminio/"$MINIO_BUCKET_SNAPSHOTS" --ignore-existing
mc mb myminio/"$MINIO_BUCKET_IMAGES" --ignore-existing

# Enable versioning for certificates bucket
mc version enable myminio/"$MINIO_BUCKET_CERTIFICATES"

echo "MinIO initialized with buckets: $MINIO_BUCKET_CERTIFICATES, $MINIO_BUCKET_SNAPSHOTS, $MINIO_BUCKET_IMAGES"

# Bring server to foreground
wait $SERVER_PID
FROM minio/minio:latest

# Копирование mc (скачайте mc вручную с https://dl.min.io/client/mc/release/linux-amd64/mc и положите в папку minio/)
COPY mc /usr/local/bin/mc
RUN chmod +x /usr/local/bin/mc

# Копирование скриптов
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
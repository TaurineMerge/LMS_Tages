#!/bin/bash
# Скрипт сборки MkDocs документации
# Запускать из директории personal-account/

set -e

echo "=== Сборка MkDocs документации ==="


# Собираем статику
echo "Сборка статических файлов..."
mkdocs build --clean --site-dir docs-tech

echo "=== Готово! ==="
echo "Статика в: docs-tech/"
echo "Для локальной разработки: mkdocs serve"

# Решение проблемы запуска на Windows

## Проблема
При запуске Docker контейнеров на Windows возникает ошибка:
```
/docker-entrypoint-initdb.d/00-init-dbs.sh: cannot execute: required file not found
```

Это происходит потому, что bash скрипты (`.sh`) на Windows сохраняются с окончаниями строк CRLF (Windows), а Docker/Linux требует LF (Unix).

## Решения

### Вариант 1: Автоматический скрипт (Рекомендуется)

Запустите PowerShell скрипт перед первым запуском:

```powershell
.\convert-line-endings.ps1
```

Затем запустите Docker:

```powershell
docker-compose down -v
docker-compose up -d --build
```

### Вариант 2: Git автоматически (Лучшее долгосрочное решение)

Файл `.gitattributes` уже настроен в проекте. При клонировании репозитория на новом ПК выполните:

```bash
# Обновить все файлы согласно .gitattributes
git add --renormalize .
git status
```

Если у вас уже склонирован репозиторий, выполните:

```bash
# Пересохранить файлы с правильными окончаниями строк
git rm --cached -r .
git reset --hard
```

### Вариант 3: Вручную через VS Code

1. Откройте любой `.sh` файл в VS Code
2. Внизу справа найдите `CRLF`
3. Нажмите на него и выберите `LF`
4. Сохраните файл
5. Повторите для всех `.sh` файлов в папке `init-sql/`

### Вариант 4: Использование dos2unix (если установлен)

```bash
# Установить dos2unix через WSL или Git Bash
find ./init-sql -name "*.sh" -exec dos2unix {} \;
```

## Проверка

После конвертации проверьте статус контейнеров:

```powershell
docker ps
```

Все контейнеры должны быть в статусе `Up`, а не `Restarting`.

## Почему это происходит?

- **Windows** использует два символа для новой строки: `\r\n` (CRLF)
- **Linux/Unix** использует один символ: `\n` (LF)
- Bash интерпретатор не распознаёт `\r` и считает скрипт невалидным
- `.gitattributes` решает эту проблему автоматически при работе через Git

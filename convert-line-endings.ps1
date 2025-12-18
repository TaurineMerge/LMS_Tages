# Скрипт для конвертации окончаний строк из CRLF (Windows) в LF (Unix/Linux)
# Используется для исправления bash скриптов перед запуском в Docker

Write-Host "🔄 Конвертирую .sh файлы из CRLF в LF формат..." -ForegroundColor Cyan

$scriptPath = $PSScriptRoot
if ([string]::IsNullOrEmpty($scriptPath)) {
    $scriptPath = (Get-Location).Path
}

# Ищем все .sh файлы в директории init-sql и её поддиректориях
$shFiles = Get-ChildItem -Path "$scriptPath\init-sql" -Recurse -Filter "*.sh"

if ($shFiles.Count -eq 0) {
    Write-Host "⚠️  Не найдено .sh файлов для конвертации" -ForegroundColor Yellow
    exit 0
}

Write-Host "📁 Найдено файлов: $($shFiles.Count)" -ForegroundColor Green

foreach ($file in $shFiles) {
    try {
        # Читаем содержимое файла
        $content = Get-Content $file.FullName -Raw
        
        # Заменяем CRLF на LF
        $content = $content -replace "`r`n", "`n"
        
        # Записываем обратно в файл с UTF-8 без BOM
        [System.IO.File]::WriteAllText($file.FullName, $content, [System.Text.UTF8Encoding]::new($false))
        
        Write-Host "  ✅ $($file.Name)" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Ошибка при обработке $($file.Name): $_" -ForegroundColor Red
    }
}

Write-Host "`n✨ Конвертация завершена!" -ForegroundColor Cyan
Write-Host "💡 Теперь можно запустить: docker-compose down -v && docker-compose up -d --build" -ForegroundColor Yellow

# AdminPanel

# Как запустить Docker с сервером

# Запуск через Docker

1. **Соберите Docker образ:**
```bash
docker build -t adminpanel .
```

2. **Запустите контейнер:**
```bash
docker run -p 4000:4000 adminpanel
```

3. **Откройте браузер:**
```
http://localhost:4000
```

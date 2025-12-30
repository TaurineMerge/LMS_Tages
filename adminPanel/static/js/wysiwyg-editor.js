/**
 * TinyMCE Editor для редактирования HTML контента уроков
 * Интеграция с MinIO/S3 для загрузки изображений
 */

// Инициализация TinyMCE с интеграцией загрузки изображений
function initTinyMCE(selector) {
    tinymce.init({
        selector: selector,
        
        // Основные настройки
        height: 600,
        menubar: true,
        
        // Плагины
        plugins: [
            'advlist', 'autolink', 'lists', 'link', 'image', 'charmap', 'preview',
            'anchor', 'searchreplace', 'visualblocks', 'code', 'fullscreen',
            'insertdatetime', 'media', 'table', 'help', 'wordcount', 'codesample'
        ],
        
        // Панель инструментов
        toolbar: 'undo redo | blocks | ' +
            'bold italic underline strikethrough | forecolor backcolor | ' +
            'alignleft aligncenter alignright alignjustify | ' +
            'bullist numlist outdent indent | ' +
            'link image media table codesample | ' +
            'removeformat code fullscreen | help',
        
        // Настройки контента
        content_style: 'body { font-family: Arial, sans-serif; font-size: 14px; }',
        
        // Автоматическое изменение размера
        autoresize_bottom_margin: 50,
        autoresize_overflow_padding: 50,
        
        // Настройки изображений
        image_advtab: true,
        image_caption: true,
        image_title: true,
        
        // Настройки загрузки файлов
        automatic_uploads: true,
        file_picker_types: 'image',
        
        // Обработчик загрузки изображений
        images_upload_handler: function (blobInfo, progress) {
            return new Promise((resolve, reject) => {
                const formData = new FormData();
                formData.append('image', blobInfo.blob(), blobInfo.filename());
                
                const xhr = new XMLHttpRequest();
                xhr.open('POST', '/admin/api/v1/upload/image', true);
                
                // Обработка прогресса загрузки
                xhr.upload.onprogress = function(e) {
                    if (e.lengthComputable) {
                        progress(e.loaded / e.total * 100);
                    }
                };
                
                xhr.onload = function() {
                    if (xhr.status === 200) {
                        try {
                            const json = JSON.parse(xhr.responseText);
                            if (json.image_url) {
                                resolve(json.image_url);
                            } else {
                                reject('Ошибка: URL изображения не получен');
                            }
                        } catch (e) {
                            reject('Ошибка парсинга ответа сервера: ' + e.message);
                        }
                    } else {
                        try {
                            const json = JSON.parse(xhr.responseText);
                            reject('Ошибка загрузки: ' + (json.error || xhr.statusText));
                        } catch (e) {
                            reject('HTTP Error: ' + xhr.status);
                        }
                    }
                };
                
                xhr.onerror = function() {
                    reject('Ошибка сети при загрузке изображения');
                };
                
                xhr.send(formData);
            });
        },
        
        // Дополнительная обработка вставки изображений по URL
        file_picker_callback: function(callback, value, meta) {
            if (meta.filetype === 'image') {
                // Создаем диалог выбора способа добавления изображения
                const input = document.createElement('input');
                input.setAttribute('type', 'file');
                input.setAttribute('accept', 'image/*');
                
                input.onchange = function() {
                    const file = this.files[0];
                    
                    if (file) {
                        // Проверка типа файла
                        const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
                        if (!validTypes.includes(file.type)) {
                            tinymce.activeEditor.notificationManager.open({
                                text: 'Неверный тип файла. Разрешены только: JPEG, PNG, GIF, WEBP',
                                type: 'error',
                                timeout: 5000
                            });
                            return;
                        }
                        
                        // Проверка размера файла (10 МБ)
                        if (file.size > 10 * 1024 * 1024) {
                            tinymce.activeEditor.notificationManager.open({
                                text: 'Файл слишком большой. Максимальный размер: 10 МБ',
                                type: 'error',
                                timeout: 5000
                            });
                            return;
                        }
                        
                        // Загружаем файл
                        const formData = new FormData();
                        formData.append('image', file);
                        
                        const notification = tinymce.activeEditor.notificationManager.open({
                            text: 'Загрузка изображения...',
                            type: 'info',
                            closeButton: false,
                            progressBar: true
                        });
                        
                        fetch('/admin/api/v1/upload/image', {
                            method: 'POST',
                            body: formData
                        })
                        .then(response => response.json())
                        .then(data => {
                            notification.close();
                            
                            if (data.image_url) {
                                callback(data.image_url, { title: file.name, alt: file.name });
                                tinymce.activeEditor.notificationManager.open({
                                    text: 'Изображение успешно загружено',
                                    type: 'success',
                                    timeout: 3000
                                });
                            } else {
                                throw new Error('URL изображения не получен');
                            }
                        })
                        .catch(error => {
                            notification.close();
                            tinymce.activeEditor.notificationManager.open({
                                text: 'Ошибка загрузки: ' + error.message,
                                type: 'error',
                                timeout: 5000
                            });
                        });
                    }
                };
                
                input.click();
            }
        },
        
        // Настройки таблиц
        table_appearance_options: false,
        table_grid: false,
        table_default_attributes: {
            border: '1'
        },
        table_default_styles: {
            'border-collapse': 'collapse',
            'width': '100%'
        },
        
        // Настройки кода
        codesample_languages: [
            { text: 'HTML/XML', value: 'markup' },
            { text: 'JavaScript', value: 'javascript' },
            { text: 'CSS', value: 'css' },
            { text: 'PHP', value: 'php' },
            { text: 'Python', value: 'python' },
            { text: 'Java', value: 'java' },
            { text: 'C', value: 'c' },
            { text: 'C++', value: 'cpp' },
            { text: 'C#', value: 'csharp' },
            { text: 'Go', value: 'go' },
            { text: 'Ruby', value: 'ruby' },
            { text: 'SQL', value: 'sql' },
            { text: 'JSON', value: 'json' }
        ],
        
        // Обработка вставки контента
        paste_preprocess: function(plugin, args) {
            // Очистка нежелательных стилей при вставке
            console.log('Вставка контента:', args.content);
        }
    });
}

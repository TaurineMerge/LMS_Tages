/**
 * WYSIWYG Editor для редактирования HTML контента уроков
 * Простой редактор с базовым функционалом форматирования
 */

class WYSIWYGEditor {
    constructor(editorId, toolbarId) {
        this.editor = document.getElementById(editorId);
        this.toolbar = document.getElementById(toolbarId);
        this.hiddenInput = null;
        
        if (!this.editor || !this.toolbar) {
            console.error('Editor or toolbar not found');
            return;
        }
        
        this.init();
    }
    
    init() {
        // Делаем редактор редактируемым
        this.editor.contentEditable = true;
        this.editor.classList.add('wysiwyg-editor__content');
        
        // Создаем скрытое поле для отправки HTML
        this.createHiddenInput();
        
        // Добавляем обработчики для панели инструментов
        this.setupToolbar();
        
        // Обновляем скрытое поле при изменении контента
        this.editor.addEventListener('input', () => this.updateHiddenInput());
        
        // Обрабатываем вставку текста (очищаем форматирование)
        this.editor.addEventListener('paste', (e) => this.handlePaste(e));
        
        // Добавляем обработчики для существующих изображений
        this.setupExistingImages();
        
        // Закрываем выпадающие меню при клике вне их
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.wysiwyg-menu') && !e.target.closest('.wysiwyg-dropdown')) {
                document.querySelectorAll('.wysiwyg-menu__dropdown.show, .wysiwyg-dropdown__list.show').forEach(d => {
                    d.classList.remove('show');
                });
            }
        });
        
        // Инициализируем начальный контент
        this.updateHiddenInput();
    }
    
    setupExistingImages() {
        const images = this.editor.querySelectorAll('img');
        images.forEach(img => {
            // Устанавливаем стили для изображений
            if (!img.style.maxWidth) img.style.maxWidth = '100%';
            if (!img.style.height) img.style.height = 'auto';
            if (!img.style.display) img.style.display = 'block';
            if (!img.style.margin) img.style.margin = '10px 0';
            img.style.cursor = 'pointer';
            img.className = 'wysiwyg-image';
            
            // Добавляем обработчик клика для редактирования
            img.onclick = (e) => {
                e.preventDefault();
                this.editImageSize(img);
            };
            
            // Обработчик ошибки загрузки
            img.onerror = () => {
                img.alt = '⚠️ Изображение не загружено';
                img.style.border = '2px dashed #ff0000';
                img.style.padding = '20px';
                img.style.background = '#fff3cd';
            };
        });
    }
    
    createHiddenInput() {
        const form = this.editor.closest('form');
        if (!form) return;
        
        this.hiddenInput = document.createElement('input');
        this.hiddenInput.type = 'hidden';
        this.hiddenInput.name = 'html_content';
        form.appendChild(this.hiddenInput);
    }
    
    updateHiddenInput() {
        if (this.hiddenInput) {
            this.hiddenInput.value = this.editor.innerHTML;
        }
    }
    
    setupToolbar() {
        // Создаем структуру меню
        const menuStructure = [
            {
                type: 'menu',
                label: 'Правка',
                items: [
                    { label: 'Отменить', icon: '↶', action: 'undo', shortcut: 'Ctrl+Z' },
                    { label: 'Повторить', icon: '↷', action: 'redo', shortcut: 'Ctrl+Y' },
                    { type: 'separator' },
                    { label: 'Вырезать', icon: '✂', action: 'cut', shortcut: 'Ctrl+X' },
                    { label: 'Копировать', icon: '📋', action: 'copy', shortcut: 'Ctrl+C' },
                    { label: 'Вставить', icon: '📄', action: 'paste', shortcut: 'Ctrl+V' },
                    { label: 'Вставить как текст', icon: '📃', action: 'pasteAsText' },
                    { type: 'separator' },
                    { label: 'Выделить все', icon: '⬚', action: 'selectAll', shortcut: 'Ctrl+A' },
                    { label: 'Найти и заменить', icon: '🔍', action: 'findReplace', shortcut: 'Ctrl+F' },
                ]
            },
            {
                type: 'menu',
                label: 'Вид',
                items: [
                    { label: 'Исходный код', icon: '<>', action: 'sourceCode' },
                    { label: 'Предпросмотр', icon: '👁', action: 'preview' },
                    { label: 'Полноэкранный режим', icon: '⛶', action: 'fullscreen' },
                ]
            },
            {
                type: 'menu',
                label: 'Вставка',
                items: [
                    { label: 'Изображение', icon: '🖼', action: 'insertImage' },
                    { label: 'Ссылка', icon: '🔗', action: 'createLink' },
                    { label: 'Таблица', icon: '⊞', action: 'insertTable' },
                    { label: 'Код', icon: '</>', action: 'insertCode' },
                    { label: 'Горизонтальная линия', icon: '―', action: 'insertHR' },
                ]
            },
            { type: 'separator' },
            { type: 'button', id: 'bold', command: 'bold', icon: '<b>B</b>', title: 'Жирный (Ctrl+B)' },
            { type: 'button', id: 'italic', command: 'italic', icon: '<i>I</i>', title: 'Курсив (Ctrl+I)' },
            { type: 'button', id: 'underline', command: 'underline', icon: '<u>U</u>', title: 'Подчеркнутый (Ctrl+U)' },
            { type: 'button', id: 'strikethrough', command: 'strikeThrough', icon: '<s>S</s>', title: 'Зачеркнутый' },
            { type: 'separator' },
            {
                type: 'dropdown',
                label: 'Формат',
                items: [
                    { label: 'Параграф', value: 'p', action: 'formatBlock' },
                    { label: 'Заголовок 1', value: 'h1', action: 'formatBlock' },
                    { label: 'Заголовок 2', value: 'h2', action: 'formatBlock' },
                    { label: 'Заголовок 3', value: 'h3', action: 'formatBlock' },
                    { label: 'Заголовок 4', value: 'h4', action: 'formatBlock' },
                    { label: 'Цитата', value: 'blockquote', action: 'formatBlock' },
                    { label: 'Код', value: 'pre', action: 'formatBlock' },
                ]
            },
            {
                type: 'dropdown',
                label: 'Шрифт',
                items: [
                    { label: 'Arial', value: 'Arial', action: 'fontName' },
                    { label: 'Times New Roman', value: 'Times New Roman', action: 'fontName' },
                    { label: 'Courier New', value: 'Courier New', action: 'fontName' },
                    { label: 'Georgia', value: 'Georgia', action: 'fontName' },
                    { label: 'Verdana', value: 'Verdana', action: 'fontName' },
                ]
            },
            {
                type: 'dropdown',
                label: 'Размер',
                items: [
                    { label: '8px', value: '8px', action: 'fontSize' },
                    { label: '10px', value: '10px', action: 'fontSize' },
                    { label: '12px', value: '12px', action: 'fontSize' },
                    { label: '14px', value: '14px', action: 'fontSize' },
                    { label: '16px', value: '16px', action: 'fontSize' },
                    { label: '18px', value: '18px', action: 'fontSize' },
                    { label: '20px', value: '20px', action: 'fontSize' },
                    { label: '22px', value: '22px', action: 'fontSize' },
                    { label: '24px', value: '24px', action: 'fontSize' },
                    { label: '28px', value: '28px', action: 'fontSize' },
                    { label: '32px', value: '32px', action: 'fontSize' },
                    { label: '36px', value: '36px', action: 'fontSize' },
                    { label: '40px', value: '40px', action: 'fontSize' },
                ]
            },
            // ...existing code...
            { type: 'separator' },
            { type: 'button', id: 'ul', command: 'insertUnorderedList', icon: '•', title: 'Маркированный список' },
            { type: 'button', id: 'ol', command: 'insertOrderedList', icon: '1.', title: 'Нумерованный список' },
            { type: 'separator' },
            { type: 'button', id: 'alignLeft', command: 'justifyLeft', icon: '⬅', title: 'По левому краю' },
            { type: 'button', id: 'alignCenter', command: 'justifyCenter', icon: '↔', title: 'По центру' },
            { type: 'button', id: 'alignRight', command: 'justifyRight', icon: '➡', title: 'По правому краю' },
            { type: 'button', id: 'alignJustify', command: 'justifyFull', icon: '≡', title: 'По ширине' },
            { type: 'separator' },
            { type: 'button', id: 'clear', command: 'removeFormat', icon: '✕', title: 'Очистить форматирование' },
        ];
        
        this.renderToolbar(menuStructure);
    }
    
    renderToolbar(structure) {
        structure.forEach(item => {
            if (item.type === 'separator') {
                const separator = document.createElement('span');
                separator.className = 'wysiwyg-toolbar__separator';
                this.toolbar.appendChild(separator);
            }
            else if (item.type === 'menu') {
                this.createMenu(item);
            }
            else if (item.type === 'dropdown') {
                this.createDropdown(item);
            }
            else if (item.type === 'button') {
                this.createButton(item);
            }
        });
    }
    
    createMenu(menuData) {
        const menuContainer = document.createElement('div');
        menuContainer.className = 'wysiwyg-menu';
        
        const menuButton = document.createElement('button');
        menuButton.type = 'button';
        menuButton.className = 'wysiwyg-menu__button';
        menuButton.textContent = menuData.label;
        
        const menuDropdown = document.createElement('div');
        menuDropdown.className = 'wysiwyg-menu__dropdown';
        
        menuData.items.forEach(item => {
            if (item.type === 'separator') {
                const sep = document.createElement('div');
                sep.className = 'wysiwyg-menu__separator';
                menuDropdown.appendChild(sep);
            } else {
                const menuItem = document.createElement('button');
                menuItem.type = 'button';
                menuItem.className = 'wysiwyg-menu__item';
                menuItem.innerHTML = `
                    <span class="wysiwyg-menu__item-icon">${item.icon}</span>
                    <span class="wysiwyg-menu__item-label">${item.label}</span>
                    ${item.shortcut ? `<span class="wysiwyg-menu__item-shortcut">${item.shortcut}</span>` : ''}
                `;
                menuItem.addEventListener('click', (e) => {
                    e.preventDefault();
                    this.executeAction(item.action);
                    menuDropdown.classList.remove('show');
                });
                menuDropdown.appendChild(menuItem);
            }
        });
        
        menuButton.addEventListener('click', (e) => {
            e.stopPropagation();
            // Закрываем все другие меню
            document.querySelectorAll('.wysiwyg-menu__dropdown.show').forEach(d => {
                if (d !== menuDropdown) d.classList.remove('show');
            });
            menuDropdown.classList.toggle('show');
        });
        
        menuContainer.appendChild(menuButton);
        menuContainer.appendChild(menuDropdown);
        this.toolbar.appendChild(menuContainer);
    }
    
    createDropdown(dropdownData) {
        const dropdownContainer = document.createElement('div');
        dropdownContainer.className = 'wysiwyg-dropdown';
        
        const dropdownButton = document.createElement('button');
        dropdownButton.type = 'button';
        dropdownButton.className = 'wysiwyg-dropdown__button';
        dropdownButton.innerHTML = `${dropdownData.label} <span class="wysiwyg-dropdown__arrow">▼</span>`;
        
        const dropdownList = document.createElement('div');
        dropdownList.className = 'wysiwyg-dropdown__list';
        
        dropdownData.items.forEach(item => {
            const listItem = document.createElement('button');
            listItem.type = 'button';
            listItem.className = 'wysiwyg-dropdown__item';
            listItem.textContent = item.label;
            listItem.addEventListener('click', (e) => {
                e.preventDefault();
                if (item.action === 'formatBlock') {
                    document.execCommand('formatBlock', false, item.value);
                } else if (item.action === 'fontName') {
                    document.execCommand('fontName', false, item.value);
                } else if (item.action === 'fontSize') {
                    // Для размеров в px используем inline style
                    const selection = window.getSelection();
                    if (selection.rangeCount > 0) {
                        const range = selection.getRangeAt(0);
                        const span = document.createElement('span');
                        span.style.fontSize = item.value;
                        range.surroundContents(span);
                    }
                }
                this.updateHiddenInput();
                dropdownList.classList.remove('show');
            });
            dropdownList.appendChild(listItem);
        });
        
        dropdownButton.addEventListener('click', (e) => {
            e.stopPropagation();
            document.querySelectorAll('.wysiwyg-dropdown__list.show, .wysiwyg-menu__dropdown.show').forEach(d => {
                if (d !== dropdownList) d.classList.remove('show');
            });
            dropdownList.classList.toggle('show');
        });
        
        dropdownContainer.appendChild(dropdownButton);
        dropdownContainer.appendChild(dropdownList);
        this.toolbar.appendChild(dropdownContainer);
    }
    
    createButton(btnData) {
        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'wysiwyg-toolbar__btn';
        button.id = `btn-${btnData.id}`;
        button.innerHTML = btnData.icon;
        button.title = btnData.title;
        
        // Добавляем color picker для кнопок цвета
        if (btnData.colorPicker) {
            button.style.position = 'relative';
            
            const colorInput = document.createElement('input');
            colorInput.type = 'color';
            colorInput.className = 'wysiwyg-color-picker';
            colorInput.style.cssText = 'position: absolute; opacity: 0; width: 100%; height: 100%; cursor: pointer; left: 0; top: 0;';
            
            colorInput.addEventListener('change', (e) => {
                const color = e.target.value;
                if (btnData.action === 'textColor') {
                    document.execCommand('foreColor', false, color);
                } else if (btnData.action === 'bgColor') {
                    document.execCommand('backColor', false, color);
                }
                this.updateHiddenInput();
            });
            
            button.appendChild(colorInput);
        } else {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                if (btnData.action) {
                    this.executeAction(btnData.action);
                } else if (btnData.command) {
                    document.execCommand(btnData.command, false, null);
                    this.updateHiddenInput();
                }
            });
        }
        
        this.toolbar.appendChild(button);
    }
    
    executeAction(action) {
        switch(action) {
            case 'undo':
                document.execCommand('undo');
                break;
            case 'redo':
                document.execCommand('redo');
                break;
            case 'cut':
                document.execCommand('cut');
                break;
            case 'copy':
                document.execCommand('copy');
                break;
            case 'paste':
                document.execCommand('paste');
                break;
            case 'pasteAsText':
                this.pasteAsText();
                break;
            case 'selectAll':
                document.execCommand('selectAll');
                break;
            case 'findReplace':
                this.showFindReplace();
                break;
            case 'sourceCode':
                this.toggleSourceCode();
                break;
            case 'preview':
                this.showPreview();
                break;
            case 'fullscreen':
                this.toggleFullscreen();
                break;
            case 'wordCount':
                this.showWordCount();
                break;
            case 'insertImage':
                this.insertImageWithDialog();
                break;
            case 'createLink':
                this.insertLink();
                break;
            case 'insertTable':
                this.insertTable();
                break;
            case 'insertCode':
                this.insertCodeBlock();
                break;
            case 'insertHR':
                document.execCommand('insertHorizontalRule');
                break;
        }
        this.updateHiddenInput();
    }
    
    insertImageWithDialog() {
        const url = prompt('Введите URL изображения:\n\n⚠️ Внимание:\n- Используйте прямые ссылки на изображения (заканчиваются на .jpg, .png, .gif)\n- Pinterest и социальные сети могут блокировать вставку\n- Рекомендуется загружать изображения на imgur.com или imgbb.com');
        
        if (!url) return;
        
        // Запрашиваем размер
        const width = prompt('Введите ширину изображения в пикселях (например: 600):', '600');
        
        // Создаем изображение
        const img = document.createElement('img');
        img.src = url;
        img.alt = 'Image';
        img.style.maxWidth = '100%';
        img.style.height = 'auto';
        img.style.display = 'block';
        img.style.margin = '10px 0';
        img.style.cursor = 'pointer';
        
        if (width && !isNaN(width) && width > 0) {
            img.style.width = width + 'px';
        }
        
        // Добавляем класс для идентификации
        img.className = 'wysiwyg-image';
        
        // Обработчик ошибки загрузки
        img.onerror = () => {
            img.alt = '⚠️ Изображение не загружено. Возможные причины:\n- Неверная ссылка\n- Сайт блокирует внешние запросы\n- Изображение удалено';
            img.style.border = '2px dashed #ff0000';
            img.style.padding = '20px';
            img.style.background = '#fff3cd';
            img.style.color = '#856404';
            img.style.fontSize = '14px';
            img.style.whiteSpace = 'pre-wrap';
            img.removeAttribute('src');
        };
        
        // Добавляем обработчик клика для редактирования размера
        img.onclick = (e) => {
            e.preventDefault();
            this.editImageSize(img);
        };
        
        // Вставляем изображение в редактор
        const selection = window.getSelection();
        if (selection.rangeCount > 0) {
            const range = selection.getRangeAt(0);
            range.deleteContents();
            range.insertNode(img);
            
            // Добавляем перенос строки после изображения
            const br = document.createElement('br');
            img.parentNode.insertBefore(br, img.nextSibling);
            
            // Перемещаем курсор после изображения
            range.setStartAfter(br);
            range.collapse(true);
            selection.removeAllRanges();
            selection.addRange(range);
        }
        
        this.updateHiddenInput();
    }
    
    editImageSize(img) {
        const currentWidth = img.style.width ? parseInt(img.style.width) : img.naturalWidth || 600;
        const newWidth = prompt('Введите новую ширину изображения в пикселях (например: 800):', currentWidth);
        
        if (newWidth && !isNaN(newWidth) && newWidth > 0) {
            img.style.width = newWidth + 'px';
            this.updateHiddenInput();
        }
    }
    
    handlePaste(e) {
        e.preventDefault();
        
        // Получаем только текст без форматирования
        const text = e.clipboardData.getData('text/plain');
        
        // Вставляем как текст
        document.execCommand('insertText', false, text);
    }
    
    getContent() {
        return this.editor.innerHTML;
    }
    
    setContent(html) {
        this.editor.innerHTML = html;
        this.updateHiddenInput();
    }
    
    clear() {
        this.editor.innerHTML = '';
        this.updateHiddenInput();
    }
    
    // Новые методы
    
    pasteAsText() {
        const text = prompt('Вставьте текст:');
        if (text) {
            document.execCommand('insertText', false, text);
        }
    }
    
    showFindReplace() {
        const searchText = prompt('Найти:');
        if (!searchText) return;
        
        const replaceText = prompt('Заменить на:');
        if (replaceText === null) return;
        
        const content = this.editor.innerHTML;
        const regex = new RegExp(searchText, 'g');
        this.editor.innerHTML = content.replace(regex, replaceText);
        this.updateHiddenInput();
        alert(`Заменено ${(content.match(regex) || []).length} вхождений`);
    }
    
    toggleSourceCode() {
        if (!this.editor.classList.contains('source-code-mode')) {
            // Включаем режим исходного кода
            const html = this.editor.innerHTML;
            this.originalHTML = html;
            const formattedHtml = this.formatHTML(html);
            this.editor.textContent = formattedHtml;
            this.editor.contentEditable = 'false';
            this.editor.classList.add('source-code-mode');
        }
    }
    
    formatHTML(html) {
        // Убираем лишние пробелы и переносы
        html = html.trim();
        
        // Простое форматирование HTML с отступами
        let formatted = '';
        let indent = 0;
        const tab = '  '; // 2 пробела для отступа
        
        // Разбиваем на теги
        const tags = html.match(/<[^>]+>|[^<]+/g) || [];
        
        tags.forEach(tag => {
            if (tag.match(/^<\/\w/)) {
                // Закрывающий тег
                indent = Math.max(0, indent - 1);
                formatted += tab.repeat(indent) + tag + '\n';
            } else if (tag.match(/^<\w[^>]*[^\/]>$/)) {
                // Открывающий тег (не самозакрывающийся)
                formatted += tab.repeat(indent) + tag + '\n';
                indent++;
            } else if (tag.match(/^<\w[^>]*\/>$/)) {
                // Самозакрывающийся тег
                formatted += tab.repeat(indent) + tag + '\n';
            } else if (tag.trim()) {
                // Текстовое содержимое
                formatted += tab.repeat(indent) + tag.trim() + '\n';
            }
        });
        
        return formatted.trim();
    }
    
    highlightHTML(html) {
        // Подсветка синтаксиса HTML
        let highlighted = html
            // Экранируем HTML
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            // Подсвечиваем теги
            .replace(/(&lt;\/?)(\w+)(.*?)(&gt;)/g, '<span class="html-tag">$1</span><span class="html-tag-name">$2</span><span class="html-attr">$3</span><span class="html-tag">$4</span>')
            // Подсвечиваем атрибуты
            .replace(/(\w+)=(".*?"|'.*?')/g, '<span class="html-attr-name">$1</span>=<span class="html-attr-value">$2</span>');
        
        return `<pre class="html-source"><code>${highlighted}</code></pre>`;
    }
    
    showPreview() {
        // Если сейчас исходный код — выходим из него и возвращаем визуальный режим
        if (this.editor.classList.contains('source-code-mode')) {
            if (this.originalHTML) {
                this.editor.innerHTML = this.originalHTML;
                this.originalHTML = null;
            }
            this.editor.contentEditable = 'true';
            this.editor.classList.remove('source-code-mode');
            this.setupExistingImages();
        }
    }
    
    toggleFullscreen() {
        const editorContainer = this.editor.closest('.wysiwyg-editor');
        if (!editorContainer.classList.contains('fullscreen')) {
            editorContainer.classList.add('fullscreen');
            document.body.style.overflow = 'hidden';
        } else {
            editorContainer.classList.remove('fullscreen');
            document.body.style.overflow = '';
        }
    }
    
    showWordCount() {
        const text = this.editor.textContent || '';
        const words = text.trim().split(/\s+/).filter(w => w.length > 0).length;
        const chars = text.length;
        const charsNoSpaces = text.replace(/\s/g, '').length;
        
        alert(`
Статистика текста:
────────────────
Слов: ${words}
Символов (с пробелами): ${chars}
Символов (без пробелов): ${charsNoSpaces}
        `.trim());
    }
    
    insertLink() {
        const url = prompt('Введите URL:');
        if (!url) return;
        
        const selection = window.getSelection();
        const selectedText = selection.toString();
        
        if (selectedText) {
            document.execCommand('createLink', false, url);
        } else {
            const text = prompt('Введите текст ссылки:', url);
            if (text) {
                const link = document.createElement('a');
                link.href = url;
                link.textContent = text;
                link.target = '_blank';
                
                const range = selection.getRangeAt(0);
                range.deleteContents();
                range.insertNode(link);
            }
        }
        this.updateHiddenInput();
    }
    
    insertTable() {
        const rows = parseInt(prompt('Количество строк:', '3'));
        const cols = parseInt(prompt('Количество столбцов:', '3'));
        
        if (!rows || !cols || rows < 1 || cols < 1) return;
        
        let tableHTML = '<table border="1" style="border-collapse: collapse; width: 100%; margin: 10px 0;">';
        
        for (let i = 0; i < rows; i++) {
            tableHTML += '<tr>';
            for (let j = 0; j < cols; j++) {
                tableHTML += '<td style="padding: 8px; border: 1px solid #ddd;"></td>';
            }
            tableHTML += '</tr>';
        }
        
        tableHTML += '</table>';
        
        document.execCommand('insertHTML', false, tableHTML);
        this.updateHiddenInput();
    }
    
    insertCodeBlock() {
        const code = prompt('Вставьте код:');
        if (!code) return;
        
        const pre = document.createElement('pre');
        pre.style.background = '#2d2d2d';
        pre.style.color = '#f8f8f2';
        pre.style.padding = '16px';
        pre.style.borderRadius = '4px';
        pre.style.overflow = 'auto';
        pre.textContent = code;
        
        const selection = window.getSelection();
        if (selection.rangeCount > 0) {
            const range = selection.getRangeAt(0);
            range.deleteContents();
            range.insertNode(pre);
        }
        
        this.updateHiddenInput();
    }
}

// Инициализация редактора при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    const editorElement = document.getElementById('lesson-content-editor');
    if (editorElement) {
        window.lessonEditor = new WYSIWYGEditor('lesson-content-editor', 'wysiwyg-toolbar');
    }
});

/**
 * WYSIWYG Editor для редактирования HTML контента уроков
 * Простой редактор с базовым функционалом форматирования
 */

class WYSIWYGEditor {
        wrapSelectionWithSpan(styleObj) {
            const selection = window.getSelection();
            if (!selection.rangeCount) return;
            const range = selection.getRangeAt(0);
            // Собираем все текстовые узлы в выделении
            const walker = document.createTreeWalker(
                range.commonAncestorContainer,
                NodeFilter.SHOW_TEXT,
                {
                    acceptNode: (node) => {
                        if (!node.nodeValue.trim()) return NodeFilter.FILTER_REJECT;
                        const nodeRange = document.createRange();
                        nodeRange.selectNodeContents(node);
                        return (range.compareBoundaryPoints(Range.END_TO_START, nodeRange) < 0 &&
                                range.compareBoundaryPoints(Range.START_TO_END, nodeRange) > 0)
                            ? NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_REJECT;
                    }
                },
                false
            );
            const textNodes = [];
            let currentNode;
            while ((currentNode = walker.nextNode())) {
                textNodes.push(currentNode);
            }
            textNodes.forEach(node => {
                // Если уже есть span с такими стилями — просто дополняем стиль
                if (node.parentNode.nodeName === 'SPAN') {
                    Object.assign(node.parentNode.style, styleObj);
                } else {
                    const span = document.createElement('span');
                    Object.assign(span.style, styleObj);
                    node.parentNode.replaceChild(span, node);
                    span.appendChild(node);
                }
            });
        }
    constructor(editorId, toolbarId) {
        this.editor = document.getElementById(editorId);
        this.toolbar = document.getElementById(toolbarId);
        this.hiddenInput = null;
        this.cursorMarker = null; // Маркер для сохранения позиции курсора
        this.draggedElement = null; // Элемент, который перетаскивается
        
        // Добавляем CSS стили для изображений
        const style = document.createElement('style');
        style.textContent = `
            #${editorId} img.wysiwyg-image {
                min-height: 100px !important;
                background: #f8f9fa !important;
                border: 1px solid #dee2e6 !important;
            }
            #${editorId} img.wysiwyg-image:not([src]), 
            #${editorId} img.wysiwyg-image[loading] {
                aspect-ratio: 4/3 !important;
                object-fit: cover !important;
                background: linear-gradient(45deg, #f8f9fa 25%, transparent 25%), 
                           linear-gradient(-45deg, #f8f9fa 25%, transparent 25%), 
                           linear-gradient(45deg, transparent 75%, #f8f9fa 75%), 
                           linear-gradient(-45deg, transparent 75%, #f8f9fa 75%) !important;
                background-size: 20px 20px !important;
                background-position: 0 0, 0 10px, 10px -10px, -10px 0px !important;
            }
            .wysiwyg-float-left {
                float: left !important;
                margin: 5px 15px 5px 0 !important;
                shape-outside: margin-box !important;
            }
            .wysiwyg-float-right {
                float: right !important;
                margin: 5px 0 5px 15px !important;
                shape-outside: margin-box !important;
            }
            .wysiwyg-block {
                float: none !important;
                display: block !important;
                margin: 10px auto !important;
                shape-outside: none !important;
            }
            #${editorId} {
                word-wrap: break-word;
                overflow-wrap: break-word;
            }
            #${editorId} p, #${editorId} div {
                clear: none !important;
                float: none !important;
                display: block !important;
                position: relative;
                overflow: visible;
                margin: 0 0 1em 0 !important;
                padding: 0 !important;
            }
            #${editorId} br {
                clear: none !important;
            }
            #${editorId} .wysiwyg-float-left + * {
                margin-left: 0;
                margin-right: 0;
            }
            #${editorId} .wysiwyg-float-right + * {
                margin-left: 0;
                margin-right: 0;
            }
            /* Улучшаем выделение текста возле float изображений */
            #${editorId} .wysiwyg-float-left,
            #${editorId} .wysiwyg-float-right {
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -ms-user-select: none;
            }
            #${editorId} .wysiwyg-float-left:hover,
            #${editorId} .wysiwyg-float-right:hover {
                outline: 2px solid #007bff;
                outline-offset: 2px;
            }
            /* Убеждаемся, что текст возле float изображений можно выделять */
            #${editorId} p, #${editorId} div, #${editorId} span {
                user-select: text;
                -webkit-user-select: text;
                -moz-user-select: text;
                -ms-user-select: text;
            }
        `;
        document.head.appendChild(style);
        
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
        
        // Инициализируем drag and drop для редактора
        this.setupDragAndDrop();
        
        // Обновляем скрытое поле при изменении контента
        this.editor.addEventListener('input', () => this.updateHiddenInput());
        
        // Обрабатываем нажатие клавиш
        this.editor.addEventListener('keydown', (e) => this.handleKeyDown(e));
        
        // Обрабатываем выделение текста
        this.editor.addEventListener('mouseup', (e) => this.handleMouseUp(e));
        
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
        this.hiddenInput.name = 'content';
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
    
    setupDragAndDrop() {
        // Добавляем обработчики drag and drop для редактора
        this.editor.addEventListener('dragover', (e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = 'move';
        });
        
        this.editor.addEventListener('drop', (e) => {
            e.preventDefault();
            
            if (this.draggedElement) {
                const range = document.caretRangeFromPoint(e.clientX, e.clientY);
                if (range) {
                    range.insertNode(this.draggedElement);
                    range.collapse(false);
                    const selection = window.getSelection();
                    selection.removeAllRanges();
                    selection.addRange(range);
                }
                this.draggedElement = null;
                this.updateHiddenInput();
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
                    this.wrapSelectionWithSpan({ fontFamily: item.value });
                } else if (item.action === 'fontSize') {
                    this.wrapSelectionWithSpan({ fontSize: item.value });
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
        // Сохраняем позицию курсора перед открытием диалога
        this.saveSelection();
        // Создаем кастомное диалоговое окно
        this.showImageUploadDialog();
    }
    
    saveSelection() {
        const selection = window.getSelection();
        if (selection.rangeCount > 0) {
            const range = selection.getRangeAt(0);
            
            // Вставляем временный маркер в позицию курсора
            const marker = document.createElement('span');
            marker.id = 'cursor-marker-' + Date.now();
            marker.style.display = 'none';
            
            range.insertNode(marker);
            this.cursorMarker = marker;
        } else {
            this.cursorMarker = null;
        }
    }
    
    restoreSelection() {
        if (this.cursorMarker) {
            const marker = this.cursorMarker;
            const range = document.createRange();
            const selection = window.getSelection();
            
            range.setStartBefore(marker);
            range.setEndBefore(marker);
            selection.removeAllRanges();
            selection.addRange(range);
            
            // Удаляем маркер
            marker.remove();
            this.cursorMarker = null;
        }
    }
    
    showImageUploadDialog() {
        // Создаем модальное окно
        const modal = document.createElement('div');
        modal.style.cssText = 'position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.7); z-index: 10000; display: flex; align-items: center; justify-content: center;';
        
        const dialog = document.createElement('div');
        dialog.style.cssText = 'background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.3); max-width: 500px; width: 90%;';
        
        dialog.innerHTML = `
            <h3 style="margin: 0 0 20px 0; font-size: 20px; color: #333;">📷 Добавить изображение</h3>
            
            <div style="margin-bottom: 20px;">
                <label style="display: block; margin-bottom: 8px; font-weight: 500; color: #555;">URL изображения:</label>
                <input type="text" id="image-url-input" placeholder="https://example.com/image.jpg" 
                    style="width: 100%; padding: 10px; border: 2px solid #ddd; border-radius: 6px; font-size: 14px; box-sizing: border-box;" />
            </div>
            
            <div style="margin-bottom: 25px; text-align: center;">
                <div style="margin-bottom: 10px; color: #666; font-size: 14px;">или</div>
                <button id="upload-from-pc-btn" type="button" 
                    style="padding: 12px 24px; background: #4CAF50; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    📁 Загрузить с компьютера
                </button>
                <input type="file" id="image-file-input" accept="image/jpeg,image/jpg,image/png,image/gif,image/webp" style="display: none;" />
            </div>
            
            <div style="display: flex; gap: 10px; justify-content: flex-end;">
                <button id="cancel-btn" type="button" 
                    style="padding: 10px 20px; background: #f5f5f5; color: #333; border: none; border-radius: 6px; cursor: pointer; font-size: 14px;">
                    Отмена
                </button>
                <button id="ok-btn" type="button" 
                    style="padding: 10px 20px; background: #2196F3; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500;">
                    ОК
                </button>
            </div>
        `;
        
        modal.appendChild(dialog);
        document.body.appendChild(modal);
        
        const urlInput = dialog.querySelector('#image-url-input');
        const fileInput = dialog.querySelector('#image-file-input');
        const uploadBtn = dialog.querySelector('#upload-from-pc-btn');
        const okBtn = dialog.querySelector('#ok-btn');
        const cancelBtn = dialog.querySelector('#cancel-btn');
        
        let selectedFile = null;
        
        // Обработчик кнопки загрузки с компьютера
        uploadBtn.addEventListener('click', () => {
            fileInput.click();
        });
        
        // Обработчик выбора файла
        fileInput.addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (file) {
                selectedFile = file;
                uploadBtn.textContent = `✅ ${file.name}`;
                uploadBtn.style.background = '#4CAF50';
                urlInput.value = '';
                urlInput.disabled = true;
            }
        });
        
        // Обработчик ввода URL
        urlInput.addEventListener('input', () => {
            if (urlInput.value) {
                selectedFile = null;
                fileInput.value = '';
                uploadBtn.textContent = '📁 Загрузить с компьютера';
                uploadBtn.style.background = '#4CAF50';
                urlInput.disabled = false;
            }
        });
        
        // Обработчик кнопки OK
        okBtn.addEventListener('click', async () => {
            if (selectedFile) {
                document.body.removeChild(modal);
                await this.uploadImageFromFileObject(selectedFile);
            } else if (urlInput.value.trim()) {
                const imageUrl = urlInput.value.trim();
                document.body.removeChild(modal);
                
                // Загружаем изображение по URL в S3
                await this.uploadImageFromURL(imageUrl);
            } else {
                alert('⚠️ Пожалуйста, введите URL или выберите файл');
            }
        });
        
        // Обработчик кнопки Отмена
        cancelBtn.addEventListener('click', () => {
            document.body.removeChild(modal);
        });
        
        // Закрытие по клику вне диалога
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                document.body.removeChild(modal);
            }
        });
        
        // Фокус на поле URL
        setTimeout(() => urlInput.focus(), 100);
    }
    
    async uploadImageFromFileObject(file) {
        if (!file) return;
            
            // Проверяем размер файла (максимум 10 МБ)
            const maxSize = 10 * 1024 * 1024; // 10 MB
            if (file.size > maxSize) {
                alert('❌ Ошибка: размер файла превышает 10 МБ');
                return;
            }
            
            // Показываем индикатор загрузки
            const loadingMsg = document.createElement('div');
            loadingMsg.style.cssText = 'position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); background: rgba(0,0,0,0.8); color: white; padding: 20px; border-radius: 8px; z-index: 10000;';
            loadingMsg.textContent = '⏳ Загрузка изображения...';
            document.body.appendChild(loadingMsg);
            
            try {
                // Создаем FormData для отправки файла
                const formData = new FormData();
                formData.append('image', file);
                
                // Отправляем на сервер
                const response = await fetch('/admin/api/v1/upload/image', {
                    method: 'POST',
                    body: formData,
                    // Добавляем заголовок авторизации, если он есть в localStorage
                    headers: {
                        // JWT токен будет добавлен автоматически из cookie или localStorage
                    }
                });
                
                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.message || 'Ошибка загрузки изображения');
                }
                
                const data = await response.json();
                
                // Убираем индикатор загрузки
                document.body.removeChild(loadingMsg);
                
                // Вставляем изображение в редактор
                this.insertImageElement(data.image_url);
                
            } catch (error) {
                // Убираем индикатор загрузки
                if (document.body.contains(loadingMsg)) {
                    document.body.removeChild(loadingMsg);
                }
                
                // Показываем ошибку
                alert(`❌ Ошибка загрузки изображения:\n${error.message}`);
                console.error('Upload error:', error);
            }
    }
    
    async uploadImageFromURL(url) {
        // Показываем индикатор загрузки
        const loadingMsg = document.createElement('div');
        loadingMsg.style.cssText = 'position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); background: rgba(0,0,0,0.8); color: white; padding: 20px; border-radius: 8px; z-index: 10000;';
        loadingMsg.textContent = '⏳ Загрузка изображения по URL...';
        document.body.appendChild(loadingMsg);
        
        try {
            // Отправляем URL на сервер
            const response = await fetch('/admin/api/v1/upload/image-from-url', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url })
            });
            
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Ошибка загрузки изображения');
            }
            
            const data = await response.json();
            
            // Убираем индикатор загрузки
            document.body.removeChild(loadingMsg);
            
            // Вставляем изображение в редактор
            this.insertImageElement(data.image_url);
            
        } catch (error) {
            // Убираем индикатор загрузки
            if (document.body.contains(loadingMsg)) {
                document.body.removeChild(loadingMsg);
            }
            
            // Показываем ошибку
            alert(`❌ Ошибка загрузки изображения по URL:\n${error.message}`);
            console.error('Upload from URL error:', error);
        }
    }

    
    insertImageElement(url, width = null) {
        // Создаем HTML для изображения как отдельного блока
        const widthAttr = (width && !isNaN(width) && width > 0) ? ` style="width: ${width}px;"` : ' style="width: 300px;"';
        const imgHTML = `<div><img src="${url}" alt="Image" class="wysiwyg-image wysiwyg-block" draggable="true" style="max-width: 100%; height: auto; cursor: move;${widthAttr.replace(' style="', '').replace(';"', '')}"></div>`;
        
        // Фокусируем редактор
        this.editor.focus();
        
        // Восстанавливаем позицию курсора если она была сохранена
        if (this.cursorMarker) {
            this.restoreSelection();
        }
        
        const selection = window.getSelection();
        
        if (selection.rangeCount === 0) {
            // Вставляем в конец
            this.editor.insertAdjacentHTML('beforeend', imgHTML);
            // Прокручиваем к новому изображению
            const newImages = this.editor.querySelectorAll('img[src="' + url + '"]');
            const lastImg = newImages[newImages.length - 1];
            if (lastImg) {
                lastImg.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            }
        } else {
            // Вставляем в позицию курсора как отдельный блок
            const range = selection.getRangeAt(0);
            
            // Если курсор находится внутри div'а с текстом, нужно разделить его
            const container = range.commonAncestorContainer;
            let parentDiv = container.nodeType === Node.TEXT_NODE ? container.parentElement : container;
            
            // Ищем ближайший div
            while (parentDiv && parentDiv !== this.editor && parentDiv.tagName !== 'DIV') {
                parentDiv = parentDiv.parentElement;
            }
            
            if (parentDiv && parentDiv !== this.editor && parentDiv.tagName === 'DIV') {
                // Разделяем div на две части: до курсора и после
                const divContent = parentDiv.innerHTML;
                const beforeCursor = divContent.substring(0, range.startOffset);
                const afterCursor = divContent.substring(range.endOffset);
                
                // Заменяем содержимое div'а на часть до курсора
                parentDiv.innerHTML = beforeCursor;
                
                // Вставляем изображение после этого div'а
                const imgDiv = document.createElement('div');
                imgDiv.innerHTML = `<img src="${url}" alt="Image" class="wysiwyg-image wysiwyg-block" draggable="true" style="max-width: 100%; height: auto; cursor: move;${widthAttr.replace(' style="', '').replace(';"', '')}">`;
                
                if (parentDiv.nextSibling) {
                    parentDiv.parentNode.insertBefore(imgDiv, parentDiv.nextSibling);
                } else {
                    parentDiv.parentNode.appendChild(imgDiv);
                }
                
                // Если есть текст после курсора, создаем новый div для него
                if (afterCursor.trim()) {
                    const afterDiv = document.createElement('div');
                    afterDiv.innerHTML = afterCursor;
                    imgDiv.parentNode.insertBefore(afterDiv, imgDiv.nextSibling);
                }
                
                // Устанавливаем курсор после изображения
                const newRange = document.createRange();
                newRange.setStartAfter(imgDiv);
                newRange.collapse(true);
                selection.removeAllRanges();
                selection.addRange(newRange);
            } else {
                // Обычная вставка
                range.deleteContents();
                const tempDiv = document.createElement('div');
                tempDiv.innerHTML = imgHTML;
                range.insertNode(tempDiv);
                
                // Устанавливаем курсор после вставленного блока
                range.setStartAfter(tempDiv);
                range.collapse(true);
                selection.removeAllRanges();
                selection.addRange(range);
            }
        }
        
        // Находим вставленное изображение и добавляем обработчики
        const images = this.editor.querySelectorAll('img[src="' + url + '"]');
        const img = images[images.length - 1]; // Берем последнее изображение с этим URL
        
        if (img) {
            // Добавляем обработчик успешной загрузки
            img.onload = () => {
                console.log('Image loaded successfully:', url);
                // Убираем placeholder стили
                img.style.background = 'transparent';
                img.style.border = 'none';
                img.style.aspectRatio = 'auto'; // Убираем фиксированные пропорции
                // Обновляем позиционирование
                this.updateHiddenInput();
            };
            
            // Если изображение уже загружено (из кэша), вызываем onload
            if (img.complete && img.naturalHeight > 0) {
                img.onload();
            }
            
            // Добавляем обработчик ошибки загрузки
            img.onerror = () => {
                console.error('Failed to load image:', url);
                img.alt = '⚠️ Изображение не загружено';
                img.style.border = '2px dashed #ff0000';
                img.style.padding = '20px';
                img.style.background = '#fff3cd';
                img.style.color = '#856404';
                img.style.fontSize = '14px';
                img.style.minHeight = '100px';
                img.style.display = 'flex';
                img.style.alignItems = 'center';
                img.style.justifyContent = 'center';
            };
            
            // Добавляем обработчик клика для редактирования размера
            img.onclick = (e) => {
                e.preventDefault();
                this.editImageSize(img);
            };
            
            // Добавляем обработчики для drag and drop
            img.addEventListener('dragstart', (e) => {
                e.dataTransfer.effectAllowed = 'move';
                e.dataTransfer.setData('text/html', img.outerHTML);
                img.style.opacity = '0.5';
                this.draggedElement = img;
            });
            
            img.addEventListener('dragend', (e) => {
                img.style.opacity = '1';
                this.draggedElement = null;
                this.updateHiddenInput();
            });
            
            // Добавляем обработчики для drop зоны
            img.addEventListener('dragover', (e) => {
                e.preventDefault();
                e.dataTransfer.dropEffect = 'move';
            });
            
            img.addEventListener('drop', (e) => {
                e.preventDefault();
                if (this.draggedElement && this.draggedElement !== img) {
                    // Меняем местами элементы
                    const draggedHTML = this.draggedElement.outerHTML;
                    const targetHTML = img.outerHTML;
                    
                    this.draggedElement.outerHTML = targetHTML;
                    img.outerHTML = draggedHTML;
                    
                    this.updateHiddenInput();
                }
            });
        }
        
        this.updateHiddenInput();
        
        console.log('Image inserted into editor:', url);
        
        // Очищаем маркер позиции курсора
        this.cursorMarker = null;
    }
    
    editImageSize(img) {
        const currentWidth = img.style.width ? parseInt(img.style.width) : img.naturalWidth || 300;
        let currentFloat = 'left';
        if (img.classList.contains('wysiwyg-float-right')) {
            currentFloat = 'right';
        } else if (img.classList.contains('wysiwyg-block')) {
            currentFloat = 'none';
        }
        
        // Создаем модальное окно для редактирования
        const modal = document.createElement('div');
        modal.style.cssText = `
            position: fixed; top: 0; left: 0; width: 100%; height: 100%;
            background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center;
            z-index: 10000;
        `;
        
        const modalContent = document.createElement('div');
        modalContent.style.cssText = `
            background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            min-width: 300px;
        `;
        
        modalContent.innerHTML = `
            <h3 style="margin-top: 0;">Настройки изображения</h3>
            <div style="margin-bottom: 15px;">
                <label>Ширина (px): <input type="number" id="img-width" value="${currentWidth}" min="50" max="800" style="width: 80px;"></label>
            </div>
            <div style="margin-bottom: 15px;">
                <label>Выравнивание: 
                    <select id="img-align">
                        <option value="left" ${currentFloat === 'left' ? 'selected' : ''}>Слева (текст справа)</option>
                        <option value="right" ${currentFloat === 'right' ? 'selected' : ''}>Справа (текст слева)</option>
                        <option value="none" ${currentFloat === 'none' ? 'selected' : ''}>По центру (блок)</option>
                    </select>
                </label>
            </div>
            <div style="text-align: right;">
                <button id="img-save" style="margin-right: 10px; padding: 5px 15px;">Сохранить</button>
                <button id="img-cancel" style="padding: 5px 15px;">Отмена</button>
            </div>
        `;
        
        modal.appendChild(modalContent);
        document.body.appendChild(modal);
        
        // Обработчики кнопок
        document.getElementById('img-save').onclick = () => {
            const newWidth = parseInt(document.getElementById('img-width').value);
            const newAlign = document.getElementById('img-align').value;
            
            if (newWidth && !isNaN(newWidth) && newWidth > 0) {
                img.style.width = newWidth + 'px';
                
                // Удаляем старые классы выравнивания
                img.classList.remove('wysiwyg-float-left', 'wysiwyg-float-right', 'wysiwyg-block');
                
                // Добавляем новый класс выравнивания
                if (newAlign === 'left') {
                    img.classList.add('wysiwyg-float-left');
                } else if (newAlign === 'right') {
                    img.classList.add('wysiwyg-float-right');
                } else {
                    img.classList.add('wysiwyg-block');
                }
                
                this.updateHiddenInput();
            }
            
            document.body.removeChild(modal);
        };
        
        document.getElementById('img-cancel').onclick = () => {
            document.body.removeChild(modal);
        };
        
        // Закрытие по клику вне модального окна
        modal.onclick = (e) => {
            if (e.target === modal) {
                document.body.removeChild(modal);
            }
        };
    }
    
    handlePaste(e) {
        e.preventDefault();
        
        // Получаем только текст без форматирования
        const text = e.clipboardData.getData('text/plain');
        
        // Вставляем как текст
        document.execCommand('insertText', false, text);
    }
    
    handleMouseUp(e) {
        // Проверяем, было ли выделение текста возле float изображений
        const selection = window.getSelection();
        if (selection.rangeCount > 0) {
            const range = selection.getRangeAt(0);
            if (range.collapsed) {
                // Если выделение свернуто (просто клик), проверяем, не кликнули ли по float изображению
                let element = e.target;
                while (element && element !== this.editor) {
                    if (element.classList && 
                        (element.classList.contains('wysiwyg-float-left') || 
                         element.classList.contains('wysiwyg-float-right'))) {
                        // Кликнули по float изображению - выделяем его
                        const imgRange = document.createRange();
                        imgRange.selectNode(element);
                        selection.removeAllRanges();
                        selection.addRange(imgRange);
                        break;
                    }
                    element = element.parentElement;
                }
            }
        }
    }
    
    handleKeyDown(e) {
        // Убираем специальную обработку Enter возле float изображений
        // Браузер сам справляется с созданием новых параграфов
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

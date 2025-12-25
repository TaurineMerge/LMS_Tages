// Глобальные переменные - будут установлены в шаблоне
let questionCounter;
let answerCounter = {};

// Функция расчета баллов
function calculateAnswerPoints(questionIndex) {
    const questionCard = document.querySelector(`[data-question-index="${questionIndex}"]`);
    if (!questionCard) return 0;
    
    const answerInputs = questionCard.querySelectorAll('.answer-points-input');
    let total = 0;
    
    answerInputs.forEach(input => {
        const value = parseFloat(input.value) || 0;
        total += value;
    });
    
    const totalDisplay = questionCard.querySelector(`#question-points-${questionIndex}`);
    if (totalDisplay) {
        totalDisplay.textContent = total.toFixed(1);
    }
    
    return total;
}

// Функция обновления всех расчетов
function updateAllPointsCalculations() {
    let totalTestPoints = 0;
    
    document.querySelectorAll('.question-editor').forEach(card => {
        const questionIndex = parseInt(card.dataset.questionIndex);
        const questionPoints = calculateAnswerPoints(questionIndex);
        totalTestPoints += questionPoints;
    });
    
    const totalPointsDisplay = document.getElementById('total-test-points');
    if (totalPointsDisplay) {
        totalPointsDisplay.textContent = totalTestPoints.toFixed(1);
    }
    
    return totalTestPoints;
}

// Настройка обработчиков баллов
function setupPointsListeners() {
    document.querySelectorAll('.answer-points-input').forEach(input => {
        input.addEventListener('input', function() {
            const questionCard = this.closest('.question-editor');
            const questionIndex = parseInt(questionCard.dataset.questionIndex);
            calculateAnswerPoints(questionIndex);
            updateAllPointsCalculations();
        });
    });
}

// Функция добавления ответа
function addAnswer(questionIndex) {
    const answersContainer = document.getElementById(`answers-${questionIndex}`);
    const answerIndex = answerCounter[questionIndex] || answersContainer.children.length;
    
    // Создаем элементы динамически
    const answerGroup = document.createElement('div');
    answerGroup.className = 'answer-group';
    
    const textInput = document.createElement('input');
    textInput.type = 'text';
    textInput.name = `question[${questionIndex}][answers][${answerIndex}][text]`;
    textInput.placeholder = `Текст ответа ${answerIndex + 1}`;
    textInput.required = true;
    
    const pointsGroup = document.createElement('div');
    pointsGroup.className = 'answer-points-group';
    
    const pointsLabel = document.createElement('label');
    pointsLabel.textContent = 'Баллы:';
    
    const pointsInput = document.createElement('input');
    pointsInput.type = 'number';
    pointsInput.className = 'answer-points-input';
    pointsInput.name = `question[${questionIndex}][answers][${answerIndex}][points]`;
    pointsInput.value = '0';
    pointsInput.min = '0';
    pointsInput.step = '1';
    pointsInput.required = true;
    
    const removeBtn = document.createElement('button');
    removeBtn.type = 'button';
    removeBtn.className = 'btn-danger remove-answer-btn';
    removeBtn.textContent = '× Удалить';
    removeBtn.disabled = answerIndex === 0;
    removeBtn.onclick = function() {
        removeAnswer(questionIndex, answerIndex);
    };
    
    // Собираем структуру
    pointsGroup.appendChild(pointsLabel);
    pointsGroup.appendChild(pointsInput);
    
    answerGroup.appendChild(textInput);
    answerGroup.appendChild(pointsGroup);
    answerGroup.appendChild(removeBtn);
    
    answersContainer.appendChild(answerGroup);
    answerCounter[questionIndex] = (answerCounter[questionIndex] || 0) + 1;
    
    // Добавляем обработчик
    pointsInput.addEventListener('input', function() {
        calculateAnswerPoints(questionIndex);
        updateAllPointsCalculations();
    });
    
    calculateAnswerPoints(questionIndex);
    updateAllPointsCalculations();
}

// Функция удаления ответа
function removeAnswer(questionIndex, answerIndex) {
    const answersContainer = document.getElementById(`answers-${questionIndex}`);
    const answerElements = answersContainer.querySelectorAll('.answer-group');
    
    if (answerElements.length <= 2) {
        alert('Вопрос должен содержать минимум 2 ответа');
        return;
    }
    
    if (answerIndex < answerElements.length) {
        answerElements[answerIndex].remove();
        
        // Обновляем индексы
        const remainingAnswers = answersContainer.querySelectorAll('.answer-group');
        remainingAnswers.forEach((element, newIndex) => {
            const textInput = element.querySelector('input[type="text"]');
            const pointsInput = element.querySelector('.answer-points-input');
            const removeBtn = element.querySelector('.remove-answer-btn');
            
            if (textInput) {
                textInput.name = `question[${questionIndex}][answers][${newIndex}][text]`;
                textInput.placeholder = `Текст ответа ${newIndex + 1}`;
            }
            
            if (pointsInput) {
                pointsInput.name = `question[${questionIndex}][answers][${newIndex}][points]`;
            }
            
            if (removeBtn) {
                removeBtn.onclick = function() {
                    removeAnswer(questionIndex, newIndex);
                };
                removeBtn.disabled = newIndex === 0;
            }
        });
        
        answerCounter[questionIndex] = remainingAnswers.length;
        calculateAnswerPoints(questionIndex);
        updateAllPointsCalculations();
    }
}

// Остальные функции перенесите аналогично, создавая элементы через document.createElement()
// вместо строковых шаблонов

// Инициализация
document.addEventListener('DOMContentLoaded', function() {
    // Установите questionCounter и answerCounter из data-атрибутов
    const container = document.getElementById('questions-container');
    if (container) {
        questionCounter = parseInt(container.dataset.questionCount) || 1;
    }
    
    // Инициализация счетчиков
    document.querySelectorAll('.question-editor').forEach((card, index) => {
        const questionIndex = parseInt(card.dataset.questionIndex);
        const answerCount = card.querySelectorAll('.answer-group').length;
        answerCounter[questionIndex] = answerCount;
    });
    
    setupPointsListeners();
    updateAllPointsCalculations();
});
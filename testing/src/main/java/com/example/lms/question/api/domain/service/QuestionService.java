package com.example.lms.question.api.domain.service;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import com.example.lms.question.api.dto.Question;
import com.example.lms.question.api.domain.model.QuestionModel;
import com.example.lms.question.api.domain.repository.QuestionRepositoryInterface;

/**
 * Сервисный слой для работы с вопросами.
 * Выполняет преобразование между DTO и доменной моделью,
 * а также содержит бизнес-логику, связанную с вопросами тестов.
 */
public class QuestionService {

    private final QuestionRepositoryInterface repository;

    /**
     * Конструктор сервиса.
     *
     * @param repository репозиторий для работы с вопросами
     */
    public QuestionService(QuestionRepositoryInterface repository) {
        this.repository = repository;
    }

    /**
     * Преобразует DTO в доменную модель.
     * <p>
     * Если {@code id == null}, считается, что это создание нового вопроса.
     * Если поле {@code order} не задано, оно может быть рассчитано как
     * следующий порядковый номер для теста.
     *
     * @param dto объект DTO
     * @return доменная модель QuestionModel
     */
    private QuestionModel toModel(Question dto) {
        int order = dto.getOrder() != null ? dto.getOrder() : 0;

        if (dto.getId() == null) {
            // Новый вопрос: если order не указан, вычисляем следующий
            if (dto.getOrder() == null && dto.getTestId() != null) {
                order = repository.getNextOrderForTest(dto.getTestId());
            }

            return new QuestionModel(
                    dto.getTestId(),
                    dto.getTextOfQuestion(),
                    order
            );
        } else {
            // Существующий вопрос
            return new QuestionModel(
                    dto.getId(),
                    dto.getTestId(),
                    dto.getTextOfQuestion(),
                    order
            );
        }
    }

    /**
     * Преобразует доменную модель в DTO.
     *
     * @param model доменная модель QuestionModel
     * @return объект DTO Question
     */
    private Question toDTO(QuestionModel model) {
        return new Question(
                model.getId(),
                model.getTestId(),
                model.getTextOfQuestion(),
                model.getOrder()
        );
    }

    /**
     * Получает все вопросы.
     *
     * @return список всех вопросов в формате DTO
     */
    public List<Question> getAllQuestions() {
        return repository.findAll().stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Получает вопрос по его идентификатору.
     *
     * @param id идентификатор вопроса
     * @return объект Question
     * @throws RuntimeException если вопрос не найден
     */
    public Question getQuestionById(UUID id) {
        QuestionModel model = repository.findById(id)
                .orElseThrow(() -> new RuntimeException("Вопрос с ID " + id + " не найден"));
        return toDTO(model);
    }

    /**
     * Получает вопрос по его идентификатору с возможностью отсутствия.
     *
     * @param id идентификатор вопроса
     * @return Optional с Question, если вопрос найден
     */
    public java.util.Optional<Question> findQuestionById(UUID id) {
        return repository.findById(id)
                .map(this::toDTO);
    }

    /**
     * Получает все вопросы для указанного теста.
     *
     * @param testId идентификатор теста
     * @return список вопросов указанного теста
     */
    public List<Question> getQuestionsByTestId(UUID testId) {
        return repository.findByTestId(testId).stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Создает новый вопрос.
     *
     * @param question объект Question с данными вопроса
     * @return созданный вопрос
     * @throws IllegalArgumentException если данные невалидны
     */
    public Question createQuestion(Question question) {
        // Проверяем обязательные поля
        if (question.getTestId() == null) {
            throw new IllegalArgumentException("Идентификатор теста обязателен");
        }
        if (question.getTextOfQuestion() == null || question.getTextOfQuestion().trim().isEmpty()) {
            throw new IllegalArgumentException("Текст вопроса обязателен");
        }

        QuestionModel model = toModel(question);
        model.validate();
        QuestionModel saved = repository.save(model);
        return toDTO(saved);
    }

    /**
     * Обновляет существующий вопрос.
     *
     * @param question объект Question с обновлёнными данными
     * @return обновлённый вопрос
     * @throws IllegalArgumentException если данные невалидны
     * @throws RuntimeException если вопрос не найден
     */
    public Question updateQuestion(Question question) {
        if (question.getId() == null) {
            throw new IllegalArgumentException("Идентификатор вопроса обязателен для обновления");
        }

        if (!repository.existsById(question.getId())) {
            throw new RuntimeException("Вопрос с ID " + question.getId() + " не найден");
        }

        if (question.getTestId() == null) {
            throw new IllegalArgumentException("Идентификатор теста обязателен");
        }
        if (question.getTextOfQuestion() == null || question.getTextOfQuestion().trim().isEmpty()) {
            throw new IllegalArgumentException("Текст вопроса обязателен");
        }

        QuestionModel model = toModel(question);
        model.validate();
        QuestionModel updated = repository.update(model);
        return toDTO(updated);
    }

    /**
     * Удаляет вопрос по его идентификатору.
     *
     * @param id идентификатор вопроса
     * @return true если вопрос был удалён, false если не найден
     */
    public boolean deleteQuestion(UUID id) {
        return repository.deleteById(id);
    }

    /**
     * Удаляет все вопросы для указанного теста.
     *
     * @param testId идентификатор теста
     * @return количество удалённых вопросов
     */
    public int deleteQuestionsByTestId(UUID testId) {
        return repository.deleteByTestId(testId);
    }

    /**
     * Проверяет существование вопроса.
     *
     * @param id идентификатор вопроса
     * @return true если вопрос существует
     */
    public boolean existsById(UUID id) {
        return repository.existsById(id);
    }

    /**
     * Подсчитывает количество вопросов в тесте.
     *
     * @param testId идентификатор теста
     * @return количество вопросов
     */
    public int countByTestId(UUID testId) {
        return repository.countByTestId(testId);
    }

    /**
     * Ищет вопросы по подстроке в тексте.
     *
     * @param text часть текста вопроса (поиск регистронезависимый в репозитории)
     * @return список подходящих вопросов
     */
    public List<Question> searchByText(String text) {
        return repository.findByTextContaining(text).stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Сдвигает порядок вопросов в тесте, начиная с указанного.
     * <p>
     * Может использоваться при вставке нового вопроса в середину
     * или при удалении вопроса с нужным перераспределением порядков.
     *
     * @param testId   идентификатор теста
     * @param fromOrder начиная с какого порядка выполнять сдвиг
     * @param shiftBy   на сколько изменить порядок (может быть отрицательным)
     * @return количество обновлённых вопросов
     */
    public int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy) {
        return repository.shiftQuestionsOrder(testId, fromOrder, shiftBy);
    }
}
package com.example.lms.answer.domain.service;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import com.example.lms.answer.api.dto.Answer;
import com.example.lms.answer.domain.model.AnswerModel;
import com.example.lms.answer.domain.repository.AnswerRepositoryInterface;

/**
 * Сервисный слой для работы с ответами.
 * Выполняет преобразование между DTO и доменной моделью,
 * а также содержит бизнес-логику.
 */
public class AnswerService {
    private final AnswerRepositoryInterface repository;

    /**
     * Конструктор сервиса.
     *
     * @param repository репозиторий для работы с ответами
     */
    public AnswerService(AnswerRepositoryInterface repository) {
        this.repository = repository;
    }

    /**
     * Преобразует DTO в доменную модель.
     *
     * @param dto объект DTO
     * @return доменная модель AnswerModel
     */
    private AnswerModel toModel(Answer dto) {
        return new AnswerModel(
            dto.getId(),
            dto.getQuestionId(),
            dto.getText(),
            dto.getScore() != null ? dto.getScore() : 0
        );
    }

    /**
     * Преобразует доменную модель в DTO.
     *
     * @param model доменная модель AnswerModel
     * @return объект DTO Answer
     */
    private Answer toDTO(AnswerModel model) {
        return new Answer(
            model.getId(),
            model.getText(),
            model.getQuestionId(),
            model.getScore()
        );
    }

    /**
     * Получает все ответы.
     *
     * @return список всех ответов в формате DTO
     */
    public List<Answer> getAllAnswers() {
        return repository.findAll().stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Получает ответ по его идентификатору.
     *
     * @param id идентификатор ответа
     * @return объект Answer
     * @throws RuntimeException если ответ не найден
     */
    public Answer getAnswerById(UUID id) {
        AnswerModel model = repository.findById(id)
                .orElseThrow(() -> new RuntimeException("Ответ с ID " + id + " не найден"));
        return toDTO(model);
    }

    /**
     * Получает ответ по его идентификатору с проверкой существования.
     *
     * @param id идентификатор ответа
     * @return Optional с Answer, если ответ найден
     */
    public java.util.Optional<Answer> findAnswerById(UUID id) {
        return repository.findById(id)
                .map(this::toDTO);
    }

    /**
     * Получает все ответы для указанного вопроса.
     *
     * @param questionId идентификатор вопроса
     * @return список ответов для указанного вопроса
     */
    public List<Answer> getAnswersByQuestionId(UUID questionId) {
        return repository.findByQuestionId(questionId).stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Получает все правильные ответы для указанного вопроса.
     * Ответ считается правильным, если score > 0.
     *
     * @param questionId идентификатор вопроса
     * @return список правильных ответов
     */
    public List<Answer> getCorrectAnswersByQuestionId(UUID questionId) {
        return repository.findCorrectAnswersByQuestionId(questionId).stream()
                .map(this::toDTO)
                .collect(Collectors.toList());
    }

    /**
     * Создает новый ответ.
     *
     * @param answer объект Answer с данными ответа
     * @return созданный ответ
     * @throws IllegalArgumentException если данные невалидны
     * @throws RuntimeException если ответ с таким текстом уже существует для этого вопроса
     */
    public Answer createAnswer(Answer answer) {
        // Проверяем обязательные поля
        if (answer.getQuestionId() == null) {
            throw new IllegalArgumentException("Идентификатор вопроса обязателен");
        }
        if (answer.getText() == null || answer.getText().trim().isEmpty()) {
            throw new IllegalArgumentException("Текст ответа обязателен");
        }
        if (answer.getScore() == null) {
            throw new IllegalArgumentException("Балл за ответ обязателен");
        }

        // Проверяем уникальность ответа для вопроса
        if (repository.existsByQuestionIdAndText(answer.getQuestionId(), answer.getText())) {
            throw new RuntimeException("Ответ с таким текстом уже существует для этого вопроса");
        }

        // Создаем и сохраняем модель
        AnswerModel model = toModel(answer);
        AnswerModel saved = repository.save(model);
        return toDTO(saved);
    }

    /**
     * Обновляет существующий ответ.
     *
     * @param answer объект Answer с обновленными данными
     * @return обновленный ответ
     * @throws IllegalArgumentException если данные невалидны
     * @throws RuntimeException если ответ не найден
     */
    public Answer updateAnswer(Answer answer) {
        if (answer.getId() == null) {
            throw new IllegalArgumentException("Идентификатор ответа обязателен для обновления");
        }

        // Проверяем существование ответа
        if (!repository.existsById(answer.getId())) {
            throw new RuntimeException("Ответ с ID " + answer.getId() + " не найден");
        }

        // Проверяем обязательные поля
        if (answer.getQuestionId() == null) {
            throw new IllegalArgumentException("Идентификатор вопроса обязателен");
        }
        if (answer.getText() == null || answer.getText().trim().isEmpty()) {
            throw new IllegalArgumentException("Текст ответа обязателен");
        }
        if (answer.getScore() == null) {
            throw new IllegalArgumentException("Балл за ответ обязателен");
        }

        // Обновляем модель
        AnswerModel model = toModel(answer);
        AnswerModel updated = repository.update(model);
        return toDTO(updated);
    }

    /**
     * Удаляет ответ по его идентификатору.
     *
     * @param id идентификатор ответа
     * @return true если ответ был удален, false если не найден
     */
    public boolean deleteAnswer(UUID id) {
        return repository.deleteById(id);
    }

    /**
     * Удаляет все ответы для указанного вопроса.
     *
     * @param questionId идентификатор вопроса
     * @return количество удаленных ответов
     */
    public int deleteAnswersByQuestionId(UUID questionId) {
        return repository.deleteByQuestionId(questionId);
    }

    /**
     * Проверяет существование ответа.
     *
     * @param id идентификатор ответа
     * @return true если ответ существует
     */
    public boolean existsById(UUID id) {
        return repository.existsById(id);
    }

    /**
     * Проверяет существование ответа с указанным текстом для вопроса.
     *
     * @param questionId идентификатор вопроса
     * @param text текст ответа
     * @return true если такой ответ существует
     */
    public boolean existsByQuestionIdAndText(UUID questionId, String text) {
        return repository.existsByQuestionIdAndText(questionId, text);
    }

    /**
     * Подсчитывает количество ответов для вопроса.
     *
     * @param questionId идентификатор вопроса
     * @return количество ответов
     */
    public int countByQuestionId(UUID questionId) {
        return repository.countByQuestionId(questionId);
    }

    /**
     * Подсчитывает количество правильных ответов для вопроса.
     *
     * @param questionId идентификатор вопроса
     * @return количество правильных ответов
     */
    public int countCorrectAnswersByQuestionId(UUID questionId) {
        return repository.countCorrectAnswersByQuestionId(questionId);
    }

    /**
     * Проверяет, является ли ответ правильным.
     * Ответ считается правильным, если score > 0.
     *
     * @param answerId идентификатор ответа
     * @return true если ответ правильный
     */
    public boolean isAnswerCorrect(UUID answerId) {
        Answer answer = getAnswerById(answerId);
        return answer.getScore() > 0;
    }

    /**
     * Проверяет, является ли ответ правильным.
     *
     * @param answer объект Answer
     * @return true если ответ правильный (score > 0)
     */
    public boolean isAnswerCorrect(Answer answer) {
        return answer.getScore() != null && answer.getScore() > 0;
    }

    /**
     * Получает баллы за ответ.
     *
     * @param answerId идентификатор ответа
     * @return количество баллов
     */
    public int getAnswerScore(UUID answerId) {
        Answer answer = getAnswerById(answerId);
        return answer.getScore() != null ? answer.getScore() : 0;
    }
}
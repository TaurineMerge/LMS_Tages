package com.example.lms.answer.api.controller;

import com.example.lms.answer.api.domain.model.AnswerModel;
import com.example.lms.answer.api.domain.repository.AnswerRepositoryInterface;
import com.example.lms.tracing.SimpleTracer;
import io.javalin.http.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.UUID;

/**
 * Контроллер для управления ответами на вопросы.
 * Предоставляет REST API для операций CRUD над ответами.
 */
public class AnswerController {

    private static final Logger logger = LoggerFactory.getLogger(AnswerController.class);
    private final AnswerRepositoryInterface answerRepository;

    /**
     * Конструктор контроллера.
     *
     * @param answerRepository репозиторий для работы с ответами
     */
    public AnswerController(AnswerRepositoryInterface answerRepository) {
        this.answerRepository = answerRepository;
    }

    /**
     * Получает список всех ответов.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void getAllAnswers(Context ctx) {
        try {
            List<AnswerModel> answers = answerRepository.findAll();
            logger.info("Получено {} ответов", answers.size());
            ctx.json(answers);
        } catch (Exception e) {
            logger.error("Ошибка при получении всех ответов", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при получении ответов"));
        }
    }

    /**
     * Получает ответ по его идентификатору.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void getAnswerById(Context ctx) {
        try {
            String idParam = ctx.pathParam("id");
            UUID id = UUID.fromString(idParam);
            
            var answer = answerRepository.findById(id);
            if (answer.isPresent()) {
                ctx.json(answer.get());
            } else {
                ctx.status(404).json(new ErrorResponse("Ответ с ID " + id + " не найден"));
            }
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора"));
        } catch (Exception e) {
            logger.error("Ошибка при получении ответа по ID", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при получении ответа"));
        }
    }

    /**
     * Получает все ответы для указанного вопроса.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void getAnswersByQuestionId(Context ctx) {
        try {
            String questionIdParam = ctx.queryParam("questionId");
            if (questionIdParam == null || questionIdParam.isEmpty()) {
                ctx.status(400).json(new ErrorResponse("Не указан идентификатор вопроса"));
                return;
            }
            
            UUID questionId = UUID.fromString(questionIdParam);
            List<AnswerModel> answers = answerRepository.findByQuestionId(questionId);
            
            logger.info("Получено {} ответов для вопроса {}", answers.size(), questionId);
            ctx.json(answers);
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID вопроса", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора вопроса"));
        } catch (Exception e) {
            logger.error("Ошибка при получении ответов по вопросу", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при получении ответов"));
        }
    }

    /**
     * Получает правильные ответы для указанного вопроса.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void getCorrectAnswersByQuestionId(Context ctx) {
        try {
            String questionIdParam = ctx.queryParam("questionId");
            if (questionIdParam == null || questionIdParam.isEmpty()) {
                ctx.status(400).json(new ErrorResponse("Не указан идентификатор вопроса"));
                return;
            }
            
            UUID questionId = UUID.fromString(questionIdParam);
            List<AnswerModel> answers = answerRepository.findCorrectAnswersByQuestionId(questionId);
            
            logger.info("Получено {} правильных ответов для вопроса {}", answers.size(), questionId);
            ctx.json(answers);
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID вопроса", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора вопроса"));
        } catch (Exception e) {
            logger.error("Ошибка при получении правильных ответов", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при получении правильных ответов"));
        }
    }

    /**
     * Создает новый ответ.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void createAnswer(Context ctx) {
        try {
            AnswerModel answer = ctx.bodyValidator(AnswerModel.class)
                    .check(a -> a.getQuestionId() != null, "Идентификатор вопроса обязателен")
                    .check(a -> a.getText() != null && !a.getText().isEmpty(), "Текст ответа обязателен")
                    .get();
            
            // Проверяем, не существует ли уже ответ с таким текстом для этого вопроса
            if (answerRepository.existsByQuestionIdAndText(answer.getQuestionId(), answer.getText())) {
                ctx.status(409).json(new ErrorResponse("Ответ с таким текстом уже существует для этого вопроса"));
                return;
            }
            
            AnswerModel savedAnswer = answerRepository.save(answer);
            logger.info("Создан новый ответ с ID: {}", savedAnswer.getId());
            ctx.status(201).json(savedAnswer);
        } catch (IllegalArgumentException e) {
            logger.error("Ошибка валидации при создании ответа", e);
            ctx.status(400).json(new ErrorResponse("Ошибка валидации: " + e.getMessage()));
        } catch (Exception e) {
            logger.error("Ошибка при создании ответа", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при создании ответа"));
        }
    }

    /**
     * Обновляет существующий ответ.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void updateAnswer(Context ctx) {
        try {
            String idParam = ctx.pathParam("id");
            UUID id = UUID.fromString(idParam);
            
            // Проверяем существование ответа
            if (!answerRepository.existsById(id)) {
                ctx.status(404).json(new ErrorResponse("Ответ с ID " + id + " не найден"));
                return;
            }
            
            AnswerModel answer = ctx.bodyValidator(AnswerModel.class)
                    .check(a -> a.getQuestionId() != null, "Идентификатор вопроса обязателен")
                    .check(a -> a.getText() != null && !a.getText().isEmpty(), "Текст ответа обязателен")
                    .get();
            
            answer.setId(id); // Устанавливаем ID из пути
            AnswerModel updatedAnswer = answerRepository.update(answer);
            
            logger.info("Обновлен ответ с ID: {}", id);
            ctx.json(updatedAnswer);
        } catch (IllegalArgumentException e) {
            logger.error("Ошибка валидации при обновлении ответа", e);
            ctx.status(400).json(new ErrorResponse("Ошибка валидации: " + e.getMessage()));
        } catch (Exception e) {
            logger.error("Ошибка при обновлении ответа", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при обновлении ответа"));
        }
    }

    /**
     * Удаляет ответ по его идентификатору.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void deleteAnswer(Context ctx) {
        try {
            String idParam = ctx.pathParam("id");
            UUID id = UUID.fromString(idParam);
            
            boolean deleted = answerRepository.deleteById(id);
            if (deleted) {
                logger.info("Удален ответ с ID: {}", id);
                ctx.status(204);
            } else {
                ctx.status(404).json(new ErrorResponse("Ответ с ID " + id + " не найден"));
            }
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора"));
        } catch (Exception e) {
            logger.error("Ошибка при удалении ответа", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при удалении ответа"));
        }
    }

    /**
     * Удаляет все ответы для указанного вопроса.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void deleteAnswersByQuestionId(Context ctx) {
        try {
            String questionIdParam = ctx.queryParam("questionId");
            if (questionIdParam == null || questionIdParam.isEmpty()) {
                ctx.status(400).json(new ErrorResponse("Не указан идентификатор вопроса"));
                return;
            }
            
            UUID questionId = UUID.fromString(questionIdParam);
            int deletedCount = answerRepository.deleteByQuestionId(questionId);
            
            logger.info("Удалено {} ответов для вопроса {}", deletedCount, questionId);
            ctx.json(new DeleteResponse(deletedCount, "Удалено ответов: " + deletedCount));
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID вопроса", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора вопроса"));
        } catch (Exception e) {
            logger.error("Ошибка при удалении ответов по вопросу", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при удалении ответов"));
        }
    }

    /**
     * Подсчитывает количество ответов для вопроса.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void countAnswersByQuestionId(Context ctx) {
        try {
            String questionIdParam = ctx.queryParam("questionId");
            if (questionIdParam == null || questionIdParam.isEmpty()) {
                ctx.status(400).json(new ErrorResponse("Не указан идентификатор вопроса"));
                return;
            }
            
            UUID questionId = UUID.fromString(questionIdParam);
            int count = answerRepository.countByQuestionId(questionId);
            
            ctx.json(new CountResponse(count));
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID вопроса", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора вопроса"));
        } catch (Exception e) {
            logger.error("Ошибка при подсчете ответов", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при подсчете ответов"));
        }
    }

    /**
     * Подсчитывает количество правильных ответов для вопроса.
     * 
     * @param ctx контекст HTTP-запроса
     */
    public void countCorrectAnswersByQuestionId(Context ctx) {
        try {
            String questionIdParam = ctx.queryParam("questionId");
            if (questionIdParam == null || questionIdParam.isEmpty()) {
                ctx.status(400).json(new ErrorResponse("Не указан идентификатор вопроса"));
                return;
            }
            
            UUID questionId = UUID.fromString(questionIdParam);
            int count = answerRepository.countCorrectAnswersByQuestionId(questionId);
            
            ctx.json(new CountResponse(count));
        } catch (IllegalArgumentException e) {
            logger.error("Неверный формат UUID вопроса", e);
            ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора вопроса"));
        } catch (Exception e) {
            logger.error("Ошибка при подсчете правильных ответов", e);
            ctx.status(500).json(new ErrorResponse("Ошибка при подсчете правильных ответов"));
        }
    }

    // Вспомогательные классы для ответов

    /**
     * Класс для передачи сообщений об ошибках.
     */
    private static class ErrorResponse {
        private final String error;

        public ErrorResponse(String error) {
            this.error = error;
        }

        public String getError() {
            return error;
        }
    }

    /**
     * Класс для передачи результатов подсчета.
     */
    private static class CountResponse {
        private final int count;

        public CountResponse(int count) {
            this.count = count;
        }

        public int getCount() {
            return count;
        }
    }

    /**
     * Класс для передачи результатов удаления.
     */
    private static class DeleteResponse {
        private final int deletedCount;
        private final String message;

        public DeleteResponse(int deletedCount, String message) {
            this.deletedCount = deletedCount;
            this.message = message;
        }

        public int getDeletedCount() {
            return deletedCount;
        }

        public String getMessage() {
            return message;
        }
    }
}
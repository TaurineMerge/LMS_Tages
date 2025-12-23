package com.example.lms.question.domain.service;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import com.example.lms.question.api.dto.Question;
import com.example.lms.question.domain.model.QuestionModel;
import com.example.lms.question.domain.repository.QuestionRepositoryInterface;

/**
 * Сервисный слой для работы с вопросами тестов и черновиков.
 */
public class QuestionService {

	private final QuestionRepositoryInterface repository;

	public QuestionService(QuestionRepositoryInterface repository) {
		this.repository = repository;
	}

	/**
	 * Преобразует DTO в доменную модель с поддержкой draftId.
	 */
	private QuestionModel toModel(Question dto) {
		int order = dto.getOrder() != null ? dto.getOrder() : 0;

		if (dto.getId() == null) {
			// Новый вопрос: вычисляем следующий порядковый номер
			if (dto.getOrder() == null) {
				if (dto.getTestId() != null) {
					order = repository.getNextOrderForTest(dto.getTestId());
				} else if (dto.getDraftId() != null) {
					order = repository.getNextOrderForDraft(dto.getDraftId());
				}
			}

			return new QuestionModel(
					dto.getTestId(),
					dto.getDraftId(),
					dto.getTextOfQuestion(),
					order);
		} else {
			// Существующий вопрос
			return new QuestionModel(
					dto.getId(),
					dto.getTestId(),
					dto.getDraftId(),
					dto.getTextOfQuestion(),
					order);
		}
	}

	/**
	 * Преобразует доменную модель в DTO с поддержкой draftId.
	 */
	private Question toDTO(QuestionModel model) {
		Question question = new Question();
		question.setId(model.getId());
		question.setTestId(model.getTestId());
		question.setDraftId(model.getDraftId());
		question.setTextOfQuestion(model.getTextOfQuestion());
		question.setOrder(model.getOrder());
		return question;
	}

	/**
	 * Получает все вопросы.
	 */
	public List<Question> getAllQuestions() {
		return repository.findAll().stream()
				.map(this::toDTO)
				.collect(Collectors.toList());
	}

	/**
	 * Получает вопрос по его идентификатору.
	 */
	public Question getQuestionById(UUID id) {
		QuestionModel model = repository.findById(id)
				.orElseThrow(() -> new RuntimeException("Вопрос с ID " + id + " не найден"));
		return toDTO(model);
	}

	/**
	 * Получает вопрос по его идентификатору с возможностью отсутствия.
	 */
	public java.util.Optional<Question> findQuestionById(UUID id) {
		return repository.findById(id)
				.map(this::toDTO);
	}

	/**
	 * Получает все вопросы для указанного теста.
	 */
	public List<Question> getQuestionsByTestId(UUID testId) {
		return repository.findByTestId(testId).stream()
				.map(this::toDTO)
				.collect(Collectors.toList());
	}

	/**
	 * Получает все вопросы для указанного черновика.
	 */
	public List<Question> getQuestionsByDraftId(UUID draftId) {
		return repository.findByDraftId(draftId).stream()
				.map(this::toDTO)
				.collect(Collectors.toList());
	}

	/**
	 * Создает новый вопрос (для теста или черновика).
	 */
	public Question createQuestion(Question question) {
		// Проверяем обязательные поля
		if (question.getTestId() == null && question.getDraftId() == null) {
			throw new IllegalArgumentException("Вопрос должен принадлежать либо тесту, либо черновику");
		}
		if (question.getTestId() != null && question.getDraftId() != null) {
			throw new IllegalArgumentException("Вопрос не может принадлежать одновременно и тесту, и черновику");
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
	 */
	public Question updateQuestion(Question question) {
		if (question.getId() == null) {
			throw new IllegalArgumentException("Идентификатор вопроса обязателен для обновления");
		}

		if (question.getTestId() == null && question.getDraftId() == null) {
			throw new IllegalArgumentException("Вопрос должен принадлежать либо тесту, либо черновику");
		}
		if (question.getTestId() != null && question.getDraftId() != null) {
			throw new IllegalArgumentException("Вопрос не может принадлежать одновременно и тесту, и черновику");
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
	 */
	public boolean deleteQuestion(UUID id) {
		return repository.deleteById(id);
	}

	/**
	 * Удаляет все вопросы черновика.
	 */
	public boolean deleteQuestionsByDraftId(UUID draftId) {
		return repository.deleteByDraftId(draftId);
	}

	/**
	 * Подсчитывает количество вопросов в тесте.
	 */
	public int countByTestId(UUID testId) {
		return repository.countByTestId(testId);
	}

	/**
	 * Подсчитывает количество вопросов в черновике.
	 */
	public int countByDraftId(UUID draftId) {
		return repository.countByDraftId(draftId);
	}

	/**
	 * Ищет вопросы по подстроке в тексте.
	 */
	public List<Question> searchByText(String text) {
		return repository.findByTextContaining(text).stream()
				.map(this::toDTO)
				.collect(Collectors.toList());
	}

	/**
	 * Сдвигает порядок вопросов в тесте.
	 */
	public int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy) {
		return repository.shiftQuestionsOrder(testId, fromOrder, shiftBy);
	}

	/**
	 * Сдвигает порядок вопросов в черновике.
	 */
	public int shiftDraftQuestionsOrder(UUID draftId, int fromOrder, int shiftBy) {
		return repository.shiftDraftQuestionsOrder(draftId, fromOrder, shiftBy);
	}

	/**
	 * Получает следующий порядковый номер для вопроса в тесте.
	 */
	public int getNextOrderForTest(UUID testId) {
		return repository.getNextOrderForTest(testId);
	}

	/**
	 * Получает следующий порядковый номер для вопроса в черновике.
	 */
	public int getNextOrderForDraft(UUID draftId) {
		return repository.getNextOrderForDraft(draftId);
	}
}
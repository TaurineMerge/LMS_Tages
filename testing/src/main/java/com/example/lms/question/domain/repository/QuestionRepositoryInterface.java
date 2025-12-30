package com.example.lms.question.domain.repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.question.domain.model.QuestionModel;

/**
 * Репозиторий для работы с вопросами тестов и черновиков.
 *
 * Соответствует таблице question_d:
 * - id UUID (PK, not null)
 * - test_id UUID (FK → test_d.id, может быть null для черновиков)
 * - draft_id UUID (FK → draft_b.id, может быть null для тестов)
 * - text_of_question TEXT (текст вопроса)
 * - order INT (порядок вопроса в тесте/черновике)
 */
public interface QuestionRepositoryInterface {

	/**
	 * Сохранить новый вопрос.
	 *
	 * @param question доменная модель вопроса (без id или с временным id)
	 * @return сохранённая модель с присвоенным ID
	 */
	QuestionModel save(QuestionModel question);

	/**
	 * Обновить существующий вопрос.
	 *
	 * @param question модель вопроса с обновлёнными данными (id должен быть
	 *                 заполнен)
	 * @return обновлённая модель
	 */
	QuestionModel update(QuestionModel question);

	/**
	 * Найти вопрос по ID.
	 *
	 * @param id идентификатор вопроса (question_d.id)
	 * @return Optional с вопросом, если найден
	 */
	Optional<QuestionModel> findById(UUID id);

	/**
	 * Найти все вопросы в системе.
	 *
	 * @return список всех вопросов
	 */
	List<QuestionModel> findAll();

	/**
	 * Найти все вопросы конкретного теста.
	 *
	 * @param testId идентификатор теста (question_d.test_id)
	 * @return список вопросов теста, отсортированный по полю order
	 */
	List<QuestionModel> findByTestId(UUID testId);

	/**
	 * Найти все вопросы конкретного черновика.
	 *
	 * @param draftId идентификатор черновика (question_d.draft_id)
	 * @return список вопросов черновика, отсортированный по полю order
	 */
	List<QuestionModel> findByDraftId(UUID draftId);

	/**
	 * Удалить вопрос по ID.
	 *
	 * @param id идентификатор вопроса
	 * @return true, если вопрос был удалён; false — если запись не найдена
	 */
	boolean deleteById(UUID id);

	/**
	 * Удалить все вопросы черновика.
	 *
	 * @param draftId идентификатор черновика
	 * @return true, если вопросы были удалены; false — если записей не найдено
	 */
	boolean deleteByDraftId(UUID draftId);

	/**
	 * Получить количество вопросов в тесте.
	 *
	 * @param testId идентификатор теста
	 * @return количество вопросов
	 */
	int countByTestId(UUID testId);

	/**
	 * Получить количество вопросов в черновике.
	 *
	 * @param draftId идентификатор черновика
	 * @return количество вопросов
	 */
	int countByDraftId(UUID draftId);

	/**
	 * Найти вопросы по тексту (регистронезависимый поиск).
	 *
	 * Использует поле question_d.text_of_question.
	 *
	 * @param text часть текста вопроса для поиска
	 * @return список найденных вопросов
	 */
	List<QuestionModel> findByTextContaining(String text);

	/**
	 * Получить следующий порядковый номер для нового вопроса в тесте.
	 *
	 * Как правило, реализуется как:
	 * MAX(order) + 1 для заданного testId.
	 *
	 * @param testId идентификатор теста
	 * @return следующий доступный порядковый номер
	 */
	int getNextOrderForTest(UUID testId);

	/**
	 * Получить следующий порядковый номер для нового вопроса в черновике.
	 *
	 * @param draftId идентификатор черновика
	 * @return следующий доступный порядковый номер
	 */
	int getNextOrderForDraft(UUID draftId);

	/**
	 * Обновить порядок вопросов в тесте.
	 *
	 * Доменный сценарий:
	 * - при вставке/удалении вопроса "посередине" теста
	 * нужно сдвинуть остальные вопросы.
	 *
	 * @param testId    идентификатор теста
	 * @param fromOrder порядок, с которого начинать сдвиг
	 * @param shiftBy   на сколько сдвигать (может быть отрицательным)
	 * @return количество обновлённых вопросов
	 */
	int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy);

	/**
	 * Обновить порядок вопросов в черновике.
	 *
	 * @param draftId   идентификатор черновика
	 * @param fromOrder порядок, с которого начинать сдвиг
	 * @param shiftBy   на сколько сдвигать (может быть отрицательным)
	 * @return количество обновлённых вопросов
	 */
	int shiftDraftQuestionsOrder(UUID draftId, int fromOrder, int shiftBy);
}
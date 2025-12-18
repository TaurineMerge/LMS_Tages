"""Centralized SQL query strings used across repositories."""

# Base repository templates -------------------------------------------------
BASE_SELECT_BY_ID = "SELECT * FROM {table} WHERE id = :id"

BASE_SELECT_ALL = """
    SELECT * FROM {table}
    ORDER BY {order_clause}
    LIMIT :limit OFFSET :offset
    """

BASE_COUNT = "SELECT COUNT(*) AS count FROM {table} {where_clause}"

BASE_DELETE = "DELETE FROM {table} WHERE id = :id"

BASE_EXISTS = "SELECT 1 FROM {table} WHERE id = :id LIMIT 1"

# Student queries -----------------------------------------------------------
STUDENT_INSERT = """
    INSERT INTO personal_account.student_s
        (name, surname, birth_date, avatar, contacts, email, phone)
    VALUES
        (:name, :surname, :birth_date, :avatar, CAST(:contacts AS jsonb), :email, :phone)
    RETURNING *
    """

STUDENT_UPDATE_TEMPLATE = """
    UPDATE personal_account.student_s
    SET {set_clause}
    WHERE id = :id
    RETURNING *
    """

STUDENT_BY_EMAIL = "SELECT * FROM personal_account.student_s WHERE email = :email"

STUDENT_PAGINATED = """
    SELECT * FROM personal_account.student_s
    ORDER BY created_at DESC
    LIMIT :limit OFFSET :offset
    """

STUDENT_COUNT = "SELECT COUNT(*) AS count FROM personal_account.student_s"

STUDENT_EMAIL_EXISTS = """
    SELECT 1 FROM personal_account.student_s
    WHERE email = :email {exclude_clause}
    LIMIT 1
    """

# Certificate queries -------------------------------------------------------
CERTIFICATE_INSERT = """
    INSERT INTO personal_account.certificate_b
        (content, student_id, course_id, test_attempt_id)
    VALUES
        (:content, :student_id, :course_id, :test_attempt_id)
    RETURNING *
    """

CERTIFICATES_BY_STUDENT = """
    SELECT * FROM personal_account.certificate_b
    WHERE student_id = :student_id
    ORDER BY created_at DESC
    """

CERTIFICATES_BY_COURSE = """
    SELECT * FROM personal_account.certificate_b
    WHERE course_id = :course_id
    ORDER BY created_at DESC
    """

CERTIFICATES_FILTERED_TEMPLATE = """
    SELECT * FROM personal_account.certificate_b
    {where_clause}
    ORDER BY created_at DESC
    """

CERTIFICATE_BY_NUMBER = "SELECT * FROM personal_account.certificate_b WHERE certificate_number = :certificate_number"

# Visit queries -------------------------------------------------------------
VISIT_INSERT = """
    INSERT INTO personal_account.visit_students_for_lessons
        (student_id, lesson_id)
    VALUES (:student_id, :lesson_id)
    ON CONFLICT (student_id, lesson_id) DO NOTHING
    RETURNING *
    """

VISIT_BY_STUDENT = "SELECT * FROM personal_account.visit_students_for_lessons WHERE student_id = :student_id"

VISIT_BY_LESSON = "SELECT * FROM personal_account.visit_students_for_lessons WHERE lesson_id = :lesson_id"

VISIT_FILTERED_TEMPLATE = """
    SELECT * FROM personal_account.visit_students_for_lessons
    {where_clause}
    """

VISIT_EXISTS = """
    SELECT 1 FROM personal_account.visit_students_for_lessons
    WHERE student_id = :student_id AND lesson_id = :lesson_id
    LIMIT 1
    """

VISIT_BY_ID = "SELECT * FROM personal_account.visit_students_for_lessons WHERE id = :id"

# Integration (raw data) queries -------------------------------------------
INSERT_RAW_USER_STATS = """
    INSERT INTO integration.raw_user_stats
        (student_id, payload, received_at, processed, error_message)
    VALUES
        (:student_id, CAST(:payload AS jsonb), :received_at, :processed, :error_message)
    RETURNING *
"""

INSERT_RAW_ATTEMPT = """
    INSERT INTO integration.raw_attempts
        (external_attempt_id, student_id, test_id, payload, received_at, processed, processing_attempts, error_message)
    VALUES
        (:external_attempt_id, :student_id, :test_id, CAST(:payload AS jsonb), :received_at, :processed, :processing_attempts, :error_message)
    ON CONFLICT (external_attempt_id) DO UPDATE
      SET payload = EXCLUDED.payload,
          received_at = EXCLUDED.received_at,
          processed = EXCLUDED.processed,
          processing_attempts = integration.raw_attempts.processing_attempts
    RETURNING *
"""

SELECT_UNPROCESSED_USER_STATS = """
    SELECT * FROM integration.raw_user_stats
    WHERE processed = FALSE
    ORDER BY received_at ASC
    LIMIT :limit
"""

SELECT_UNPROCESSED_ATTEMPTS = """
    SELECT * FROM integration.raw_attempts
    WHERE processed = FALSE
    ORDER BY received_at ASC
    LIMIT :limit
"""

MARK_RAW_USER_STATS_PROCESSED = """
    UPDATE integration.raw_user_stats
    SET processed = TRUE, updated_at = CURRENT_TIMESTAMP
    WHERE id = :id
    RETURNING *
"""

MARK_RAW_ATTEMPT_PROCESSED = """
    UPDATE integration.raw_attempts
    SET processed = TRUE, updated_at = CURRENT_TIMESTAMP
    WHERE id = :id
    RETURNING *
"""

INCREMENT_RAW_ATTEMPT_PROCESSING = """
    UPDATE integration.raw_attempts
    SET processing_attempts = processing_attempts + 1, updated_at = CURRENT_TIMESTAMP
    WHERE id = :id
    RETURNING processing_attempts
"""

# Business table upsert for attempts (used by StatsRepository)
TEST_ATTEMPT_UPSERT = """
INSERT INTO tests.test_attempt_b
(
    id,
    student_id,
    test_id,
    date_of_attempt,
    point,
    result,
    completed,
    passed,
    certificate_id,
    attempt_snapshot_s3,
    attempt_version,
    meta,
    created_at,
    updated_at
)
VALUES (
    :id,
    :student_id,
    :test_id,
    :date_of_attempt,
    :point,
    CAST(:result AS jsonb),
    :completed,
    :passed,
    :certificate_id,
    :snapshot,
    CAST(:version AS jsonb),
    CAST(:meta AS jsonb),
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    point = EXCLUDED.point,
    result = EXCLUDED.result,
    completed = EXCLUDED.completed,
    passed = EXCLUDED.passed,
    certificate_id = EXCLUDED.certificate_id,
    attempt_snapshot_s3 = EXCLUDED.attempt_snapshot_s3,
    attempt_version = EXCLUDED.attempt_version,
    meta = EXCLUDED.meta,
    updated_at = NOW();
"""

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
GET_PASSING_ATTEMPTS_WITHOUT_CERTIFICATES = """
SELECT 
    ta.id,
    ta.student_id,
    ta.test_id as course_id,
    ta.point as score,
    td.min_point as max_score,
    cb.title as course_name
FROM tests.test_attempt_b ta
JOIN tests.test_d td ON ta.test_id = td.id
JOIN knowledge_base.course_b cb ON td.course_id = cb.id
WHERE ta.passed = true
AND ta.completed = true
AND ta.certificate_id IS NULL  -- certificates are generated here, so testing doesn't set it
ORDER BY ta.date_of_attempt DESC
LIMIT 100
"""

GET_PASSING_ATTEMPTS_WITHOUT_CERTIFICATES_FOR_STUDENT = """
SELECT 
    ta.id,
    ta.student_id,
    ta.test_id as course_id,
    ta.point as score,
    td.min_point as max_score,
    cb.title as course_name
FROM tests.test_attempt_b ta
JOIN tests.test_d td ON ta.test_id = td.id
JOIN knowledge_base.course_b cb ON td.course_id = cb.id
WHERE ta.passed = true
AND ta.completed = true
AND ta.certificate_id IS NULL  -- certificates are generated here, so testing doesn't set it
AND ta.student_id = :student_id
ORDER BY ta.date_of_attempt DESC
LIMIT 100
"""

CERTIFICATE_INSERT = """
    INSERT INTO personal_account.certificate_b
        (student_id, certificate_number, pdf_s3_key, snapshot_s3_key)
    VALUES
        (:student_id, :certificate_number, :pdf_s3_key, :snapshot_s3_key)
    RETURNING *
    """

GET_MAX_CERTIFICATE_NUMBER = "SELECT MAX(certificate_number) FROM personal_account.certificate_b"

CERTIFICATE_UPDATE_S3_KEY = """
    UPDATE personal_account.certificate_b
    SET pdf_s3_key = :pdf_s3_key  -- Или snapshot_s3_key
    WHERE id = :id
    RETURNING *
    """

UPDATE_TEST_ATTEMPT_CERTIFICATE_ID = """
    UPDATE tests.test_attempt_b
    SET certificate_id = :certificate_id
    WHERE id = :test_attempt_id
    """

CERTIFICATES_BY_STUDENT = """
SELECT
    c.*,
    ta.test_id,
    t.course_id,
    kb.title as course_name
FROM personal_account.certificate_b c
LEFT JOIN tests.test_attempt_b ta ON ta.certificate_id = c.id
LEFT JOIN tests.test_d t ON t.id = ta.test_id
LEFT JOIN knowledge_base.course_b kb ON kb.id = t.course_id
WHERE c.student_id = :student_id
ORDER BY c.created_at DESC
"""

CERTIFICATES_BY_STUDENT_WITH_COURSE = """
SELECT
    c.*,
    ta.test_id,
    t.course_id,
    kb.title as course_name
FROM personal_account.certificate_b c
LEFT JOIN tests.test_attempt_b ta ON ta.certificate_id = c.id
LEFT JOIN tests.test_d t ON t.id = ta.test_id
LEFT JOIN knowledge_base.course_b kb ON kb.id = t.course_id
WHERE c.student_id = :student_id
ORDER BY c.created_at DESC
"""


CERTIFICATES_BY_COURSE = (
    "SELECT * FROM personal_account.certificate_b WHERE course_id = :course_id"  # Удалить, если course_id нет
)


CERTIFICATES_FILTERED_TEMPLATE = """
    SELECT * FROM personal_account.certificate_b
    {where_clause}
    ORDER BY created_at DESC
    """

CERTIFICATE_BY_NUMBER = "SELECT * FROM personal_account.certificate_b WHERE certificate_number = :certificate_number"

# Student attempts queries --------------------------------------------------
STUDENT_ATTEMPTS = """
SELECT
    ta.id,
    ta.student_id,
    ta.test_id,
    ta.date_of_attempt,
    ta.point,
    ta.passed,
    ta.completed,
    ta.result,
    td.title as test_title,
    cb.title as course_name,
    cb.id as course_id
FROM tests.test_attempt_b ta
JOIN tests.test_d td ON ta.test_id = td.id
JOIN knowledge_base.course_b cb ON td.course_id = cb.id
WHERE ta.student_id = :student_id
ORDER BY ta.date_of_attempt DESC
LIMIT 50
"""

STUDENT_ATTEMPTS_WITH_CERTIFICATES = """
SELECT
    ta.id,
    ta.student_id,
    ta.test_id,
    ta.date_of_attempt,
    ta.point,
    ta.passed,
    ta.completed,
    ta.result,
    ta.certificate_id,
    td.title as test_title,
    cb.title as course_name,
    cb.id as course_id
FROM tests.test_attempt_b ta
JOIN tests.test_d td ON ta.test_id = td.id
JOIN knowledge_base.course_b cb ON td.course_id = cb.id
WHERE ta.student_id = :student_id
ORDER BY ta.date_of_attempt DESC
LIMIT 50
"""

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

# Integration (raw data) queries ---------------------------------------MARK_RAW_ATTEMPT_PROCESSED----
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

# Aggregated stats table query
STUDENT_STATS_AGGREGATED_UPSERT = """
INSERT INTO stats.student_stats_aggregated (
    student_id,
    total_attempts,
    passed_attempts,
    failed_attempts,
    avg_score,
    total_tests_taken,
    last_attempt_at,
    stats_json,
    created_at,
    updated_at
)
VALUES (
    :student_id,
    :total_attempts,
    :passed_attempts,
    :failed_attempts,
    :avg_score,
    :total_tests_taken,
    :last_attempt_at,
    CAST(:stats_json AS jsonb),
    NOW(),
    NOW()
)
ON CONFLICT (student_id) DO UPDATE SET
    total_attempts = EXCLUDED.total_attempts,
    passed_attempts = EXCLUDED.passed_attempts,
    failed_attempts = EXCLUDED.failed_attempts,
    avg_score = EXCLUDED.avg_score,
    total_tests_taken = EXCLUDED.total_tests_taken,
    last_attempt_at = EXCLUDED.last_attempt_at,
    stats_json = EXCLUDED.stats_json,
    updated_at = NOW();
"""

STUDENT_STATS_AGGREGATED_SELECT = """
SELECT * FROM stats.student_stats_aggregated WHERE student_id = :student_id
"""

STUDENT_STATS_CALCULATE = """
SELECT
    student_id,
    COUNT(*) as total_attempts,
    COUNT(*) FILTER (WHERE passed = true) as passed_attempts,
    COUNT(*) FILTER (WHERE passed = false) as failed_attempts,
    COALESCE(AVG(point), 0) as avg_score,
    COUNT(DISTINCT test_id) as total_tests_taken,
    MAX(date_of_attempt) as last_attempt_at
FROM tests.test_attempt_b
WHERE student_id = :student_id
GROUP BY student_id
"""

# Raw data queries for charts and analytics
RAW_USER_STATS_BY_STUDENT = """
SELECT
    id,
    student_id,
    payload,
    received_at,
    processed,
    error_message
FROM integration.raw_user_stats
WHERE student_id = :student_id
ORDER BY received_at DESC
LIMIT 100
"""

RAW_ATTEMPTS_BY_STUDENT = """
SELECT
    id,
    external_attempt_id,
    student_id,
    test_id,
    payload,
    received_at,
    processed,
    processing_attempts,
    error_message
FROM integration.raw_attempts
WHERE student_id = :student_id
ORDER BY received_at DESC
LIMIT 200
"""

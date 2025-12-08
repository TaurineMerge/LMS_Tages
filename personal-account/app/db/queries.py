"""Centralized SQL query strings used across repositories."""

# Base repository templates -------------------------------------------------
BASE_SELECT_BY_ID = "SELECT * FROM {table} WHERE id = :id"

BASE_SELECT_ALL = (
    """
    SELECT * FROM {table}
    ORDER BY {order_clause}
    LIMIT :limit OFFSET :offset
    """
)

BASE_COUNT = "SELECT COUNT(*) AS count FROM {table} {where_clause}"

BASE_DELETE = "DELETE FROM {table} WHERE id = :id"

BASE_EXISTS = "SELECT 1 FROM {table} WHERE id = :id LIMIT 1"

# Student queries -----------------------------------------------------------
STUDENT_INSERT = (
    """
    INSERT INTO personal_account.student_s
        (name, surname, birth_date, avatar, contacts, email, phone)
    VALUES
        (:name, :surname, :birth_date, :avatar, CAST(:contacts AS jsonb), :email, :phone)
    RETURNING *
    """
)

STUDENT_UPDATE_TEMPLATE = (
    """
    UPDATE personal_account.student_s
    SET {set_clause}
    WHERE id = :id
    RETURNING *
    """
)

STUDENT_BY_EMAIL = (
    "SELECT * FROM personal_account.student_s WHERE email = :email"
)

STUDENT_PAGINATED = (
    """
    SELECT * FROM personal_account.student_s
    ORDER BY created_at DESC
    LIMIT :limit OFFSET :offset
    """
)

STUDENT_COUNT = "SELECT COUNT(*) AS count FROM personal_account.student_s"

STUDENT_EMAIL_EXISTS = (
    """
    SELECT 1 FROM personal_account.student_s
    WHERE email = :email {exclude_clause}
    LIMIT 1
    """
)

# Certificate queries -------------------------------------------------------
CERTIFICATE_INSERT = (
    """
    INSERT INTO personal_account.certificate_b
        (content, student_id, course_id, test_attempt_id)
    VALUES
        (:content, :student_id, :course_id, :test_attempt_id)
    RETURNING *
    """
)

CERTIFICATES_BY_STUDENT = (
    """
    SELECT * FROM personal_account.certificate_b
    WHERE student_id = :student_id
    ORDER BY created_at DESC
    """
)

CERTIFICATES_BY_COURSE = (
    """
    SELECT * FROM personal_account.certificate_b
    WHERE course_id = :course_id
    ORDER BY created_at DESC
    """
)

CERTIFICATES_FILTERED_TEMPLATE = (
    """
    SELECT * FROM personal_account.certificate_b
    {where_clause}
    ORDER BY created_at DESC
    """
)

CERTIFICATE_BY_NUMBER = (
    "SELECT * FROM personal_account.certificate_b WHERE certificate_number = :certificate_number"
)

# Visit queries -------------------------------------------------------------
VISIT_INSERT = (
    """
    INSERT INTO personal_account.visit_students_for_lessons
        (student_id, lesson_id)
    VALUES (:student_id, :lesson_id)
    ON CONFLICT (student_id, lesson_id) DO NOTHING
    RETURNING *
    """
)

VISIT_BY_STUDENT = (
    "SELECT * FROM personal_account.visit_students_for_lessons WHERE student_id = :student_id"
)

VISIT_BY_LESSON = (
    "SELECT * FROM personal_account.visit_students_for_lessons WHERE lesson_id = :lesson_id"
)

VISIT_FILTERED_TEMPLATE = (
    """
    SELECT * FROM personal_account.visit_students_for_lessons
    {where_clause}
    """
)

VISIT_EXISTS = (
    """
    SELECT 1 FROM personal_account.visit_students_for_lessons
    WHERE student_id = :student_id AND lesson_id = :lesson_id
    LIMIT 1
    """
)

VISIT_BY_ID = (
    "SELECT * FROM personal_account.visit_students_for_lessons WHERE id = :id"
)

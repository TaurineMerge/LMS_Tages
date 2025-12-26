# Integration Repository (raw data)

Repository for storing raw integration payloads received from external services
(such as the `testing` service). The repository writes into the `integration`
schema tables created by `init-sql/migrate-002.sql`.

::: app.repositories.integration.integration_repository
    options:
      show_root_heading: true
      show_source: true
      members_order: source


## Purpose

This repository is intended to be used by webhooks or client adapters that
receive payloads from external services. It stores the data "as is" into
`integration.raw_user_stats` and `integration.raw_attempts` for asynchronous
processing by workers.

## Usage example

```python
from app.repositories.integration import integration_repository

# insert a user stats payload
await integration_repository.insert_raw_user_stats(student_id, payload)

# insert/update an attempt
await integration_repository.insert_raw_attempt(external_attempt_id, student_id, test_id, payload)
```

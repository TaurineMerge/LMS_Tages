# Certificate Generation System - Best Practice Implementation

## Overview

This system implements automatic certificate generation for students who successfully complete tests. It follows best practice architecture with clear separation of concerns:

- **Storage Layer** (`storage_service.py`): Manages S3/MinIO operations
- **Certificate Generation** (`cert_service.py`): Generates PDF certificates
- **Integration** (`stats_processor.py`): Integrates certificates with statistics processing
- **Worker** (`stats_worker.py`): Manages background job scheduling

## Architecture

### 1. Storage Service (`app/services/storage_service.py`)

**Purpose**: Universal S3-compatible storage interface for all file operations.

**Key Features**:
- Certificate storage and retrieval from MinIO
- Image management (logos, signatures, backgrounds)
- Presigned URL generation for downloads
- Metadata management with tagging
- Error handling and logging

**Usage**:
```python
from app.services.storage_service import get_storage_service

storage = get_storage_service()

# Store certificate
s3_key = await storage.store_certificate(
    certificate_id=cert_id,
    student_id=student_id,
    course_id=course_id,
    pdf_content=pdf_bytes
)

# Generate download URL
download_url = await storage.get_certificate_url(s3_key)

# Store images
image_key = await storage.store_image("logo.png", image_bytes, "image/png")
```

**Bucket Structure**:
```
myminio/
├── certificates/
│   └── {course_id}/
│       └── {student_id}/
│           └── {certificate_id}.pdf
└── images/
    └── python/
        ├── logo.png
        ├── signature.png
        └── backgrounds/
```

### 2. Certificate Service (`app/services/cert_service.py`)

**Purpose**: Handles certificate generation logic and PDF creation.

**Key Features**:
- Fetches student data from Keycloak
- Generates beautiful PDF certificates with ReportLab
- Stores PDFs in S3 via storage service
- Updates database with S3 keys
- Professional certificate design with:
  - Student name and email
  - Course name
  - Score and percentage
  - Issue date
  - Decorative borders and signatures

**Usage**:
```python
from app.services.cert_service import get_certificate_service

cert_service = get_certificate_service()

# Generate certificate
cert_id, s3_key = await cert_service.generate_certificate(
    student_id=UUID("..."),
    course_id=UUID("..."),
    course_name="Advanced Python",
    test_attempt_id=UUID("..."),
    score=85,
    max_score=100
)
```

**Error Handling**:
- `CertificateGenerationError`: Raised for generation failures
- Automatic rollback if S3 storage fails
- Comprehensive logging with tracing

### 3. Stats Processor Integration (`app/services/stats_processor.py`)

**Purpose**: Integrates certificate generation with statistics processing.

**Key Method**:
```python
async def check_and_generate_certificates(self) -> dict[str, int]:
    """
    Check for passing attempts without certificates and generate them.
    Called by the worker scheduler.
    """
```

**Workflow**:
1. Fetch passing attempts (score >= passing_score)
2. For each attempt without a certificate:
   - Validate student and course exist
   - Generate PDF certificate
   - Store in S3
   - Update database with S3 key
   - Handle errors gracefully

### 4. Stats Worker (`app/services/stats_worker.py`)

**Purpose**: Manages background job scheduling using APScheduler.

**Jobs**:
- `fetch_testing`: Fetches data from testing service (every 60s)
- `process_raws`: Processes raw statistics (every 15s)
- `generate_certificates`: Checks and generates certificates (every 60s)

**Usage**:
```python
worker = StatsWorker()
worker.start()  # Start all scheduled jobs

# Later...
worker.stop()   # Gracefully shutdown
```

## Database Schema Updates

Add these columns to `certificate_b` table:

```sql
-- Store S3 key instead of content
ALTER TABLE certificate_b
ADD COLUMN s3_key VARCHAR(500),
DROP COLUMN content;

-- Add index for S3 key lookups
CREATE INDEX idx_certificate_s3_key ON certificate_b(s3_key);

-- Track certificate generation state
ALTER TABLE certificate_b
ADD COLUMN generated_at TIMESTAMP,
ADD COLUMN last_error VARCHAR(500);
```

## Configuration

**Environment Variables** (`.env`):
```bash
# MinIO S3 Configuration
MINIO_ENDPOINT_URL=http://minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_REGION=us-east-1

# Keycloak (for student data)
KEYCLOAK_SERVER_URL=http://keycloak:8080/auth
KEYCLOAK_REALM=student
KEYCLOAK_ADMIN_USERNAME=admin
KEYCLOAK_ADMIN_PASSWORD=admin

# Stats Worker
STATS_WORKER_FETCH_INTERVAL=60
STATS_WORKER_PROCESS_INTERVAL=15
```

## Dependencies

Added to `pyproject.toml`:
```
boto3>=1.28.0,<2.0.0          # AWS S3/MinIO client
reportlab>=4.0.0,<5.0.0       # PDF generation
Pillow>=10.0.0,<11.0.0        # Image processing
```

## API Endpoints (Frontend Integration)

Create these endpoints in your routers:

```python
from fastapi import APIRouter, HTTPException
from app.services.storage_service import StorageError, get_storage_service

router = APIRouter(prefix="/certificates", tags=["certificates"])

@router.get("/{certificate_id}/download")
async def download_certificate(certificate_id: UUID):
    """Download certificate PDF with presigned URL."""
    # Get S3 key from database
    certificate = await certificate_service.get_certificate(certificate_id)
    
    if not certificate.s3_key:
        raise HTTPException(status_code=404, detail="Certificate not found")
    
    try:
        download_url = await storage_service.get_certificate_url(
            certificate.s3_key,
            expiration=3600  # 1 hour
        )
        return {"download_url": download_url}
    except StorageError as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/{student_id}")
async def list_certificates(student_id: UUID):
    """List all certificates for a student."""
    certificates = await certificate_service.get_certificates(student_id=student_id)
    return certificates
```

## Best Practices Implemented

### 1. **Separation of Concerns**
- Storage logic isolated in `storage_service`
- PDF generation isolated in `cert_service`
- Processing orchestration in `stats_processor`
- Job scheduling in `stats_worker`

### 2. **Error Handling**
- Custom exceptions for different failure scenarios
- Transaction rollback on S3 failure
- Comprehensive logging with correlation IDs
- Graceful degradation

### 3. **Performance Optimization**
- S3 presigned URLs avoid direct file downloads
- Metadata stored in database for quick lookups
- Batch processing in statistics processor
- Configurable job intervals for resource management

### 4. **Security**
- MinIO credentials from environment variables
- Presigned URLs with expiration
- Student data validation against Keycloak
- S3 object metadata for audit trail

### 5. **Scalability**
- Stateless service design (can run multiple instances)
- Efficient batch processing
- Configurable job intervals
- S3 as unlimited storage backend

### 6. **Observability**
- Distributed tracing with `@traced` decorators
- Structured logging with contextual information
- Metrics for certificates generated/failed
- Error tracking and alerting

## Testing

### Unit Tests Example:
```python
@pytest.mark.asyncio
async def test_generate_certificate():
    # Mock dependencies
    mock_storage = AsyncMock()
    mock_repo = AsyncMock()
    mock_keycloak = AsyncMock()
    
    cert_service = CertificateService(
        storage_service=mock_storage,
        certificate_repository_inst=mock_repo,
        keycloak_service_inst=mock_keycloak
    )
    
    # Test generation
    cert_id, s3_key = await cert_service.generate_certificate(...)
    
    # Verify calls
    mock_storage.store_certificate.assert_called_once()
    mock_repo.update_s3_key.assert_called_once()
```

## Troubleshooting

### Issue: MinIO connection refused
- Check MinIO is running: `docker-compose ps`
- Verify endpoint URL and credentials
- Check firewall rules

### Issue: Certificates not generating
- Check stats_worker logs: `docker-compose logs personal-account`
- Verify student exists in Keycloak
- Check S3 permissions
- Monitor job execution: `STATS_WORKER_PROCESS_INTERVAL`

### Issue: S3 keys not saving to database
- Verify database connection
- Check certificate table schema
- Review transaction logs

## Future Enhancements

1. **Template System**: Allow customizable certificate templates
2. **Digital Signatures**: Add cryptographic signatures
3. **QR Codes**: Include verifiable QR codes
4. **Multi-language**: Support certificates in different languages
5. **Webhooks**: Notify external systems of certificate generation
6. **Analytics**: Track certificate downloads and validity
7. **Batch Operations**: Mass certificate generation for past attempts

## References

- [ReportLab Documentation](https://www.reportlab.com/)
- [Boto3 S3 Reference](https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/s3.html)
- [MinIO Documentation](https://min.io/)
- [APScheduler Documentation](https://apscheduler.readthedocs.io/)

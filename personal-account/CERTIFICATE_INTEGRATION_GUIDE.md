# Certificate Generation System - Integration Guide

## 🎯 Quick Start

### 1. Install Dependencies
```bash
cd personal-account
poetry add boto3 reportlab Pillow
poetry lock
```

### 2. Update Configuration
Add to `.env` or `.env.local`:
```bash
# MinIO Configuration
MINIO_ENDPOINT_URL=http://minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_REGION=us-east-1
```

### 3. Run Database Migration
```bash
psql -h localhost -U appuser -d appdb < init-sql/migrate-certificate-s3.sql
```

### 4. Add Router to Your Application
```python
# In main.py or app.py
from app.routers.certificates import router as certificates_router

app.include_router(
    certificates_router,
    prefix="/api/v1",
)
```

### 5. Update Stats Worker
The stats_worker automatically generates certificates! It runs:
- `fetch_testing` - every 60 seconds
- `process_raws` - every 15 seconds  
- `generate_certificates` - every 60 seconds (NEW)

```python
# In your application startup
from app.services.stats_worker import StatsWorker

worker = StatsWorker()
worker.start()

# On shutdown
worker.stop()
```

## 📁 File Structure

```
personal-account/
├── app/
│   ├── services/
│   │   ├── storage_service.py      # S3/MinIO operations
│   │   ├── cert_service.py         # Certificate generation
│   │   ├── stats_processor.py      # Updated with cert generation
│   │   └── stats_worker.py         # Updated with cert job
│   │
│   ├── repositories/
│   │   └── certificate.py          # Updated with update_s3_key
│   │
│   ├── schemas/
│   │   └── certificate.py          # Updated schemas
│   │
│   └── routers/
│       └── certificates.py.example # Example endpoints
│
├── init-sql/
│   └── migrate-certificate-s3.sql  # Database migration
│
├── CERTIFICATE_SYSTEM.md           # Full documentation
└── FRONTEND_CERTIFICATES_EXAMPLE.md # Frontend examples
```

## 🔄 Certificate Generation Flow

```
Test Attempt Passed
        ↓
StatsWorker (every 60s)
    check_and_generate_certificates()
        ↓
1. Query passing attempts without certificates
        ↓
2. For each attempt:
   - Fetch student data from Keycloak
   - Validate student/course exist
        ↓
3. Generate PDF using ReportLab:
   - Student name, course name
   - Score and percentage
   - Issue date
   - Decorative design
        ↓
4. Store PDF in S3:
   - Bucket: myminio
   - Path: certificates/{course_id}/{student_id}/{cert_id}.pdf
        ↓
5. Update database:
   - certificate_b.s3_key = S3 path
   - certificate_b.generated_at = timestamp
   - certificate_generation_log = audit entry
        ↓
6. Generate presigned URL for download
```

## 📊 Database Changes

### New Columns
```sql
-- certificate_b table
ALTER TABLE certificate_b
ADD COLUMN s3_key VARCHAR(500) UNIQUE,
ADD COLUMN generated_at TIMESTAMP,
ADD COLUMN last_error TEXT;
```

### New Tables
- `certificate_generation_log` - Audit trail
- `certificate_templates` - Template customization

### New Views
- `v_certificates_with_metadata` - Certificates with status

## 🔐 Configuration Options

### Environment Variables
```bash
# S3/MinIO
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
STATS_WORKER_FETCH_INTERVAL=60        # seconds
STATS_WORKER_PROCESS_INTERVAL=15      # seconds
```

## 🎨 Certificate Customization

The certificate design is controlled in `cert_service.py`:

```python
# Modify colors
colors.HexColor("#1F4788")  # Primary color

# Modify fonts
pdf_canvas.setFont("Helvetica-Bold", 36)

# Modify layout
_draw_border(pdf_canvas, width, height)
```

For advanced customization:
1. Create certificate templates in `certificate_templates` table
2. Store logo/signature images in S3
3. Update `_generate_pdf` to use template styling

## 🧪 Testing

### Unit Test Example
```python
import pytest
from unittest.mock import AsyncMock

@pytest.mark.asyncio
async def test_generate_certificate():
    # Mock S3 storage
    mock_storage = AsyncMock()
    mock_storage.store_certificate.return_value = "certificates/course/student/cert.pdf"
    
    # Mock database
    mock_repo = AsyncMock()
    mock_repo.create.return_value = {"id": UUID("...")}
    
    # Create service
    from app.services.cert_service import CertificateService
    service = CertificateService(
        storage_service=mock_storage,
        certificate_repository_inst=mock_repo,
    )
    
    # Test
    cert_id, s3_key = await service.generate_certificate(
        student_id=UUID("..."),
        course_id=UUID("..."),
        course_name="Test Course",
        test_attempt_id=UUID("..."),
        score=95,
        max_score=100,
    )
    
    assert s3_key == "certificates/course/student/cert.pdf"
    mock_storage.store_certificate.assert_called_once()
```

### Integration Test Example
```python
@pytest.mark.asyncio
async def test_certificate_download_flow():
    """Test complete certificate generation and download flow"""
    
    # 1. Generate certificate
    cert_id, s3_key = await cert_service.generate_certificate(...)
    
    # 2. Get download URL
    download_url = await storage_service.get_certificate_url(s3_key)
    
    # 3. Verify URL is valid
    response = requests.head(download_url)
    assert response.status_code == 200
    
    # 4. Download and verify PDF
    pdf_content = requests.get(download_url).content
    assert pdf_content.startswith(b'%PDF')
```

## 📱 Frontend Integration

### Vue.js Example
See `FRONTEND_CERTIFICATES_EXAMPLE.md` for complete Vue.js component

Key points:
- Use presigned URLs to avoid authentication issues
- Store download URL expiration (default 1 hour)
- Handle download errors gracefully
- Show loading state during generation

### API Endpoints

```
GET    /api/v1/certificates/{student_id}
       List all certificates for student

GET    /api/v1/certificates/{certificate_id}
       Get certificate details

GET    /api/v1/certificates/{certificate_id}/download
       Get presigned download URL

GET    /api/v1/certificates/course/{course_id}
       List all certificates for course

GET    /api/v1/certificates/{student_id}/latest
       Get most recent certificate
```

## 🚀 Performance Optimization

### Batch Processing
The stats processor fetches up to 100 pending certificates per run:
```python
passing_attempts = await self.repo.get_passing_attempts_without_certificates()
```

### Caching
Consider adding Redis caching for:
- Certificate list per student
- Presigned URLs (with TTL matching expiration)
- Student data from Keycloak

```python
# Example
cache_key = f"certificates:{student_id}"
cached = await redis.get(cache_key)
if cached:
    return cached
```

### S3 Optimization
- Use multipart upload for large PDFs
- Set appropriate expiration on presigned URLs
- Monitor S3 bucket size and cleanup old certificates

## 🔍 Monitoring & Debugging

### Logs
Watch for these logs:
```
Certificate generated successfully: certificate_id=..., s3_key=...
Failed to generate certificate for student: ...
Certificate generation completed: generated=10, failed=2
```

### Database Queries
```sql
-- Check pending certificates
SELECT * FROM get_certificates_pending_generation();

-- View generation history
SELECT * FROM certificate_generation_log
ORDER BY created_at DESC
LIMIT 20;

-- Check for errors
SELECT * FROM certificate_b
WHERE last_error IS NOT NULL;

-- Certificate status overview
SELECT status, COUNT(*) FROM v_certificates_with_metadata
GROUP BY status;
```

### Metrics to Monitor
- Certificates generated per hour
- Certificate generation success rate
- Average generation time
- S3 storage usage
- S3 API errors

## 🐛 Troubleshooting

### Issue: Certificates not generating
**Check**:
1. Stats worker is running: `ps aux | grep python`
2. Database connection works: test query
3. Keycloak accessible: `curl http://keycloak:8080`
4. MinIO accessible: `curl http://minio:9000`
5. Check logs: `docker-compose logs personal-account`

### Issue: Download URL expires
**Solution**: Increase `expiration` parameter
```python
await storage_service.get_certificate_url(s3_key, expiration=86400)  # 24 hours
```

### Issue: PDF generation fails
**Check**:
1. ReportLab installed: `pip show reportlab`
2. Course name doesn't exceed limits
3. Student data valid: check Keycloak
4. Disk space available for temp files

### Issue: S3 connection error
**Check**:
1. MinIO running: `docker-compose ps`
2. Credentials correct: compare with `.env`
3. Endpoint URL format: `http://minio:9000` (no trailing slash)
4. Network: `docker-compose exec personal-account curl http://minio:9000`

## 📚 Additional Resources

- [ReportLab Documentation](https://www.reportlab.com/)
- [Boto3 S3 Guide](https://boto3.amazonaws.com/v1/documentation/api/latest/guide/s3.html)
- [MinIO Client Guide](https://min.io/docs/minio/docker/index.html)
- [APScheduler Documentation](https://apscheduler.readthedocs.io/)
- [Keycloak Admin API](https://www.keycloak.org/docs/latest/server_admin/#admin-rest-api)

## 🎓 Best Practices Summary

✅ **DO**:
- Store certificates in S3, not database
- Use presigned URLs for downloads
- Validate student data before generation
- Log all certificate operations
- Handle errors gracefully
- Monitor certificate generation metrics

❌ **DON'T**:
- Store large PDFs in database
- Expose S3 credentials in code
- Generate certificates synchronously
- Skip input validation
- Ignore storage errors
- Hardcode certificate design

## 🔄 Future Enhancements

1. **Template System**: Allow custom certificate designs per course
2. **Digital Signatures**: Add cryptographic signatures
3. **QR Codes**: Include verifiable QR codes
4. **Bulk Operations**: Generate certificates for past attempts
5. **Webhook Notifications**: Notify students when ready
6. **Analytics**: Track certificate downloads and verification
7. **Multi-language**: Support certificates in different languages
8. **Expiration**: Set certificate validity periods

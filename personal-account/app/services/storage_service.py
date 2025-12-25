"""Storage service for S3-compatible storage (MinIO).

Handles all operations with S3 storage including:
- Certificate storage and retrieval
- Image storage and management
- Metadata management
"""

import logging
from uuid import UUID

import boto3
from botocore.exceptions import ClientError

from app.config import get_settings
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StorageService:
    """Service for managing S3-compatible storage operations.

    Provides unified interface for certificate and image storage operations.
    Supports certificate generation, retrieval, deletion and image management.
    """

    # Storage bucket and paths configuration
    CERTIFICATES_BUCKET = "certificates"
    CERTIFICATES_PATH = ""
    IMAGES_BUCKET = "images"
    IMAGES_PATH = "python"

    def __init__(self):
        """Initialize S3 client with MinIO credentials."""
        settings = get_settings()
        self.s3_client = boto3.client(
            "s3",
            endpoint_url=settings.MINIO_ENDPOINT_URL,
            aws_access_key_id=settings.MINIO_ACCESS_KEY,
            aws_secret_access_key=settings.MINIO_SECRET_KEY,
            region_name=settings.MINIO_REGION,
        )
        self.settings = settings

    # ==================== Certificate Operations ====================

    @traced("storage_service.store_certificate", record_args=True, record_result=True)
    async def store_certificate(
        self,
        certificate_id: UUID,
        student_id: UUID,
        course_id: UUID,
        pdf_content: bytes,
    ) -> str:
        """Store generated certificate PDF in S3.

        Args:
            certificate_id: Unique certificate identifier
            student_id: Student UUID
            course_id: Course UUID
            pdf_content: PDF file content as bytes

        Returns:
            S3 object key for the certificate

        Raises:
            StorageError: If upload fails
        """
        s3_key = self._build_certificate_key(certificate_id, student_id, course_id)

        try:
            self.s3_client.put_object(
                Bucket=self.CERTIFICATES_BUCKET,
                Key=s3_key,
                Body=pdf_content,
                ContentType="application/pdf",
                Metadata={
                    "certificate-id": str(certificate_id),
                    "student-id": str(student_id),
                    "course-id": str(course_id),
                },
            )
            logger.info("Certificate stored: %s", s3_key)
            return s3_key
        except ClientError as e:
            logger.error("Failed to store certificate: %s", e)
            raise StorageError(f"Failed to store certificate: {e}") from e

    @traced("storage_service.get_certificate", record_args=True, record_result=True)
    async def get_certificate(self, s3_key: str) -> bytes:
        """Retrieve certificate PDF from S3.

        Args:
            s3_key: S3 object key for the certificate

        Returns:
            PDF file content as bytes

        Raises:
            StorageError: If retrieval fails or object not found
        """
        try:
            response = self.s3_client.get_object(
                Bucket=self.CERTIFICATES_BUCKET,
                Key=s3_key,
            )
            content = response["Body"].read()
            logger.info("Certificate retrieved: %s", s3_key)
            return content
        except ClientError as e:
            if e.response["Error"]["Code"] == "NoSuchKey":
                logger.warning("Certificate not found: %s", s3_key)
                raise StorageError(f"Certificate not found: {s3_key}") from e
            logger.error("Failed to retrieve certificate: %s", e)
            raise StorageError(f"Failed to retrieve certificate: {e}") from e

    @traced("storage_service.delete_certificate", record_args=True, record_result=True)
    async def delete_certificate(self, s3_key: str) -> bool:
        """Delete certificate PDF from S3.

        Args:
            s3_key: S3 object key for the certificate

        Returns:
            True if deletion successful

        Raises:
            StorageError: If deletion fails
        """
        try:
            self.s3_client.delete_object(
                Bucket=self.CERTIFICATES_BUCKET,
                Key=s3_key,
            )
            logger.info("Certificate deleted: %s", s3_key)
            return True
        except ClientError as e:
            logger.error("Failed to delete certificate: %s", e)
            raise StorageError(f"Failed to delete certificate: {e}") from e

    @traced("storage_service.get_certificate_url", record_args=True, record_result=True)
    async def get_certificate_url(self, s3_key: str, expiration: int = 3600) -> str:
        """Generate presigned URL for certificate download.

        Args:
            s3_key: S3 object key for the certificate
            expiration: URL expiration time in seconds (default: 1 hour)

        Returns:
            Presigned URL for downloading certificate

        Raises:
            StorageError: If URL generation fails
        """
        try:
            url = self.s3_client.generate_presigned_url(
                "get_object",
                Params={
                    "Bucket": self.CERTIFICATES_BUCKET,
                    "Key": s3_key,
                },
                ExpiresIn=expiration,
            )
            logger.info("Certificate URL generated: %s", s3_key)
            return url
        except ClientError as e:
            logger.error("Failed to generate certificate URL: %s", e)
            raise StorageError(f"Failed to generate certificate URL: {e}") from e

    # ==================== Image Operations ====================

    @traced("storage_service.store_image", record_args=True, record_result=True)
    async def store_image(
        self,
        image_name: str,
        image_content: bytes,
        content_type: str = "image/png",
    ) -> str:
        """Store image in S3.

        Args:
            image_name: Name of the image file
            image_content: Image file content as bytes
            content_type: MIME type of the image (default: image/png)

        Returns:
            S3 object key for the image

        Raises:
            StorageError: If upload fails
        """
        s3_key = self._build_image_key(image_name)

        try:
            self.s3_client.put_object(
                Bucket=self.IMAGES_BUCKET,
                Key=s3_key,
                Body=image_content,
                ContentType=content_type,
            )
            logger.info("Image stored: %s", s3_key)
            return s3_key
        except ClientError as e:
            logger.error("Failed to store image: %s", e)
            raise StorageError(f"Failed to store image: {e}") from e

    @traced("storage_service.get_image", record_args=True, record_result=True)
    async def get_image(self, s3_key: str) -> bytes:
        """Retrieve image from S3.

        Args:
            s3_key: S3 object key for the image

        Returns:
            Image file content as bytes

        Raises:
            StorageError: If retrieval fails or object not found
        """
        try:
            response = self.s3_client.get_object(
                Bucket=self.IMAGES_BUCKET,
                Key=s3_key,
            )
            content = response["Body"].read()
            logger.info("Image retrieved: %s", s3_key)
            return content
        except ClientError as e:
            if e.response["Error"]["Code"] == "NoSuchKey":
                logger.warning("Image not found: %s", s3_key)
                raise StorageError(f"Image not found: {s3_key}") from e
            logger.error("Failed to retrieve image: %s", e)
            raise StorageError(f"Failed to retrieve image: {e}") from e

    @traced("storage_service.delete_image", record_args=True, record_result=True)
    async def delete_image(self, s3_key: str) -> bool:
        """Delete image from S3.

        Args:
            s3_key: S3 object key for the image

        Returns:
            True if deletion successful

        Raises:
            StorageError: If deletion fails
        """
        try:
            self.s3_client.delete_object(
                Bucket=self.IMAGES_BUCKET,
                Key=s3_key,
            )
            logger.info("Image deleted: %s", s3_key)
            return True
        except ClientError as e:
            logger.error("Failed to delete image: %s", e)
            raise StorageError(f"Failed to delete image: {e}") from e

    @traced("storage_service.get_image_url", record_args=True, record_result=True)
    async def get_image_url(self, s3_key: str, expiration: int = 3600) -> str:
        """Generate presigned URL for image download.

        Args:
            s3_key: S3 object key for the image
            expiration: URL expiration time in seconds (default: 1 hour)

        Returns:
            Presigned URL for downloading image

        Raises:
            StorageError: If URL generation fails
        """
        try:
            url = self.s3_client.generate_presigned_url(
                "get_object",
                Params={
                    "Bucket": self.IMAGES_BUCKET,
                    "Key": s3_key,
                },
                ExpiresIn=expiration,
            )
            logger.info("Image URL generated: %s", s3_key)
            return url
        except ClientError as e:
            logger.error("Failed to generate image URL: %s", e)
            raise StorageError(f"Failed to generate image URL: {e}") from e

    @traced("storage_service.update_image", record_args=True, record_result=True)
    async def update_image(
        self,
        s3_key: str,
        image_content: bytes,
        content_type: str = "image/png",
    ) -> str:
        """Update existing image in S3.

        Args:
            s3_key: S3 object key for the image
            image_content: New image file content as bytes
            content_type: MIME type of the image (default: image/png)

        Returns:
            S3 object key for the image

        Raises:
            StorageError: If update fails
        """
        try:
            self.s3_client.put_object(
                Bucket=self.IMAGES_BUCKET,
                Key=s3_key,
                Body=image_content,
                ContentType=content_type,
            )
            logger.info("Image updated: %s", s3_key)
            return s3_key
        except ClientError as e:
            logger.error("Failed to update image: %s", e)
            raise StorageError(f"Failed to update image: {e}") from e

    # ==================== Utility Methods ====================

    @staticmethod
    def _build_certificate_key(
        certificate_id: UUID,
        student_id: UUID,
        course_id: UUID,
    ) -> str:
        """Build S3 key for certificate storage.

        Format: certificates/{course_id}/{student_id}/{certificate_id}.pdf

        Args:
            certificate_id: Unique certificate identifier
            student_id: Student UUID
            course_id: Course UUID

        Returns:
            S3 object key
        """
        return f"{StorageService.CERTIFICATES_PATH}/{course_id}/{student_id}/{certificate_id}.pdf"

    @staticmethod
    def _build_image_key(image_name: str) -> str:
        """Build S3 key for image storage.

        Format: images/python/{image_name}

        Args:
            image_name: Name of the image file

        Returns:
            S3 object key
        """
        return f"{StorageService.IMAGES_PATH}/{image_name}"

    @traced("storage_service.list_certificates", record_args=True, record_result=True)
    async def list_certificates(self, prefix: str = "") -> list[str]:
        """List all certificates in storage with optional prefix filter.

        Args:
            prefix: Optional prefix to filter certificates

        Returns:
            List of S3 object keys

        Raises:
            StorageError: If listing fails
        """
        try:
            if not prefix:
                prefix = self.CERTIFICATES_PATH

            response = self.s3_client.list_objects_v2(
                Bucket=self.CERTIFICATES_BUCKET,
                Prefix=prefix,
            )

            keys = []
            if "Contents" in response:
                keys = [obj["Key"] for obj in response["Contents"]]

            logger.info("Listed %d certificates with prefix: %s", len(keys), prefix)
            return keys
        except ClientError as e:
            logger.error("Failed to list certificates: %s", e)
            raise StorageError(f"Failed to list certificates: {e}") from e

    @traced("storage_service.list_images", record_args=True, record_result=True)
    async def list_images(self, prefix: str = "") -> list[str]:
        """List all images in storage with optional prefix filter.

        Args:
            prefix: Optional prefix to filter images

        Returns:
            List of S3 object keys

        Raises:
            StorageError: If listing fails
        """
        try:
            if not prefix:
                prefix = self.IMAGES_PATH

            response = self.s3_client.list_objects_v2(
                Bucket=self.IMAGES_BUCKET,
                Prefix=prefix,
            )

            keys = []
            if "Contents" in response:
                keys = [obj["Key"] for obj in response["Contents"]]

            logger.info("Listed %d images with prefix: %s", len(keys), prefix)
            return keys
        except ClientError as e:
            logger.error("Failed to list images: %s", e)
            raise StorageError(f"Failed to list images: {e}") from e


class StorageError(Exception):
    """Custom exception for storage operations."""

    pass


# Singleton instance
_storage_service_instance: StorageService | None = None


def get_storage_service() -> StorageService:
    """Get or create storage service instance.

    Returns:
        StorageService instance
    """
    global _storage_service_instance
    if _storage_service_instance is None:
        _storage_service_instance = StorageService()
    return _storage_service_instance


storage_service = StorageService()

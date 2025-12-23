
// ======================================================================
// UploadResult.java
// ======================================================================
package com.example.lms.shared.storage.dto;

/**
 * Результат загрузки файла в хранилище.
 */
public class UploadResult {
    private final String objectPath;
    private final String url;
    private final long size;
    private final String contentType;
    private final boolean success;
    private final String errorMessage;

    /**
     * Конструктор для успешной загрузки.
     */
    public UploadResult(String objectPath, String url, long size, String contentType, boolean success) {
        this.objectPath = objectPath;
        this.url = url;
        this.size = size;
        this.contentType = contentType;
        this.success = success;
        this.errorMessage = null;
    }

    /**
     * Конструктор для неудачной загрузки.
     */
    public UploadResult(String objectPath, String errorMessage) {
        this.objectPath = objectPath;
        this.url = null;
        this.size = 0;
        this.contentType = null;
        this.success = false;
        this.errorMessage = errorMessage;
    }

    public String getObjectPath() {
        return objectPath;
    }

    public String getUrl() {
        return url;
    }

    public long getSize() {
        return size;
    }

    public String getContentType() {
        return contentType;
    }

    public boolean isSuccess() {
        return success;
    }

    public String getErrorMessage() {
        return errorMessage;
    }

    @Override
    public String toString() {
        if (success) {
            return "UploadResult{" +
                    "objectPath='" + objectPath + '\'' +
                    ", url='" + url + '\'' +
                    ", size=" + size +
                    ", contentType='" + contentType + '\'' +
                    ", success=true" +
                    '}';
        } else {
            return "UploadResult{" +
                    "objectPath='" + objectPath + '\'' +
                    ", success=false" +
                    ", errorMessage='" + errorMessage + '\'' +
                    '}';
        }
    }
}
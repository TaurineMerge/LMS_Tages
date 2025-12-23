package com.example.lms.shared.storage.dto;

import java.time.Instant;

/**
 * Метаданные файла в хранилище.
 */
public class FileMetadata {
    private final String path;
    private final long size;
    private final String contentType;
    private final Instant lastModified;

    public FileMetadata(String path, long size, String contentType, Instant lastModified) {
        this.path = path;
        this.size = size;
        this.contentType = contentType;
        this.lastModified = lastModified;
    }

    public String getPath() {
        return path;
    }

    public long getSize() {
        return size;
    }

    public String getContentType() {
        return contentType;
    }

    public Instant getLastModified() {
        return lastModified;
    }

    @Override
    public String toString() {
        return "FileMetadata{" +
                "path='" + path + '\'' +
                ", size=" + size +
                ", contentType='" + contentType + '\'' +
                ", lastModified=" + lastModified +
                '}';
    }
}
package com.example.lms.shared.storage.exceptions;

/**
 * Исключение, возникающее при работе с хранилищем.
 * <p>
 * Обертка над исключениями MinIO для единообразной обработки ошибок.
 */
public class StorageException extends RuntimeException {

    public StorageException(String message) {
        super(message);
    }

    public StorageException(String message, Throwable cause) {
        super(message, cause);
    }

    public StorageException(Throwable cause) {
        super(cause);
    }
}
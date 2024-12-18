package com.fanyamin.reminder.model;

import jakarta.validation.constraints.*;
import java.time.Instant;

public record Book(
        Long id,
        
        @NotBlank(message = "The book ISBN must be defined.")
        @Pattern(regexp = "^([0-9]{10}|[0-9]{13})$", message = "The ISBN format must be valid.")
        String isbn,

        @NotBlank(message = "The book title must be defined.")
        String title,

        @NotBlank(message = "The book author must be defined.")
        String author,
        
        @NotNull(message = "The book price must be defined.")
        @Positive(message = "The book price must be greater than zero.")
        Double price,

        Instant borrowTime,
        Instant returnTime,
        Instant createdDate,
        Instant lastModifiedDate,
        int version
) {
    public static Book of(String isbn, String title, String author, Double price, Instant borrowTime, Instant returnTime) {
        return new Book(null, isbn, title, author, price, borrowTime, returnTime, null, null, 0);
    }

    public Book withBorrowTime(Instant borrowTime) {
        return new Book(id, isbn, title, author, price, borrowTime, null, createdDate, lastModifiedDate, version);
    }

    public Book withReturnTime(Instant returnTime) {
        return new Book(id, isbn, title, author, price, borrowTime, returnTime, createdDate, lastModifiedDate, version);
    }

    public boolean isAvailable() {
        return borrowTime == null || returnTime != null;
    }

    public boolean isBorrowed() {
        return borrowTime != null && returnTime == null;
    }
}
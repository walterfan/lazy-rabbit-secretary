package com.fanyamin.reminder.service;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.fanyamin.reminder.exception.BookNotFoundException;
import com.fanyamin.reminder.exception.BookOperationException;
import com.fanyamin.reminder.mapper.BookMapper;
import com.fanyamin.reminder.model.Book;

import java.time.Instant;
import java.util.List;

@Service
public class BookService {
    private final BookMapper bookMapper;
    
    public BookService(BookMapper bookMapper) {
        this.bookMapper = bookMapper;
    }
    
    public List<Book> getAllBooks() {
        return bookMapper.findAll();
    }
    
    public Book getBookById(Long id) {
        Book book = bookMapper.findById(id);
        if (book == null) {
            throw new BookNotFoundException(id);
        }
        return book;
    }
    
    @Transactional
    public Book createBook(Book book) {
        bookMapper.insert(book);
        return book;
    }
    
    @Transactional
    public Book borrowBook(Long id) {
        Book book = getBookById(id);
        
        if (!book.isAvailable()) {
            throw new BookOperationException("Book is already borrowed");
        }
        
        Book updatedBook = book.withBorrowTime(Instant.now());
        
        if (bookMapper.borrow(updatedBook) != 1) {
            throw new BookOperationException("Failed to borrow book - concurrent modification detected");
        }
        
        return updatedBook;
    }
    
    @Transactional
    public Book returnBook(Long id) {
        Book book = getBookById(id);
        
        if (!book.isBorrowed()) {
            throw new BookOperationException("Book is not borrowed");
        }
        
        Book updatedBook = book.withReturnTime(Instant.now());
        
        if (bookMapper.return_(updatedBook) != 1) {
            throw new BookOperationException("Failed to return book - concurrent modification detected");
        }
        
        return updatedBook;
    }
}
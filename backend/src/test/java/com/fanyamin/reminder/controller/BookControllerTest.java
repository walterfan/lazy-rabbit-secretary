
package com.fanyamin.reminder.controller;

import com.fanyamin.reminder.model.Book;
import com.fanyamin.reminder.service.BookService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import java.util.Arrays;
import java.util.List;

import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@WebMvcTest(BookController.class)
public class BookControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private BookService bookService;

    private ObjectMapper objectMapper;

    @BeforeEach
    public void setUp() {
        objectMapper = new ObjectMapper();
    }

    @Test
    public void testGetAllBooks() throws Exception {
        List<Book> books = Arrays.asList(
                new Book(1L, "978-7-115-63056-1", "Author One", null, null, null, null, null, null, 0),
                new Book(2L, "978-7-115-63056-2", "Author Two", null, null, null, null, null, null, 0)
        );

        when(bookService.getAllBooks()).thenReturn(books);

        mockMvc.perform(get("/api/books"))
                .andExpect(status().isOk())
                .andExpect(content().json(objectMapper.writeValueAsString(books)));
    }

    @Test
    public void testGetBookById() throws Exception {
        Book book = new Book(1L, "978-7-115-63056-3", "Author One", null, null, null, null, null, null, 0);

        when(bookService.getBookById(1L)).thenReturn(book);

        mockMvc.perform(get("/api/books/1"))
                .andExpect(status().isOk())
                .andExpect(content().json(objectMapper.writeValueAsString(book)));
    }

    @Test
    public void testCreateBook() throws Exception {
        Book book = new Book(1L, "978-7-115-63056-1", "Author One", null, null, null, null, null, null, 0);

        when(bookService.createBook(any(Book.class))).thenReturn(book);

        mockMvc.perform(post("/api/books")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(book)))
                .andExpect(status().isOk())
                .andExpect(content().json(objectMapper.writeValueAsString(book)));
    }

    @Test
    public void testBorrowBook() throws Exception {
        Book book = new Book(1L, "978-7-115-63056-5", "Author One", null, null, null, null, null, null, 0);

        when(bookService.borrowBook(1L)).thenReturn(book);

        mockMvc.perform(post("/api/books/1/borrow"))
                .andExpect(status().isOk())
                .andExpect(content().json(objectMapper.writeValueAsString(book)));
    }

    @Test
    public void testReturnBook() throws Exception {
        Book book = new Book(1L, "978-7-115-63056-6", "Author One", null, null, null, null, null, null, 0);

        when(bookService.returnBook(1L)).thenReturn(book);

        mockMvc.perform(post("/api/books/1/return"))
                .andExpect(status().isOk())
                .andExpect(content().json(objectMapper.writeValueAsString(book)));
    }
}
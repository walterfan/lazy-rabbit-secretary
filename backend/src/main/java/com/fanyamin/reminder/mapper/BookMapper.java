package com.fanyamin.reminder.mapper;

import org.apache.ibatis.annotations.*;

import com.fanyamin.reminder.model.Book;

import java.util.List;

@Mapper
public interface BookMapper {
    @Select("SELECT * FROM book")
    List<Book> findAll();
    
    @Select("SELECT * FROM book WHERE id = #{id}")
    Book findById(Long id);
    
    @Insert("INSERT INTO book (isbn, title, author, price, borrow_time, return_time, created_date, last_modified_date, version) " +
           "VALUES (#{isbn}, #{title}, #{author}, #{price}, #{borrowTime}, #{returnTime}, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)")
    @Options(useGeneratedKeys = true, keyProperty = "id")
    void insert(Book book);
    
    @Update("UPDATE book SET borrow_time = #{borrowTime}, last_modified_date = CURRENT_TIMESTAMP, version = version + 1 " +
           "WHERE id = #{id} AND version = #{version}")
    int borrow(Book book);
    
    @Update("UPDATE book SET return_time = #{returnTime}, last_modified_date = CURRENT_TIMESTAMP, version = version + 1 " +
           "WHERE id = #{id} AND version = #{version}")
    int return_(Book book);
}
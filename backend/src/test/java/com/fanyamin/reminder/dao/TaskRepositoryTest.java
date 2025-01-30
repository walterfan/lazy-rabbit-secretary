package com.fanyamin.reminder.dao;

import com.fanyamin.reminder.domain.Task;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.autoconfigure.orm.jpa.TestEntityManager;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
public class TaskRepositoryTest {

    @Autowired
    private TaskRepository taskRepository;

    @Autowired
    private TestEntityManager entityManager;

    private Task task;

    @BeforeEach
    public void setUp() {
        task = new Task();
        task.setName("Sample Task");
        task.setDesc("This is a sample task.");
        entityManager.persist(task);
        entityManager.flush();
    }

    @Test
    public void testFindByName() {
        Task foundTask = taskRepository.findByName(task.getName());
        assertThat(foundTask).isNotNull();
        assertThat(foundTask.getName()).isEqualTo(task.getName());
        assertThat(foundTask.getDesc()).isEqualTo(task.getDesc());
    }

    @Test
    public void testFindById() {
        Task foundTask = taskRepository.findById(task.getTaskId()).orElse(null);
        assertThat(foundTask).isNotNull();
        assertThat(foundTask.getName()).isEqualTo(task.getName());
        assertThat(foundTask.getDesc()).isEqualTo(task.getDesc());
    }

    @Test
    public void testSave() {
        Task newTask = new Task();
        newTask.setName("New Task");
        newTask.setDesc("This is a new task.");

        Task savedTask = taskRepository.save(newTask);
        assertThat(savedTask.getTaskId()).isNotNull();
        assertThat(savedTask.getName()).isEqualTo(newTask.getName());
        assertThat(savedTask.getDesc()).isEqualTo(newTask.getDesc());
    }

    @Test
    public void testDelete() {
        taskRepository.deleteById(task.getTaskId());
        Task deletedTask = taskRepository.findById(task.getTaskId()).orElse(null);
        assertThat(deletedTask).isNull();
    }
}
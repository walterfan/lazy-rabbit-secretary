import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;

import com.fanyamin.reminder.controller.TaskController;
import com.fanyamin.reminder.domain.Task;
import com.fanyamin.reminder.service.TaskService;

import java.util.Arrays;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
public class TaskControllerTest {

    @Mock
    private TaskService taskService;

    @InjectMocks
    private TaskController taskController;

    @Test
    public void testGetAllTasksWithPagination() {
        // 模拟任务服务返回的分页任务列表
        Page<Task> page = new PageImpl<>(Arrays.asList(
                Task.builder().taskId(UUID.randomUUID()).name("task 1").build(),
                Task.builder().taskId(UUID.randomUUID()).name("task 2").build()
        ), PageRequest.of(0, 10), 2);

        when(taskService.findAll(any(Pageable.class))).thenReturn(page);

        // 调用控制器方法获取分页任务列表
        Page<Task> result = taskController.getAllTasks(PageRequest.of(0, 10));

        // 验证结果
        assertNotNull(result);
        assertEquals(2, result.getTotalElements());
        assertEquals(10, result.getSize());
        assertEquals("task 1", result.getContent().get(0).getName());
        assertEquals("task 2", result.getContent().get(1).getName());
    }

    @Test
    public void testGetAllTasksEmpty() {
        // 模拟任务服务返回空任务列表
        when(taskService.findAll(Pageable.unpaged())).thenReturn(Page.empty());

        // 调用控制器方法获取任务列表
        Page<Task> result = taskController.getAllTasks(Pageable.unpaged());

        // 验证结果
        assertNotNull(result);
        assertEquals(0, result.getSize());
    }
}
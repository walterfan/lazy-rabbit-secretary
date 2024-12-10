package com.github.walterfan.bjava.reminder.web;

import com.github.walterfan.bjava.reminder.config.ReminderProperties;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HomeController {
    
    private final ReminderProperties reminderProperties;
    
    public HomeController(ReminderProperties reminderProperties) {
        this.reminderProperties = reminderProperties;
    }

    @GetMapping("/")
    public String getGreeting() {
        return reminderProperties.getGreeting();
    }

}

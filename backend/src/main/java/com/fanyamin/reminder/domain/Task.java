package com.fanyamin.reminder.domain;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.Builder;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.Id;
import java.util.UUID;

@Entity
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Task extends BaseEntity {
    @Id
    @GeneratedValue
    private UUID taskId;
    private String name;
    private String desc;
    private String tags;
    private String startTime;
    private String endTime;
    private String deadline;
    private String tenantId;
}
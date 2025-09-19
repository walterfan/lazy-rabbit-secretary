Feature: Task Reminding
  As a user
  I want to create tasks with automatic reminders
  So that I can be notified before my tasks are due

  Background:
    Given I am authenticated as user "john.doe@example.com"
    And I belong to realm "test-realm"

  Scenario: Create a simple task with email reminder
    Given I want to create a task with the following details:
      | name         | Team Meeting                    |
      | description  | Weekly team sync meeting        |
      | schedule_time| 2025-09-17T09:00:00Z           |
      | minutes      | 60                              |
      | deadline     | 2025-09-17T10:00:00Z           |
      | priority     | 3                               |
      | difficulty   | 2                               |
      | tags         | meeting,team                    |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 30                      |
      | reminder_methods        | email                   |
      | reminder_targets        | john.doe@example.com    |
    When I create the task
    Then the task should be created successfully
    And a reminder should be created for "2025-09-17T08:30:00Z"
    And the reminder should have method "email"
    And the reminder should target "john.doe@example.com"
    And the reminder content should contain "Team Meeting"
    And the reminder content should contain "📋 Task:"

  Scenario: Create a repeating daily task with multiple reminder methods
    Given I want to create a task with the following details:
      | name         | Daily Standup                   |
      | description  | Daily team standup meeting      |
      | schedule_time| 2025-09-17T09:00:00Z           |
      | minutes      | 30                              |
      | deadline     | 2025-09-17T09:30:00Z           |
    And I want to set up repeat settings:
      | is_repeating    | true                            |
      | repeat_pattern  | daily                           |
      | repeat_interval | 1                               |
      | repeat_count    | 5                               |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 15                      |
      | reminder_methods        | email,webhook           |
      | reminder_targets        | team@example.com,https://hooks.slack.com/webhook |
    When I create the task
    Then the task should be created successfully
    And the task should be marked as repeating
    And 5 task instances should be created
    And 5 reminders should be created
    And each reminder should have methods "email,webhook"
    And each reminder should be scheduled 15 minutes before its task

  Scenario: Create a weekly repeating task with specific days
    Given I want to create a task with the following details:
      | name         | Gym Workout                     |
      | description  | Regular workout session         |
      | schedule_time| 2025-09-17T18:00:00Z           |
      | minutes      | 90                              |
      | deadline     | 2025-09-17T19:30:00Z           |
    And I want to set up repeat settings:
      | is_repeating     | true                           |
      | repeat_pattern   | weekly                         |
      | repeat_interval  | 1                              |
      | repeat_days_of_week| mon,wed,fri                  |
      | repeat_end_date  | 2025-12-31T00:00:00Z          |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 60                      |
      | reminder_methods        | email                   |
      | reminder_targets        | john.doe@example.com    |
    When I create the task
    Then the task should be created successfully
    And task instances should be created for Monday, Wednesday, and Friday
    And reminders should be created 60 minutes before each workout

  Scenario: Create a monthly task with specific day
    Given I want to create a task with the following details:
      | name         | Monthly Report                  |
      | description  | Generate monthly status report  |
      | schedule_time| 2025-09-15T14:00:00Z           |
      | minutes      | 120                             |
      | deadline     | 2025-09-15T16:00:00Z           |
    And I want to set up repeat settings:
      | is_repeating     | true                           |
      | repeat_pattern   | monthly                        |
      | repeat_interval  | 1                              |
      | repeat_day_of_month| 15                           |
      | repeat_count     | 12                             |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 120                     |
      | reminder_methods        | email                   |
      | reminder_targets        | manager@example.com     |
    When I create the task
    Then the task should be created successfully
    And task instances should be created for the 15th of each month
    And reminders should be created 2 hours before each report task

  Scenario: Task without reminders should not create reminders
    Given I want to create a task with the following details:
      | name         | Simple Task                     |
      | description  | A task without reminders        |
      | schedule_time| 2025-09-17T10:00:00Z           |
      | minutes      | 30                              |
      | deadline     | 2025-09-17T10:30:00Z           |
    And I do not want to generate reminders:
      | generate_reminders| false                          |
    When I create the task
    Then the task should be created successfully
    And no reminders should be created

  Scenario: Validate reminder methods
    Given I want to create a task with reminder settings:
      | generate_reminders      | true                    |
      | reminder_methods        | invalid_method          |
    When I create the task
    Then the task creation should fail
    And I should receive an error "invalid reminder method: invalid_method"

  Scenario: Validate reminder advance time
    Given I want to create a task with reminder settings:
      | generate_reminders      | true                    |
      | reminder_methods        | email                   |
      | reminder_advance_minutes| -30                     |
    When I create the task
    Then the task creation should fail
    And I should receive an error "reminder_advance_minutes cannot be negative"

  Scenario: Reminder content formatting for recurring task
    Given I want to create a repeating task with the following details:
      | name         | Weekly Review                   |
      | description  | Review weekly progress          |
      | schedule_time| 2025-09-19T16:00:00Z           |
      | minutes      | 45                              |
      | deadline     | 2025-09-19T16:45:00Z           |
      | priority     | 4                               |
      | difficulty   | 3                               |
      | tags         | review,weekly,important         |
    And I want to set up repeat settings:
      | is_repeating    | true                            |
      | repeat_pattern  | weekly                          |
      | repeat_interval | 1                               |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 30                      |
      | reminder_methods        | email                   |
      | reminder_targets        | john.doe@example.com    |
    When I create the task
    Then the task should be created successfully
    And the reminder content should contain:
      | 📋 Task: Weekly Review                            |
      | 📝 Description: Review weekly progress            |
      | ⏰ Scheduled: 2025-09-19 16:00:00                |
      | ⏱️  Duration: 45 minutes                         |
      | 📅 Deadline: 2025-09-19 16:45:00                |
      | 🎯 Priority: 4/5                                 |
      | 🔧 Difficulty: 3/5                               |
      | 🏷️  Tags: review,weekly,important                |
      | 🔄 This is a recurring task instance.            |
      | Generated automatically by Lazy Rabbit Secretary System |

  Scenario: Reminder tags formatting
    Given I want to create a task with the following details:
      | name         | Project Meeting                 |
      | description  | Discuss project milestones      |
      | tags         | project,milestone,urgent        |
    And I want to set up repeat settings:
      | is_repeating    | true                            |
      | repeat_pattern  | weekly                          |
    And I want to generate reminders:
      | generate_reminders| true                           |
      | reminder_methods  | email                          |
    When I create the task
    Then the reminder tags should contain "task,auto-generated,recurring,project,milestone,urgent"

  Scenario: Skip reminder creation for past times
    Given I want to create a task with the following details:
      | name         | Past Task                       |
      | schedule_time| 2024-01-01T10:00:00Z           |
      | minutes      | 30                              |
      | deadline     | 2024-01-01T10:30:00Z           |
    And I want to generate reminders with the following settings:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 60                      |
      | reminder_methods        | email                   |
      | reminder_targets        | john.doe@example.com    |
    When I create the task
    Then the task should be created successfully
    But no reminder should be created
    And I should receive a warning "skipping reminder creation for task"

  Scenario: Handle reminder service unavailable
    Given the reminder service is unavailable
    And I want to create a task with reminder settings:
      | generate_reminders| true                           |
      | reminder_methods  | email                          |
    When I create the task
    Then the task should be created successfully
    But reminder creation should fail
    And I should receive a warning "reminder service not available"

  Scenario Outline: Validate repeat patterns with reminders
    Given I want to create a task with repeat pattern "<pattern>"
    And I want to set repeat interval to <interval>
    And I want to generate reminders
    When I create the task
    Then the task creation should <result>

    Examples:
      | pattern  | interval | result  |
      | daily    | 1        | succeed |
      | weekly   | 2        | succeed |
      | monthly  | 1        | succeed |
      | yearly   | 1        | succeed |
      | invalid  | 1        | fail    |
      | daily    | 0        | fail    |
      | daily    | -1       | fail    |

  Scenario: Complex weekly pattern with multiple days
    Given I want to create a task with the following details:
      | name         | Team Standup                    |
      | schedule_time| 2025-09-16T09:00:00Z           |
      | minutes      | 15                              |
      | deadline     | 2025-09-16T09:15:00Z           |
    And I want to set up repeat settings:
      | is_repeating     | true                           |
      | repeat_pattern   | weekly                         |
      | repeat_interval  | 1                              |
      | repeat_days_of_week| tue,thu                      |
      | repeat_count     | 10                             |
    And I want to generate reminders:
      | generate_reminders      | true                    |
      | reminder_advance_minutes| 10                      |
      | reminder_methods        | webhook                 |
      | reminder_targets        | https://hooks.slack.com/standup |
    When I create the task
    Then the task should be created successfully
    And task instances should be created only for Tuesdays and Thursdays
    And each instance should have a webhook reminder 10 minutes before
    And the total number of instances should not exceed 10

  Scenario: End date limits for repeating tasks
    Given I want to create a task with the following details:
      | name         | Limited Task                    |
      | schedule_time| 2025-09-17T10:00:00Z           |
      | minutes      | 30                              |
      | deadline     | 2025-09-17T10:30:00Z           |
    And I want to set up repeat settings:
      | is_repeating    | true                            |
      | repeat_pattern  | daily                           |
      | repeat_interval | 1                               |
      | repeat_end_date | 2025-09-20T00:00:00Z           |
    And I want to generate reminders:
      | generate_reminders| true                           |
      | reminder_methods  | email                          |
    When I create the task
    Then the task should be created successfully
    And only 3 task instances should be created
    And only 3 reminders should be created
    And no instances should be created after 2025-09-20

  Scenario: Reminder service integration validation
    Given I want to create a task with valid reminder settings
    But the reminder service returns an error
    When I create the task
    Then the task should be created successfully
    But I should receive a warning about reminder creation failure
    And the task should still be functional without reminders

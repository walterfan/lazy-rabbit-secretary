Feature: User Registration
  As a website visitor
  I want to register an account
  So that I can use the website

  Scenario: Successful registration
    Given a clean user repository
    And an email "test@example.com" doesn't exist
    When I register with email "test@example.com" and password "SecurePass123!"
    Then the response status should be 201
    And the user "test@example.com" should exist with status "pending"

  Scenario: Registration with existing email
    Given a clean user repository
    And an email "existing@example.com" already exists
    When I register with email "existing@example.com" and password "SecurePass123!"
    Then the response status should be 409
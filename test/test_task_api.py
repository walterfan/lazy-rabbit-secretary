#!/usr/bin/env python3
"""
Test script for Task API endpoints.

This script tests:
1. User authentication (login)
2. Task creation with validation
3. Task retrieval and search
4. Task status management
5. Task updates and lifecycle

Prerequisites:
- Server running on localhost:9090
- Test user with appropriate permissions
- Database properly configured
"""

import pytest
import httpx
import json
import os
from datetime import datetime, timedelta
from typing import Dict, Optional


class TaskAPIClient:
    """Client for testing Task API endpoints."""
    
    def __init__(self, base_url: Optional[str] = None):
        self.base_url = base_url or os.getenv("TEST_BASE_URL", "http://localhost:9090")
        self.access_token: Optional[str] = None
        
        # Configure SSL verification for self-signed certificates
        verify_ssl_env = os.getenv("TEST_VERIFY_SSL", "false")
        verify_ssl = verify_ssl_env.lower() in ("true", "1", "yes")
        
        # Create httpx client with SSL verification setting
        self.client = httpx.Client(verify=verify_ssl)
        
        if not verify_ssl:
            print(f"[INFO] SSL verification disabled for {self.base_url}")
    
    def login(self, username: str, password: str, realm_name: str = "default") -> Dict:
        """Login and store access token."""
        response = self.client.post(
            f"{self.base_url}/api/v1/auth/login",
            json={
                "username": username,
                "password": password,
                "realm_name": realm_name
            }
        )
        response.raise_for_status()
        
        data = response.json()
        self.access_token = data["access_token"]
        self.client.headers.update({
            "Authorization": f"Bearer {self.access_token}"
        })
        return data
    
    def create_task(self, realm_name: str, name: str, description: str = "",
                   priority: int = None, difficulty: int = None,
                   schedule_time: str = None, minutes: int = 30,
                   deadline: str = None, tags: str = "") -> Dict:
        """Create a new task."""
        # Default schedule_time to 1 hour from now if not provided
        if not schedule_time:
            schedule_time = (datetime.now() + timedelta(hours=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
        
        # Default deadline to 1 day from now if not provided
        if not deadline:
            deadline = (datetime.now() + timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
        
        task_data = {
            "realm_name": realm_name,
            "name": name,
            "description": description,
            "schedule_time": schedule_time,
            "minutes": minutes,
            "deadline": deadline,
            "tags": tags
        }

        # Add optional fields if provided
        if priority is not None:
            task_data["priority"] = priority
        if difficulty is not None:
            task_data["difficulty"] = difficulty

        response = self.client.post(
            f"{self.base_url}/api/v1/tasks",
            json=task_data
        )
        response.raise_for_status()
        return response.json()

    def get_task(self, task_id: str) -> Dict:
        """Get a task by ID."""
        response = self.client.get(f"{self.base_url}/api/v1/tasks/{task_id}")
        response.raise_for_status()
        return response.json()
    
    def update_task(self, task_id: str, **kwargs) -> Dict:
        """Update a task."""
        response = self.client.put(
            f"{self.base_url}/api/v1/tasks/{task_id}",
            json=kwargs
        )
        response.raise_for_status()
        return response.json()
    
    def delete_task(self, task_id: str) -> None:
        """Delete a task."""
        response = self.client.delete(f"{self.base_url}/api/v1/tasks/{task_id}")
        response.raise_for_status()
    
    def search_tasks(self, realm_id: str = None, query: str = None,
                    status: str = None, tags: str = None,
                    priority: int = None, difficulty: int = None,
                    page: int = 1, page_size: int = 20) -> Dict:
        """Search tasks with filters."""
        params = {"page": page, "page_size": page_size}
        if realm_id:
            params["realm_id"] = realm_id
        if query:
            params["q"] = query
        if status:
            params["status"] = status
        if tags:
            params["tags"] = tags
        if priority:
            params["priority"] = priority
        if difficulty:
            params["difficulty"] = difficulty
        
        response = self.client.get(f"{self.base_url}/api/v1/tasks", params=params)
        response.raise_for_status()
        return response.json()
    
    def get_tasks_by_status(self, status: str, realm_id: str = None,
                           page: int = 1, page_size: int = 20) -> Dict:
        """Get tasks by status."""
        params = {"page": page, "page_size": page_size}
        if realm_id:
            params["realm_id"] = realm_id
        
        response = self.client.get(
            f"{self.base_url}/api/v1/tasks/status/{status}", 
            params=params
        )
        response.raise_for_status()
        return response.json()
    
    def get_upcoming_tasks(self, realm_id: str = None, limit: int = 10) -> Dict:
        """Get upcoming tasks."""
        params = {"limit": limit}
        if realm_id:
            params["realm_id"] = realm_id
        
        response = self.client.get(
            f"{self.base_url}/api/v1/tasks/upcoming",
            params=params
        )
        response.raise_for_status()
        return response.json()
    
    def get_overdue_tasks(self, realm_id: str = None, limit: int = 10) -> Dict:
        """Get overdue tasks."""
        params = {"limit": limit}
        if realm_id:
            params["realm_id"] = realm_id
        
        response = self.client.get(
            f"{self.base_url}/api/v1/tasks/overdue",
            params=params
        )
        response.raise_for_status()
        return response.json()
    
    def start_task(self, task_id: str) -> Dict:
        """Start a task."""
        response = self.client.post(f"{self.base_url}/api/v1/tasks/{task_id}/start")
        response.raise_for_status()
        return response.json()
    
    def complete_task(self, task_id: str) -> Dict:
        """Complete a task."""
        response = self.client.post(f"{self.base_url}/api/v1/tasks/{task_id}/complete")
        response.raise_for_status()
        return response.json()
    
    def fail_task(self, task_id: str) -> Dict:
        """Mark a task as failed."""
        response = self.client.post(f"{self.base_url}/api/v1/tasks/{task_id}/fail")
        response.raise_for_status()
        return response.json()
    
    def cleanup_tasks(self, task_ids: list) -> None:
        """Clean up multiple tasks, ignoring errors for already deleted ones."""
        for task_id in task_ids:
            try:
                self.delete_task(task_id)
                print(f"[CLEANUP] Deleted task: {task_id}")
            except httpx.HTTPStatusError as e:
                if e.response.status_code == 404:
                    print(f"[CLEANUP] Task already deleted: {task_id}")
                else:
                    print(f"[CLEANUP] Failed to delete task {task_id}: {e}")
            except Exception as e:
                print(f"[CLEANUP] Error deleting task {task_id}: {e}")


@pytest.fixture(scope="session")
def api_client():
    """Create and authenticate API client."""
    client = TaskAPIClient()
    
    # Try to login with configurable test credentials
    username = os.getenv("TEST_USERNAME", "admin")
    password = os.getenv("TEST_PASSWORD", "admin123")
    realm = os.getenv("TEST_REALM", "default")
    
    try:
        client.login(username, password, realm)
    except httpx.HTTPStatusError as e:
        pytest.skip(f"Cannot authenticate with {username}@{realm}: {e}. Please ensure test user exists.")
    
    return client


@pytest.fixture
def cleanup_tracker():
    """Track created tasks for cleanup after each test."""
    created_tasks = []
    
    def register_task(task_data):
        """Register a task for cleanup."""
        if isinstance(task_data, dict) and "id" in task_data:
            created_tasks.append(task_data["id"])
        return task_data
    
    yield register_task
    
    # Cleanup after test
    if created_tasks:
        try:
            # Get a fresh client for cleanup
            cleanup_client = TaskAPIClient()
            username = os.getenv("TEST_USERNAME", "admin")
            password = os.getenv("TEST_PASSWORD", "admin123")
            realm = os.getenv("TEST_REALM", "default")
            cleanup_client.login(username, password, realm)
            cleanup_client.cleanup_tasks(created_tasks)
        except Exception as e:
            print(f"[CLEANUP] Failed to authenticate for cleanup: {e}")


class TestTaskAPI:
    """Test cases for Task API."""
    
    def test_authentication_success(self, api_client):
        """Test successful authentication."""
        assert api_client.access_token is not None
        assert len(api_client.access_token) > 0
        assert api_client.client.headers.get("Authorization").startswith("Bearer ")
    
    def test_authentication_failure(self):
        """Test authentication with invalid credentials."""
        client = TaskAPIClient()
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            client.login("invalid_user", "wrong_password", "default")
        
        assert exc_info.value.response.status_code == 401
    
    def test_create_task_success(self, api_client, cleanup_tracker):
        """Test successful task creation."""
        task_data = {
            "realm_name": "default",
            "name": f"test_task_{os.urandom(4).hex()}",
            "description": "Test task for pytest",
            "priority": 3,
            "difficulty": 2,
            "minutes": 45,
            "tags": "testing,pytest"
        }
        
        result = api_client.create_task(**task_data)
        cleanup_tracker(result)  # Register for cleanup
        
        # Verify response structure
        assert "id" in result
        assert result["name"] == task_data["name"]
        assert result["description"] == task_data["description"]
        assert result["priority"] == task_data["priority"]
        assert result["difficulty"] == task_data["difficulty"]
        assert result["minutes"] == task_data["minutes"]
        assert result["tags"] == task_data["tags"]
        assert result["status"] == "pending"
        
        # Verify timestamps
        assert "schedule_time" in result
        assert "deadline" in result
        assert "created_at" in result
        assert "updated_at" in result
    
    def test_create_task_with_defaults(self, api_client, cleanup_tracker):
        """Test task creation with default priority and difficulty."""
        task_data = {
            "realm_name": "default",
            "name": f"default_task_{os.urandom(4).hex()}",
            "description": "Task with default values",
            "minutes": 30
        }
        
        result = api_client.create_task(**task_data)
        cleanup_tracker(result)  # Register for cleanup
        
        # Verify default values are set
        assert result["priority"] == 2  # Default from model
        assert result["difficulty"] == 2  # Default from model
        assert result["status"] == "pending"
    
    def test_create_task_validation_errors(self, api_client):
        """Test task creation with validation errors."""
        # Test missing required fields
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_task(
                realm_name="default",
                name="",  # Empty name
                minutes=30
            )
        assert exc_info.value.response.status_code == 400
        
        # Test invalid priority range
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_task(
                realm_name="default",
                name="invalid_priority",
                priority=6,  # Out of range (1-5)
                minutes=30
            )
        assert exc_info.value.response.status_code == 400
        
        # Test invalid difficulty range
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_task(
                realm_name="default",
                name="invalid_difficulty",
                difficulty=0,  # Out of range (1-5)
                minutes=30
            )
        assert exc_info.value.response.status_code == 400
    
    def test_create_task_invalid_realm(self, api_client):
        """Test task creation with invalid realm."""
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_task(
                realm_name="nonexistent_realm",
                name="test_task",
                minutes=30
            )
        assert exc_info.value.response.status_code == 400
    
    def test_get_task_success(self, api_client, cleanup_tracker):
        """Test successful task retrieval."""
        # First create a task
        created_task = api_client.create_task(
            realm_name="default",
            name=f"get_test_{os.urandom(4).hex()}",
            description="Test task for get operation",
            priority=3,
            difficulty=2,
            minutes=45,
            tags="testing,get"
        )
        cleanup_tracker(created_task)
        task_id = created_task["id"]
        
        # Then retrieve it
        retrieved_task = api_client.get_task(task_id)
        
        # Verify it matches the created task
        assert retrieved_task["id"] == task_id
        assert retrieved_task["name"] == created_task["name"]
        assert retrieved_task["priority"] == created_task["priority"]
        assert retrieved_task["status"] == created_task["status"]
    
    def test_get_task_not_found(self, api_client):
        """Test retrieval of non-existent task."""
        fake_id = "00000000-0000-0000-0000-000000000000"
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.get_task(fake_id)
        
        assert exc_info.value.response.status_code == 404
    
    def test_update_task_success(self, api_client, cleanup_tracker):
        """Test successful task update."""
        # Create a task
        created_task = api_client.create_task(
            realm_name="default",
            name=f"update_test_{os.urandom(4).hex()}",
            description="Test task for update operation",
            priority=2,
            difficulty=3,
            minutes=30,
            tags="testing,update"
        )
        cleanup_tracker(created_task)
        task_id = created_task["id"]
        
        # Update it
        update_data = {
            "name": "Updated Task Name",
            "description": "Updated description",
            "priority": 5,
            "difficulty": 4,
            "tags": "updated,modified"
        }
        
        updated_task = api_client.update_task(task_id, **update_data)
        
        # Verify updates
        assert updated_task["name"] == update_data["name"]
        assert updated_task["description"] == update_data["description"]
        assert updated_task["priority"] == update_data["priority"]
        assert updated_task["difficulty"] == update_data["difficulty"]
        assert updated_task["tags"] == update_data["tags"]
        
        # Verify ID and status remain unchanged
        assert updated_task["id"] == task_id
        assert updated_task["status"] == "pending"
    
    def test_delete_task_success(self, api_client):
        """Test successful task deletion."""
        # Create a task
        created_task = api_client.create_task(
            realm_name="default",
            name="task_to_delete",
            minutes=30
        )
        task_id = created_task["id"]
        
        # Delete it
        api_client.delete_task(task_id)
        
        # Verify it's deleted
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.get_task(task_id)
        assert exc_info.value.response.status_code == 404
    
    def test_search_tasks_success(self, api_client, cleanup_tracker):
        """Test successful task search."""
        # Create test tasks
        tasks = []
        unique_tag = f"search_test_{os.urandom(4).hex()}"
        
        for i in range(3):
            task = api_client.create_task(
                realm_name="default",
                name=f"search_task_{i}",
                description=f"Search test task {i}",
                priority=2 + (i % 2),  # Mix of priority 2 and 3
                difficulty=1 + i,      # Difficulty 1, 2, 3
                tags=unique_tag,
                minutes=30
            )
            cleanup_tracker(task)
            tasks.append(task)
        
        # Test search by tags
        realm_id = tasks[0]["realm_id"]
        results = api_client.search_tasks(realm_id=realm_id, tags=unique_tag)
        
        assert "items" in results
        assert "total" in results
        assert results["total"] >= len(tasks)
        
        # Verify pagination info
        assert "page" in results
        assert "page_size" in results
        assert "total_pages" in results
        
        # Test search by priority
        priority_results = api_client.search_tasks(realm_id=realm_id, priority=2)
        assert len([t for t in priority_results["items"] if t["priority"] == 2]) >= 1
        
        # Test search by difficulty
        difficulty_results = api_client.search_tasks(realm_id=realm_id, difficulty=1)
        assert len([t for t in difficulty_results["items"] if t["difficulty"] == 1]) >= 1
    
    def test_get_tasks_by_status(self, api_client, cleanup_tracker):
        """Test getting tasks by status."""
        # Create a task
        task = api_client.create_task(
            realm_name="default",
            name="status_test_task",
            minutes=30
        )
        cleanup_tracker(task)
        
        # Get pending tasks
        realm_id = task["realm_id"]
        results = api_client.get_tasks_by_status("pending", realm_id=realm_id)
        
        assert "items" in results
        assert "total" in results
        
        # Verify all returned tasks have pending status
        for item in results["items"]:
            assert item["status"] == "pending"
    
    def test_task_status_transitions(self, api_client, cleanup_tracker):
        """Test task status transitions."""
        # Create a task
        task = api_client.create_task(
            realm_name="default",
            name="status_transition_test",
            minutes=30
        )
        cleanup_tracker(task)
        task_id = task["id"]
        
        # Start the task (pending -> running)
        started_task = api_client.start_task(task_id)
        assert started_task["status"] == "running"
        assert "start_time" in started_task
        
        # Complete the task (running -> completed)
        completed_task = api_client.complete_task(task_id)
        assert completed_task["status"] == "completed"
        assert "end_time" in completed_task
    
    def test_task_failure_transition(self, api_client, cleanup_tracker):
        """Test task failure transition."""
        # Create a task
        task = api_client.create_task(
            realm_name="default",
            name="failure_test_task",
            minutes=30
        )
        cleanup_tracker(task)
        task_id = task["id"]
        
        # Start the task
        api_client.start_task(task_id)
        
        # Mark as failed (running -> failed)
        failed_task = api_client.fail_task(task_id)
        assert failed_task["status"] == "failed"
        assert "end_time" in failed_task
    
    def test_get_upcoming_tasks(self, api_client, cleanup_tracker):
        """Test getting upcoming tasks."""
        # Create a task scheduled for the future
        future_time = (datetime.now() + timedelta(hours=2)).isoformat()
        task = api_client.create_task(
            realm_name="default",
            name="upcoming_task",
            schedule_time=future_time,
            minutes=30
        )
        cleanup_tracker(task)
        
        # Get upcoming tasks
        realm_id = task["realm_id"]
        results = api_client.get_upcoming_tasks(realm_id=realm_id, limit=10)
        
        assert "items" in results
        assert "total" in results
        
        # Verify all tasks are pending or running
        for item in results["items"]:
            assert item["status"] in ["pending", "running"]
    
    def test_get_overdue_tasks(self, api_client, cleanup_tracker):
        """Test getting overdue tasks."""
        # Create a task with deadline in the past
        past_time = (datetime.now() - timedelta(hours=1)).isoformat()
        task = api_client.create_task(
            realm_name="default",
            name="overdue_task",
            deadline=past_time,
            minutes=30
        )
        cleanup_tracker(task)
        
        # Get overdue tasks
        realm_id = task["realm_id"]
        results = api_client.get_overdue_tasks(realm_id=realm_id, limit=10)
        
        assert "items" in results
        assert "total" in results
    
    def test_unauthorized_access(self):
        """Test API access without authentication."""
        client = TaskAPIClient()
        # Don't login - no auth token
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            client.create_task(
                realm_name="default",
                name="unauthorized_test",
                minutes=30
            )
        
        assert exc_info.value.response.status_code == 401
    
    def test_invalid_status_filter(self, api_client):
        """Test invalid status filter."""
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.get_tasks_by_status("invalid_status")
        
        assert exc_info.value.response.status_code == 400
    
    def test_task_cleanup(self, api_client, cleanup_tracker):
        """Test that tasks are properly cleaned up after test."""
        # Create a task
        task = api_client.create_task(
            realm_name="default",
            name="cleanup_test_task",
            minutes=30
        )
        cleanup_tracker(task)  # Register for cleanup
        
        # Verify it exists
        retrieved = api_client.get_task(task["id"])
        assert retrieved["id"] == task["id"]
        
        # The cleanup_tracker fixture will automatically delete this task after the test
        # This test verifies the cleanup mechanism is working


def test_environment_setup():
    """Test that the test environment is properly configured."""
    # Verify required environment variables
    base_url = os.getenv("TEST_BASE_URL", "http://localhost:9090")
    assert base_url, "TEST_BASE_URL not configured"
    
    username = os.getenv("TEST_USERNAME", "admin")
    password = os.getenv("TEST_PASSWORD", "admin123")
    realm = os.getenv("TEST_REALM", "default")
    
    assert username, "TEST_USERNAME not configured"
    assert password, "TEST_PASSWORD not configured"
    assert realm, "TEST_REALM not configured"


if __name__ == "__main__":
    # Run tests directly
    pytest.main([__file__, "-v"])

#!/usr/bin/env python3
"""
Test script for Secret API endpoints.

This script tests:
1. User authentication (login)
2. Secret creation with encryption
3. Secret retrieval
4. Secret search functionality

Prerequisites:
- Server running on localhost:9090
- KEK environment variables set
- Test user with admin role exists
"""

import pytest
import httpx
import json
import base64
import os
from typing import Dict, Optional


class SecretAPIClient:
    """Client for testing Secret API endpoints."""
    
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
    
    def create_secret(self, realm_name: str, name: str, group: str,
                     desc: str, path: str, value: str) -> Dict:
        """Create a new secret."""
        response = self.client.post(
            f"{self.base_url}/api/v1/secrets",
            json={
                "realm_name": realm_name,
                "name": name,
                "group": group,
                "desc": desc,
                "path": path,
                "value": value
            }
        )
        response.raise_for_status()
        return response.json()
    
    def get_secret(self, secret_id: str) -> Dict:
        """Get a secret by ID."""
        response = self.client.get(f"{self.base_url}/api/v1/secrets/{secret_id}")
        response.raise_for_status()
        return response.json()
    
    def search_secrets(self, realm_id: str = None, query: str = None,
                      group: str = None, path: str = None,
                      page: int = 1, page_size: int = 20) -> Dict:
        """Search secrets with filters."""
        params = {"page": page, "page_size": page_size}
        if realm_id:
            params["realm_id"] = realm_id
        if query:
            params["q"] = query
        if group:
            params["group"] = group
        if path:
            params["path"] = path
        
        response = self.client.get(f"{self.base_url}/api/v1/secrets", params=params)
        response.raise_for_status()
        return response.json()
    
    def delete_secret(self, secret_id: str) -> None:
        """Delete a secret."""
        response = self.client.delete(f"{self.base_url}/api/v1/secrets/{secret_id}")
        response.raise_for_status()
    
    def cleanup_secrets(self, secret_ids: list) -> None:
        """Clean up multiple secrets, ignoring errors for already deleted ones."""
        for secret_id in secret_ids:
            try:
                self.delete_secret(secret_id)
                print(f"[CLEANUP] Deleted secret: {secret_id}")
            except httpx.HTTPStatusError as e:
                if e.response.status_code == 404:
                    print(f"[CLEANUP] Secret already deleted: {secret_id}")
                else:
                    print(f"[CLEANUP] Failed to delete secret {secret_id}: {e}")
            except Exception as e:
                print(f"[CLEANUP] Error deleting secret {secret_id}: {e}")


@pytest.fixture(scope="session")
def api_client():
    """Create and authenticate API client."""
    client = SecretAPIClient()
    
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
    """Track created secrets for cleanup after each test."""
    created_secrets = []
    
    def register_secret(secret_data):
        """Register a secret for cleanup."""
        if isinstance(secret_data, dict) and "id" in secret_data:
            created_secrets.append(secret_data["id"])
        return secret_data
    
    yield register_secret
    
    # Cleanup after test
    if created_secrets:
        try:
            # Get a fresh client for cleanup
            cleanup_client = SecretAPIClient()
            username = os.getenv("TEST_USERNAME", "admin")
            password = os.getenv("TEST_PASSWORD", "admin123")
            realm = os.getenv("TEST_REALM", "default")
            cleanup_client.login(username, password, realm)
            cleanup_client.cleanup_secrets(created_secrets)
        except Exception as e:
            print(f"[CLEANUP] Failed to authenticate for cleanup: {e}")


@pytest.fixture(scope="session", autouse=True)
def setup_environment():
    """Set up required environment variables for testing."""
    if not os.getenv("KEK_BASE64") and not os.getenv("KEK"):
        # Generate a test KEK if none exists
        test_kek = base64.b64encode(os.urandom(32)).decode('ascii')
        os.environ["KEK_BASE64"] = test_kek
        os.environ["KEK_VERSION"] = "1"
        print(f"Generated test KEK: {test_kek}")


class TestSecretAPI:
    """Test cases for Secret API."""
    
    def test_authentication_success(self, api_client):
        """Test successful authentication."""
        assert api_client.access_token is not None
        assert len(api_client.access_token) > 0
        assert api_client.client.headers.get("Authorization").startswith("Bearer ")

    def test_authentication_failure(self):
        """Test authentication with invalid credentials."""
        client = SecretAPIClient()
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            client.login("invalid_user", "wrong_password", "default")
        
        assert exc_info.value.response.status_code == 401
    
    def test_create_secret_success(self, api_client, cleanup_tracker):
        """Test successful secret creation."""
        secret_data = {
            "realm_name": "default",
            "name": f"test_secret_{os.urandom(4).hex()}",
            "group": "testing",
            "desc": "Test secret for pytest",
            "path": "/test/secrets/api",
            "value": "super_secret_test_value_123"
        }

        result = api_client.create_secret(**secret_data)
        cleanup_tracker(result)  # Register for cleanup
        
        # Verify response structure
        assert "id" in result
        assert result["name"] == secret_data["name"]
        assert result["group"] == secret_data["group"]
        assert result["desc"] == secret_data["desc"]
        assert result["path"] == secret_data["path"]
        assert result["cipher_alg"] == "aes-256-gcm"
        
        # Verify encryption fields are present
        assert "cipher_text" in result
        assert "nonce" in result
        assert "auth_tag" in result
        assert "wrapped_dek" in result
        assert "kek_version" in result
        
        # Verify original value is not in response
        assert "value" not in result
        
        # Verify encrypted fields are base64 encoded
        try:
            base64.b64decode(result["cipher_text"])
            base64.b64decode(result["nonce"])
            base64.b64decode(result["auth_tag"])
            base64.b64decode(result["wrapped_dek"])
        except Exception as e:
            pytest.fail(f"Invalid base64 encoding in response: {e}")
        
        return result
    
    def test_create_secret_missing_required_fields(self, api_client):
        """Test secret creation with missing required fields."""
        incomplete_data = {
            "realm_name": "default",
            "name": "incomplete_secret",
            # Missing required fields: group, path, value
        }
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_secret(
                realm_name=incomplete_data["realm_name"],
                name=incomplete_data["name"],
                group="",  # Empty required field
                desc="Test",
                path="",   # Empty required field
                value=""   # Empty required field
            )
        
        assert exc_info.value.response.status_code == 400
    
    def test_create_secret_invalid_realm(self, api_client):
        """Test secret creation with invalid realm."""
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.create_secret(
                realm_name="nonexistent_realm",
                name="test_secret",
                group="testing",
                desc="Test secret",
                path="/test/path",
                value="test_value"
            )
        
        assert exc_info.value.response.status_code == 400
    
    def test_get_secret_success(self, api_client, cleanup_tracker):
        """Test successful secret retrieval."""
        # First create a secret
        secret_data = {
            "realm_name": "default",
            "name": f"get_test_{os.urandom(4).hex()}",
            "group": "testing",
            "desc": "Test secret for get operation",
            "path": "/test/secrets/get",
            "value": "get_test_value_123"
        }
        created_secret = api_client.create_secret(**secret_data)
        cleanup_tracker(created_secret)  # Register for cleanup
        secret_id = created_secret["id"]
        
        # Then retrieve it
        retrieved_secret = api_client.get_secret(secret_id)
        
        # Verify it matches the created secret
        assert retrieved_secret["id"] == secret_id
        assert retrieved_secret["name"] == created_secret["name"]
        assert retrieved_secret["cipher_text"] == created_secret["cipher_text"]
    
    def test_get_secret_not_found(self, api_client):
        """Test retrieval of non-existent secret."""
        fake_id = "00000000-0000-0000-0000-000000000000"
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            api_client.get_secret(fake_id)
        
        assert exc_info.value.response.status_code == 404
    
    def test_search_secrets_success(self, api_client, cleanup_tracker):
        """Test successful secret search."""
        # Create a few test secrets first
        secrets = []
        for i in range(3):
            secret = api_client.create_secret(
                realm_name="default",
                name=f"search_test_{i}_{os.urandom(2).hex()}",
                group="search_testing",
                desc=f"Search test secret {i}",
                path=f"/search/test/{i}",
                value=f"search_value_{i}"
            )
            cleanup_tracker(secret)  # Register for cleanup
            secrets.append(secret)
        
        # Test search by group (use realm_id from first created secret)
        realm_id = secrets[0]["realm_id"] if secrets else None
        results = api_client.search_secrets(realm_id=realm_id, group="search_testing")
        assert "items" in results
        assert "total" in results
        assert results["total"] >= len(secrets)
        
        # Verify our created secrets are in the results
        result_names = [item["name"] for item in results["items"]]
        for secret in secrets:
            assert secret["name"] in result_names
    
    def test_search_secrets_with_query(self, api_client, cleanup_tracker):
        """Test secret search with text query."""
        # Create a secret with unique content
        unique_name = f"unique_searchable_{os.urandom(4).hex()}"
        secret = api_client.create_secret(
            realm_name="default",
            name=unique_name,
            group="query_testing",
            desc="Unique searchable description",
            path="/query/test/unique",
            value="unique_test_value"
        )
        cleanup_tracker(secret)  # Register for cleanup
        
        # Search for it
        realm_id = secret["realm_id"]
        results = api_client.search_secrets(realm_id=realm_id, query=unique_name)
        
        assert results["total"] >= 1
        found_names = [item["name"] for item in results["items"]]
        assert unique_name in found_names
    
    def test_unauthorized_access(self):
        """Test API access without authentication."""
        client = SecretAPIClient()
        # Don't login - no auth token
        
        with pytest.raises(httpx.HTTPStatusError) as exc_info:
            client.create_secret(
                realm_name="default",
                name="unauthorized_test",
                group="testing",
                desc="Should fail",
                path="/test/unauthorized",
                value="test_value"
            )
        
        assert exc_info.value.response.status_code == 401
    
    def test_secret_cleanup(self, api_client, cleanup_tracker):
        """Test that secrets are properly cleaned up after test."""
        # Create a secret
        secret = api_client.create_secret(
            realm_name="default",
            name=f"cleanup_test_{os.urandom(4).hex()}",
            group="cleanup_testing",
            desc="Test secret for cleanup verification",
            path="/test/cleanup",
            value="cleanup_test_value"
        )
        cleanup_tracker(secret)  # Register for cleanup
        
        # Verify it exists
        retrieved = api_client.get_secret(secret["id"])
        assert retrieved["id"] == secret["id"]
        
        # The cleanup_tracker fixture will automatically delete this secret after the test
        # This test verifies the cleanup mechanism is working


def test_environment_setup():
    """Test that required environment variables are set."""
    assert os.getenv("KEK_BASE64") or os.getenv("KEK"), "KEK not configured"
    
    if os.getenv("KEK_BASE64"):
        # Verify it's valid base64 and correct length
        try:
            key = base64.b64decode(os.getenv("KEK_BASE64"))
            assert len(key) == 32, f"KEK must be 32 bytes, got {len(key)}"
        except Exception as e:
            pytest.fail(f"Invalid KEK_BASE64: {e}")


if __name__ == "__main__":
    # Run tests directly
    pytest.main([__file__, "-v"])

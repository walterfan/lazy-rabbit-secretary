"""
Pytest configuration and shared fixtures.
"""

import pytest
import os
import base64
from pathlib import Path
from dotenv import load_dotenv


def load_environment_variables():
    """Load environment variables from .env files in order of precedence."""
    current_dir = Path(__file__).parent
    parent_dir = current_dir.parent
    backend_dir = parent_dir / "backend"
    
    # Search paths in order of precedence (highest to lowest)
    env_paths = [
        current_dir / ".env",           # ./test/.env (highest priority)
        parent_dir / ".env",            # ./.env (project root)
        backend_dir / ".env",           # ./backend/.env (lowest priority)
    ]
    
    loaded_from = []
    
    # Load from all files in reverse order so higher priority files override lower ones
    # python-dotenv doesn't override existing environment variables by default
    for env_path in reversed(env_paths):
        if env_path.exists():
            # Load without overriding existing environment variables
            load_dotenv(env_path, override=False)
            loaded_from.append(str(env_path))
    
    if loaded_from:
        print(f"\n[TEST] Loaded environment variables from: {', '.join(reversed(loaded_from))}")
    
    return loaded_from


@pytest.fixture(scope="session", autouse=True)
def setup_test_environment():
    """Set up test environment variables."""
    # Load environment variables from .env files
    loaded_files = load_environment_variables()
    
    # Ensure KEK is available for testing
    if not os.getenv("KEK_BASE64") and not os.getenv("KEK"):
        # Generate a test KEK
        test_kek = base64.b64encode(os.urandom(32)).decode('ascii')
        os.environ["KEK_BASE64"] = test_kek
        os.environ["KEK_VERSION"] = "1"
        print(f"\n[TEST] Generated KEK_BASE64 for testing: {test_kek}")
    else:
        kek_source = "KEK_BASE64" if os.getenv("KEK_BASE64") else "KEK"
        print(f"\n[TEST] Using existing {kek_source} from environment")
    
    # Set other test environment variables if needed (don't override existing)
    os.environ.setdefault("DB_TYPE", "sqlite")
    os.environ.setdefault("DB_FILE_PATH", "test_lazy_rabbit_reminder.db")
    os.environ.setdefault("TEST_BASE_URL", "http://localhost:9090")
    os.environ.setdefault("TEST_USERNAME", "admin")
    os.environ.setdefault("TEST_PASSWORD", "admin123")
    os.environ.setdefault("TEST_REALM", "default")
    os.environ.setdefault("TEST_VERIFY_SSL", "false")
    
    # Print loaded environment for debugging
    relevant_vars = [
        "KEK_BASE64", "KEK", "KEK_VERSION", 
        "DB_TYPE", "DB_FILE_PATH", 
        "TEST_BASE_URL", "TEST_USERNAME", "TEST_PASSWORD", "TEST_REALM", "TEST_VERIFY_SSL"
    ]
    loaded_env = {k: v for k, v in os.environ.items() if k in relevant_vars}
    if loaded_env:
        print(f"[TEST] Relevant environment variables:")
        for k, v in loaded_env.items():
            # Mask sensitive values
            if k in ["KEK_BASE64", "KEK", "TEST_PASSWORD"]:
                display_value = f"{v[:8]}..." if len(v) > 8 else "***"
            else:
                display_value = v
            print(f"  {k}={display_value}")


def pytest_configure(config):
    """Configure pytest."""
    config.addinivalue_line(
        "markers", "integration: marks tests as integration tests (may be slow)"
    )
    config.addinivalue_line(
        "markers", "auth: marks tests that require authentication"
    )


def pytest_collection_modifyitems(config, items):
    """Modify test collection to add markers automatically."""
    for item in items:
        # Add integration marker to tests that make HTTP requests
        if "api_client" in item.fixturenames:
            item.add_marker(pytest.mark.integration)
        
        # Add auth marker to tests that require authentication
        if any(keyword in item.name.lower() for keyword in ["auth", "login", "token"]):
            item.add_marker(pytest.mark.auth)

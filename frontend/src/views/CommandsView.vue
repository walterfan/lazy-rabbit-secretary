<template>
  <div class="container mt-4 mb-5">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2 class="mb-0">AI Command Interface</h2>
      <div class="badge bg-secondary">v0.2</div>
    </div>
    
    <hr class="mb-4" />
    
    <div class="row g-3">
      <!-- Command Configuration Card -->
      <div class="col-md-6">
        <div class="card mb-3">
          <div class="card-header bg-primary text-white">
            <h5 class="mb-0"><i class="fas fa-terminal"></i> Command Configuration</h5>
          </div>
          <div class="card-body">
            <div class="mb-3">
              <label class="form-label">Command Name</label>
              <select v-model="formData.name" class="form-select">
                <option value="">Select a command...</option>
                <option v-for="(cmd, index) in commands" :key="cmd.name" :value="cmd.name">
                  {{ index + 1 }}. {{ cmd.name }} — {{ cmd.desc }}
                </option>
              </select>
            </div>

            <div class="mb-3">
              <label class="form-label">Question or Request</label>
              <textarea
                v-model="formData.question"
                rows="3"
                class="form-control"
                placeholder="Enter your question or request..."
              ></textarea>
            </div>

            <div class="mb-3">
              <label class="form-label">Additional Parameters (Optional)</label>
              <textarea
                v-model="formData.parameters"
                rows="4"
                class="form-control"
                placeholder="Any specific requirements or context..."
              ></textarea>
              <small class="form-text text-muted">
                <i class="fas fa-info-circle"></i> JSON format parameters for the command
              </small>
            </div>
          </div>
        </div>
      </div>

      <!-- Settings Configuration Card -->
      <div class="col-md-6">
        <div class="card mb-3">
          <div class="card-header bg-success text-white">
            <h5 class="mb-0"><i class="fas fa-cog"></i> Application Settings</h5>
          </div>
          <div class="card-body">
            <div class="row g-3">
              <div class="col-md-6">
                <div class="d-flex flex-column align-items-center text-center">
                  <label class="form-label mb-2">
                    <i class="fas fa-stream"></i> Stream
                  </label>
                  <div class="form-check form-switch">
                    <input
                      class="form-check-input"
                      type="checkbox"
                      v-model="formData.stream"
                      id="streamSwitch"
                    />
                    <label class="form-check-label" for="streamSwitch">
                      Response
                    </label>
                  </div>
                </div>
              </div>

              <div class="col-md-6">
                <div class="d-flex flex-column align-items-center text-center">
                  <label class="form-label mb-2">
                    <i class="fas fa-memory"></i> Memory
                  </label>
                  <div class="form-check form-switch">
                    <input
                      class="form-check-input"
                      type="checkbox"
                      v-model="formData.remember"
                      id="rememberSwitch"
                    />
                    <label class="form-check-label" for="rememberSwitch">
                      Remember
                    </label>
                  </div>
                </div>
              </div>

              <div class="col-12">
                <div class="d-flex flex-column">
                  <label class="form-label mb-2">
                    <i class="fas fa-language"></i> Output Languages
                  </label>
                  <div class="d-flex gap-3 align-items-center">
                    <div class="form-check">
                      <input
                        class="form-check-input"
                        type="checkbox"
                        v-model="formData.output_languages"
                        value="english"
                        id="langEnglish"
                      />
                      <label class="form-check-label" for="langEnglish">
                        🇺🇸 English
                      </label>
                    </div>
                    <div class="form-check">
                      <input
                        class="form-check-input"
                        type="checkbox"
                        v-model="formData.output_languages"
                        value="chinese"
                        id="langChinese"
                      />
                      <label class="form-check-label" for="langChinese">
                        🇨🇳 Chinese
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="d-grid mb-4">
      <button @click="handleSubmit" class="btn btn-primary btn-lg" :disabled="isLoading || !formData.name">
        <span v-if="isLoading">
          <i class="fas fa-spinner fa-spin"></i> {{ loadingMessage }}
        </span>
        <span v-else>
          <i class="fas fa-paper-plane"></i> Submit Request
        </span>
      </button>
    </div>

    <!-- AI Assistant Response -->
    <div class="card mb-4">
      <div class="card-header bg-secondary text-white d-flex justify-content-between align-items-center">
        <h5 class="mb-0"><i class="fas fa-robot"></i> AI Assistant Response</h5>
      </div>
      <div class="card-body">
        <pre id="answer-display" class="bg-light p-3 rounded border">{{ answer }}</pre>
      </div>

      <div class="mx-auto" style="width: fit-content">
        <button @click="drawImage" class="btn btn-info btn-sm me-2" :disabled="isDrawing || !answer">
          <span v-if="isDrawing">
            <i class="fas fa-spinner fa-spin"></i> Generating...
          </span>
          <span v-else>
            <i class="fas fa-image"></i> Draw Image
          </span>
        </button>
        <button @click="saveConversation" class="btn btn-success btn-sm me-2" :disabled="isSaving || !answer">
          <span v-if="isSaving">
            <i class="fas fa-spinner fa-spin"></i> Saving...
          </span>
          <span v-else>
            <i class="fas fa-save"></i> Save Answer
          </span>
        </button>
        <button @click="copyAnswer" class="btn btn-secondary btn-sm" :disabled="!answer">
          <i class="fas fa-copy"></i> Copy Answer
        </button>
      </div>

      <div id="imageContainer" class="mb-4"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/authStore';

const authStore = useAuthStore();

// Helper function to get headers with optional authentication
const getHeaders = () => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  if (authStore.token) {
    headers.Authorization = `Bearer ${authStore.token}`;
  }
  
  return headers;
};

// Form data model
const formData = ref({
  question: '',
  name: '',
  parameters: '',
  stream: false,
  remember: false,
  output_languages: [] as string[]
});

// Response and loading states
const answer = ref('');
const isLoading = ref(false);
const isDrawing = ref(false);
const isSaving = ref(false);
const loadingMessage = ref('Processing...');

const commands = ref<{ name: string; desc: string }[]>([]);

// Fetch available commands for demo
const fetchCommands = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/api/v1/commands`, {
      headers: getHeaders()
    });
    const data = await response.json();
    commands.value = data.commands || [];
  } catch (error) {
    console.error('Failed to fetch commands:', error);
    // Fallback demo commands if API is not available
    commands.value = [
      { name: 'start_pomodoro', desc: 'Start a 25-minute focus session' },
      { name: 'create_task', desc: 'Create a new task with AI assistance' },
      { name: 'plan_day', desc: 'Generate an optimized daily schedule' },
      { name: 'analyze_productivity', desc: 'Get productivity insights and suggestions' },
      { name: 'set_reminder', desc: 'Set intelligent reminders for tasks' },
      { name: 'gtd_process', desc: 'Process tasks using GTD methodology' }
    ];
  }
};

// Handle form submission
const handleSubmit = async () => {
  if (!formData.value.name) {
    answer.value = 'Please select a command to execute.';
    return;
  }

  isLoading.value = true;
  loadingMessage.value = 'Processing your request...';
  answer.value = '';
  
  const apiBase = import.meta.env.VITE_API_BASE_URL;
  
  try {
    if (formData.value.stream) {
      // Use WebSocket for streaming
      await handleWebSocketRequest();
    } else {
      // Use HTTP for regular request
      const response = await fetch(`${apiBase}/api/v1/commands`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({
          name: formData.value.name,
          question: formData.value.question,
          parameters: formData.value.parameters
        })
      });

      const result = await response.json();
      answer.value = result.result || JSON.stringify(result, null, 2);
    }
  } catch (error) {
    console.error('Request failed:', error);
    answer.value = `Error: ${error instanceof Error ? error.message : 'Request failed'}`;
  } finally {
    isLoading.value = false;
  }
};

// Handle WebSocket request for streaming
const handleWebSocketRequest = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;
  const wsUrl = apiBase.replace('http', 'ws') + '/api/v1/commands/ws';
  
  return new Promise<void>((resolve, reject) => {
    const ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
      ws.send(JSON.stringify({
        name: formData.value.name,
        question: formData.value.question,
        parameters: formData.value.parameters
      }));
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.result) {
        answer.value = data.result;
      } else if (data.error) {
        answer.value = `Error: ${data.error}`;
      }
    };
    
    ws.onclose = () => {
      resolve();
    };
    
    ws.onerror = (error) => {
      reject(error);
    };
  });
};

// Action handlers
const drawImage = async () => {
  isDrawing.value = true;
  try {
    // TODO: Implement image drawing functionality
    console.log('Drawing image from answer:', answer.value);
  } catch (error) {
    console.error('Failed to draw image:', error);
  } finally {
    isDrawing.value = false;
  }
};

const saveConversation = async () => {
  isSaving.value = true;
  try {
    // TODO: Implement conversation saving functionality
    console.log('Saving conversation:', { formData: formData.value, answer: answer.value });
  } catch (error) {
    console.error('Failed to save conversation:', error);
  } finally {
    isSaving.value = false;
  }
};

const copyAnswer = async () => {
  try {
    await navigator.clipboard.writeText(answer.value);
    // TODO: Show success message
    console.log('Answer copied to clipboard');
  } catch (error) {
    console.error('Failed to copy answer:', error);
  }
};

// Initialize component
onMounted(() => {
  fetchCommands();
});
</script>

<style scoped>
/* Card styling */
.card {
  border-radius: 10px;
  border: none;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.card-header {
  border-radius: 10px 10px 0 0 !important;
  border-bottom: none;
}

/* Form enhancements */
.form-control, .form-select {
  border-radius: 8px;
  border: 1px solid #dee2e6;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

.form-control:focus, .form-select:focus {
  border-color: #0d6efd;
  box-shadow: 0 0 0 0.2rem rgba(13, 110, 253, 0.25);
}

.form-label {
  font-weight: 500;
  color: #495057;
  margin-bottom: 0.5rem;
}

.form-label i {
  margin-right: 0.5rem;
  color: #6c757d;
}

/* Button enhancements */
.btn-lg {
  padding: 0.75rem 2rem;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #0d6efd 0%, #6610f2 100%);
  border: none;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(13, 110, 253, 0.3);
}

/* Response display */
#answer-display {
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 400px;
  overflow-y: auto;
}

/* Action buttons */
.btn-sm {
  border-radius: 6px;
  font-size: 0.875rem;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .btn-lg {
    padding: 0.5rem 1.5rem;
    font-size: 1rem;
  }
  
  .card-body {
    padding: 1rem;
  }
  
  #answer-display {
    font-size: 0.8rem;
    max-height: 300px;
  }
}

/* Loading animation */
.fa-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Form switch styling */
.form-check-input:checked {
  background-color: #0d6efd;
  border-color: #0d6efd;
}

/* Badge styling */
.badge {
  font-size: 0.75rem;
  padding: 0.375rem 0.75rem;
}
</style>

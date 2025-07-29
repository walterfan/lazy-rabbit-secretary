<template>
  <div class="container mt-4">

    <!-- News Prompt Section -->
    <div class="alert alert-info mb-4" v-if="newsPrompt">
      {{ newsPrompt }}
    </div>

    <!-- Form Section -->
    <div class="card mb-4">
      <div class="card-body">
        <h2>What can I do for you?</h2>
        <form @submit.prevent="handleSubmit">


          <!-- Select Box -->
          <div class="mb-3">
            <label for="selectOption" class="form-label">Command</label>
            <select id="selectOption" v-model="formData.selectedOption" class="form-select">
              <option v-for="cmd in commands" :key="cmd.name" :value="cmd.name">
                {{ cmd.name }} — {{ cmd.desc }}
              </option>
            </select>
          </div>

          <!-- Input Text -->
          <div class="mb-3">
            <label for="inputText" class="form-label">Question</label>
            <input type="text" id="inputText" v-model="formData.inputText" class="form-control" />
          </div>

          <!-- Input TextArea -->
          <div class="mb-3">
            <label for="inputTextArea" class="form-label">Parameters</label>
            <textarea id="inputTextArea" v-model="formData.additionalNotes" class="form-control" rows="3"></textarea>
          </div>

          <!-- Output TextArea -->
          <div class="mb-3">
            <label for="outputTextArea" class="form-label">Response</label>
            <textarea id="outputTextArea" v-model="outputResponse" class="form-control" rows="4" readonly></textarea>
          </div>

          <!-- Submit Button -->
          <button type="submit" class="btn btn-primary">Submit</button>
        </form>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue';

// Form data model
const formData = ref({
  inputText: '',
  selectedOption: '',
  additionalNotes: ''
});

// Output response area
const outputResponse = ref('');

// News prompt from API
const newsPrompt = ref<string | null>(null);

const commands = ref<{ name: string; desc: string }[]>([]);
// Simulated API call to fetch news prompt
const fetchNewsPrompt = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/news`); // Replace with your actual API endpoint
    const data = await response.json();
    newsPrompt.value = data.message || 'Welcome to our service!';
  } catch (error) {
    console.error('Failed to fetch news prompt:', error);
    newsPrompt.value = 'Unable to load latest news.';
  }
};

const fetchCommands = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/api/v1/commands`);
    const data = await response.json();

    // Assuming the response is { commands: [...] }
    commands.value = data.commands || [];
  } catch (error) {
    console.error('Failed to fetch commands:', error);
  }
};

// Simulated POST request to backend
const handleSubmit = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;
  try {
    const response = await fetch(`${apiBase}/api/v1/commands`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData.value)
    });

    const result = await response.json();
    outputResponse.value = JSON.stringify(result, null, 2); // Pretty print JSON
  } catch (error) {
    console.error('Request failed:', error);
    outputResponse.value = 'Error occurred while submitting the request.';
  }
};

// Fetch news prompt when component mounts
onMounted(() => {
  fetchNewsPrompt();
  fetchCommands();
});
</script>
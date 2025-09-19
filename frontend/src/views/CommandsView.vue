<template>
  <div class="container mt-4">


    <!-- Demo Section -->
    <div class="demo-section card mb-5">
      <div class="card-body">

        <form @submit.prevent="handleSubmit">
          <!-- Command Selection -->
          <div class="mb-3">
            <label for="selectOption" class="form-label">
              <i class="fas fa-terminal"></i> Choose a Command
            </label>
            <select id="selectOption" v-model="formData.selectedOption" class="form-select">
              <option value="">Select a demo command...</option>
              <option v-for="cmd in commands" :key="cmd.name" :value="cmd.name">
                {{ cmd.name }} — {{ cmd.desc }}
              </option>
            </select>
          </div>

          <!-- Input Text -->
          <div class="mb-3">
            <label for="inputText" class="form-label">
              <i class="fas fa-question-circle"></i> Your Question or Request
            </label>
            <input 
              type="text" 
              id="inputText" 
              v-model="formData.inputText" 
              class="form-control"
              placeholder="e.g., 'Help me plan my day' or 'Start a 25-minute focus session'"
            />
          </div>

          <!-- Parameters -->
          <div class="mb-3">
            <label for="inputTextArea" class="form-label">
              <i class="fas fa-cog"></i> Additional Parameters (Optional)
            </label>
            <textarea 
              id="inputTextArea" 
              v-model="formData.additionalNotes" 
              class="form-control" 
              rows="3"
              placeholder="Any specific requirements or context..."
            ></textarea>
          </div>

          <!-- Response Area -->
          <div class="mb-3">
            <label for="outputTextArea" class="form-label">
              <i class="fas fa-robot"></i> AI Response
            </label>
            <textarea 
              id="outputTextArea" 
              v-model="outputResponse" 
              class="form-control" 
              rows="6" 
              readonly
              placeholder="AI response will appear here..."
            ></textarea>
          </div>

          <!-- Submit Button -->
          <div class="text-center">
            <button type="submit" class="btn btn-primary btn-lg" :disabled="!formData.selectedOption">
              <i class="fas fa-paper-plane"></i> Send Request
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Available Commands Info -->
    <div class="commands-info-section mb-5">
      <h2 class="text-center mb-4">📋 Available Commands</h2>
      <div class="row">
        <div class="col-md-6">
          <div class="card h-100">
            <div class="card-body">
              <h5 class="card-title">
                <i class="fas fa-clock text-danger"></i> Productivity Commands
              </h5>
              <ul class="list-group list-group-flush">
                <li class="list-group-item">
                  <strong>start_pomodoro</strong> - Start a 25-minute focus session
                </li>
                <li class="list-group-item">
                  <strong>create_task</strong> - Create a new task with AI assistance
                </li>
                <li class="list-group-item">
                  <strong>plan_day</strong> - Generate an optimized daily schedule
                </li>
                <li class="list-group-item">
                  <strong>analyze_productivity</strong> - Get productivity insights and suggestions
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="col-md-6">
          <div class="card h-100">
            <div class="card-body">
              <h5 class="card-title">
                <i class="fas fa-cogs text-primary"></i> System Commands
              </h5>
              <ul class="list-group list-group-flush">
                <li class="list-group-item">
                  <strong>set_reminder</strong> - Set intelligent reminders for tasks
                </li>
                <li class="list-group-item">
                  <strong>gtd_process</strong> - Process tasks using GTD methodology
                </li>
                <li class="list-group-item">
                  <strong>get_weather</strong> - Get current weather information
                </li>
                <li class="list-group-item">
                  <strong>help</strong> - Get help and command documentation
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Usage Tips -->
    <div class="tips-section mb-5">
      <h2 class="text-center mb-4">💡 Usage Tips</h2>
      <div class="row">
        <div class="col-lg-4 col-md-6 mb-3">
          <div class="tip-card text-center">
            <div class="tip-icon mb-3">
              <i class="fas fa-lightbulb text-warning"></i>
            </div>
            <h6>Be Specific</h6>
            <p class="text-muted small">Provide clear, specific requests for better AI responses</p>
          </div>
        </div>
        <div class="col-lg-4 col-md-6 mb-3">
          <div class="tip-card text-center">
            <div class="tip-icon mb-3">
              <i class="fas fa-comments text-info"></i>
            </div>
            <h6>Use Context</h6>
            <p class="text-muted small">Add context in the additional parameters field</p>
          </div>
        </div>
        <div class="col-lg-4 col-md-6 mb-3">
          <div class="tip-card text-center">
            <div class="tip-icon mb-3">
              <i class="fas fa-sync text-success"></i>
            </div>
            <h6>Try Different Commands</h6>
            <p class="text-muted small">Experiment with various commands to see different capabilities</p>
          </div>
        </div>
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

const commands = ref<{ name: string; desc: string }[]>([]);

// Fetch available commands for demo
const fetchCommands = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/api/v1/commands`);
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

// Handle demo form submission
const handleSubmit = async () => {
  if (!formData.value.selectedOption) {
    outputResponse.value = 'Please select a command to try the demo.';
    return;
  }

  outputResponse.value = 'Processing your request...';
  
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
    outputResponse.value = JSON.stringify(result, null, 2);
  } catch (error) {
    console.error('Request failed:', error);
    
    // Provide demo responses when API is not available
    const demoResponses = {
      'start_pomodoro': {
        status: 'success',
        message: 'Pomodoro timer started! 🍅',
        data: {
          duration: '25 minutes',
          task: formData.value.inputText || 'Focus session',
          break_reminder: 'You\'ll get a break reminder in 25 minutes',
          tips: ['Minimize distractions', 'Focus on one task', 'Take notes of interruptions']
        }
      },
      'create_task': {
        status: 'success',
        message: 'Task created successfully! ✅',
        data: {
          task_id: 'task_' + Math.random().toString(36).substr(2, 9),
          title: formData.value.inputText || 'New Task',
          priority: 'Medium',
          estimated_time: '30 minutes',
          suggested_schedule: 'Tomorrow at 10:00 AM',
          ai_insights: 'This task pairs well with your morning productivity peak'
        }
      },
      'plan_day': {
        status: 'success',
        message: 'Daily plan generated! 📅',
        data: {
          schedule: [
            { time: '09:00', task: 'Morning review & planning', duration: '15 min' },
            { time: '09:15', task: 'High-priority tasks', duration: '90 min' },
            { time: '10:45', task: 'Break', duration: '15 min' },
            { time: '11:00', task: 'Communication & emails', duration: '45 min' },
            { time: '11:45', task: 'Creative work', duration: '75 min' }
          ],
          productivity_score: '87%',
          ai_recommendation: 'Schedule demanding tasks during your peak hours (9-11 AM)'
        }
      },
      'analyze_productivity': {
        status: 'success',
        message: 'Productivity analysis complete! 📊',
        data: {
          weekly_score: '78%',
          completed_tasks: 23,
          total_focus_time: '12.5 hours',
          peak_hours: '9:00 AM - 11:00 AM',
          suggestions: [
            'Consider batching similar tasks together',
            'Your productivity drops after 3 PM - schedule lighter tasks then',
            'You complete 23% more tasks when using Pomodoro technique'
          ]
        }
      },
      'set_reminder': {
        status: 'success',
        message: 'Reminder set successfully! 🔔',
        data: {
          reminder_id: 'reminder_' + Math.random().toString(36).substr(2, 9),
          task: formData.value.inputText || 'Important task',
          reminder_time: 'Tomorrow at 2:00 PM',
          notification_type: 'Email + Push',
          ai_suggestion: 'This reminder is optimally timed based on your productivity patterns'
        }
      },
      'gtd_process': {
        status: 'success',
        message: 'GTD processing complete! 📋',
        data: {
          inbox_items: 5,
          processed_items: 5,
          next_actions: 3,
          waiting_for: 1,
          someday_maybe: 1,
          gtd_insights: 'Your inbox is clear! Consider reviewing your projects weekly.'
        }
      }
    };

    const selectedCommand = formData.value.selectedOption;
    const demoResponse = demoResponses[selectedCommand as keyof typeof demoResponses] || {
      status: 'demo',
      message: `Demo response for: ${selectedCommand}`,
      data: { note: 'This is a demo response. Connect to the backend for full functionality.' }
    };

    outputResponse.value = JSON.stringify(demoResponse, null, 2);
  }
};

// Initialize component
onMounted(() => {
  fetchCommands();
});
</script>

<style scoped>
/* Page Header */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 3rem 2rem;
  border-radius: 20px;
  margin-bottom: 2rem;
}

.page-header .display-5 {
  font-weight: 700;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
}

/* Demo Section */
.demo-section {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 20px;
  border: none;
}

.demo-section .card-body {
  padding: 2.5rem;
}

/* Button enhancements */
.btn-lg {
  padding: 0.75rem 2rem;
  font-weight: 500;
  border-radius: 50px;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

/* Form enhancements */
.form-control, .form-select {
  border-radius: 10px;
  border: 2px solid #e9ecef;
  padding: 0.75rem 1rem;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

.form-control:focus, .form-select:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 0.2rem rgba(102, 126, 234, 0.25);
}

.form-label {
  font-weight: 600;
  color: #495057;
  margin-bottom: 0.75rem;
}

.form-label i {
  margin-right: 0.5rem;
  color: #667eea;
}

/* Commands Info Cards */
.commands-info-section .card {
  border-radius: 15px;
  border: none;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}

/* Tips Cards */
.tip-card {
  padding: 2rem 1rem;
  border-radius: 15px;
  background: white;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  height: 100%;
}

.tip-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.tip-icon i {
  font-size: 2.5rem;
  transition: transform 0.3s ease;
}

.tip-card:hover .tip-icon i {
  transform: scale(1.1);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .page-header {
    padding: 2rem 1rem;
  }
  
  .page-header .display-5 {
    font-size: 2rem;
  }
  
  .demo-section .card-body {
    padding: 1.5rem;
  }
  
  .btn-lg {
    padding: 0.5rem 1.5rem;
    font-size: 1rem;
  }
  
  .tip-card {
    padding: 1.5rem 1rem;
  }
}

/* Animation for tip cards */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tip-card {
  animation: fadeInUp 0.6s ease-out;
}

.tip-card:nth-child(2) {
  animation-delay: 0.1s;
}

.tip-card:nth-child(3) {
  animation-delay: 0.2s;
}
</style>

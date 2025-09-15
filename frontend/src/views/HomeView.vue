<template>
  <div class="container mt-4">
    <!-- Hero Section -->
    <div class="hero-section text-center mb-5">
      <h1 class="display-4 mb-3">🐰 Lazy Rabbit Reminder</h1>
      <p class="lead text-muted mb-4">
        Your AI-powered productivity companion to boost efficiency and never miss important tasks
      </p>
      <div class="d-flex justify-content-center gap-3">
        <button class="btn btn-primary btn-lg" @click="scrollToFeatures">
          <i class="fas fa-rocket"></i> Explore Features
        </button>
        <button class="btn btn-outline-primary btn-lg" @click="scrollToDemo">
          <i class="fas fa-play"></i> Try Demo
        </button>
      </div>
    </div>

    <!-- Features Section -->
    <div id="features" class="features-section mb-5">
      <h2 class="text-center mb-5">🚀 Boost Your Efficiency With</h2>
      
      <div class="row g-4">
        <!-- Pomodoro Timer -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-clock text-danger" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">🍅 Pomodoro Timer</h5>
              <p class="card-text text-muted">
                Focus with time-boxed work sessions. 25-minute focused work periods followed by short breaks to maximize productivity and prevent burnout.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Enhanced focus</small><br>
                <small class="text-success">✓ Reduced distractions</small><br>
                <small class="text-success">✓ Better time awareness</small>
              </div>
            </div>
          </div>
        </div>

        <!-- GTD Flow Helper -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-project-diagram text-primary" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">📋 GTD Flow Helper</h5>
              <p class="card-text text-muted">
                Get Things Done methodology implementation. Capture, clarify, organize, reflect, and engage with your tasks systematically.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Stress-free productivity</small><br>
                <small class="text-success">✓ Clear mind organization</small><br>
                <small class="text-success">✓ Action-oriented workflow</small>
              </div>
            </div>
          </div>
        </div>

        <!-- Smart Todo Lists -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-list-check text-success" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">✅ Smart Todo Lists</h5>
              <p class="card-text text-muted">
                Intelligent task management with priority levels, difficulty ratings, deadlines, and automatic scheduling based on your preferences.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Priority-based sorting</small><br>
                <small class="text-success">✓ Deadline tracking</small><br>
                <small class="text-success">✓ Progress visualization</small>
              </div>
            </div>
          </div>
        </div>

        <!-- Task Reminding -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-bell text-warning" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">🔔 Smart Reminders</h5>
              <p class="card-text text-muted">
                Never miss important tasks with intelligent reminders. Customizable notifications based on urgency, location, and your daily patterns.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Contextual notifications</small><br>
                <small class="text-success">✓ Multiple reminder types</small><br>
                <small class="text-success">✓ Smart scheduling</small>
              </div>
            </div>
          </div>
        </div>

        <!-- AI Agent -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-robot text-info" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">🤖 AI Agent Assistant</h5>
              <p class="card-text text-muted">
                Your personal AI assistant that can execute tasks, analyze patterns, suggest optimizations, and automate routine workflows.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Task automation</small><br>
                <small class="text-success">✓ Pattern analysis</small><br>
                <small class="text-success">✓ Smart suggestions</small>
              </div>
            </div>
          </div>
        </div>

        <!-- Productivity Analytics -->
        <div class="col-lg-4 col-md-6">
          <div class="feature-card card h-100 shadow-sm border-0">
            <div class="card-body text-center p-4">
              <div class="feature-icon mb-3">
                <i class="fas fa-chart-line text-purple" style="font-size: 3rem;"></i>
              </div>
              <h5 class="card-title">📊 Productivity Insights</h5>
              <p class="card-text text-muted">
                Track your efficiency with detailed analytics. Understand your peak hours, task completion patterns, and productivity trends.
              </p>
              <div class="feature-benefits">
                <small class="text-success">✓ Performance tracking</small><br>
                <small class="text-success">✓ Habit formation</small><br>
                <small class="text-success">✓ Goal achievement</small>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- How It Works Section -->
    <div class="how-it-works-section mb-5">
      <h2 class="text-center mb-5">🎯 How It Works</h2>
      <div class="row">
        <div class="col-lg-3 col-md-6 mb-4">
          <div class="text-center">
            <div class="step-number mb-3">1</div>
            <h6>Capture Tasks</h6>
            <p class="text-muted small">Quickly add tasks, ideas, and reminders to your inbox</p>
          </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-4">
          <div class="text-center">
            <div class="step-number mb-3">2</div>
            <h6>AI Organization</h6>
            <p class="text-muted small">Let AI categorize, prioritize, and schedule your tasks</p>
          </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-4">
          <div class="text-center">
            <div class="step-number mb-3">3</div>
            <h6>Focus & Execute</h6>
            <p class="text-muted small">Use Pomodoro sessions to focus and complete tasks efficiently</p>
          </div>
        </div>
        <div class="col-lg-3 col-md-6 mb-4">
          <div class="text-center">
            <div class="step-number mb-3">4</div>
            <h6>Review & Improve</h6>
            <p class="text-muted small">Analyze your productivity and continuously optimize</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Demo Section -->
    <div id="demo" class="demo-section card mb-5">
      <div class="card-body">
        <h3 class="text-center mb-4">🎮 Try Our Interactive Demo</h3>
        <p class="text-center text-muted mb-4">
          Experience the power of AI-assisted productivity management
        </p>

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

    <!-- News/Updates Section -->
    <div class="alert alert-info" v-if="newsPrompt">
      <i class="fas fa-info-circle"></i> <strong>Latest Update:</strong> {{ newsPrompt }}
    </div>

    <!-- Call to Action -->
    <div class="cta-section text-center py-5 mb-4">
      <h3 class="mb-3">Ready to Transform Your Productivity?</h3>
      <p class="text-muted mb-4">Join thousands of users who have already boosted their efficiency</p>
      <div class="d-flex justify-content-center gap-3">
        <button class="btn btn-success btn-lg">
          <i class="fas fa-user-plus"></i> Get Started Free
        </button>
        <button class="btn btn-outline-secondary btn-lg">
          <i class="fas fa-play-circle"></i> Watch Demo Video
        </button>
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

// Scroll functions for smooth navigation
const scrollToFeatures = () => {
  const element = document.getElementById('features');
  element?.scrollIntoView({ behavior: 'smooth' });
};

const scrollToDemo = () => {
  const element = document.getElementById('demo');
  element?.scrollIntoView({ behavior: 'smooth' });
};

// API call to fetch news prompt
const fetchNewsPrompt = async () => {
  const apiBase = import.meta.env.VITE_API_BASE_URL;

  try {
    const response = await fetch(`${apiBase}/news`);
    const data = await response.json();
    newsPrompt.value = data.message || 'Welcome to Lazy Rabbit Reminder! Your AI-powered productivity companion is ready to help you achieve more with less effort.';
  } catch (error) {
    console.error('Failed to fetch news prompt:', error);
    newsPrompt.value = 'Welcome to Lazy Rabbit Reminder! Boost your productivity with AI-powered task management.';
  }
};

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
  fetchNewsPrompt();
  fetchCommands();
});
</script>

<style scoped>
/* Hero Section */
.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 4rem 2rem;
  border-radius: 20px;
  margin-bottom: 3rem;
}

.hero-section .display-4 {
  font-weight: 700;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
}

.hero-section .lead {
  font-size: 1.25rem;
  opacity: 0.95;
}

/* Feature Cards */
.feature-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border-radius: 15px;
  overflow: hidden;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0,0,0,0.1) !important;
}

.feature-icon i {
  transition: transform 0.3s ease;
}

.feature-card:hover .feature-icon i {
  transform: scale(1.1);
}

.feature-benefits {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e9ecef;
}

/* How It Works Section */
.step-number {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0 auto;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
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

/* Call to Action */
.cta-section {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border-radius: 20px;
  margin: 2rem 0;
}

.cta-section h3 {
  font-weight: 600;
}

/* Custom colors */
.text-purple {
  color: #6f42c1 !important;
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

.btn-success {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.3);
}

.btn-success:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(40, 167, 69, 0.4);
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

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero-section {
    padding: 2rem 1rem;
  }
  
  .hero-section .display-4 {
    font-size: 2.5rem;
  }
  
  .feature-card .card-body {
    padding: 1.5rem;
  }
  
  .demo-section .card-body {
    padding: 1.5rem;
  }
  
  .btn-lg {
    padding: 0.5rem 1.5rem;
    font-size: 1rem;
  }
  
  .d-flex.gap-3 {
    flex-direction: column;
    gap: 1rem !important;
  }
}

/* Animation for feature cards */
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

.feature-card {
  animation: fadeInUp 0.6s ease-out;
}

.feature-card:nth-child(2) {
  animation-delay: 0.1s;
}

.feature-card:nth-child(3) {
  animation-delay: 0.2s;
}

.feature-card:nth-child(4) {
  animation-delay: 0.3s;
}

.feature-card:nth-child(5) {
  animation-delay: 0.4s;
}

.feature-card:nth-child(6) {
  animation-delay: 0.5s;
}
</style>
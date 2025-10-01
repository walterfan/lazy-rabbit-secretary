<template>
  <div class="ai-platform-view">
    <div class="container-fluid">
      <!-- Header -->
      <div class="d-flex justify-content-between align-items-center mb-4">
        <div>
          <h1 class="h2 mb-0">
            <i class="bi bi-cpu me-2 text-primary"></i>
            AI Platform
          </h1>
          <p class="text-muted mb-0">Build and manage AI agents, tools, and workflows</p>
        </div>
        <div class="d-flex gap-2">
          <button
            class="btn btn-outline-info"
            @click="showOverview = !showOverview"
          >
            <i class="bi bi-info-circle me-2"></i>
            {{ showOverview ? 'Hide' : 'Show' }} Overview
          </button>
          <button
            class="btn btn-success"
            @click="showCreateModal = true"
          >
            <i class="bi bi-plus me-2"></i>
            Create New
          </button>
        </div>
      </div>

      <!-- Platform Overview -->
      <div v-if="showOverview" class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">
            <i class="bi bi-diagram-3 me-2"></i>
            AI Platform Overview
          </h5>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="col-md-4">
              <div class="overview-item">
                <div class="overview-icon">
                  <i class="bi bi-gear-wide-connected text-primary"></i>
                </div>
                <h6>AI Tools</h6>
                <p class="small text-muted">Reusable AI-powered utilities for specific tasks like text processing, data analysis, and content generation.</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="overview-item">
                <div class="overview-icon">
                  <i class="bi bi-robot text-success"></i>
                </div>
                <h6>AI Agents</h6>
                <p class="small text-muted">Intelligent autonomous entities that can perform complex tasks and make decisions based on context and goals.</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="overview-item">
                <div class="overview-icon">
                  <i class="bi bi-diagram-2 text-info"></i>
                </div>
                <h6>AI Workflows</h6>
                <p class="small text-muted">Orchestrated sequences of AI tools and agents working together to accomplish complex multi-step processes.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="mb-4">
        <ul class="nav nav-tabs ai-platform-tabs">
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'prompts' }"
              @click="activeTab = 'prompts'"
            >
              <i class="bi bi-chat-dots me-2"></i>
              Prompts
              <span class="badge bg-primary ms-2">{{ promptsCount }}</span>
            </button>
          </li>
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'tools' }"
              @click="activeTab = 'tools'"
            >
              <i class="bi bi-gear-wide-connected me-2"></i>
              AI Tools
              <span class="badge bg-secondary ms-2">{{ toolsCount }}</span>
            </button>
          </li>
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'agents' }"
              @click="activeTab = 'agents'"
            >
              <i class="bi bi-robot me-2"></i>
              AI Agents
              <span class="badge bg-secondary ms-2">{{ agentsCount }}</span>
            </button>
          </li>
          <li class="nav-item">
            <button
              class="nav-link"
              :class="{ active: activeTab === 'workflows' }"
              @click="activeTab = 'workflows'"
            >
              <i class="bi bi-diagram-2 me-2"></i>
              Workflows
              <span class="badge bg-secondary ms-2">{{ workflowsCount }}</span>
            </button>
          </li>
        </ul>
      </div>

      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Prompts Tab -->
        <div v-if="activeTab === 'prompts'" class="tab-pane fade show active">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>
              <i class="bi bi-chat-dots me-2"></i>
              Prompt Library
            </h4>
            <router-link to="/prompts" class="btn btn-primary">
              <i class="bi bi-arrow-right me-2"></i>
              Manage Prompts
            </router-link>
          </div>

          <div class="row g-4">
            <!-- Existing Prompts Preview -->
            <div class="col-lg-6">
              <div class="feature-card h-100">
                <div class="card-header bg-primary text-white">
                  <h6 class="mb-0">
                    <i class="bi bi-collection me-2"></i>
                    Available Prompts
                  </h6>
                </div>
                <div class="card-body">
                  <p class="text-muted">Manage your collection of AI prompts for various tasks and scenarios.</p>
                  <ul class="list-unstyled">
                    <li class="mb-2">
                      <i class="bi bi-check-circle text-success me-2"></i>
                      Create and edit custom prompts
                    </li>
                    <li class="mb-2">
                      <i class="bi bi-check-circle text-success me-2"></i>
                      Organize by categories and tags
                    </li>
                    <li class="mb-2">
                      <i class="bi bi-check-circle text-success me-2"></i>
                      Version control and history
                    </li>
                    <li class="mb-2">
                      <i class="bi bi-check-circle text-success me-2"></i>
                      Test and validate prompts
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            <div class="col-lg-6">
              <div class="feature-card h-100">
                <div class="card-header bg-info text-white">
                  <h6 class="mb-0">
                    <i class="bi bi-lightning me-2"></i>
                    Quick Actions
                  </h6>
                </div>
                <div class="card-body">
                  <div class="d-grid gap-2">
                    <router-link to="/prompts" class="btn btn-outline-primary">
                      <i class="bi bi-plus me-2"></i>
                      Create New Prompt
                    </router-link>
                    <router-link to="/prompts" class="btn btn-outline-secondary">
                      <i class="bi bi-search me-2"></i>
                      Browse Prompt Library
                    </router-link>
                    <button class="btn btn-outline-info" @click="importPrompts">
                      <i class="bi bi-upload me-2"></i>
                      Import Prompts
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- AI Tools Tab -->
        <div v-if="activeTab === 'tools'" class="tab-pane fade show active">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>
              <i class="bi bi-gear-wide-connected me-2"></i>
              AI Tools
            </h4>
            <button class="btn btn-success" @click="createTool">
              <i class="bi bi-plus me-2"></i>
              Create Tool
            </button>
          </div>

          <!-- Coming Soon Content -->
          <div class="coming-soon-section">
            <div class="row g-4">
              <div class="col-lg-4">
                <div class="tool-preview-card">
                  <div class="tool-icon">
                    <i class="bi bi-file-text text-primary"></i>
                  </div>
                  <h6>Text Processing Tools</h6>
                  <p class="small text-muted">Summarization, translation, sentiment analysis, and content optimization tools.</p>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
              
              <div class="col-lg-4">
                <div class="tool-preview-card">
                  <div class="tool-icon">
                    <i class="bi bi-graph-up text-success"></i>
                  </div>
                  <h6>Data Analysis Tools</h6>
                  <p class="small text-muted">Statistical analysis, pattern recognition, and data visualization tools.</p>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
              
              <div class="col-lg-4">
                <div class="tool-preview-card">
                  <div class="tool-icon">
                    <i class="bi bi-code-slash text-info"></i>
                  </div>
                  <h6>Code Generation Tools</h6>
                  <p class="small text-muted">Code generation, documentation, testing, and optimization tools.</p>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- AI Agents Tab -->
        <div v-if="activeTab === 'agents'" class="tab-pane fade show active">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>
              <i class="bi bi-robot me-2"></i>
              AI Agents
            </h4>
            <button class="btn btn-success" @click="createAgent">
              <i class="bi bi-plus me-2"></i>
              Create Agent
            </button>
          </div>

          <!-- Coming Soon Content -->
          <div class="coming-soon-section">
            <div class="row g-4">
              <div class="col-lg-4">
                <div class="agent-preview-card">
                  <div class="agent-avatar">
                    <i class="bi bi-person-workspace text-primary"></i>
                  </div>
                  <h6>Research Assistant</h6>
                  <p class="small text-muted">Autonomous agent that can gather information, analyze sources, and compile research reports.</p>
                  <div class="agent-capabilities">
                    <span class="capability-tag">Web Search</span>
                    <span class="capability-tag">Analysis</span>
                    <span class="capability-tag">Reporting</span>
                  </div>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
              
              <div class="col-lg-4">
                <div class="agent-preview-card">
                  <div class="agent-avatar">
                    <i class="bi bi-pencil-square text-success"></i>
                  </div>
                  <h6>Content Creator</h6>
                  <p class="small text-muted">Intelligent agent for creating and optimizing content across multiple formats and platforms.</p>
                  <div class="agent-capabilities">
                    <span class="capability-tag">Writing</span>
                    <span class="capability-tag">SEO</span>
                    <span class="capability-tag">Multi-format</span>
                  </div>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
              
              <div class="col-lg-4">
                <div class="agent-preview-card">
                  <div class="agent-avatar">
                    <i class="bi bi-calendar-check text-info"></i>
                  </div>
                  <h6>Task Manager</h6>
                  <p class="small text-muted">Smart agent that helps organize, prioritize, and track tasks and projects automatically.</p>
                  <div class="agent-capabilities">
                    <span class="capability-tag">Planning</span>
                    <span class="capability-tag">Scheduling</span>
                    <span class="capability-tag">Tracking</span>
                  </div>
                  <div class="status-badge coming-soon">Coming Soon</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Workflows Tab -->
        <div v-if="activeTab === 'workflows'" class="tab-pane fade show active">
          <div class="d-flex justify-content-between align-items-center mb-3">
            <h4>
              <i class="bi bi-diagram-2 me-2"></i>
              AI Workflows
            </h4>
            <button class="btn btn-success" @click="createWorkflow">
              <i class="bi bi-plus me-2"></i>
              Create Workflow
            </button>
          </div>

          <!-- Coming Soon Content -->
          <div class="coming-soon-section">
            <div class="workflow-builder-preview">
              <div class="workflow-canvas">
                <div class="workflow-node start-node">
                  <i class="bi bi-play-circle"></i>
                  <span>Start</span>
                </div>
                <div class="workflow-arrow">→</div>
                <div class="workflow-node tool-node">
                  <i class="bi bi-gear"></i>
                  <span>AI Tool</span>
                </div>
                <div class="workflow-arrow">→</div>
                <div class="workflow-node agent-node">
                  <i class="bi bi-robot"></i>
                  <span>AI Agent</span>
                </div>
                <div class="workflow-arrow">→</div>
                <div class="workflow-node end-node">
                  <i class="bi bi-check-circle"></i>
                  <span>Complete</span>
                </div>
              </div>
              <div class="workflow-description">
                <h6>Visual Workflow Builder</h6>
                <p class="text-muted">Drag-and-drop interface to create complex AI workflows by connecting tools, agents, and decision points.</p>
              </div>
              <div class="status-badge coming-soon large">Coming Soon</div>
            </div>

            <div class="row g-4 mt-4">
              <div class="col-lg-6">
                <div class="workflow-feature-card">
                  <h6>
                    <i class="bi bi-diagram-3 me-2"></i>
                    Workflow Features
                  </h6>
                  <ul class="list-unstyled">
                    <li><i class="bi bi-check text-success me-2"></i>Visual workflow designer</li>
                    <li><i class="bi bi-check text-success me-2"></i>Conditional logic and branching</li>
                    <li><i class="bi bi-check text-success me-2"></i>Error handling and retries</li>
                    <li><i class="bi bi-check text-success me-2"></i>Real-time monitoring</li>
                    <li><i class="bi bi-check text-success me-2"></i>Version control</li>
                  </ul>
                </div>
              </div>
              
              <div class="col-lg-6">
                <div class="workflow-feature-card">
                  <h6>
                    <i class="bi bi-lightning me-2"></i>
                    Integration Options
                  </h6>
                  <ul class="list-unstyled">
                    <li><i class="bi bi-check text-success me-2"></i>API endpoints and webhooks</li>
                    <li><i class="bi bi-check text-success me-2"></i>Schedule-based execution</li>
                    <li><i class="bi bi-check text-success me-2"></i>Event-triggered workflows</li>
                    <li><i class="bi bi-check text-success me-2"></i>External service integration</li>
                    <li><i class="bi bi-check text-success me-2"></i>Human-in-the-loop approval</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Create Modal -->
      <div
        class="modal fade"
        :class="{ show: showCreateModal }"
        :style="{ display: showCreateModal ? 'block' : 'none' }"
        tabindex="-1"
      >
        <div class="modal-dialog modal-lg">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">
                <i class="bi bi-plus-circle me-2"></i>
                Create New AI Component
              </h5>
              <button
                type="button"
                class="btn-close"
                @click="showCreateModal = false"
              ></button>
            </div>
            <div class="modal-body">
              <div class="row g-3">
                <div class="col-md-6">
                  <div class="create-option-card" @click="goToPrompts">
                    <div class="create-icon">
                      <i class="bi bi-chat-dots text-primary"></i>
                    </div>
                    <h6>Create Prompt</h6>
                    <p class="small text-muted">Design custom AI prompts for specific tasks and scenarios.</p>
                    <div class="status-badge available">Available</div>
                  </div>
                </div>
                
                <div class="col-md-6">
                  <div class="create-option-card disabled">
                    <div class="create-icon">
                      <i class="bi bi-gear-wide-connected text-secondary"></i>
                    </div>
                    <h6>Create AI Tool</h6>
                    <p class="small text-muted">Build reusable AI-powered tools for specific functions.</p>
                    <div class="status-badge coming-soon">Coming Soon</div>
                  </div>
                </div>
                
                <div class="col-md-6">
                  <div class="create-option-card disabled">
                    <div class="create-icon">
                      <i class="bi bi-robot text-secondary"></i>
                    </div>
                    <h6>Create AI Agent</h6>
                    <p class="small text-muted">Design intelligent agents with specific roles and capabilities.</p>
                    <div class="status-badge coming-soon">Coming Soon</div>
                  </div>
                </div>
                
                <div class="col-md-6">
                  <div class="create-option-card disabled">
                    <div class="create-icon">
                      <i class="bi bi-diagram-2 text-secondary"></i>
                    </div>
                    <h6>Create Workflow</h6>
                    <p class="small text-muted">Build complex workflows connecting tools and agents.</p>
                    <div class="status-badge coming-soon">Coming Soon</div>
                  </div>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="showCreateModal = false">
                Close
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Backdrop -->
      <div
        v-if="showCreateModal"
        class="modal-backdrop fade show"
        @click="showCreateModal = false"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// Local state
const showOverview = ref(true);
const showCreateModal = ref(false);
const activeTab = ref('prompts');

// Mock data for counts (replace with real data later)
const promptsCount = computed(() => 12); // This should come from prompts store
const toolsCount = computed(() => 0);
const agentsCount = computed(() => 0);
const workflowsCount = computed(() => 0);

// Methods
const createTool = () => {
  // Future implementation
  alert('AI Tools creation will be available in a future update!');
};

const createAgent = () => {
  // Future implementation
  alert('AI Agents creation will be available in a future update!');
};

const createWorkflow = () => {
  // Future implementation
  alert('AI Workflows creation will be available in a future update!');
};

const importPrompts = () => {
  // Future implementation
  alert('Prompt import functionality will be available soon!');
};

const goToPrompts = () => {
  showCreateModal.value = false;
  router.push('/prompts');
};
</script>

<style scoped>
.ai-platform-view {
  min-height: 100vh;
  padding: 2rem 0;
}

/* Overview Section */
.overview-item {
  text-align: center;
  padding: 1.5rem;
}

.overview-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.overview-icon i {
  display: block;
}

/* Navigation Tabs */
.ai-platform-tabs {
  border-bottom: 2px solid #e9ecef;
}

.ai-platform-tabs .nav-link {
  border: none;
  border-bottom: 3px solid transparent;
  background: none;
  color: #6c757d;
  font-weight: 500;
  padding: 1rem 1.5rem;
  transition: all 0.3s ease;
}

.ai-platform-tabs .nav-link:hover {
  border-bottom-color: #dee2e6;
  color: #495057;
  background-color: #f8f9fa;
}

.ai-platform-tabs .nav-link.active {
  color: #007bff;
  border-bottom-color: #007bff;
  background-color: #fff;
}

.ai-platform-tabs .badge {
  font-size: 0.7rem;
}

/* Feature Cards */
.feature-card {
  border: 1px solid #e9ecef;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

/* Coming Soon Section */
.coming-soon-section {
  padding: 2rem;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 12px;
  position: relative;
}

/* Tool Preview Cards */
.tool-preview-card,
.agent-preview-card {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  position: relative;
  transition: all 0.3s ease;
}

.tool-preview-card:hover,
.agent-preview-card:hover {
  border-color: #007bff;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 123, 255, 0.15);
}

.tool-icon,
.agent-avatar {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.agent-capabilities {
  margin: 1rem 0;
}

.capability-tag {
  display: inline-block;
  background: #e9ecef;
  color: #495057;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  margin: 0.25rem;
}

/* Workflow Builder Preview */
.workflow-builder-preview {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  position: relative;
}

.workflow-canvas {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.workflow-node {
  background: #f8f9fa;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  padding: 1rem;
  min-width: 100px;
  text-align: center;
  transition: all 0.3s ease;
}

.workflow-node i {
  display: block;
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.start-node {
  border-color: #28a745;
  background-color: #d4edda;
}

.tool-node {
  border-color: #007bff;
  background-color: #d1ecf1;
}

.agent-node {
  border-color: #6f42c1;
  background-color: #e2d9f3;
}

.end-node {
  border-color: #dc3545;
  background-color: #f8d7da;
}

.workflow-arrow {
  font-size: 1.5rem;
  color: #6c757d;
  font-weight: bold;
}

.workflow-feature-card {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 1.5rem;
}

/* Status Badges */
.status-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-badge.available {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status-badge.coming-soon {
  background: #fff3cd;
  color: #856404;
  border: 1px solid #ffeaa7;
}

.status-badge.large {
  font-size: 1rem;
  padding: 0.5rem 1rem;
}

/* Create Modal */
.create-option-card {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  height: 100%;
}

.create-option-card:hover:not(.disabled) {
  border-color: #007bff;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 123, 255, 0.15);
}

.create-option-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.create-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

/* Modal styling */
.modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.modal.show {
  display: block !important;
}

.modal-backdrop {
  background-color: rgba(0, 0, 0, 0.5);
}

/* Card styling */
.card {
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
  border: 1px solid rgba(0, 0, 0, 0.125);
}

.card:hover {
  box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15);
  transition: box-shadow 0.15s ease-in-out;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .ai-platform-view {
    padding: 1rem 0;
  }
  
  .workflow-canvas {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .workflow-arrow {
    transform: rotate(90deg);
  }
  
  .d-flex.gap-2 {
    flex-direction: column;
    gap: 0.5rem !important;
  }
  
  .ai-platform-tabs .nav-link {
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }
}

/* Animation for cards */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tool-preview-card,
.agent-preview-card,
.create-option-card {
  animation: fadeInUp 0.6s ease-out;
}

.tool-preview-card:nth-child(2),
.agent-preview-card:nth-child(2) {
  animation-delay: 0.1s;
}

.tool-preview-card:nth-child(3),
.agent-preview-card:nth-child(3) {
  animation-delay: 0.2s;
}
</style>
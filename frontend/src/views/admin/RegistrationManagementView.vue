<template>
  <div class="container-fluid mt-4">
    <div class="row">
      <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
          <h2>Registration Management</h2>
          <div class="d-flex gap-2">
            <button class="btn btn-outline-primary" @click="refreshData">
              <i class="fas fa-refresh"></i> Refresh
            </button>
            <button class="btn btn-info" @click="showStats = !showStats">
              <i class="fas fa-chart-bar"></i> {{ showStats ? 'Hide' : 'Show' }} Stats
            </button>
          </div>
        </div>

        <!-- Statistics Card -->
        <div v-if="showStats" class="card mb-4">
          <div class="card-header">
            <h5 class="mb-0">Registration Statistics</h5>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-2 col-sm-4 col-6 mb-3">
                <div class="text-center">
                  <div class="h3 text-warning">{{ stats.pending || 0 }}</div>
                  <small class="text-muted">Pending</small>
                </div>
              </div>
              <div class="col-md-2 col-sm-4 col-6 mb-3">
                <div class="text-center">
                  <div class="h3 text-success">{{ stats.approved || 0 }}</div>
                  <small class="text-muted">Approved</small>
                </div>
              </div>
              <div class="col-md-2 col-sm-4 col-6 mb-3">
                <div class="text-center">
                  <div class="h3 text-danger">{{ stats.denied || 0 }}</div>
                  <small class="text-muted">Denied</small>
                </div>
              </div>
              <div class="col-md-2 col-sm-4 col-6 mb-3">
                <div class="text-center">
                  <div class="h3 text-secondary">{{ stats.suspended || 0 }}</div>
                  <small class="text-muted">Suspended</small>
                </div>
              </div>
              <div class="col-md-2 col-sm-4 col-6 mb-3">
                <div class="text-center">
                  <div class="h3 text-primary">{{ stats.total || 0 }}</div>
                  <small class="text-muted">Total</small>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Filters -->
        <div class="card mb-4">
          <div class="card-body">
            <div class="row align-items-end">
              <div class="col-md-3">
                <label class="form-label">Status Filter</label>
                <select v-model="filters.status" @change="loadRegistrations" class="form-select">
                  <option value="">All Statuses</option>
                  <option value="pending">Pending</option>
                  <option value="approved">Approved</option>
                  <option value="denied">Denied</option>
                  <option value="suspended">Suspended</option>
                </select>
              </div>
              <div class="col-md-3">
                <label class="form-label">Realm</label>
                <select v-model="filters.realmName" @change="loadRegistrations" class="form-select">
                  <option value="">All Realms</option>
                  <option value="default">Default</option>
                </select>
              </div>
              <div class="col-md-2">
                <label class="form-label">Page Size</label>
                <select v-model="filters.pageSize" @change="loadRegistrations" class="form-select">
                  <option value="10">10</option>
                  <option value="25">25</option>
                  <option value="50">50</option>
                  <option value="100">100</option>
                </select>
              </div>
              <div class="col-md-4">
                <button class="btn btn-outline-secondary" @click="clearFilters">
                  Clear Filters
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Registrations Table -->
        <div class="card">
          <div class="card-header">
            <h5 class="mb-0">User Registrations</h5>
          </div>
          <div class="card-body">
            <div v-if="loading" class="text-center py-4">
              <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>
            
            <div v-else-if="registrations.length === 0" class="text-center py-4 text-muted">
              No registrations found.
            </div>
            
            <div v-else>
              <div class="table-responsive">
                <table class="table table-striped table-hover">
                  <thead>
                    <tr>
                      <th>Username</th>
                      <th>Email</th>
                      <th>Status</th>
                      <th>Created</th>
                      <th>Updated</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="user in registrations" :key="user.id">
                      <td>
                        <strong>{{ user.username }}</strong>
                      </td>
                      <td>{{ user.email }}</td>
                      <td>
                        <span class="badge" :class="getStatusBadgeClass(user.status)">
                          {{ user.status.charAt(0).toUpperCase() + user.status.slice(1) }}
                        </span>
                      </td>
                      <td>
                        <small>{{ formatDate(user.created_at) }}</small>
                      </td>
                      <td>
                        <small>{{ formatDate(user.updated_at) }}</small>
                      </td>
                      <td>
                        <div class="btn-group btn-group-sm" role="group">
                          <button 
                            v-if="user.status === 'pending'" 
                            class="btn btn-success"
                            @click="approveUser(user)"
                            :disabled="processing"
                          >
                            <i class="fas fa-check"></i> Approve
                          </button>
                          <button 
                            v-if="user.status === 'pending'" 
                            class="btn btn-danger"
                            @click="denyUser(user)"
                            :disabled="processing"
                          >
                            <i class="fas fa-times"></i> Deny
                          </button>
                          <button 
                            v-if="user.status === 'approved'" 
                            class="btn btn-warning"
                            @click="suspendUser(user)"
                            :disabled="processing"
                          >
                            <i class="fas fa-ban"></i> Suspend
                          </button>
                          <button 
                            class="btn btn-info"
                            @click="viewUserDetails(user)"
                          >
                            <i class="fas fa-eye"></i> View
                          </button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Pagination -->
              <nav v-if="pagination.totalPages > 1" aria-label="Registration pagination">
                <ul class="pagination justify-content-center">
                  <li class="page-item" :class="{ disabled: pagination.page <= 1 }">
                    <button class="page-link" @click="changePage(pagination.page - 1)" :disabled="pagination.page <= 1">
                      Previous
                    </button>
                  </li>
                  <li 
                    v-for="page in getVisiblePages()" 
                    :key="page" 
                    class="page-item" 
                    :class="{ active: page === pagination.page }"
                  >
                    <button class="page-link" @click="changePage(page)">{{ page }}</button>
                  </li>
                  <li class="page-item" :class="{ disabled: pagination.page >= pagination.totalPages }">
                    <button class="page-link" @click="changePage(pagination.page + 1)" :disabled="pagination.page >= pagination.totalPages">
                      Next
                    </button>
                  </li>
                </ul>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Approval Modal -->
    <div class="modal fade" id="approvalModal" tabindex="-1" aria-labelledby="approvalModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="approvalModalLabel">
              {{ modalData.action === 'approve' ? 'Approve' : modalData.action === 'deny' ? 'Deny' : 'Suspend' }} Registration
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>
              <strong>User:</strong> {{ modalData.user?.username }} ({{ modalData.user?.email }})
            </p>
            <div class="mb-3">
              <label for="reason" class="form-label">
                {{ modalData.action === 'approve' ? 'Approval Note' : 'Reason' }}
                <span v-if="modalData.action !== 'approve'" class="text-danger">*</span>
              </label>
              <textarea 
                id="reason" 
                v-model="modalData.reason" 
                class="form-control" 
                rows="3"
                :placeholder="modalData.action === 'approve' ? 'Optional approval note...' : 'Please provide a reason...'"
                :required="modalData.action !== 'approve'"
              ></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
            <button 
              type="button" 
              class="btn"
              :class="{
                'btn-success': modalData.action === 'approve',
                'btn-danger': modalData.action === 'deny',
                'btn-warning': modalData.action === 'suspend'
              }"
              @click="confirmAction"
              :disabled="processing || (modalData.action !== 'approve' && !modalData.reason)"
            >
              {{ processing ? 'Processing...' : 'Confirm' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- User Details Modal -->
    <div class="modal fade" id="userDetailsModal" tabindex="-1" aria-labelledby="userDetailsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="userDetailsModalLabel">User Details</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div v-if="selectedUser" class="row">
              <div class="col-md-6">
                <h6>Basic Information</h6>
                <table class="table table-sm">
                  <tbody>
                    <tr>
                      <td><strong>ID:</strong></td>
                      <td><small class="font-monospace">{{ selectedUser.id }}</small></td>
                    </tr>
                    <tr>
                      <td><strong>Username:</strong></td>
                      <td>{{ selectedUser.username }}</td>
                    </tr>
                    <tr>
                      <td><strong>Email:</strong></td>
                      <td>{{ selectedUser.email }}</td>
                    </tr>
                    <tr>
                      <td><strong>Realm ID:</strong></td>
                      <td><small class="font-monospace">{{ selectedUser.realm_id }}</small></td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="col-md-6">
                <h6>Status Information</h6>
                <table class="table table-sm">
                  <tbody>
                    <tr>
                      <td><strong>Status:</strong></td>
                      <td>
                        <span class="badge" :class="getStatusBadgeClass(selectedUser.status)">
                          {{ selectedUser.status.charAt(0).toUpperCase() + selectedUser.status.slice(1) }}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td><strong>Active:</strong></td>
                      <td>
                        <span class="badge" :class="selectedUser.is_active ? 'bg-success' : 'bg-secondary'">
                          {{ selectedUser.is_active ? 'Yes' : 'No' }}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td><strong>Created:</strong></td>
                      <td>{{ formatDate(selectedUser.created_at) }}</td>
                    </tr>
                    <tr>
                      <td><strong>Updated:</strong></td>
                      <td>{{ formatDate(selectedUser.updated_at) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'

interface User {
  id: string
  username: string
  email: string
  realm_id: string
  is_active: boolean
  status: 'pending' | 'approved' | 'denied' | 'suspended'
  created_by: string
  created_at: string
  updated_by: string
  updated_at: string
}

interface RegistrationResponse {
  users: User[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

interface Stats {
  pending: number
  approved: number
  denied: number
  suspended: number
  total: number
}

const authStore = useAuthStore()

const loading = ref(false)
const processing = ref(false)
const showStats = ref(true)
const registrations = ref<User[]>([])
const stats = ref<Stats>({
  pending: 0,
  approved: 0,
  denied: 0,
  suspended: 0,
  total: 0
})

const filters = ref({
  status: 'pending',
  realmName: 'default',
  pageSize: 10
})

const pagination = ref({
  page: 1,
  totalPages: 1,
  total: 0
})

const modalData = ref({
  action: '',
  user: null as User | null,
  reason: ''
})

const selectedUser = ref<User | null>(null)

onMounted(() => {
  loadRegistrations()
  loadStats()
})

const loadRegistrations = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      page_size: filters.value.pageSize.toString()
    })
    
    if (filters.value.status) params.append('status', filters.value.status)
    if (filters.value.realmName) params.append('realm_name', filters.value.realmName)

    const response = await fetch(`/api/v1/admin/registrations?${params}`, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to load registrations')
    }

    const data: RegistrationResponse = await response.json()
    registrations.value = data.users
    pagination.value = {
      page: data.page,
      totalPages: data.total_pages,
      total: data.total
    }
  } catch (error) {
    console.error('Error loading registrations:', error)
    alert('Failed to load registrations')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const params = new URLSearchParams()
    if (filters.value.realmName) params.append('realm_name', filters.value.realmName)

    const response = await fetch(`/api/v1/admin/registrations/stats?${params}`, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to load stats')
    }

    const data = await response.json()
    stats.value = data.stats
  } catch (error) {
    console.error('Error loading stats:', error)
  }
}

const refreshData = () => {
  loadRegistrations()
  loadStats()
}

const clearFilters = () => {
  filters.value = {
    status: '',
    realmName: '',
    pageSize: 10
  }
  pagination.value.page = 1
  loadRegistrations()
}

const changePage = (page: number) => {
  pagination.value.page = page
  loadRegistrations()
}

const getVisiblePages = () => {
  const current = pagination.value.page
  const total = pagination.value.totalPages
  const delta = 2
  const range = []
  const start = Math.max(1, current - delta)
  const end = Math.min(total, current + delta)

  for (let i = start; i <= end; i++) {
    range.push(i)
  }

  return range
}

const getStatusBadgeClass = (status: string) => {
  const classes = {
    pending: 'bg-warning text-dark',
    approved: 'bg-success',
    denied: 'bg-danger',
    suspended: 'bg-secondary'
  }
  return classes[status as keyof typeof classes] || 'bg-secondary'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const approveUser = (user: User) => {
  modalData.value = {
    action: 'approve',
    user: user,
    reason: ''
  }
  const modal = new (window as any).bootstrap.Modal(document.getElementById('approvalModal'))
  modal.show()
}

const denyUser = (user: User) => {
  modalData.value = {
    action: 'deny',
    user: user,
    reason: ''
  }
  const modal = new (window as any).bootstrap.Modal(document.getElementById('approvalModal'))
  modal.show()
}

const suspendUser = (user: User) => {
  modalData.value = {
    action: 'suspend',
    user: user,
    reason: ''
  }
  const modal = new (window as any).bootstrap.Modal(document.getElementById('approvalModal'))
  modal.show()
}

const viewUserDetails = (user: User) => {
  selectedUser.value = user
  const modal = new (window as any).bootstrap.Modal(document.getElementById('userDetailsModal'))
  modal.show()
}

const confirmAction = async () => {
  if (!modalData.value.user) return

  processing.value = true
  try {
    const approved = modalData.value.action === 'approve'
    const response = await fetch('/api/v1/admin/registrations/approve', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: JSON.stringify({
        user_id: modalData.value.user.id,
        approved: approved,
        reason: modalData.value.reason
      })
    })

    if (!response.ok) {
      throw new Error('Failed to process registration')
    }

    const modal = (window as any).bootstrap.Modal.getInstance(document.getElementById('approvalModal'))
    modal.hide()
    
    // Refresh data
    await refreshData()
    
    alert(`Registration ${approved ? 'approved' : 'denied'} successfully`)
  } catch (error) {
    console.error('Error processing registration:', error)
    alert('Failed to process registration')
  } finally {
    processing.value = false
  }
}
</script>

<style scoped>
.font-monospace {
  font-family: 'Courier New', Courier, monospace;
  font-size: 0.85em;
}

.table th {
  border-top: none;
  font-weight: 600;
}

.btn-group-sm .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.spinner-border {
  width: 3rem;
  height: 3rem;
}

.badge {
  font-size: 0.75em;
}

@media (max-width: 768px) {
  .table-responsive {
    font-size: 0.875rem;
  }
  
  .btn-group-sm .btn {
    padding: 0.2rem 0.4rem;
    font-size: 0.75rem;
  }
}
</style>


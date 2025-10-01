import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { makeAuthenticatedRequest } from '@/utils/httpInterceptor';
import { getApiUrl } from '@/utils/apiConfig';

// Types
export interface Policy {
  id: string;
  realm_id: string;
  name: string;
  description: string;
  version: string;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface Statement {
  id: string;
  policy_id: string;
  sid: string; // Statement ID (optional)
  effect: 'Allow' | 'Deny';
  actions: string; // JSON array as string
  resources: string; // JSON array as string
  conditions: string; // JSON string (optional)
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface PoliciesResponse {
  policies: Policy[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface PolicyWithStatements {
  policy: Policy;
  statements: Statement[];
}

export interface CreatePolicyRequest {
  realm_name: string;
  name: string;
  description?: string;
  version?: string;
}

export interface UpdatePolicyRequest {
  name?: string;
  description?: string;
  version?: string;
}

export interface CreateStatementRequest {
  policy_id: string;
  sid?: string;
  effect: 'Allow' | 'Deny';
  actions: string[];
  resources: string[];
  conditions?: Record<string, any>;
}

export interface UpdateStatementRequest {
  sid?: string;
  effect?: 'Allow' | 'Deny';
  actions?: string[];
  resources?: string[];
  conditions?: Record<string, any>;
}

export interface PolicyFilters {
  realm_name?: string;
  page?: number;
  page_size?: number;
}

export const usePolicyStore = defineStore('policy', () => {
  // State
  const policies = ref<Policy[]>([]);
  const currentPolicy = ref<Policy | null>(null);
  const statements = ref<Statement[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const pagination = ref({
    total: 0,
    page: 1,
    page_size: 10,
    total_pages: 0,
  });

  // Computed
  const hasPolicies = computed(() => policies.value.length > 0);
  const totalPolicies = computed(() => pagination.value.total);
  const currentPage = computed(() => pagination.value.page);
  const totalPages = computed(() => pagination.value.total_pages);

  // Actions
  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading;
  };

  const setError = (errorMessage: string | null) => {
    error.value = errorMessage;
  };

  const clearError = () => {
    error.value = null;
  };

  /**
   * Get policies with pagination and filtering
   */
  const getPolicies = async (filters: PolicyFilters = {}) => {
    setLoading(true);
    clearError();

    try {
      const params = new URLSearchParams();
      
      if (filters.realm_name) params.append('realm_name', filters.realm_name);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.page_size) params.append('page_size', filters.page_size.toString());

      const url = `${getApiUrl('/api/v1/admin/policies')}${params.toString() ? `?${params.toString()}` : ''}`;
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch policies: ${response.statusText}`);
      }

      const data: PoliciesResponse = await response.json();
      
      policies.value = data.policies;
      pagination.value = {
        total: data.total,
        page: data.page,
        page_size: data.page_size,
        total_pages: data.total_pages,
      };

      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get policies';
      setError(errorMessage);
      console.error('Failed to get policies:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Get a single policy by ID with its statements
   */
  const getPolicy = async (policyId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/policies/${policyId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch policy: ${response.statusText}`);
      }

      const data: PolicyWithStatements = await response.json();
      currentPolicy.value = data.policy;
      statements.value = data.statements;
      
      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get policy';
      setError(errorMessage);
      console.error('Failed to get policy:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Create a new policy
   */
  const createPolicy = async (policyData: CreatePolicyRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl('/api/v1/admin/policies');
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'POST',
        body: JSON.stringify(policyData),
      });

      if (!response.ok) {
        throw new Error(`Failed to create policy: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Add the new policy to the list
      policies.value.unshift(data.policy);
      pagination.value.total += 1;
      
      return data.policy;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create policy';
      setError(errorMessage);
      console.error('Failed to create policy:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Update an existing policy
   */
  const updatePolicy = async (policyId: string, policyData: UpdatePolicyRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/policies/${policyId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'PUT',
        body: JSON.stringify(policyData),
      });

      if (!response.ok) {
        throw new Error(`Failed to update policy: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Update the policy in the list
      const index = policies.value.findIndex(policy => policy.id === policyId);
      if (index !== -1) {
        policies.value[index] = data.policy;
      }
      
      // Update current policy if it's the same policy
      if (currentPolicy.value?.id === policyId) {
        currentPolicy.value = data.policy;
      }
      
      return data.policy;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update policy';
      setError(errorMessage);
      console.error('Failed to update policy:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Delete a policy
   */
  const deletePolicy = async (policyId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/policies/${policyId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error(`Failed to delete policy: ${response.statusText}`);
      }

      // Remove the policy from the list
      policies.value = policies.value.filter(policy => policy.id !== policyId);
      pagination.value.total -= 1;
      
      // Clear current policy if it's the deleted policy
      if (currentPolicy.value?.id === policyId) {
        currentPolicy.value = null;
        statements.value = [];
      }
      
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete policy';
      setError(errorMessage);
      console.error('Failed to delete policy:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Create a new statement for a policy
   */
  const createStatement = async (statementData: CreateStatementRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/policies/${statementData.policy_id}/statements`);
      
      // Convert arrays to JSON strings as expected by backend
      const requestData = {
        ...statementData,
        actions: JSON.stringify(statementData.actions),
        resources: JSON.stringify(statementData.resources),
        conditions: statementData.conditions ? JSON.stringify(statementData.conditions) : undefined,
      };
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'POST',
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error(`Failed to create statement: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Add the new statement to the list
      statements.value.push(data.statement);
      
      return data.statement;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create statement';
      setError(errorMessage);
      console.error('Failed to create statement:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Update an existing statement
   */
  const updateStatement = async (statementId: string, statementData: UpdateStatementRequest) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/statements/${statementId}`);
      
      // Convert arrays to JSON strings as expected by backend
      const requestData = {
        ...statementData,
        actions: statementData.actions ? JSON.stringify(statementData.actions) : undefined,
        resources: statementData.resources ? JSON.stringify(statementData.resources) : undefined,
        conditions: statementData.conditions ? JSON.stringify(statementData.conditions) : undefined,
      };
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'PUT',
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error(`Failed to update statement: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Update the statement in the list
      const index = statements.value.findIndex(statement => statement.id === statementId);
      if (index !== -1) {
        statements.value[index] = data.statement;
      }
      
      return data.statement;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update statement';
      setError(errorMessage);
      console.error('Failed to update statement:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Delete a statement
   */
  const deleteStatement = async (statementId: string) => {
    setLoading(true);
    clearError();

    try {
      const url = getApiUrl(`/api/v1/admin/statements/${statementId}`);
      
      const response = await makeAuthenticatedRequest(url, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error(`Failed to delete statement: ${response.statusText}`);
      }

      // Remove the statement from the list
      statements.value = statements.value.filter(statement => statement.id !== statementId);
      
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete statement';
      setError(errorMessage);
      console.error('Failed to delete statement:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * Search policies by name or description
   */
  const searchPolicies = (searchTerm: string) => {
    if (!searchTerm.trim()) {
      return policies.value;
    }
    
    const term = searchTerm.toLowerCase();
    return policies.value.filter(policy =>
      policy.name.toLowerCase().includes(term) ||
      policy.description.toLowerCase().includes(term)
    );
  };

  /**
   * Get policy by ID from current list
   */
  const getPolicyById = (policyId: string) => {
    return policies.value.find(policy => policy.id === policyId);
  };

  /**
   * Parse JSON string fields from statements
   */
  const parseStatementActions = (actionsJson: string): string[] => {
    try {
      return JSON.parse(actionsJson);
    } catch {
      return [];
    }
  };

  const parseStatementResources = (resourcesJson: string): string[] => {
    try {
      return JSON.parse(resourcesJson);
    } catch {
      return [];
    }
  };

  const parseStatementConditions = (conditionsJson: string): Record<string, any> => {
    try {
      return JSON.parse(conditionsJson);
    } catch {
      return {};
    }
  };

  /**
   * Clear all policy data
   */
  const clearPolicies = () => {
    policies.value = [];
    currentPolicy.value = null;
    statements.value = [];
    pagination.value = {
      total: 0,
      page: 1,
      page_size: 10,
      total_pages: 0,
    };
    clearError();
  };

  /**
   * Refresh current page of policies
   */
  const refreshPolicies = async () => {
    return getPolicies({
      page: pagination.value.page,
      page_size: pagination.value.page_size,
    });
  };

  return {
    // State
    policies,
    currentPolicy,
    statements,
    loading,
    error,
    pagination,
    
    // Computed
    hasPolicies,
    totalPolicies,
    currentPage,
    totalPages,
    
    // Actions
    getPolicies,
    getPolicy,
    createPolicy,
    updatePolicy,
    deletePolicy,
    createStatement,
    updateStatement,
    deleteStatement,
    searchPolicies,
    getPolicyById,
    parseStatementActions,
    parseStatementResources,
    parseStatementConditions,
    clearPolicies,
    refreshPolicies,
    setLoading,
    setError,
    clearError,
  };
});

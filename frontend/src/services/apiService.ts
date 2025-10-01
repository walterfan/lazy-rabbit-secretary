// This file is deprecated. Use the userStore instead.
// The userStore provides better state management and automatic authentication handling.

export { useUserStore } from '@/stores/userStore';
export type { 
  User, 
  UsersResponse, 
  CreateUserRequest, 
  UpdateUserRequest, 
  UserFilters,
  RegistrationStats 
} from '@/stores/userStore';

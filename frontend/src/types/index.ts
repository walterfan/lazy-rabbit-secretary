export interface Book {
  id: string;
  realm_id: string;
  isbn: string;
  name: string;
  author: string;
  description: string;
  price: number;
  borrow_time?: Date;
  return_time?: Date;
  deadline?: Date;
  tags: string[];
  created_by: string;
  created_at: Date;
  updated_by: string;
  updated_time: Date;
}

export interface CreateBookRequest {
  isbn: string;
  name: string;
  author?: string;
  description?: string;
  price?: number;
  tags?: string;
  deadline?: Date;
}

export interface UpdateBookRequest {
  isbn?: string;
  name?: string;
  author?: string;
  description?: string;
  price?: number;
  tags?: string;
  deadline?: Date;
}

export interface BookListResponse {
  books: Book[];
  total: number;
  page: number;
  limit: number;
}

export interface Prompt {
  id: string;
  name: string;
  description: string;
  system_prompt: string;
  user_prompt: string;
  tags: string;
  created_by: string;
  created_at: Date;
  updated_by: string;
  updated_at: Date;
}

export interface CreatePromptRequest {
  name: string;
  description?: string;
  system_prompt: string;
  user_prompt: string;
  tags?: string;
}

export interface UpdatePromptRequest {
  name?: string;
  description?: string;
  system_prompt?: string;
  user_prompt?: string;
  tags?: string;
}

export interface PromptListResponse {
  prompts: Prompt[];
  total: number;
  page: number;
  limit: number;
}

export interface Task {
  id: string;
  name: string;
  description: string;
  priority: 'low' | 'medium' | 'high';
  difficulty: 'easy' | 'medium' | 'hard';
  status: 'pending' | 'running' | 'completed' | 'failed';
  minutes: number;
  deadline: Date;
  schedule_time: Date;
  start_time?: Date;
  end_time?: Date;
  tags: string[];
}

export interface Secret {
  id: string;
  realm_id: string;
  name: string;
  group: string;
  desc: string;
  path: string;
  cipher_alg: string;
  cipher_text: string;
  nonce: string;
  auth_tag: string;
  wrapped_dek: string;
  kek_version: number;
  created_by: string;
  created_at: Date;
  updated_by: string;
  updated_at: Date;
}

export interface CreateSecretRequest {
  name: string;
  group: string;
  desc: string;
  path: string;
  value: string;
  kek: string; // Custom KEK (empty string means use default)
}

export interface UpdateSecretRequest {
  name: string;
  group: string;
  desc: string;
  path: string;
  value: string;
  current_value: string; // Current secret value for verification
  kek: string; // Custom KEK (empty string means use default)
}

export interface Reminder {
  id: string;
  realm_id: string;
  name: string;
  content: string;
  remind_time: Date;
  status: 'pending' | 'active' | 'completed' | 'cancelled';
  tags: string;
  remind_methods: string; // Comma-separated: email,im,webhook
  remind_targets: string; // JSON or comma-separated targets
  created_by: string;
  created_at: Date;
  updated_by: string;
  updated_at: Date;
}

export interface CreateReminderRequest {
  name: string;
  content: string;
  remind_time: Date;
  tags: string;
  remind_methods: string; // Comma-separated: email,im,webhook
  remind_targets: string; // JSON or comma-separated targets
}

export interface UpdateReminderRequest {
  name?: string;
  content?: string;
  status?: 'pending' | 'active' | 'completed' | 'cancelled';
  remind_time?: Date;
  tags?: string;
  remind_methods?: string; // Comma-separated: email,im,webhook
  remind_targets?: string; // JSON or comma-separated targets
}

export type ValidationErrors = {
  [key: string]: string[];
}

// Post-related types
export interface Post {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  status: 'draft' | 'pending' | 'published' | 'private' | 'trash' | 'scheduled';
  type: 'post' | 'page' | 'attachment' | 'revision' | 'custom';
  format: 'standard' | 'aside' | 'gallery' | 'link' | 'image' | 'quote' | 'status' | 'video' | 'audio' | 'chat';
  password?: string;
  meta_title: string;
  meta_description: string;
  meta_keywords: string;
  featured_image: string;
  categories: string[];
  tags: string[];
  published_at?: string;
  scheduled_for?: string;
  view_count: number;
  comment_count: number;
  parent_id?: string;
  menu_order: number;
  is_sticky: boolean;
  allow_pings: boolean;
  comment_status: 'open' | 'closed' | 'registration_required';
  language: string;
  custom_fields?: Record<string, any>;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface CreatePostRequest {
  title: string;
  slug?: string;
  content: string;
  excerpt?: string;
  status?: 'draft' | 'pending' | 'published' | 'private' | 'scheduled';
  type?: 'post' | 'page';
  format?: string;
  password?: string;
  meta_title?: string;
  meta_description?: string;
  meta_keywords?: string;
  featured_image?: string;
  categories?: string[];
  tags?: string[];
  parent_id?: string;
  menu_order?: number;
  is_sticky?: boolean;
  allow_pings?: boolean;
  comment_status?: string;
  scheduled_for?: string;
  custom_fields?: Record<string, any>;
}

export interface UpdatePostRequest extends Partial<CreatePostRequest> {}

export interface PostListResponse {
  posts: Post[];
  total: number;
  page: number;
  limit: number;
}
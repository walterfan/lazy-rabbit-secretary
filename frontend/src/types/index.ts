export interface Book {
  isbn: string;
  title: string;
  author: string;
  price: number;
  borrower: string;
  borrowTime?: Date;
  returnTime?: Date;
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
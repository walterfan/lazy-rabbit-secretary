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
}

export type ValidationErrors = {
  [key: string]: string[];
}
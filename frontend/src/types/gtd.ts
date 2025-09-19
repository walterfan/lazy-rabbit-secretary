export interface InboxItem {
  id: string;
  title: string;
  description: string;
  priority: 'low' | 'normal' | 'high' | 'urgent';
  status: 'pending' | 'processing' | 'completed' | 'archived';
  tags: string;
  context: string;
  created_by: string;
  created_at: Date;
  updated_at: Date;
}

export interface DailyChecklistItem {
  id: string;
  title: string;
  description: string;
  priority: 'A' | 'B+' | 'B' | 'C' | 'D';
  estimated_time: number; // in minutes
  deadline?: Date;
  context: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  completion_time?: Date;
  actual_time: number; // in minutes
  notes: string;
  inbox_item_id?: string;
  date: Date;
  created_by: string;
  created_at: Date;
  updated_at: Date;
}

export interface CreateInboxItemRequest {
  title: string;
  description?: string;
  priority?: 'low' | 'normal' | 'high' | 'urgent';
  tags?: string;
  context?: string;
}

export interface UpdateInboxItemRequest {
  title?: string;
  description?: string;
  priority?: 'low' | 'normal' | 'high' | 'urgent';
  status?: 'pending' | 'processing' | 'completed' | 'archived';
  tags?: string;
  context?: string;
}

export interface CreateDailyItemRequest {
  title: string;
  description?: string;
  priority?: 'A' | 'B+' | 'B' | 'C' | 'D';
  estimated_time?: number;
  deadline?: Date;
  context?: string;
  notes?: string;
  inbox_item_id?: string;
  date: Date;
}

export interface UpdateDailyItemRequest {
  title?: string;
  description?: string;
  priority?: 'A' | 'B+' | 'B' | 'C' | 'D';
  estimated_time?: number;
  deadline?: Date;
  context?: string;
  status?: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  actual_time?: number;
  notes?: string;
}

export interface InboxListResponse {
  items: InboxItem[];
  total: number;
  page: number;
  limit: number;
}

export interface DailyListResponse {
  items: DailyChecklistItem[];
  total: number;
  page: number;
  limit: number;
}

export interface DailyStatsResponse {
  date: Date;
  total_items: number;
  completed_items: number;
  pending_items: number;
  in_progress_items: number;
  cancelled_items: number;
  completion_rate: number;
  total_estimated_time: number;
  total_actual_time: number;
  priority_breakdown: Record<string, number>;
}

export interface InboxStatsResponse {
  total: number;
  status_pending: number;
  status_processing: number;
  status_completed: number;
  status_archived: number;
  priority_urgent: number;
  priority_high: number;
  priority_normal: number;
  priority_low: number;
}

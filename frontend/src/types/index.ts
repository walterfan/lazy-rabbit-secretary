export interface Book {
  isbn: string;
  title: string;
  author: string;
  price: number;
  borrowTime?: Date;
  returnTime?: Date;
}

export interface Task {
  id: string;
  name: string;
  description: string;
  priority: 'low' | 'medium' | 'high';
  minutes: number;
  deadline: Date;
  start_time?: Date;
  end_time?: Date;
  tags: string[];
}

export type ValidationErrors = {
  [key: string]: string[];
}
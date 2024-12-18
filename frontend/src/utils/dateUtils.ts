import { format } from 'date-fns';

export const formatDate = (date?: Date): string => {
  if (!date) return '-';
  return format(new Date(date), 'yyyy-MM-dd HH:mm');
};

export const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(price);
};
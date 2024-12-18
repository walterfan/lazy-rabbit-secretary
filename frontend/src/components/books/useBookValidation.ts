import { required, minLength } from '@vuelidate/validators';

export const useBookValidation = () => {
  const rules = {
    isbn: { required, minLength: minLength(10) },
    title: { required },
    author: { required },
    price: { required }
  };

  return {
    rules
  };
};
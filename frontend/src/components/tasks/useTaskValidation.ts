import { required, minValue } from '@vuelidate/validators';

export const useTaskValidation = () => {
  const rules = {
    name: { required },
    description: { required },
    priority: { required },
    minutes: { required, minValue: minValue(1) },
    deadline: { required }
  };

  return {
    rules
  };
};
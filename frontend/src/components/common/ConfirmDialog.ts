export const useConfirmDialog = () => {
  const confirm = (message: string): boolean => {
    return window.confirm(message);
  };

  return {
    confirm
  };
};
import { ref } from 'vue';

export const useLoading = () => {
  const loading = ref(false);

  const showLoading = () => {
    loading.value = true;
  };

  const hideLoading = () => {
    loading.value = false;
  };

  return {
    loading,
    showLoading,
    hideLoading,
  };
};

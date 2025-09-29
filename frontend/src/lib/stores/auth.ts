import { writable } from 'svelte/store';
import type { User } from '$lib/types';
import { api } from '$lib/api';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  });

  return {
    subscribe,
    login: async (username: string, password: string) => {
      update(state => ({ ...state, isLoading: true }));
      
      const response = await api.login({ username, password });
      
      if (response.data) {
        set({
          user: response.data.user,
          isAuthenticated: true,
          isLoading: false,
        });
        return { success: true };
      } else {
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
        return { success: false, error: response.error };
      }
    },
    
    logout: async () => {
      await api.logout();
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });
    },
    
    checkAuth: async () => {
      update(state => ({ ...state, isLoading: true }));
      
      const response = await api.getCurrentUser();
      
      if (response.data) {
        set({
          user: response.data,
          isAuthenticated: true,
          isLoading: false,
        });
      } else {
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      }
    },
    
    setUser: (user: User) => {
      set({
        user,
        isAuthenticated: true,
        isLoading: false,
      });
    },
  };
}

export const auth = createAuthStore(); 
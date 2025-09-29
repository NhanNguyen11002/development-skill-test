import { writable } from 'svelte/store';
import type { Alert } from '$lib/types';
import { api } from '$lib/api';

interface AlertsState {
  alerts: Alert[];
  isLoading: boolean;
  error: string | null;
}

function createAlertsStore() {
  const { subscribe, set, update } = writable<AlertsState>({
    alerts: [],
    isLoading: false,
    error: null,
  });

  return {
    subscribe,
    
    fetchAlerts: async () => {
      update(state => ({ ...state, isLoading: true, error: null }));
      
      const response = await api.getAlerts();
      
      if (response.data) {
        set({
          alerts: response.data,
          isLoading: false,
          error: null,
        });
      } else {
        set({
          alerts: [],
          isLoading: false,
          error: response.error || 'Failed to fetch alerts',
        });
      }
    },
    
    addAlert: (alert: Alert) => {
      update(state => ({
        ...state,
        alerts: [alert, ...state.alerts],
      }));
    },
    
    updateAlert: (updatedAlert: Alert) => {
      update(state => ({
        ...state,
        alerts: state.alerts.map(alert => 
          alert.id === updatedAlert.id ? updatedAlert : alert
        ),
      }));
    },
    
    removeAlert: (alertId: string) => {
      update(state => ({
        ...state,
        alerts: state.alerts.filter(alert => alert.id !== alertId),
      }));
    },
    
    acknowledgeAlert: async (alertId: string) => {
      const response = await api.acknowledgeAlert(alertId);
      
      if (response.data) {
        update(state => ({
          ...state,
          alerts: state.alerts.map(alert => 
            alert.id === alertId ? response.data! : alert
          ),
        }));
        return { success: true };
      } else {
        return { success: false, error: response.error };
      }
    },
    
    assignAlert: async (alertId: string, guardId: string) => {
      const response = await api.assignAlert(alertId, guardId);
      
      if (response.data) {
        update(state => ({
          ...state,
          alerts: state.alerts.map(alert => 
            alert.id === alertId ? response.data!.alert : alert
          ),
        }));
        return { success: true, incident: response.data.incident };
      } else {
        return { success: false, error: response.error };
      }
    },
    
    clearError: () => {
      update(state => ({ ...state, error: null }));
    },
  };
}

export const alerts = createAlertsStore(); 
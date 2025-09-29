import type { 
  User, 
  Premise, 
  Camera, 
  Alert, 
  Incident, 
  IncidentUpdate,
  LoginRequest,
  LoginResponse,
  CameraGrid,
  StreamResponse,
  ApiResponse
} from './types';

const API_BASE_URL = 'http://localhost:8081/api';

class ApiClient {
  private token: string | null = null;

  setToken(token: string) {
    this.token = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  getToken(): string | null {
    if (!this.token && typeof window !== 'undefined') {
      this.token = localStorage.getItem('token');
    }
    return this.token;
  }

  clearToken() {
    this.token = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${API_BASE_URL}${endpoint}`;
    const token = this.getToken();

    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        return { error: data.error || 'An error occurred' };
      }

      return { data };
    } catch (error) {
      return { error: error instanceof Error ? error.message : 'Network error' };
    }
  }

  // Authentication
  async login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    const response = await this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });

    if (response.data) {
      this.setToken(response.data.token);
    }

    return response;
  }

  async logout(): Promise<ApiResponse<void>> {
    const response = await this.request<void>('/auth/logout', {
      method: 'POST',
    });
    this.clearToken();
    return response;
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.request<User>('/auth/me');
  }

  // Premises
  async getPremises(): Promise<ApiResponse<Premise[]>> {
    return this.request<Premise[]>('/premises');
  }

  async getPremise(id: string): Promise<ApiResponse<Premise>> {
    return this.request<Premise>(`/premises/${id}`);
  }

  async getPremiseCameras(id: string): Promise<ApiResponse<Camera[]>> {
    return this.request<Camera[]>(`/premises/${id}/cameras`);
  }

  // Cameras
  async getCameras(): Promise<ApiResponse<Camera[]>> {
    return this.request<Camera[]>('/cameras');
  }

  async getCamera(id: string): Promise<ApiResponse<Camera>> {
    return this.request<Camera>(`/cameras/${id}`);
  }

  async getCameraGrid(): Promise<ApiResponse<CameraGrid>> {
    return this.request<CameraGrid>('/cameras/grid');
  }

  async getCameraStream(id: string): Promise<ApiResponse<StreamResponse>> {
    return this.request<StreamResponse>(`/cameras/${id}/stream`);
  }

  async updateCameraStatus(id: string, status: Camera['status']): Promise<ApiResponse<Camera>> {
    return this.request<Camera>(`/cameras/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  }

  // Alerts
  async getAlerts(): Promise<ApiResponse<Alert[]>> {
    return this.request<Alert[]>('/alerts');
  }

  async getAlert(id: string): Promise<ApiResponse<Alert>> {
    return this.request<Alert>(`/alerts/${id}`);
  }

  async acknowledgeAlert(id: string): Promise<ApiResponse<Alert>> {
    return this.request<Alert>(`/alerts/${id}/acknowledge`, {
      method: 'POST',
    });
  }

  async assignAlert(id: string, guardId: string): Promise<ApiResponse<{ alert: Alert; incident: Incident }>> {
    return this.request<{ alert: Alert; incident: Incident }>(`/alerts/${id}/assign`, {
      method: 'POST',
      body: JSON.stringify({ guard_id: guardId }),
    });
  }

  async createAlert(alert: Partial<Alert>): Promise<ApiResponse<Alert>> {
    return this.request<Alert>('/alerts', {
      method: 'POST',
      body: JSON.stringify(alert),
    });
  }

  // Incidents
  async getIncidents(): Promise<ApiResponse<Incident[]>> {
    return this.request<Incident[]>('/incidents');
  }

  async getIncident(id: string): Promise<ApiResponse<Incident>> {
    return this.request<Incident>(`/incidents/${id}`);
  }

  async updateIncident(id: string, status: Incident['status']): Promise<ApiResponse<Incident>> {
    return this.request<Incident>(`/incidents/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  }

  async addIncidentUpdate(id: string, update: Partial<IncidentUpdate>): Promise<ApiResponse<IncidentUpdate>> {
    return this.request<IncidentUpdate>(`/incidents/${id}/updates`, {
      method: 'POST',
      body: JSON.stringify(update),
    });
  }

  // Guards
  async getGuards(): Promise<ApiResponse<User[]>> {
    return this.request<User[]>('/guards');
  }

  // Health check
  async healthCheck(): Promise<ApiResponse<{ status: string; service: string }>> {
    return this.request<{ status: string; service: string }>('/health');
  }
}

export const api = new ApiClient(); 
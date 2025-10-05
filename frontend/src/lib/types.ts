export interface User {
  id: string;
  username: string;
  email: string;
  role: 'scs_operator' | 'security_guard';
  firstName: string;
  lastName: string;
  phone: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Premise {
  id: string;
  name: string;
  address: string;
  type: 'office' | 'substation';
  floorPlans: string;
  description: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  cameras?: Camera[];
}

export interface Camera {
  id: string;
  name: string;
  location: string;
  streamUrl: string;
  status: 'active' | 'inactive' | 'maintenance';
  premiseId: string;
  createdAt: string;
  updatedAt: string;
  premise?: Premise;
  assignedGuards?: User[];
}

export interface Alert {
  id: string;
  type: 'unauthorized_access' | 'suspicious_activity' | 'equipment_damage' | 'system_failure';
  severity: 'low' | 'medium' | 'high' | 'critical';
  title: string;
  description: string;
  location: string;
  status: 'pending' | 'acknowledged' | 'assigned' | 'resolved' | 'closed';
  cameraId?: string;
  premiseId: string;
  assignedGuardId?: string;
  createdAt: string;
  updatedAt: string;
  camera?: Camera;
  premise?: Premise;
  assignedGuard?: User;
  incident?: Incident;
}

export interface Incident {
  id: string;
  alertId: string;
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  assignedGuardId: string;
  location: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  alert?: Alert;
  assignedGuard?: User;
  updates?: IncidentUpdate[];
}

export interface IncidentUpdate {
  id: string;
  incidentId: string;
  guardId: string;
  type: 'arrival' | 'investigation' | 'resolution' | 'photo' | 'video';
  message: string;
  mediaUrl?: string;
  location?: string;
  createdAt: string;
  incident?: Incident;
  guard?: User;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface WebSocketMessage {
  type: string;
  payload: any;
  userId?: string;
}

export interface CameraGrid {
  grid: Camera[][];
  total: number;
  userRole: string;
}

export interface StreamResponse {
  cameraId: string;
  streamUrl: string;
  status: string;
  webrtcUrl: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  status: string;
  message?: string;
} 
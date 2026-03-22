export interface Device {
  id: number;
  hostname: string;
  location: string;
  ip: string;
  is_active: boolean;
  is_deleted: boolean;
  created_at: Date;
}

export interface CreateDevicePayload {
  hostname: string;
  ip: string;
  location: string;
  is_active: boolean;
}

export interface UpdateDevicePayload {
  hostname: string;
  ip: string;
  location: string;
  is_active: boolean;
}

export interface Device {
  id: number;
  hostname: string;
  location: string;
  ip: string;
  is_active: boolean;
  is_deleted: boolean;
  created_at: Date;
}

import type { Device } from "../types/device";

export async function getDevices(search?: string): Promise<Device[]> {
  const params = new URLSearchParams();
  if (search) params.set("search", search);

  const query = params.size ? `?${params}` : "";
  const response = await fetch(`/api/v1/devices${query}`);

  if (!response.ok) {
    throw new Error(`Не удалось загрузить устройства: ${response.statusText}`);
  }

  return response.json();
}

export interface CreateDevicePayload {
  hostname: string;
  ip: string;
  location: string;
  is_active: boolean;
}

export async function createDevice(
  payload: CreateDevicePayload,
): Promise<Device> {
  const response = await fetch("/api/v1/devices", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message);
  }

  return response.json();
}

export async function deleteDevice(id: number): Promise<void> {
  const response = await fetch(`/api/v1/devices/${id}`, { method: "DELETE" });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message);
  }
}

export async function updateDevice(
  id: number,
  payload: CreateDevicePayload,
): Promise<Device> {
  const response = await fetch(`/api/v1/devices/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message);
  }
  return response.json();
}

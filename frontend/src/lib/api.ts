import type {
  CreateDevicePayload,
  Device,
  UpdateDevicePayload,
} from "../types/device";

const API_BASE = import.meta.env.VITE_API_BASE;

export async function getDevices(
  search?: string,
  onlyActive?: boolean,
): Promise<Device[]> {
  const params = new URLSearchParams();
  if (search) params.set("search", search);
  if (onlyActive) params.set("is_active", "true");

  const query = params.size ? `?${params}` : "";
  const response = await fetch(API_BASE + `/v1/devices${query}`);

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message);
  }

  return response.json();
}

export async function createDevice(
  payload: CreateDevicePayload,
): Promise<Device> {
  const response = await fetch(API_BASE + "/v1/devices", {
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
  const response = await fetch(API_BASE + `/v1/devices/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message);
  }
}

export async function updateDevice(
  id: number,
  payload: UpdateDevicePayload,
): Promise<Device> {
  const response = await fetch(API_BASE + `/v1/devices/${id}`, {
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

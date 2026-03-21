package dto

type UpdateDeviceRequest struct {
	IP       *string `json:"ip"`
	Hostname *string `json:"hostname"`
	IsActive *bool   `json:"is_active"`
	Location *string `json:"location"`
}

type GetDevicesRequest struct {
	IsActive *string
	Search   *string
}

type CreateDeviceRequest struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	IsActive bool   `json:"is_active"`
	Location string `json:"location"`
}

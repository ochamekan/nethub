package dto

type UpdateDeviceRequest struct {
	Hostname *string `json:"hostname"`
	IP       *string `json:"ip"`
	IsActive *bool   `json:"is_active"`
}

package payload

// TenantRequest adalah payload untuk membuat tenant baru
type TenantRequest struct {
	TenantID string                 `json:"tenant_id"` // ID unik tenant
	Payload  map[string]interface{} `json:"payload"`   // Payload untuk pesan awal
}

// UpdateConcurrencyRequest adalah payload untuk memperbarui jumlah worker
type UpdateConcurrencyRequest struct {
	Workers int `json:"workers"` // Jumlah worker
}

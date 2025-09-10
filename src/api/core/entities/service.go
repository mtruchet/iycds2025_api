package entities

import "time"

// Service representa un servicio en el sistema
type Service struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	UserID       int64     `json:"user_id"`
	Category     string    `json:"category"`
	Price        float64   `json:"price"`
	Availability string    `json:"availability"` // JSON string con días y horarios
	Zones        string    `json:"zones"`        // JSON string con zonas de cobertura
	Status       string    `json:"status"`       // "active" o "inactive"
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ServiceCreate representa la solicitud de creación de servicio
type ServiceCreate struct {
	Title        string                 `json:"title" validate:"required,min=3,max=100"`
	Description  string                 `json:"description" validate:"required,min=100,max=1000"`
	Category     string                 `json:"category" validate:"required,min=2,max=50"`
	Price        float64                `json:"price" validate:"required,min=0"`
	Availability map[string]interface{} `json:"availability" validate:"required"`
	Zones        []Zone                 `json:"zones" validate:"required,min=1"`
	ImageURL     string                 `json:"image_url" validate:"omitempty,url"`
}

// ServiceUpdate representa la solicitud de actualización de servicio
type ServiceUpdate struct {
	Title        string                 `json:"title" validate:"omitempty,min=3,max=100"`
	Description  string                 `json:"description" validate:"omitempty,min=100,max=1000"`
	Category     string                 `json:"category" validate:"omitempty,min=2,max=50"`
	Price        *float64               `json:"price" validate:"omitempty,min=0"`
	Availability map[string]interface{} `json:"availability" validate:"omitempty"`
	Zones        []Zone                 `json:"zones" validate:"omitempty,min=1"`
	ImageURL     string                 `json:"image_url" validate:"omitempty,url"`
}

// Zone representa una zona de cobertura
type Zone struct {
	Province  string `json:"province" validate:"required,min=2,max=50"`
	Locality  string `json:"locality" validate:"required,min=2,max=50"`
	Neighborhood string `json:"neighborhood" validate:"omitempty,min=2,max=50"`
}

// ServiceResponse representa la respuesta de un servicio
type ServiceResponse struct {
	ID           int64                  `json:"id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	Price        float64                `json:"price"`
	Availability map[string]interface{} `json:"availability"`
	Zones        []Zone                 `json:"zones"`
	Status       string                 `json:"status"`
	ImageURL     string                 `json:"image_url"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ServiceListResponse representa la respuesta de lista de servicios
type ServiceListResponse struct {
	Services []ServiceResponse `json:"services"`
	Total    int               `json:"total"`
}

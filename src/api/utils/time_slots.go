package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DayAvailability representa la disponibilidad de un día
type DayAvailability struct {
	Available bool   `json:"available"`
	Start     string `json:"start"`
	End       string `json:"end"`
}

// GenerateTimeSlots genera slots de 30 minutos entre start y end
func GenerateTimeSlots(startTime, endTime string) ([]string, error) {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %s", startTime)
	}

	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %s", endTime)
	}

	if start.After(end) || start.Equal(end) {
		return nil, fmt.Errorf("start time must be before end time")
	}

	var slots []string
	current := start

	// Generar slots de 30 minutos
	for current.Before(end) {
		slotEnd := current.Add(30 * time.Minute)
		if slotEnd.After(end) {
			break // No crear slot si se pasa del horario final
		}

		slot := fmt.Sprintf("%s-%s", 
			current.Format("15:04"), 
			slotEnd.Format("15:04"))
		slots = append(slots, slot)
		
		current = slotEnd
	}

	return slots, nil
}

// GetDayAvailabilityFromJSON extrae la disponibilidad de un día específico del JSON
func GetDayAvailabilityFromJSON(availabilityJSON string, dayName string) (*DayAvailability, error) {
	var availability map[string]interface{}
	err := json.Unmarshal([]byte(availabilityJSON), &availability)
	if err != nil {
		return nil, fmt.Errorf("failed to parse availability JSON: %w", err)
	}

	dayData, exists := availability[strings.ToLower(dayName)]
	if !exists {
		return &DayAvailability{Available: false}, nil
	}

	dayMap, ok := dayData.(map[string]interface{})
	if !ok {
		return &DayAvailability{Available: false}, nil
	}

	// Verificar si está marcado como no disponible
	if available, exists := dayMap["available"]; exists {
		if avail, ok := available.(bool); ok && !avail {
			return &DayAvailability{Available: false}, nil
		}
	}

	// Extraer start y end
	start, startExists := dayMap["start"]
	end, endExists := dayMap["end"]

	if !startExists || !endExists {
		return &DayAvailability{Available: false}, nil
	}

	startStr, ok1 := start.(string)
	endStr, ok2 := end.(string)

	if !ok1 || !ok2 {
		return &DayAvailability{Available: false}, nil
	}

	return &DayAvailability{
		Available: true,
		Start:     startStr,
		End:       endStr,
	}, nil
}

// ValidateTimeFormat valida que el formato de hora sea HH:MM
func ValidateTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// ValidateDateFormat valida que el formato de fecha sea YYYY-MM-DD
func ValidateDateFormat(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// GetDayOfWeek retorna el nombre del día de la semana en inglés para una fecha
func GetDayOfWeek(dateStr string) (string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	dayNames := []string{
		"sunday", "monday", "tuesday", "wednesday", 
		"thursday", "friday", "saturday",
	}

	return dayNames[date.Weekday()], nil
}

// GetDayOfWeekInSpanish retorna el nombre del día de la semana en español
func GetDayOfWeekInSpanish(weekday time.Weekday) string {
	dayNames := []string{
		"domingo", "lunes", "martes", "miércoles", 
		"jueves", "viernes", "sábado",
	}
	return dayNames[weekday]
}

// IsDateInFuture verifica si una fecha es futura (no permite fechas pasadas)
func IsDateInFuture(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	return date.After(today) || date.Equal(today)
}



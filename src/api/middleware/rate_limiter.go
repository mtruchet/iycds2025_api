package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implementa un simple rate limiter en memoria
type RateLimiter struct {
	// Mapa que almacena las solicitudes por IP con sus timestamps
	ips map[string][]time.Time
	mu  sync.Mutex
	// Número máximo de solicitudes en el periodo
	maxRequests int
	// Periodo de tiempo para el rate limiting
	window time.Duration
}

// NewRateLimiter crea un nuevo rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:         make(map[string][]time.Time),
		maxRequests: maxRequests,
		window:      window,
	}
}

// cleanupOldRequests elimina las solicitudes que están fuera de la ventana de tiempo
func (rl *RateLimiter) cleanupOldRequests(ip string, now time.Time) {
	cutoff := now.Add(-rl.window)
	newRequests := []time.Time{}

	for _, timestamp := range rl.ips[ip] {
		if timestamp.After(cutoff) {
			newRequests = append(newRequests, timestamp)
		}
	}

	if len(newRequests) > 0 {
		rl.ips[ip] = newRequests
	} else {
		delete(rl.ips, ip) // Eliminar la IP si no tiene solicitudes recientes
	}
}

// RateLimit crea un middleware de Gin para aplicar rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.Request.RemoteAddr // Fallback si no se puede separar el puerto
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		rl.cleanupOldRequests(ip, now)

		// Verificar si la IP ha excedido el límite
		if len(rl.ips[ip]) >= rl.maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Registrar esta solicitud
		rl.ips[ip] = append(rl.ips[ip], now)

		c.Next()
	}
}

// Implementa diferentes limitadores para diferentes niveles de seguridad
var (
	// Limitar a 5 solicitudes cada 60 segundos (para endpoints críticos como login)
	StrictRateLimiter = NewRateLimiter(5, 60*time.Second)

	// Limitar a 10 solicitudes cada 60 segundos (para endpoints menos críticos)
	StandardRateLimiter = NewRateLimiter(10, 60*time.Second)
)

// Middleware convenientes para diferentes niveles de protección
func StrictRateLimit() gin.HandlerFunc {
	return StrictRateLimiter.RateLimit()
}

func StandardRateLimit() gin.HandlerFunc {
	return StandardRateLimiter.RateLimit()
}

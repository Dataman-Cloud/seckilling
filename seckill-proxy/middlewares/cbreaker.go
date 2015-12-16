package middlewares

import (
	"log"
	"net/http"

	"github.com/Dataman-Cloud/seckilling/seckill-proxy/oxy/cbreaker"
)

// CircuitBreaker holds the oxy circuit breaker.
type CircuitBreaker struct {
	circuitBreaker *cbreaker.CircuitBreaker
}

// NewCircuitBreaker returns a new CircuitBreaker.
func NewCircuitBreaker(next http.Handler, expression string, options ...cbreaker.CircuitBreakerOption) *CircuitBreaker {
	circuitBreaker, err := cbreaker.New(next, expression, options...)
	if err != nil {
		log.Panicln("shit happened", err)
	}
	return &CircuitBreaker{circuitBreaker}
}

func (cb *CircuitBreaker) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cb.circuitBreaker.ServeHTTP(rw, r)
}

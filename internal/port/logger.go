package port

// Logger define um contrato mínimo para logging na aplicação.
// Implementações concretas ficam em pkg/logger e são injetadas nas camadas internas.
// Mantemos uma interface simples para não vazar dependências externas.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

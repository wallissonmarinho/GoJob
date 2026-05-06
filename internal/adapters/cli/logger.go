package cli

import "log"

// Logger provides logging capabilities
type Logger struct {
	verbose bool
}

// NewLogger creates a new Logger
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

// Info logs an info message if verbose is enabled
func (l *Logger) Info(msg string) {
	if l.verbose {
		log.Printf("ℹ️  %s", msg)
	}
}

// Success logs a success message
func (l *Logger) Success(msg string) {
	log.Printf("✅ %s", msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	log.Printf("❌ %s", msg)
}

// Debugf logs a debug message if verbose is enabled
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.verbose {
		log.Printf("🔍 "+format, args...)
	}
}

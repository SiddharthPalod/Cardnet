package gates

// Config holds behavior gate configuration flags
type Config struct {
	MLEnabled      bool `yaml:"ml_enabled"`
	IssuerRequired  bool `yaml:"issuer_required"`
	AnalyticsEnabled bool `yaml:"analytics_enabled"`
}

// DefaultConfig returns default gate configuration
func DefaultConfig() Config {
	return Config{
		MLEnabled:       true,
		IssuerRequired:  true,
		AnalyticsEnabled: true,
	}
}

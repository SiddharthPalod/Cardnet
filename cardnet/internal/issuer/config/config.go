package config

import (
	"cardnet/internal/issuer/model"
	"os"

	"gopkg.in/yaml.v3"
)

type IssuerYAML struct {
	Issuers []struct {
		Name        string   `yaml:"name"`
		Bins        []string `yaml:"bins"`
		LatencyMs   int      `yaml:"latency_ms"`
		DeclineRate float64  `yaml:"decline_rate"`
		TimeoutRate float64  `yaml:"timeout_rate"`
		Chaos       struct {
			ForceTimeout   bool `yaml:"force_timeout"`
			ForceDecline   bool `yaml:"force_decline"`
			ExtraLatencyMs int  `yaml:"extra_latency_ms"`
		} `yaml:"chaos"`
	} `yaml:"issuers"`
}

func LoadIssuers(path string) (map[string]model.IssuerProfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg IssuerYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	table := make(map[string]model.IssuerProfile)
	for _, issuer := range cfg.Issuers {
		profile := model.IssuerProfile{
			Name:        issuer.Name,
			LatencyMs:   issuer.LatencyMs,
			DeclineRate: issuer.DeclineRate,
			TimeoutRate: issuer.TimeoutRate,
			Chaos: model.ChaosConfig{
				ForceTimeout:   issuer.Chaos.ForceTimeout,
				ForceDecline:   issuer.Chaos.ForceDecline,
				ExtraLatencyMs: issuer.Chaos.ExtraLatencyMs,
			},
		}

		for _, bin := range issuer.Bins {
			table[bin] = profile
		}
	}

	return table, nil
}

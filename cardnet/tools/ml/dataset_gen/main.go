package main

import (
	"cardnet/internal/analytics/features"
	"cardnet/internal/analytics/ml/training"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 1. Generate Synthetic Feature Snapshot
	merchants := make([]features.MerchantFeatures, 0)
	labels := make(map[string]int)

	for i := 0; i < 1000; i++ {
		mid := fmt.Sprintf("M%05d", i)

		// Random features
		txCount := rand.Intn(100)
		approvalRate := rand.Float64()
		avgAmount := rand.Float64() * 500

		// Simple logic for label generation (correlation)
		// If approval rate is low and tx count is high -> Risk?
		// Just random logic for demo
		isHighRisk := 0
		if approvalRate < 0.5 && txCount > 50 {
			isHighRisk = 1
		} else if rand.Float64() < 0.05 {
			isHighRisk = 1 // 5% noise
		}

		merchants = append(merchants, features.MerchantFeatures{
			MerchantID:      mid,
			TxCount24h:      txCount,
			ApprovalRate24h: approvalRate,
			AvgAmount24h:    avgAmount,
			DeclineRate24h:  1.0 - approvalRate,
		})
		labels[mid] = isHighRisk
	}

	snapshot := &features.FeatureSnapshot{
		Merchants: merchants,
	}

	// 2. Build Dataset
	rows := training.BuildMerchantDataset(snapshot, labels)
	fmt.Printf("Generated %d rows\n", len(rows))

	// 3. Export to Tools/ML/Datasets
	// Assumption: run from project root or navigate relative
	// d:\New folder (4)\SidFiles\Projectd\SystemDesign\GoLang\cardnet\tools\ml\datasets
	// We will use absolute path based on user info if possible, or relative to this file?
	// This file is in tools/ml/dataset_gen/main.go
	// So ../datasets/ is the target.

	featurePath, _ := filepath.Abs("tools/ml/datasets/merchant_features.csv")
	labelPath, _ := filepath.Abs("tools/ml/datasets/merchant_labels.csv")

	fmt.Printf("Exporting to %s and %s\n", featurePath, labelPath)

	err := training.ExportMerchantsToCSV(rows, featurePath, labelPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Dataset successfully exported!")

	// 4. Generate BINs
	bins := make([]features.BINFeatures, 0)
	for i := 0; i < 50; i++ {
		bins = append(bins, features.BINFeatures{
			BIN:             fmt.Sprintf("4111%02d", i),
			TxCount24h:      rand.Intn(1000),
			ApprovalRate24h: rand.Float64(),
		})
	}
	binPath, _ := filepath.Abs("tools/ml/datasets/bin_features.csv")
	if err := training.ExportBINsToCSV(bins, binPath); err != nil {
		panic(err)
	}

	// 5. Generate Rules
	rules := make([]training.RuleExportRow, 0)
	for i := 0; i < 20; i++ {
		rules = append(rules, training.RuleExportRow{
			RuleID:            fmt.Sprintf("R%d", i),
			FalsePositiveRate: rand.Float64(),
		})
	}
	rulePath, _ := filepath.Abs("tools/ml/datasets/rule_effectiveness.csv")
	if err := training.ExportRulesToCSV(rules, rulePath); err != nil {
		panic(err)
	}
	fmt.Println("All datasets generated.")
}

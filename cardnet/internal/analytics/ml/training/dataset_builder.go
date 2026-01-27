package training

import (
	"cardnet/internal/analytics/features"
	"encoding/csv"
	"fmt"
	"os"
)

type TrainingRow struct {
	EntityID string
	Features map[string]float64
	Label    int
}

func BuildMerchantDataset(
	snapshot *features.FeatureSnapshot,
	labels map[string]int,
) []TrainingRow {

	var rows []TrainingRow
	for _, m := range snapshot.Merchants {
		label := labels[m.MerchantID]

		rows = append(rows, TrainingRow{
			EntityID: m.MerchantID,
			Features: map[string]float64{
				"tx_count":      float64(m.TxCount24h),
				"approval_rate": m.ApprovalRate24h,
				"avg_amount":    m.AvgAmount24h,
			},
			Label: label,
		})
	}
	return rows
}

func ExportMerchantsToCSV(rows []TrainingRow, featurePath, labelPath string) error {
	// Write Features
	fFile, err := os.Create(featurePath)
	if err != nil {
		return err
	}
	defer fFile.Close()

	fWriter := csv.NewWriter(fFile)
	// Header for features
	if err := fWriter.Write([]string{"merchant_id", "tx_count_24h", "approval_rate_24h", "avg_amount_24h"}); err != nil {
		return err
	}

	// Write Labels
	lFile, err := os.Create(labelPath)
	if err != nil {
		return err
	}
	defer lFile.Close()

	lWriter := csv.NewWriter(lFile)
	// Header for labels
	if err := lWriter.Write([]string{"merchant_id", "label"}); err != nil {
		return err
	}

	for _, row := range rows {
		// Features
		err := fWriter.Write([]string{
			row.EntityID,
			fmt.Sprintf("%f", row.Features["tx_count"]),
			fmt.Sprintf("%f", row.Features["approval_rate"]),
			fmt.Sprintf("%f", row.Features["avg_amount"]),
		})
		if err != nil {
			return err
		}

		// Labels
		err = lWriter.Write([]string{
			row.EntityID,
			fmt.Sprintf("%d", row.Label),
		})
		if err != nil {
			return err
		}
	}

	fWriter.Flush()
	lWriter.Flush()
	return nil
}

func ExportBINsToCSV(bins []features.BINFeatures, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"bin", "tx_count", "approval_rate_24h"}); err != nil {
		return err
	}

	for _, b := range bins {
		err := writer.Write([]string{
			b.BIN,
			fmt.Sprintf("%d", b.TxCount24h),
			fmt.Sprintf("%f", b.ApprovalRate24h),
		})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// Defining a wrapper for RuleStats for export if needed, or using aggregation type
// Assuming accessing aggregations.RuleStats from here
// We need to synthesize FPR since it's not in RuleStats
type RuleExportRow struct {
	RuleID            string
	FalsePositiveRate float64
}

func ExportRulesToCSV(rules []RuleExportRow, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"rule_id", "false_positive_rate"}); err != nil {
		return err
	}

	for _, r := range rules {
		err := writer.Write([]string{
			r.RuleID,
			fmt.Sprintf("%f", r.FalsePositiveRate),
		})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

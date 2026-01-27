package reporting

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ExportCSV(report *BatchReport, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"batch_id",
		"batch_date",
		"status",
		"total_count",
		"total_amount",
		"currency",
	})

	return writer.Write([]string{
		report.BatchID,
		report.BatchDate,
		report.Status,
		fmt.Sprintf("%d", report.TotalCount),
		fmt.Sprintf("%d", report.TotalAmount),
		report.Currency,
	})
}

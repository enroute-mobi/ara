package audit

import "testing"

func TestDatasetFormat(t *testing.T) {
	bq := NewBigQuery("dataset-with_hyphen").(*BigQueryClient)

	if bq.dataset != "dataset_with_hyphen" {
		t.Errorf("Wrong format for BQ Client: expected 'dataset_with_hyphen' got '%s'", bq.dataset)
	}
}

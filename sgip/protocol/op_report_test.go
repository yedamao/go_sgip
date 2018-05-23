package protocol

import (
	"testing"
)

func TestReport(t *testing.T) {
	report, err := NewReport(
		[...]uint32{0, 0, 0},
		[...]uint32{0, 0, 0},
		0, "17611000000", 2, 9,
	)
	if err != nil {
		t.Error(err)
	}

	raw := report.Serialize()
	parsedOp, err := ParseOperation(raw)
	if err != nil {
		t.Error(err)
	}

	parsedReport := parsedOp.(*Report)

	if report.Length != parsedReport.Length ||
		report.SubmitSequence[0] != parsedReport.SubmitSequence[0] ||
		report.SubmitSequence[1] != parsedReport.SubmitSequence[1] ||
		report.SubmitSequence[2] != parsedReport.SubmitSequence[2] ||
		report.ReportType != parsedReport.ReportType ||
		report.UserNumber.String() != parsedReport.UserNumber.String() ||
		report.State != parsedReport.State ||
		report.ErrorCode != parsedReport.ErrorCode {

		t.Error("report not equal")
	}
}

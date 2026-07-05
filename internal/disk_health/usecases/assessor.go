package usecases

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
)

func (s *Service) AssessReport(report *domain.SmartReport) domain.Assessment {
	assessment := domain.Assessment{
		Status:        domain.StatusHealthy,
		Score:         100,
		Summary:       "SMART情報に重大な問題は見つかりませんでした。",
		OverallHealth: report.OverallHealth,
		Model:         report.Model,
		SerialNumber:  report.SerialNumber,
		Attributes:    report.Attributes,
	}

	if report.OverallHealth == "" || len(report.Attributes) == 0 {
		assessment.Status = domain.StatusUnknown
		assessment.Score = 0
		assessment.Summary = "SMART情報が不足しているため健康度を判定できません。"
		return assessment
	}

	if strings.Contains(report.OverallHealth, "FAILED") {
		assessment.Findings = append(assessment.Findings, domain.Finding{
			Severity: domain.SeverityCritical,
			Message:  "SMART全体ヘルスがFAILEDです",
		})
	}

	for _, attribute := range report.Attributes {
		s.assessAttribute(&assessment, attribute)
	}

	s.finalizeAssessment(&assessment)
	return assessment
}

func (s *Service) assessAttribute(assessment *domain.Assessment, attribute domain.SmartAttribute) {
	if attribute.WhenFailed != "-" {
		s.addFinding(assessment, attribute, domain.SeverityCritical, fmt.Sprintf("%s はしきい値を下回っています", attribute.Name))
		return
	}

	switch {
	case s.isAttribute(attribute, 197, "Current_Pending_Sector") && attribute.RawValue > 0:
		s.addFinding(assessment, attribute, domain.SeverityCritical, "代替待ちセクタを検出しました")
	case s.isAttribute(attribute, 198, "Offline_Uncorrectable") && attribute.RawValue > 0:
		s.addFinding(assessment, attribute, domain.SeverityCritical, "オフライン訂正不能セクタを検出しました")
	case s.isAttribute(attribute, 5, "Reallocated_Sector_Ct") && attribute.RawValue > 0:
		s.addFinding(assessment, attribute, domain.SeverityWarning, "代替処理済みセクタを検出しました")
	case s.isAttribute(attribute, 196, "Reallocated_Event_Count") && attribute.RawValue > 0:
		s.addFinding(assessment, attribute, domain.SeverityWarning, "セクタ代替イベントを検出しました")
	case s.isAttribute(attribute, 199, "UDMA_CRC_Error_Count") && attribute.RawValue > 0:
		s.addFinding(assessment, attribute, domain.SeverityWarning, "UDMA CRCエラーを検出しました。ケーブルや接続状態を確認してください")
	case s.isTemperatureAttribute(attribute) && attribute.RawValue >= 60:
		s.addFinding(assessment, attribute, domain.SeverityCritical, "ディスク温度が危険域です")
	case s.isTemperatureAttribute(attribute) && attribute.RawValue >= 50:
		s.addFinding(assessment, attribute, domain.SeverityWarning, "ディスク温度が高めです")
	}
}

func (s *Service) isAttribute(attribute domain.SmartAttribute, id int, name string) bool {
	return attribute.ID == id || attribute.Name == name
}

func (s *Service) isTemperatureAttribute(attribute domain.SmartAttribute) bool {
	return s.isAttribute(attribute, 194, "Temperature_Celsius") || s.isAttribute(attribute, 190, "Airflow_Temperature_Cel")
}

func (s *Service) addFinding(assessment *domain.Assessment, attribute domain.SmartAttribute, severity domain.Severity, message string) {
	assessment.Findings = append(assessment.Findings, domain.Finding{
		Severity:      severity,
		AttributeID:   attribute.ID,
		AttributeName: attribute.Name,
		RawValue:      attribute.RawValue,
		Message:       message,
	})
}

func (s *Service) finalizeAssessment(assessment *domain.Assessment) {
	hasCritical := false
	hasWarning := false
	for _, finding := range assessment.Findings {
		switch finding.Severity {
		case domain.SeverityCritical:
			hasCritical = true
		case domain.SeverityWarning:
			hasWarning = true
		}
	}

	switch {
	case hasCritical:
		assessment.Status = domain.StatusCritical
		assessment.Score = 20
		assessment.Summary = "重大な SMART 欠損指標を検出しました。速やかなバックアップと交換を推奨します。"
	case hasWarning:
		assessment.Status = domain.StatusWarning
		assessment.Score = 60
		assessment.Summary = "注意が必要な SMART 指標を検出しました。バックアップと経過観察を推奨します。"
	default:
		assessment.Status = domain.StatusHealthy
		assessment.Score = 100
		assessment.Summary = "SMART情報に重大な問題は見つかりませんでした。"
	}
}

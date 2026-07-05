package usecases

import "github.com/landmaster135/devbox/internal/disk_health/domain"

func (s *Service) AssessDmesgEvents(events []domain.DmesgEvent) domain.DmesgAssessment {
	assessment := domain.DmesgAssessment{
		Status:   domain.StatusHealthy,
		Score:    100,
		Summary:  "dmesgログにディスクI/Oの重大な問題は見つかりませんでした。",
		Findings: events,
	}

	if len(events) == 0 {
		return assessment
	}

	hasCritical := false
	for _, event := range events {
		if event.Severity == domain.SeverityCritical {
			hasCritical = true
			break
		}
	}

	if hasCritical {
		assessment.Status = domain.StatusCritical
		assessment.Score = 20
		assessment.Summary = "dmesgログからディスクI/Oの重大エラーを検出しました。速やかなバックアップと交換を推奨します。"
		return assessment
	}

	assessment.Status = domain.StatusWarning
	assessment.Score = 60
	assessment.Summary = "dmesgログからディスクI/Oの注意イベントを検出しました。バックアップと経過観察を推奨します。"
	return assessment
}

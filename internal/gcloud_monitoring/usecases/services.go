package usecases

import (
	"context"
	"fmt"
	"log"

	dashboard "cloud.google.com/go/monitoring/dashboard/apiv1"
	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

// Service はGoogle Cloud Monitoringサービスの操作を提供する
type Service struct {
	project          string
	location         string
	serviceName      string
	serviceAccountID string
}

// NewService は新しいServiceインスタンスを作成する
func NewService(project, location, serviceName, serviceAccountID string) *Service {
	return &Service{
		project:          project,
		location:         location,
		serviceName:      serviceName,
		serviceAccountID: serviceAccountID,
	}
}

// CreateDashboardForCloudRun はCloud Runサービス用のモニタリングダッシュボードを作成する
func (s *Service) CreateDashboardForCloudRun() (string, error) {
	ctx := context.Background()

	// Cloud Runサービスの存在確認
	log.Printf("Cloud Runサービスの存在確認中: %s (プロジェクト: %s, ロケーション: %s)", s.serviceName, s.project, s.location)
	exists, err := s.verifyCloudRunService(ctx)
	if err != nil {
		return "", fmt.Errorf("Cloud Runサービスの確認に失敗しました: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("指定されたCloud Runサービスが見つかりません: %s (プロジェクト: %s, ロケーション: %s)", s.serviceName, s.project, s.location)
	}

	log.Printf("Cloud Runサービスが確認されました: %s", s.serviceName)

	// ダッシュボードの作成
	log.Printf("モニタリングダッシュボードを作成中...")
	dashboardName, err := s.createMonitoringDashboard(ctx)
	if err != nil {
		return "", fmt.Errorf("ダッシュボードの作成に失敗しました: %v", err)
	}

	return fmt.Sprintf("ダッシュボードが正常に作成されました: %s", dashboardName), nil
}

// verifyCloudRunService はCloud Runサービスの存在を確認する
func (s *Service) verifyCloudRunService(ctx context.Context) (bool, error) {
	// クライアントオプションの設定
	opts, err := s.getClientOptions(ctx)
	if err != nil {
		return false, fmt.Errorf("クライアントオプションの設定に失敗しました: %v", err)
	}

	// Cloud Run クライアントの作成
	client, err := run.NewServicesClient(ctx, opts...)
	if err != nil {
		return false, fmt.Errorf("Cloud Runクライアントの作成に失敗しました: %v", err)
	}
	defer client.Close()

	// サービス名の構築
	servicePath := fmt.Sprintf("projects/%s/locations/%s/services/%s", s.project, s.location, s.serviceName)

	// サービスの取得
	req := &runpb.GetServiceRequest{
		Name: servicePath,
	}

	_, err = client.GetService(ctx, req)
	if err != nil {
		// サービスが見つからない場合
		return false, nil
	}

	return true, nil
}

// createMonitoringDashboard はモニタリングダッシュボードを作成する
func (s *Service) createMonitoringDashboard(ctx context.Context) (string, error) {
	// クライアントオプションの設定
	opts, err := s.getClientOptions(ctx)
	if err != nil {
		return "", fmt.Errorf("クライアントオプションの設定に失敗しました: %v", err)
	}

	// Dashboard クライアントの作成
	client, err := dashboard.NewDashboardsClient(ctx, opts...)
	if err != nil {
		return "", fmt.Errorf("Dashboardクライアントの作成に失敗しました: %v", err)
	}
	defer client.Close()

	// ダッシュボードの設定
	dashboardConfig := s.buildDashboardConfig()

	// ダッシュボード作成リクエスト
	req := &dashboardpb.CreateDashboardRequest{
		Parent:    fmt.Sprintf("projects/%s", s.project),
		Dashboard: dashboardConfig,
	}

	// ダッシュボードの作成
	dashboard, err := client.CreateDashboard(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ダッシュボードの作成リクエストに失敗しました: %v", err)
	}

	return dashboard.Name, nil
}

// getClientOptions はGoogle Cloudクライアントのオプションを取得する
func (s *Service) getClientOptions(ctx context.Context) ([]option.ClientOption, error) {
	var opts []option.ClientOption

	// サービスアカウントが指定されている場合はimpersonationを使用
	if s.serviceAccountID != "" {
		serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", s.serviceAccountID, s.project)
		log.Printf("サービスアカウントを使用します: %s", serviceAccountEmail)

		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: serviceAccountEmail,
			Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		})
		if err != nil {
			return nil, fmt.Errorf("サービスアカウントのimpersonationに失敗しました: %v", err)
		}
		opts = append(opts, option.WithTokenSource(ts))
	} else {
		log.Printf("Application Default Credentials (ADC) を使用します")
	}

	return opts, nil
}

// buildDashboardConfig はダッシュボードの設定を構築する
func (s *Service) buildDashboardConfig() *dashboardpb.Dashboard {
	return &dashboardpb.Dashboard{
		DisplayName: fmt.Sprintf("Cloud Run Monitoring - %s", s.serviceName),
		Layout: &dashboardpb.Dashboard_GridLayout{
			GridLayout: &dashboardpb.GridLayout{
				Columns: 12,
				Widgets: s.buildDashboardWidgets(),
			},
		},
	}
}

// buildDashboardWidgets はダッシュボードのウィジェットを構築する
func (s *Service) buildDashboardWidgets() []*dashboardpb.Widget {
	var widgets []*dashboardpb.Widget

	// 行1: リクエスト概要 (4つのウィジェット、各3列幅)
	widgets = append(widgets, s.createTextWidget("リクエスト数 (req/sec)", "リクエスト数/秒の監視"))
	widgets = append(widgets, s.createTextWidget("ステータス別リクエスト (2xx/4xx/5xx)", "ステータスコード別リクエスト数"))
	widgets = append(widgets, s.createTextWidget("累積リクエスト数 (24h total)", "24時間の累積リクエスト数"))
	widgets = append(widgets, s.createTextWidget("リクエスト時間分布 (heatmap)", "時間別リクエストパターン"))

	// 行2: パフォーマンス指標 (4つのウィジェット、各3列幅)
	widgets = append(widgets, s.createTextWidget("リクエストレイテンシ (P50,P95,P99)", "レスポンス時間のパーセンタイル"))
	widgets = append(widgets, s.createTextWidget("エラー率", "エラーリクエストの割合"))
	widgets = append(widgets, s.createTextWidget("最大同時リクエスト", "同時処理リクエスト数"))
	widgets = append(widgets, s.createTextWidget("レスポンス時間 (Mean/Max)", "平均・最大レスポンス時間"))

	// 行3: コンテナ指標 (3つのウィジェット、各4列幅)
	widgets = append(widgets, s.createTextWidget("インスタンス数", "コンテナインスタンス数"))
	widgets = append(widgets, s.createTextWidget("起動レイテンシ", "コンテナ起動時間"))
	widgets = append(widgets, s.createTextWidget("課金時間", "課金対象インスタンス時間"))

	// 行4: リソース使用状況 (3つのウィジェット、各4列幅)
	widgets = append(widgets, s.createTextWidget("CPU使用率", "CPU使用率の監視"))
	widgets = append(widgets, s.createTextWidget("メモリ使用率", "メモリ使用率の監視"))
	widgets = append(widgets, s.createTextWidget("メモリ使用量", "メモリ使用量の監視"))

	// 行5: ネットワーク (2つのウィジェット、各6列幅)
	widgets = append(widgets, s.createTextWidget("送信バイト数", "ネットワーク送信量"))
	widgets = append(widgets, s.createTextWidget("受信バイト数", "ネットワーク受信量"))

	return widgets
}

// createTextWidget はテキストウィジェットを作成する
func (s *Service) createTextWidget(title, content string) *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: title,
		Content: &dashboardpb.Widget_Text{
			Text: &dashboardpb.Text{
				Content: fmt.Sprintf("**%s**\n\n%s\n\nサービス: %s", title, content, s.serviceName),
				Format:  dashboardpb.Text_MARKDOWN,
			},
		},
	}
}

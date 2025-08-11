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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return "", fmt.Errorf("cloud Runサービスの確認に失敗しました: %v", err)
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
		return false, fmt.Errorf("cloud Runクライアントの作成に失敗しました: %v", err)
	}
	defer client.Close()

	// サービス名の構築
	servicePath := fmt.Sprintf("projects/%s/locations/%s/services/%s", s.project, s.location, s.serviceName)
	log.Printf("サービスパスを確認中: %s", servicePath)

	// サービスの取得
	req := &runpb.GetServiceRequest{
		Name: servicePath,
	}

	service, err := client.GetService(ctx, req)
	if err != nil {
		// gRPCステータスコードを確認
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				log.Printf("Cloud Runサービスが見つかりません: %s", servicePath)
				return false, nil
			case codes.PermissionDenied:
				return false, fmt.Errorf("cloud Runサービスへのアクセス権限がありません: %v", err)
			case codes.Unauthenticated:
				return false, fmt.Errorf("認証に失敗しました: %v", err)
			default:
				return false, fmt.Errorf("cloud Runサービスの取得中にエラーが発生しました: %v", err)
			}
		}
		// gRPCエラーでない場合
		return false, fmt.Errorf("予期しないエラーが発生しました: %v", err)
	}

	log.Printf("Cloud Runサービスが見つかりました: %s (状態: %s)", service.Name, service.GetConditions())
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
		return "", fmt.Errorf("dashboardクライアントの作成に失敗しました: %v", err)
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

func (s *Service) createDisplayTitleOfDashboard() string {
	return fmt.Sprintf("CloudRun ダッシュボード: %s", s.serviceName)
}

// buildDashboardConfig はダッシュボードの設定を構築する
func (s *Service) buildDashboardConfig() *dashboardpb.Dashboard {
	return &dashboardpb.Dashboard{
		DisplayName: s.createDisplayTitleOfDashboard(),
		Layout: &dashboardpb.Dashboard_MosaicLayout{
			MosaicLayout: &dashboardpb.MosaicLayout{
				Columns: 12,
				Tiles:   s.buildDashboardTiles(),
			},
		},
	}
}

// buildDashboardTiles はダッシュボードのタイルを構築する
func (s *Service) buildDashboardTiles() []*dashboardpb.MosaicLayout_Tile {
	var tiles []*dashboardpb.MosaicLayout_Tile

	// 行1: リクエスト概要 (4つのウィジェット、各3列幅)
	tiles = append(tiles, s.createTile(s.createRequestCountWidget(), 0, 0, 3, 4))
	tiles = append(tiles, s.createTile(s.createRequestsByStatusWidget(), 3, 0, 3, 4))
	tiles = append(tiles, s.createTile(s.createTotalRequestsWidget(), 6, 0, 3, 4))
	tiles = append(tiles, s.createTile(s.createLogByteByHourWidget(), 9, 0, 3, 4))

	// 行2: パフォーマンス指標 (4つのウィジェット、各3列幅)
	tiles = append(tiles, s.createTile(s.createRequestLatencyWidget(), 0, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createErrorRateWidget(), 3, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createMaxConcurrentRequestsWidget(), 6, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createResponseTimeWidget(), 9, 4, 3, 4))

	// 行3: コンテナ指標 (3つのウィジェット、各4列幅)
	tiles = append(tiles, s.createTile(s.createContainerInstanceCountWidget(), 0, 8, 4, 4))
	tiles = append(tiles, s.createTile(s.createContainerStartupLatenciesWidget(), 4, 8, 4, 4))
	tiles = append(tiles, s.createTile(s.createContainerBillableInstanceTimeWidget(), 8, 8, 4, 4))

	// 行4: リソース使用状況 (3つのウィジェット、各4列幅)
	tiles = append(tiles, s.createTile(s.createContainerCPUUtilizationsWidget(), 0, 12, 4, 4))
	tiles = append(tiles, s.createTile(s.createContainerMemoryUtilizationsWidget(), 4, 12, 4, 4))
	tiles = append(tiles, s.createTile(s.createContainerMemoryUsageTimeWidget(), 8, 12, 4, 4))

	// 行5: ネットワーク (2つのウィジェット、各6列幅)
	tiles = append(tiles, s.createTile(s.createNetworkSentBytesWidget(), 0, 16, 6, 4))
	tiles = append(tiles, s.createTile(s.createNetworkReceivedBytesWidget(), 6, 16, 6, 4))

	return tiles
}

// createTile はウィジェットとレイアウト情報からタイルを作成する
func (s *Service) createTile(widget *dashboardpb.Widget, xPos, yPos, width, height int32) *dashboardpb.MosaicLayout_Tile {
	return &dashboardpb.MosaicLayout_Tile{
		XPos:   xPos,
		YPos:   yPos,
		Width:  width,
		Height: height,
		Widget: widget,
	}
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

const (
	containerWidgetTitle01                  = "Requests per Second"
	containerWidgetTitle02                  = "Requests by Status Code"
	containerWidgetTitle03                  = "Total Requests (24h)"
	containerWidgetTitle04                  = "Log Bytes by Hour"
	containerWidgetTitle05                  = "Request Latency (P50, P95, P99)"
	containerWidgetTitle06                  = "Error Rate (%)"
	containerWidgetTitle07                  = "Max Concurrent Requests (P50/P95/P99)"
	containerWidgetTitle08                  = "Response Time (P50/P95/P99)"
	containerWidgetTitle11                  = "Container Instance Count"
	containerWidgetTitle12                  = "Container Startup Latency (P50/P95/P99)"
	containerWidgetTitle13                  = "Billable Instance Time"
	containerWidgetTitle14                  = "CPU Utilization (P50/P95/P99)"
	containerWidgetTitle15                  = "Memory Utilization (P50/P95/P99)"
	containerWidgetTitle16                  = "Memory Usage (P50/P95/P99)"
	containerWidgetTitle51                  = "Network Sent Bytes"
	containerWidgetTitle52                  = "Network Received Bytes"
	datasetLegendTemplateOfP50              = "P50"
	datasetLegendTemplateOfP95              = "P95"
	datasetLegendTemplateOfP99              = "P99"
	yAxisLabelOfRequestsPerSecond           = "requests/second"
	yAxisLabelOfBytesPerSecond              = "bytes/second"
	yAxisLabelOfLatencyMilliSec             = "latency (ms)"
	yAxisLabelOfErrorRatePercentage         = "error rate (%)"
	yAxisLabelOfConcurrentRequests          = "concurrent requests"
	yAxisLabelOfResponseTimeMilliSec        = "response time (ms)"
	yAxisLabelOfInstanceCount               = "instances"
	yAxisLabelOfStartupLatencyMilliSec      = "startup latency (ms)"
	yAxisLabelOfBillableInstanceTimeSeconds = "billable time (seconds)"
	yAxisLabelOfCPUUtilizationPercentage    = "CPU utilization (%)"
	yAxisLabelOfMemoryUtilizationPercentage = "memory utilization (%)"
	yAxisLabelOfMemoryUsageBytes            = "memory usage (bytes)"
	yAxisLabelOfNetworkBytesPerSecond       = "bytes/second"
	metricOfRequestCount                    = "run.googleapis.com/request_count"
	metricOfRequestLatencies                = "run.googleapis.com/request_latencies"
	metricOfLoggingByteCount                = "logging.googleapis.com/byte_count"
	metricOfMaxRequestConcurrencies         = "run.googleapis.com/container/max_request_concurrencies"
	metricOfContainerInstanceCount          = "run.googleapis.com/container/instance_count"
	metricOfContainerStartupLatencies       = "run.googleapis.com/container/startup_latencies"
	metricOfContainerBillableInstance       = "run.googleapis.com/container/billable_instance_time"
	metricOfContainerCPUUtilizations        = "run.googleapis.com/container/cpu/utilizations"
	metricOfContainerMemoryUtilizations     = "run.googleapis.com/container/memory/utilizations"
	metricOfContainerMemoryUsage            = "run.googleapis.com/container/memory/usage"
	metricOfNetworkSentBytesCount           = "run.googleapis.com/container/network/sent_bytes_count"
	metricOfNetworkReceivedBytesCount       = "run.googleapis.com/container/network/received_bytes_count"
)

func (s *Service) createPromQLForRequestCount() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfRequestCount)
}

func (s *Service) createPromQLForRequestLatencies() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfRequestLatencies)
}

func (s *Service) createPromQLForLogging() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfLoggingByteCount)
}

func (s *Service) createPromQLForRequestCountByResponseCode() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s" metric.label.response_code_class!="2xx"`, s.serviceName, metricOfRequestCount)
}

func (s *Service) createPromQLMaxRequestConcurrencies() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfMaxRequestConcurrencies)
}

func (s *Service) createPromQLForContainerInstanceCount() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerInstanceCount)
}

func (s *Service) createPromQLForContainerStartupLatencies() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerStartupLatencies)
}

func (s *Service) createPromQLForContainerBillableInstanceTime() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerBillableInstance)
}

func (s *Service) createPromQLForContainerCPUUtilizations() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerCPUUtilizations)
}

func (s *Service) createPromQLForContainerMemoryUtilizations() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerMemoryUtilizations)
}

func (s *Service) createPromQLForContainerMemoryUsageTime() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfContainerMemoryUsage)
}

func (s *Service) createPromQLForNetworkSentBytes() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfNetworkSentBytesCount)
}

func (s *Service) createPromQLForNetworkReceivedBytes() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`, s.serviceName, metricOfNetworkReceivedBytesCount)
}

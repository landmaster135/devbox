package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	dashboard "cloud.google.com/go/monitoring/dashboard/apiv1"
	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
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
	tiles = append(tiles, s.createTile(s.createRequestHeatmapWidget(), 9, 0, 3, 4))

	// 行2: パフォーマンス指標 (4つのウィジェット、各3列幅)
	tiles = append(tiles, s.createTile(s.createRequestLatencyWidget(), 0, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createErrorRateWidget(), 3, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createMaxConcurrentRequestsWidget(), 6, 4, 3, 4))
	tiles = append(tiles, s.createTile(s.createResponseTimeWidget(), 9, 4, 3, 4))

	// 行3: コンテナ指標 (3つのウィジェット、各4列幅)
	tiles = append(tiles, s.createTile(s.createTextWidget("インスタンス数", "コンテナインスタンス数"), 0, 8, 4, 4))
	tiles = append(tiles, s.createTile(s.createTextWidget("起動レイテンシ", "コンテナ起動時間"), 4, 8, 4, 4))
	tiles = append(tiles, s.createTile(s.createTextWidget("課金時間", "課金対象インスタンス時間"), 8, 8, 4, 4))

	// 行4: リソース使用状況 (3つのウィジェット、各4列幅)
	tiles = append(tiles, s.createTile(s.createTextWidget("CPU使用率", "CPU使用率の監視"), 0, 12, 4, 4))
	tiles = append(tiles, s.createTile(s.createTextWidget("メモリ使用率", "メモリ使用率の監視"), 4, 12, 4, 4))
	tiles = append(tiles, s.createTile(s.createTextWidget("メモリ使用量", "メモリ使用量の監視"), 8, 12, 4, 4))

	// 行5: ネットワーク (2つのウィジェット、各6列幅)
	tiles = append(tiles, s.createTile(s.createTextWidget("送信バイト数", "ネットワーク送信量"), 0, 16, 6, 4))
	tiles = append(tiles, s.createTile(s.createTextWidget("受信バイト数", "ネットワーク受信量"), 6, 16, 6, 4))

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

func (s *Service) createPromQL() string {
	return fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_count"`, s.serviceName)
}

// createRequestCountWidget はリクエスト数（req/sec）ウィジェットを作成する
func (s *Service) createRequestCountWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Requests per Second",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQL(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_LINE,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "requests/second",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createRequestsByStatusWidget はステータス別リクエスト数ウィジェットを作成する
func (s *Service) createRequestsByStatusWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Requests by Status Code",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQL(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										GroupByFields:      []string{"metric.label.response_code_class"},
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_STACKED_AREA,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "requests/second",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createTotalRequestsWidget は累積リクエスト数ウィジェットを作成する
func (s *Service) createTotalRequestsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Total Requests (24h)",
		Content: &dashboardpb.Widget_Scorecard{
			Scorecard: &dashboardpb.Scorecard{
				TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
					Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
						TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
							Filter: s.createPromQL(),
							Aggregation: &dashboardpb.Aggregation{
								AlignmentPeriod:    durationpb.New(86400 * time.Second), // 24時間
								PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_SUM,
								CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
							},
						},
					},
				},
				DataView: &dashboardpb.Scorecard_SparkChartView_{
					SparkChartView: &dashboardpb.Scorecard_SparkChartView{
						SparkChartType: dashboardpb.SparkChartType_SPARK_LINE,
					},
				},
			},
		},
	}
}

// createRequestHeatmapWidget はリクエスト時間分布ウィジェットを作成する
func (s *Service) createRequestHeatmapWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Request Pattern by Hour",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQL(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(3600 * time.Second), // 1時間
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_HEATMAP,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "requests/second",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createRequestLatencyWidget はリクエストレイテンシ (P50,P95,P99) ウィジェットを作成する
func (s *Service) createRequestLatencyWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Request Latency (P50, P95, P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P50",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P95",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P99",
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "latency (ms)",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createErrorRateWidget はエラー率ウィジェットを作成する
func (s *Service) createErrorRateWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Error Rate (%)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilterRatio{
								TimeSeriesFilterRatio: &dashboardpb.TimeSeriesFilterRatio{
									Numerator: &dashboardpb.TimeSeriesFilterRatio_RatioPart{
										Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_count" metric.label.response_code_class!="2xx"`, s.serviceName),
										Aggregation: &dashboardpb.Aggregation{
											AlignmentPeriod:    durationpb.New(60 * time.Second),
											PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
											CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										},
									},
									Denominator: &dashboardpb.TimeSeriesFilterRatio_RatioPart{
										Filter: s.createPromQL(),
										Aggregation: &dashboardpb.Aggregation{
											AlignmentPeriod:    durationpb.New(60 * time.Second),
											PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
											CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										},
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_LINE,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "error rate (%)",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createMaxConcurrentRequestsWidget は最大同時リクエストウィジェットを作成する
func (s *Service) createMaxConcurrentRequestsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Max Concurrent Requests (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/container/max_request_concurrencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P50",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/container/max_request_concurrencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P95",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/container/max_request_concurrencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P99",
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "concurrent requests",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createResponseTimeWidget はレスポンス時間 (P50/P95/P99) ウィジェットを作成する
func (s *Service) createResponseTimeWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Response Time (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P50",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P95",
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: fmt.Sprintf(`resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="run.googleapis.com/request_latencies"`, s.serviceName),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: "P99",
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: "response time (ms)",
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// internal/report/html.go
package report

import (
	"bytes"
	"html/template"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/code_analyzer/models"
)

// HTMLReporter はHTML形式のレポート生成機能を提供します
type HTMLReporter struct{}

// NewHTMLReporter は新しいHTMLReporterインスタンスを作成します
func NewHTMLReporter() *HTMLReporter {
	return &HTMLReporter{}
}

// Generate はHTML形式のレポートを生成します
func (r *HTMLReporter) Generate(metrics models.ProjectMetrics, history []models.HistoricalData) (string, error) {
	// レポートテンプレート
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Code Metrics Report</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .summary {
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            margin-bottom: 30px;
        }
        .metric-card {
            background-color: #f5f5f5;
            border-radius: 5px;
            padding: 15px;
            margin-bottom: 15px;
            width: calc(25% - 20px);
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .metric-value {
            font-size: 24px;
            font-weight: bold;
            margin: 10px 0;
        }
        .trend-up {
            color: #e74c3c;
        }
        .trend-down {
            color: #2ecc71;
        }
        .trend-neutral {
            color: #3498db;
        }
        .chart-container {
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            margin-bottom: 30px;
        }
        .chart {
            width: calc(50% - 20px);
            margin-bottom: 30px;
            background-color: #fff;
            border-radius: 5px;
            padding: 15px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
        }
        th, td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .high-complexity {
            background-color: #ffdddd;
        }
        .medium-complexity {
            background-color: #ffffdd;
        }
        .clone-section {
            margin-top: 40px;
        }
        .clone-item {
            background-color: #f9f9f9;
            border-left: 4px solid #3498db;
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 0 5px 5px 0;
        }
        .code-preview {
            font-family: monospace;
            background-color: #272822;
            color: #f8f8f2;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Code Metrics Report</h1>
        <p>Project: {{.ProjectPath}}</p>
        <p>Analyzed at: {{formatTime .AnalyzedAt}}</p>
    </div>

    <div class="summary">
        <div class="metric-card">
            <h3>Total Lines</h3>
            <div class="metric-value">{{.TotalLines}}</div>
            {{if .Trends.total_lines}}
                {{if gt .Trends.total_lines.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.total_lines.Change}} ({{printf "%.1f" .Trends.total_lines.ChangeRate}}%)</div>
                {{else if lt .Trends.total_lines.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.total_lines.Change -1)}} ({{printf "%.1f" .Trends.total_lines.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Code Lines</h3>
            <div class="metric-value">{{.TotalCodeLines}}</div>
            {{if .Trends.code_lines}}
                {{if gt .Trends.code_lines.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.code_lines.Change}} ({{printf "%.1f" .Trends.code_lines.ChangeRate}}%)</div>
                {{else if lt .Trends.code_lines.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.code_lines.Change -1)}} ({{printf "%.1f" .Trends.code_lines.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Comment Ratio</h3>
            <div class="metric-value">{{printf "%.1f" .CommentRatio}}%</div>
            {{if .Trends.comment_ratio}}
                {{if gt .Trends.comment_ratio.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.comment_ratio.Change}} ({{printf "%.1f" .Trends.comment_ratio.ChangeRate}}%)</div>
                {{else if lt .Trends.comment_ratio.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.comment_ratio.Change -1)}} ({{printf "%.1f" .Trends.comment_ratio.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Average Complexity</h3>
            <div class="metric-value">{{printf "%.1f" .AvgComplexity}}</div>
            {{if .Trends.avg_complexity}}
                {{if gt .Trends.avg_complexity.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.avg_complexity.Change}} ({{printf "%.1f" .Trends.avg_complexity.ChangeRate}}%)</div>
                {{else if lt .Trends.avg_complexity.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.avg_complexity.Change -1)}} ({{printf "%.1f" .Trends.avg_complexity.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Files</h3>
            <div class="metric-value">{{.FileCount}}</div>
        </div>
        <div class="metric-card">
            <h3>Max Complexity</h3>
            <div class="metric-value">{{.MaxComplexity}}</div>
            {{if .Trends.max_complexity}}
                {{if gt .Trends.max_complexity.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.max_complexity.Change}} ({{printf "%.1f" .Trends.max_complexity.ChangeRate}}%)</div>
                {{else if lt .Trends.max_complexity.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.max_complexity.Change -1)}} ({{printf "%.1f" .Trends.max_complexity.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Code Clones</h3>
            <div class="metric-value">{{len .Clones}}</div>
            {{if .Trends.clone_ratio}}
                {{if gt .Trends.clone_ratio.Change 0}}
                    <div class="trend-up">↑ {{printf "%.1f" .Trends.clone_ratio.Change}}% ({{printf "%.1f" .Trends.clone_ratio.ChangeRate}}%)</div>
                {{else if lt .Trends.clone_ratio.Change 0}}
                    <div class="trend-down">↓ {{printf "%.1f" (mul .Trends.clone_ratio.Change -1)}}% ({{printf "%.1f" .Trends.clone_ratio.ChangeRate}}%)</div>
                {{else}}
                    <div class="trend-neutral">↔ No change</div>
                {{end}}
            {{end}}
        </div>
        <div class="metric-card">
            <h3>Function Count</h3>
            <div class="metric-value" id="total-functions">0</div>
            <script>
                // 集計はJSで行う
                document.addEventListener('DOMContentLoaded', function() {
                    let totalFunctions = 0;
                    {{range .Files}}
                        totalFunctions += {{.FunctionCount}};
                    {{end}}
                    document.getElementById('total-functions').innerText = totalFunctions;
                });
            </script>
        </div>
    </div>

    <div class="chart-container">
        <div class="chart">
            <canvas id="lineDistributionChart"></canvas>
        </div>
        <div class="chart">
            <canvas id="complexityDistributionChart"></canvas>
        </div>
        <div class="chart">
            <canvas id="trendChart"></canvas>
        </div>
        <div class="chart">
            <canvas id="cloneChart"></canvas>
        </div>
    </div>

    <h2>Most Complex Files</h2>
    <table>
        <tr>
            <th>File</th>
            <th>Lines</th>
            <th>Code</th>
            <th>Comments</th>
            <th>Complexity</th>
            <th>Functions</th>
            <th>Comment Ratio</th>
        </tr>
        {{range .Files}}
            {{if gt .Complexity 10}}
                <tr class="high-complexity">
            {{else if gt .Complexity 5}}
                <tr class="medium-complexity">
            {{else}}
                <tr>
            {{end}}
                <td>{{.Path}}</td>
                <td>{{.TotalLines}}</td>
                <td>{{.CodeLines}}</td>
                <td>{{.CommentLines}}</td>
                <td>{{.Complexity}}</td>
                <td>{{.FunctionCount}}</td>
                <td>{{printf "%.1f" .CommentRatio}}%</td>
            </tr>
        {{end}}
    </table>

    {{if .Clones}}
    <div class="clone-section">
        <h2>Code Clones</h2>
        {{range .Clones}}
            <div class="clone-item">
                <p><strong>Source:</strong> {{.SourceFile}} (line {{.SourceLine}})</p>
                <p><strong>Target:</strong> {{.TargetFile}} (line {{.TargetLine}})</p>
                <p><strong>Size:</strong> {{.LineCount}} lines</p>
                <p><strong>Similarity:</strong> {{printf "%.1f" (mul .Similarity 100)}}%</p>
                <div class="code-preview">{{.Code}}</div>
            </div>
        {{end}}
    </div>
    {{end}}

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            // 行数分布グラフ
            const lineCtx = document.getElementById('lineDistributionChart').getContext('2d');
            const lineChart = new Chart(lineCtx, {
                type: 'doughnut',
                data: {
                    labels: ['Code', 'Comments', 'Blank'],
                    datasets: [{
                        data: [{{.TotalCodeLines}}, {{.TotalComments}}, {{.TotalBlankLines}}],
                        backgroundColor: ['#3498db', '#2ecc71', '#95a5a6']
                    }]
                },
                options: {
                    responsive: true,
                    plugins: {
                        title: {
                            display: true,
                            text: 'Line Distribution'
                        }
                    }
                }
            });

            // 複雑度分布グラフ
            const complexityData = [0, 0, 0, 0]; // 0-5, 6-10, 11-20, 21+
            {{range .Files}}
                {{if le .Complexity 5}}
                    complexityData[0]++;
                {{else if le .Complexity 10}}
                    complexityData[1]++;
                {{else if le .Complexity 20}}
                    complexityData[2]++;
                {{else}}
                    complexityData[3]++;
                {{end}}
            {{end}}

            const complexityCtx = document.getElementById('complexityDistributionChart').getContext('2d');
            const complexityChart = new Chart(complexityCtx, {
                type: 'bar',
                data: {
                    labels: ['Low (0-5)', 'Medium (6-10)', 'High (11-20)', 'Very High (21+)'],
                    datasets: [{
                        label: 'Number of Files',
                        data: complexityData,
                        backgroundColor: ['#2ecc71', '#f1c40f', '#e67e22', '#e74c3c']
                    }]
                },
                options: {
                    responsive: true,
                    plugins: {
                        title: {
                            display: true,
                            text: 'Complexity Distribution'
                        }
                    }
                }
            });

            // トレンドグラフ（過去の履歴からデータ抽出）
            const trendLabels = [];
            const codeData = [];
            const complexityData = [];
            const commentRatioData = [];

            {{range $index, $element := .HistoricalData}}
                trendLabels.push('{{formatTimeShort $element.Date}}');
                codeData.push({{$element.CodeLines}});
                complexityData.push({{$element.AvgComplexity}});
                commentRatioData.push({{$element.CommentRatio}});
            {{end}}

            // 現在のデータを追加
            trendLabels.push('Current');
            codeData.push({{.TotalCodeLines}});
            complexityData.push({{.AvgComplexity}});
            commentRatioData.push({{.CommentRatio}});

            const trendCtx = document.getElementById('trendChart').getContext('2d');
            const trendChart = new Chart(trendCtx, {
                type: 'line',
                data: {
                    labels: trendLabels,
                    datasets: [
                        {
                            label: 'Code Lines',
                            data: codeData,
                            borderColor: '#3498db',
                            tension: 0.1
                        },
                        {
                            label: 'Avg Complexity',
                            data: complexityData,
                            borderColor: '#e74c3c',
                            tension: 0.1
                        },
                        {
                            label: 'Comment Ratio',
                            data: commentRatioData,
                            borderColor: '#2ecc71',
                            tension: 0.1
                        }
                    ]
                },
                options: {
                    responsive: true,
                    plugins: {
                        title: {
                            display: true,
                            text: 'Metrics Trend'
                        }
                    }
                }
            });

            // クローングラフ（クローンの分布）
            const cloneByFile = {};
            {{range .Clones}}
                if (!cloneByFile['{{.SourceFile}}']) {
                    cloneByFile['{{.SourceFile}}'] = 1;
                } else {
                    cloneByFile['{{.SourceFile}}']++;
                }

                if (!cloneByFile['{{.TargetFile}}']) {
                    cloneByFile['{{.TargetFile}}'] = 1;
                } else {
                    cloneByFile['{{.TargetFile}}']++;
                }
            {{end}}

            const cloneLabels = [];
            const cloneByFile = {};
           {{range .Clones}}
               if (!cloneByFile['{{.SourceFile}}']) {
                   cloneByFile['{{.SourceFile}}'] = 1;
               } else {
                   cloneByFile['{{.SourceFile}}']++;
               }

               if (!cloneByFile['{{.TargetFile}}']) {
                   cloneByFile['{{.TargetFile}}'] = 1;
               } else {
                   cloneByFile['{{.TargetFile}}']++;
               }
           {{end}}

           const cloneLabels = [];
           const cloneCounts = [];

           const sortedFiles = Object.keys(cloneByFile).sort((a, b) => cloneByFile[b] - cloneByFile[a]);
           const topFiles = sortedFiles.slice(0, 10); // 上位10ファイル

           for (const file of topFiles) {
               // パスを短く表示
               const shortPath = file.split('/').pop();
               cloneLabels.push(shortPath);
               cloneCounts.push(cloneByFile[file]);
           }

           const cloneCtx = document.getElementById('cloneChart').getContext('2d');
           const cloneChart = new Chart(cloneCtx, {
               type: 'bar',
               data: {
                   labels: cloneLabels,
                   datasets: [{
                       label: 'Clone Count',
                       data: cloneCounts,
                       backgroundColor: '#9b59b6'
                   }]
               },
               options: {
                   indexAxis: 'y',
                   responsive: true,
                   plugins: {
                       title: {
                           display: true,
                           text: 'Files with Most Clones'
                       }
                   }
               }
           });
       });
   </script>
</body>
</html>
`

	// テンプレート関数
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatTimeShort": func(t time.Time) string {
			return t.Format("01/02")
		},
	}

	// データを構造化
	data := struct {
		models.ProjectMetrics
		HistoricalData []models.HistoricalData
	}{
		ProjectMetrics: metrics,
		HistoricalData: history,
	}

	// テンプレート処理
	t, err := template.New("report").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

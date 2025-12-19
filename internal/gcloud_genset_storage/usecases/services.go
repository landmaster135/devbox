package usecases

import (
	"fmt"
	"strings"
)

// Service generates gsutil/gcloud commands for Cloud Storage operations.
type Service struct{}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{}
}

const (
	discordWebhookEnvVarName = "DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD"
	discordCLIPath           = "$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook"
	successEmbedType         = "google-cloud-storage-success"
	failureEmbedType         = "google-cloud-storage-failed"
)

// UploadFilesParams holds parameters for upload operation.
type UploadFilesParams struct {
	LocalPath string
	BucketURL string
}

// DownloadFilesParams holds parameters for download operation.
type DownloadFilesParams struct {
	Sources     []string
	Destination string
}

// CreateBucketParams holds parameters for create bucket.
type CreateBucketParams struct {
	BucketURL    string
	StorageClass string
	Location     string
}

// ListContentsParams holds parameters for listing contents.
type ListContentsParams struct {
	Target string
}

// ShowDetailsParams holds parameters for detailed listing.
type ShowDetailsParams struct {
	Target string
}

// DeleteObjectParams holds parameters for deleting an object.
type DeleteObjectParams struct {
	Target string
}

// ACLParams holds parameters for ACL operations.
type ACLParams struct {
	Target  string
	ACLFile string
}

// UploadFiles command.
func (s *Service) BuildUploadFilesCommand(params UploadFilesParams) (string, error) {
	if params.LocalPath == "" || params.BucketURL == "" {
		return "", fmt.Errorf("localPath と bucketURL は必須です")
	}
	return fmt.Sprintf("gsutil -m cp -r %s %s", shellQuote(params.LocalPath), shellQuote(params.BucketURL)), nil
}

// DownloadFiles command.
func (s *Service) BuildDownloadFilesCommand(params DownloadFilesParams) (string, error) {
	if len(params.Sources) == 0 {
		return "", fmt.Errorf("sources は必須です")
	}
	if params.Destination == "" {
		return "", fmt.Errorf("destination は必須です")
	}
	quotedSources := make([]string, 0, len(params.Sources))
	for _, src := range params.Sources {
		trimmed := strings.TrimSpace(src)
		if trimmed != "" {
			quotedSources = append(quotedSources, shellQuote(trimmed))
		}
	}
	if len(quotedSources) == 0 {
		return "", fmt.Errorf("sources の形式が不正です")
	}
	command := fmt.Sprintf("gsutil -m cp %s %s", strings.Join(quotedSources, " "), shellQuote(params.Destination))
	return command, nil
}

// CreateBucket command.
func (s *Service) BuildCreateBucketCommand(params CreateBucketParams) (string, error) {
	if params.BucketURL == "" || params.StorageClass == "" || params.Location == "" {
		return "", fmt.Errorf("bucketURL, storageClass, location は必須です")
	}
	return fmt.Sprintf("gsutil mb -c %s -l %s %s", shellQuote(params.StorageClass), shellQuote(params.Location), shellQuote(params.BucketURL)), nil
}

// ListContents command (gsutil ls or ls).
func (s *Service) BuildListContentsCommand(params ListContentsParams) (string, error) {
	target := strings.TrimSpace(params.Target)
	if target == "" {
		return "", fmt.Errorf("target は必須です")
	}
	if isGCSPath(target) {
		return fmt.Sprintf("gsutil ls %s", shellQuote(target)), nil
	}
	return fmt.Sprintf("ls %s", shellQuote(target)), nil
}

// ShowDetails command (gsutil ls -Lb or -L).
func (s *Service) BuildShowDetailsCommand(params ShowDetailsParams) (string, error) {
	target := strings.TrimSpace(params.Target)
	if !isGCSPath(target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	if isBucketPath(target) {
		return fmt.Sprintf("gsutil ls -Lb %s", shellQuote(target)), nil
	}
	return fmt.Sprintf("gsutil ls -L %s", shellQuote(target)), nil
}

// BuildListNamesCommand lists physical names (gsutil ls).
func (s *Service) BuildListNamesCommand(params ListContentsParams) (string, error) {
	target := strings.TrimSpace(params.Target)
	if !isGCSPath(target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil ls %s", shellQuote(target)), nil
}

// DeleteObject command.
func (s *Service) BuildDeleteObjectCommand(params DeleteObjectParams) (string, error) {
	if !isGCSPath(params.Target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil rm %s", shellQuote(params.Target)), nil
}

// GetACL command.
func (s *Service) BuildGetACLCommand(params ACLParams) (string, error) {
	if !isGCSPath(params.Target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil acl get %s", shellQuote(params.Target)), nil
}

// SetACL command.
func (s *Service) BuildSetACLCommand(params ACLParams) (string, error) {
	if params.ACLFile == "" {
		return "", fmt.Errorf("aclFile は必須です")
	}
	if !isGCSPath(params.Target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil acl set %s %s", shellQuote(params.ACLFile), shellQuote(params.Target)), nil
}

// GrantReadAll command.
func (s *Service) BuildGrantReadAllCommand(params ACLParams) (string, error) {
	if !isGCSPath(params.Target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil acl ch -u AllUsers:R %s", shellQuote(params.Target)), nil
}

// RemoveReadAll command.
func (s *Service) BuildRemoveReadAllCommand(params ACLParams) (string, error) {
	if !isGCSPath(params.Target) {
		return "", fmt.Errorf("target には gs:// で始まるパスを指定してください")
	}
	return fmt.Sprintf("gsutil acl ch -d AllUsers %s", shellQuote(params.Target)), nil
}

// PrintHighlightedCommand outputs command nicely.
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成されたコマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

// BuildNotificationWrappedCommand creates a shell snippet with discord notifications and the gcloud command.
func (s *Service) BuildNotificationWrappedCommand(operation string, gcloudCommand string) (string, bool) {
	template, ok := notificationTemplates[operation]
	if !ok {
		return "", false
	}

	var lines []string
	if template.startContent != "" {
		lines = append(lines, buildSimpleNotificationCommand(template.startContent))
	}

	successCommand := ""
	if template.successContent != "" {
		successCommand = buildEmbedNotificationCommand(template.successContent, template.successEmbedText, successEmbedType)
	}
	failureCommand := ""
	if template.failureContent != "" {
		failureCommand = buildEmbedNotificationCommand(template.failureContent, template.failureEmbedText, failureEmbedType)
	}

	lines = append(lines, fmt.Sprintf("if %s; then", gcloudCommand))
	if successCommand != "" {
		lines = append(lines, indentCommand(successCommand, "  "))
	}
	lines = append(lines, "else")
	if failureCommand != "" {
		lines = append(lines, indentCommand(failureCommand, "  "))
	}
	lines = append(lines, "fi")

	return strings.Join(lines, "\n"), true
}

// PrintNotificationScript prints the wrapped notification script.
func (s *Service) PrintNotificationScript(script string) {
	if strings.TrimSpace(script) == "" {
		return
	}

	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("通知付きシェルコマンド")
	fmt.Println("==============================")
	fmt.Println(script)
	fmt.Println("==============================")
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

func buildSimpleNotificationCommand(content string) string {
	return buildDiscordWebhookCommand(content, "none", "")
}

func buildEmbedNotificationCommand(content, embedText, embedType string) string {
	return buildDiscordWebhookCommand(content, embedType, embedText)
}

func buildDiscordWebhookCommand(content, embedType, embedText string) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("%s \\", discordCLIPath))
	lines = append(lines, fmt.Sprintf("  -webhook-url \"$%s\" \\", discordWebhookEnvVarName))
	lines = append(lines, fmt.Sprintf("  -content-text %s \\", shellQuote(content)))
	embedLine := fmt.Sprintf("  -embed-type %s", shellQuote(embedType))
	if embedText != "" {
		lines = append(lines, embedLine+" \\")
		lines = append(lines, fmt.Sprintf("  -embed-text %s", shellQuote(embedText)))
	} else {
		lines = append(lines, embedLine)
	}
	return strings.Join(lines, "\n")
}

func indentCommand(command, indent string) string {
	if command == "" {
		return ""
	}
	parts := strings.Split(command, "\n")
	for i, part := range parts {
		parts[i] = indent + part
	}
	return strings.Join(parts, "\n")
}

func isGCSPath(target string) bool {
	return strings.HasPrefix(target, "gs://")
}

func isBucketPath(target string) bool {
	if !isGCSPath(target) {
		return false
	}
	trimmed := strings.TrimPrefix(target, "gs://")
	trimmed = strings.TrimSuffix(trimmed, "/")
	return trimmed != "" && !strings.Contains(trimmed, "/")
}

type notificationTemplate struct {
	startContent     string
	successContent   string
	successEmbedText string
	failureContent   string
	failureEmbedText string
}

var notificationTemplates = map[string]notificationTemplate{
	"upload-files": {
		startContent:     "ファイルをアップロードするよ！",
		successContent:   "アップしたよ！",
		successEmbedText: "ファイルをアップロードしたよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ファイルのアップロードに失敗したよ…",
	},
	"download-files": {
		startContent:     "ファイルをダウンロードするよ！",
		successContent:   "ダウンロードしたよ！",
		successEmbedText: "ファイルをダウンロードしたよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ファイルのダウンロードに失敗したよ…",
	},
	"create-bucket": {
		startContent:     "バケットを作るよ！",
		successContent:   "作ったよ！",
		successEmbedText: "バケットを作ったよ！",
		failureContent:   "失敗…",
		failureEmbedText: "バケットを作れなかったよ…",
	},
	"list-contents": {
		startContent:     "バケットのリストかファイルの詳細を並べるよ！",
		successContent:   "並べたよ！",
		successEmbedText: "バケットとファイルのリストを並べたよ！",
		failureContent:   "失敗…",
		failureEmbedText: "バケットとファイルを取れなかったよ…",
	},
	"show-details": {
		startContent:     "バケットかオブジェクトの詳細を並べるよ！",
		successContent:   "表示したよ！",
		successEmbedText: "バケットかオブジェクトの詳細を表示したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "バケットとオブジェクトを取れなかったよ…",
	},
	"list-names": {
		startContent:     "バケットかオブジェクトの詳細を並べるよ！",
		successContent:   "表示したよ！",
		successEmbedText: "バケットかオブジェクトの詳細を表示したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "バケットとオブジェクトを取れなかったよ…",
	},
	"delete-object": {
		startContent:     "フォルダもしくはオブジェクトを削除するよ！",
		successContent:   "削除したよ！",
		successEmbedText: "フォルダもしくはオブジェクトを削除したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "フォルダもしくはオブジェクトを削除できなかったよ…",
	},
	"get-acl": {
		startContent:     "ACLを確認するよ！",
		successContent:   "権限を取得したよ！",
		successEmbedText: "ACLを確認したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ACLを確認できなかったよ…",
	},
	"set-acl": {
		startContent:     "ACLを設定するよ！",
		successContent:   "権限を設定したよ！",
		successEmbedText: "ACLを設定したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ACLを設定できなかったよ…",
	},
	"grant-read-all": {
		startContent:     "ACLを全てのユーザに付与するよ！",
		successContent:   "権限を付与したよ！",
		successEmbedText: "ACL(READ)を全てのユーザに付与したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ACLを全てのユーザに付与できなかったよ…",
	},
	"remove-read-all": {
		startContent:     "ACLを全てのユーザから剥奪するよ！",
		successContent:   "権限を付与したよ！",
		successEmbedText: "ACLを全てのユーザから剥奪したよ！",
		failureContent:   "失敗…",
		failureEmbedText: "ACLを全てのユーザから剥奪できなかったよ…",
	},
}

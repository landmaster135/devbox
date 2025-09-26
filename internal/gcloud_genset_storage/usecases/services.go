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

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
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

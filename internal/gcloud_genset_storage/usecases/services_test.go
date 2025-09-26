package usecases

import (
	"strings"
	"testing"
)

func TestBuildUploadFilesCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildUploadFilesCommand(UploadFilesParams{
		LocalPath: "/tmp/data",
		BucketURL: "gs://bucket/path/",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil -m cp -r '/tmp/data' 'gs://bucket/path/'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildDownloadFilesCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDownloadFilesCommand(DownloadFilesParams{
		Sources:     []string{"gs://bucket/file1", "gs://bucket/file2"},
		Destination: "./dest",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil -m cp 'gs://bucket/file1' 'gs://bucket/file2' './dest'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildCreateBucketCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildCreateBucketCommand(CreateBucketParams{
		BucketURL:    "gs://new-bucket",
		StorageClass: "STANDARD",
		Location:     "US",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil mb -c 'STANDARD' -l 'US' 'gs://new-bucket'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildShowDetailsCommandBucket(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildShowDetailsCommand(ShowDetailsParams{Target: "gs://bucket/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil ls -Lb 'gs://bucket/'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildShowDetailsCommandObject(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildShowDetailsCommand(ShowDetailsParams{Target: "gs://bucket/path/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil ls -L 'gs://bucket/path/file'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildDeleteObjectCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDeleteObjectCommand(DeleteObjectParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil rm 'gs://bucket/file'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildSetACLCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildSetACLCommand(ACLParams{ACLFile: "acl.json", Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil acl set 'acl.json' 'gs://bucket/file'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildGrantReadAllCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildGrantReadAllCommand(ACLParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil acl ch -u AllUsers:R 'gs://bucket/file'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildRemoveReadAllCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildRemoveReadAllCommand(ACLParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gsutil acl ch -d AllUsers 'gs://bucket/file'"
	if cmd != expected {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildNotificationWrappedCommand(t *testing.T) {
	service := NewService()
	command := "gsutil -m cp -r './data' 'gs://bucket/'"
	script, ok := service.BuildNotificationWrappedCommand("upload-files", command)
	if !ok {
		t.Fatalf("expected notification script")
	}

	expected := strings.Join([]string{
		"$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \\",
		"  -webhook-url \"$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD\" \\",
		"  -content-text 'ファイルをアップロードするよ！' \\",
		"  -embed-type 'none'",
		"if gsutil -m cp -r './data' 'gs://bucket/'; then",
		"  $HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \\",
		"    -webhook-url \"$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD\" \\",
		"    -content-text 'アップしたよ！' \\",
		"    -embed-type 'google-cloud-storage-success' \\",
		"    -embed-text 'ファイルをアップロードしたよ！'",
		"else",
		"  $HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \\",
		"    -webhook-url \"$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD\" \\",
		"    -content-text '失敗…' \\",
		"    -embed-type 'google-cloud-storage-failed' \\",
		"    -embed-text 'ファイルのアップロードに失敗したよ…'",
		"fi",
	}, "\n")

	if script != expected {
		t.Fatalf("unexpected notification script:\n%s", script)
	}
}

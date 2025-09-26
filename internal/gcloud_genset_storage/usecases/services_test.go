package usecases

import (
	"bytes"
	"io"
	"os"
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

func TestBuildUploadFilesCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildUploadFilesCommand(UploadFilesParams{BucketURL: "gs://bucket"}); err == nil {
		t.Fatal("expected error when local path missing")
	}
	if _, err := service.BuildUploadFilesCommand(UploadFilesParams{LocalPath: "./data"}); err == nil {
		t.Fatal("expected error when bucket url missing")
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

func TestBuildDownloadFilesCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildDownloadFilesCommand(DownloadFilesParams{Destination: "./out"}); err == nil {
		t.Fatal("expected error when sources missing")
	}
	if _, err := service.BuildDownloadFilesCommand(DownloadFilesParams{Sources: []string{""}, Destination: "./out"}); err == nil {
		t.Fatal("expected error when sources invalid")
	}
	if _, err := service.BuildDownloadFilesCommand(DownloadFilesParams{Sources: []string{"gs://a"}}); err == nil {
		t.Fatal("expected error when destination missing")
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

func TestBuildCreateBucketCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildCreateBucketCommand(CreateBucketParams{StorageClass: "STANDARD", Location: "US"}); err == nil {
		t.Fatal("expected error when bucket missing")
	}
	if _, err := service.BuildCreateBucketCommand(CreateBucketParams{BucketURL: "gs://bucket", Location: "US"}); err == nil {
		t.Fatal("expected error when storage class missing")
	}
	if _, err := service.BuildCreateBucketCommand(CreateBucketParams{BucketURL: "gs://bucket", StorageClass: "STANDARD"}); err == nil {
		t.Fatal("expected error when location missing")
	}
}

func TestBuildListContentsCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildListContentsCommand(ListContentsParams{Target: "gs://bucket/path"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "gsutil ls 'gs://bucket/path'" {
		t.Fatalf("unexpected command: %s", cmd)
	}

	cmd, err = service.BuildListContentsCommand(ListContentsParams{Target: "./local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "ls './local'" {
		t.Fatalf("unexpected command for local: %s", cmd)
	}
}

func TestBuildShowDetailsCommand(t *testing.T) {
	service := NewService()
	bucketCmd, err := service.BuildShowDetailsCommand(ShowDetailsParams{Target: "gs://bucket/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bucketCmd != "gsutil ls -Lb 'gs://bucket/'" {
		t.Fatalf("unexpected bucket command: %s", bucketCmd)
	}

	objCmd, err := service.BuildShowDetailsCommand(ShowDetailsParams{Target: "gs://bucket/obj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if objCmd != "gsutil ls -L 'gs://bucket/obj'" {
		t.Fatalf("unexpected object command: %s", objCmd)
	}
}

func TestBuildShowDetailsCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildShowDetailsCommand(ShowDetailsParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error")
	}
}

func TestBuildListNamesCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildListNamesCommand(ListContentsParams{Target: "gs://bucket/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "gsutil ls 'gs://bucket/'" {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildListNamesCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildListNamesCommand(ListContentsParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error")
	}
}

func TestBuildDeleteObjectCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDeleteObjectCommand(DeleteObjectParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "gsutil rm 'gs://bucket/file'" {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestBuildDeleteObjectCommand_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildDeleteObjectCommand(DeleteObjectParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error")
	}
}

func TestACLCommands(t *testing.T) {
	service := NewService()

	getCmd, err := service.BuildGetACLCommand(ACLParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if getCmd != "gsutil acl get 'gs://bucket/file'" {
		t.Fatalf("unexpected get acl command: %s", getCmd)
	}

	setCmd, err := service.BuildSetACLCommand(ACLParams{ACLFile: "acl.json", Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if setCmd != "gsutil acl set 'acl.json' 'gs://bucket/file'" {
		t.Fatalf("unexpected set acl command: %s", setCmd)
	}

	grantCmd, err := service.BuildGrantReadAllCommand(ACLParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grantCmd != "gsutil acl ch -u AllUsers:R 'gs://bucket/file'" {
		t.Fatalf("unexpected grant command: %s", grantCmd)
	}

	removeCmd, err := service.BuildRemoveReadAllCommand(ACLParams{Target: "gs://bucket/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removeCmd != "gsutil acl ch -d AllUsers 'gs://bucket/file'" {
		t.Fatalf("unexpected remove command: %s", removeCmd)
	}
}

func TestACLCommands_Errors(t *testing.T) {
	service := NewService()
	if _, err := service.BuildGetACLCommand(ACLParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error for get acl")
	}
	if _, err := service.BuildSetACLCommand(ACLParams{Target: "gs://bucket"}); err == nil {
		t.Fatal("expected acl-file error")
	}
	if _, err := service.BuildSetACLCommand(ACLParams{ACLFile: "acl.json", Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error for set acl")
	}
	if _, err := service.BuildGrantReadAllCommand(ACLParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error grant")
	}
	if _, err := service.BuildRemoveReadAllCommand(ACLParams{Target: "./local"}); err == nil {
		t.Fatal("expected gs:// error remove")
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

func TestBuildNotificationWrappedCommand_Unknown(t *testing.T) {
	service := NewService()
	if _, ok := service.BuildNotificationWrappedCommand("unknown", "echo noop"); ok {
		t.Fatal("expected false for unknown operation")
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()
	output := captureStdout(func() {
		service.PrintHighlightedCommand("gsutil ls")
	})

	if !strings.Contains(output, "生成されたコマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gsutil ls") {
		t.Fatalf("expected command in output: %s", output)
	}
}

func TestPrintNotificationScript(t *testing.T) {
	service := NewService()
	output := captureStdout(func() {
		service.PrintNotificationScript("echo hi")
	})
	if !strings.Contains(output, "通知付きシェルコマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "echo hi") {
		t.Fatalf("expected script content in output: %s", output)
	}

	if captureStdout(func() { service.PrintNotificationScript("   ") }) != "" {
		t.Fatal("expected no output for blank script")
	}
}

func TestIndentCommand(t *testing.T) {
	indented := indentCommand("line1\nline2", "  ")
	if indented != "  line1\n  line2" {
		t.Fatalf("unexpected indent result: %s", indented)
	}
	if indentCommand("", "  ") != "" {
		t.Fatal("expected empty string for empty command")
	}
}

func TestIsBucketPath(t *testing.T) {
	if !isBucketPath("gs://bucket/") {
		t.Fatal("expected bucket path to be true")
	}
	if isBucketPath("gs://bucket/object") {
		t.Fatal("expected object path to be false")
	}
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}

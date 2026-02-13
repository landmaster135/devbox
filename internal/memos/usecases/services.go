package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
)

const defaultTimeout = 30 * time.Second

// HTTPClient は http.Client と互換なインターフェース。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL    string
	APIToken   string
	Timeout    time.Duration
	HTTPClient HTTPClient
	FileSystem infrastructures.FileSystem
}

// Service は Memos API 呼び出しのユースケースを提供する。
type Service struct {
	baseURL    string
	apiToken   string
	client     HTTPClient
	fileSystem infrastructures.FileSystem
}

// Memo は CLI/上位層に返すメモ情報。
type Memo struct {
	Name        string `json:"name,omitempty"`
	UID         string `json:"uid,omitempty"`
	ID          int64  `json:"id,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
	DisplayTime string `json:"displayTime,omitempty"`
	Content     string `json:"content,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	State       string `json:"state,omitempty"`
	Pinned      bool   `json:"pinned,omitempty"`
}

// ListMemosOutput は ListMemos のレスポンス。
type ListMemosOutput struct {
	Memos         []Memo `json:"memos,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	TotalSize     int64  `json:"totalSize,omitempty"`
}

// Attachment は Memos API の添付情報。
type Attachment struct {
	Name         string `json:"name,omitempty"`
	Filename     string `json:"filename,omitempty"`
	ExternalLink string `json:"externalLink,omitempty"`
	Type         string `json:"type,omitempty"`
	Memo         string `json:"memo,omitempty"`
}

// ListMemoAttachmentsOutput は ListMemoAttachments のレスポンス。
type ListMemoAttachmentsOutput struct {
	Attachments   []Attachment `json:"attachments,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}

// SetMemoAttachmentsOutput は SetMemoAttachments のレスポンス。
type SetMemoAttachmentsOutput struct {
	Name        string       `json:"name,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type memoMutationRequest struct {
	Content     string `json:"content,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	State       string `json:"state,omitempty"`
	Pinned      *bool  `json:"pinned,omitempty"`
	DisplayTime string `json:"displayTime,omitempty"`
}

type createAttachmentRequest struct {
	Filename     string `json:"filename,omitempty"`
	Content      []byte `json:"content,omitempty"`
	ExternalLink string `json:"externalLink,omitempty"`
	Type         string `json:"type,omitempty"`
	Memo         string `json:"memo,omitempty"`
}

type setMemoAttachmentsRequest struct {
	Name        string       `json:"name,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	client := opts.HTTPClient
	if client == nil {
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		client = &http.Client{Timeout: timeout}
	}

	fileSystem := opts.FileSystem
	if fileSystem == nil {
		fileSystem = infrastructures.NewOSFileSystem()
	}

	return &Service{
		baseURL:    normalizeBaseURL(opts.BaseURL),
		apiToken:   strings.TrimSpace(opts.APIToken),
		client:     client,
		fileSystem: fileSystem,
	}
}

// CreateMemo は CreateMemo API を呼び出す。
func (s *Service) CreateMemo(
	ctx context.Context,
	memoID string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	displayTime string,
) (*Memo, error) {
	resolvedContent, err := s.resolveContent(content, contentFile)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if memoID != "" {
		query.Set("memoId", memoID)
	}

	payload := memoMutationRequest{
		Content:     resolvedContent,
		Visibility:  visibility,
		State:       state,
		Pinned:      pinned,
		DisplayTime: displayTime,
	}

	var result Memo
	if err := s.doJSON(ctx, http.MethodPost, "/memos", query, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMemo は GetMemo API を呼び出す。
func (s *Service) GetMemo(ctx context.Context, memo string) (*Memo, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	var result Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.doJSON(ctx, http.MethodGet, requestPath, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMemos は ListMemos API を呼び出す。
func (s *Service) ListMemos(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
) (*ListMemosOutput, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}
	if state != "" {
		query.Set("state", state)
	}
	if orderBy != "" {
		query.Set("orderBy", orderBy)
	}

	var result ListMemosOutput
	if err := s.doJSON(ctx, http.MethodGet, "/memos", query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateMemo は UpdateMemo API を呼び出す。
func (s *Service) UpdateMemo(
	ctx context.Context,
	memo string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	updateMask []string,
) (*Memo, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	resolvedContent, err := s.resolveContent(content, contentFile)
	if err != nil {
		return nil, err
	}

	finalMask := buildUpdateMask(resolvedContent, visibility, state, pinned, updateMask)
	if len(finalMask) == 0 {
		return nil, fmt.Errorf("updateMask が空です")
	}

	query := url.Values{}
	query.Set("updateMask", strings.Join(finalMask, ","))

	payload := memoMutationRequest{
		Content:    resolvedContent,
		Visibility: visibility,
		State:      state,
		Pinned:     pinned,
	}

	var result Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.doJSON(ctx, http.MethodPatch, requestPath, query, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PatchFiles はファイルを添付として作成し、メモに紐づける。
func (s *Service) PatchFiles(
	ctx context.Context,
	memo string,
	filePaths []string,
	replaces bool,
) (*SetMemoAttachmentsOutput, error) {
	created := make([]Attachment, 0, len(filePaths))
	for _, filePath := range filePaths {
		attachmentFile, err := s.fileSystem.ReadAttachmentFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("files の読み込みに失敗しました (%s): %w", filePath, err)
		}
		attachment, err := s.CreateAttachment(
			ctx,
			attachmentFile.Filename,
			attachmentFile.Content,
			attachmentFile.ContentType,
			memo,
		)
		if err != nil {
			return nil, err
		}
		if attachment == nil {
			return nil, fmt.Errorf("CreateAttachment の結果が空です")
		}
		created = append(created, *attachment)
	}

	finalAttachments := created
	if !replaces {
		existing, err := s.listAllMemoAttachments(ctx, memo)
		if err != nil {
			return nil, err
		}
		finalAttachments = mergeAttachmentsByName(existing, created)
	}

	return s.SetMemoAttachments(ctx, memo, finalAttachments)
}

// CreateAttachment は CreateAttachment API を呼び出す。
func (s *Service) CreateAttachment(
	ctx context.Context,
	filename string,
	content []byte,
	attachmentType string,
	memo string,
) (*Attachment, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("filename が空です")
	}

	attachmentType = strings.TrimSpace(attachmentType)
	if attachmentType == "" {
		return nil, fmt.Errorf("type が空です")
	}

	payload := createAttachmentRequest{
		Filename: filename,
		Content:  content,
		Type:     attachmentType,
	}
	if memoName := buildMemoResourceName(memo); memoName != "" {
		payload.Memo = memoName
	}

	var result Attachment
	if err := s.doJSON(ctx, http.MethodPost, "/attachments", nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMemoAttachments は ListMemoAttachments API を呼び出す。
func (s *Service) ListMemoAttachments(
	ctx context.Context,
	memo string,
	pageSize int,
	pageToken string,
) (*ListMemoAttachmentsOutput, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "attachments")
	var result ListMemoAttachmentsOutput
	if err := s.doJSON(ctx, http.MethodGet, requestPath, query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMemoAttachments は SetMemoAttachments API を呼び出す。
func (s *Service) SetMemoAttachments(
	ctx context.Context,
	memo string,
	attachments []Attachment,
) (*SetMemoAttachmentsOutput, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "attachments")
	payload := setMemoAttachmentsRequest{
		Name:        buildMemoResourceName(memoID),
		Attachments: attachments,
	}

	var result SetMemoAttachmentsOutput
	if err := s.doJSON(ctx, http.MethodPatch, requestPath, nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) listAllMemoAttachments(ctx context.Context, memo string) ([]Attachment, error) {
	const pageSize = 100

	var all []Attachment
	pageToken := ""

	for {
		result, err := s.ListMemoAttachments(ctx, memo, pageSize, pageToken)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return all, nil
		}
		all = append(all, result.Attachments...)
		if strings.TrimSpace(result.NextPageToken) == "" {
			return all, nil
		}
		pageToken = result.NextPageToken
	}
}

func (s *Service) doJSON(ctx context.Context, method, requestPath string, query url.Values, payload any, out any) error {
	req, err := s.newRequest(ctx, method, requestPath, query, payload)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Memos API %s %s の呼び出しに失敗しました: %w", method, requestPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Memos API %s %s が失敗しました: status=%d body=%s", method, requestPath, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("Memos API %s %s のレスポンスデコードに失敗しました: %w", method, requestPath, err)
	}
	return nil
}

func (s *Service) newRequest(ctx context.Context, method, requestPath string, query url.Values, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("リクエストボディのエンコードに失敗しました: %w", err)
		}
		body = buf
	}

	requestURL := s.baseURL + requestPath
	if len(query) > 0 {
		requestURL = requestURL + "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if s.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiToken)
	}
	return req, nil
}

func normalizeBaseURL(raw string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(raw), "/")
	if baseURL == "" {
		return "/api/v1"
	}
	if strings.HasSuffix(baseURL, "/api/v1") {
		return baseURL
	}
	return baseURL + "/api/v1"
}

func normalizeMemoIdentifier(memo string) string {
	trimmed := strings.TrimSpace(memo)
	if trimmed == "" {
		return ""
	}

	if strings.Contains(trimmed, "://") {
		if parsed, err := url.Parse(trimmed); err == nil {
			trimmed = strings.Trim(parsed.Path, "/")
		}
	}

	trimmed = strings.Trim(trimmed, "/")
	trimmed = strings.TrimPrefix(trimmed, "api/v1/")
	trimmed = strings.TrimPrefix(trimmed, "memos/")
	return strings.TrimSpace(trimmed)
}

func buildMemoResourceName(memo string) string {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return ""
	}
	return "memos/" + memoID
}

func buildUpdateMask(content string, visibility string, state string, pinned *bool, updateMask []string) []string {
	if len(updateMask) > 0 {
		return cleanMaskFields(updateMask)
	}

	mask := make([]string, 0, 4)
	if content != "" {
		mask = append(mask, "content")
	}
	if visibility != "" {
		mask = append(mask, "visibility")
	}
	if state != "" {
		mask = append(mask, "state")
	}
	if pinned != nil {
		mask = append(mask, "pinned")
	}
	return mask
}

func cleanMaskFields(raw []string) []string {
	seen := make(map[string]struct{})
	fields := make([]string, 0, len(raw))
	for _, item := range raw {
		for _, token := range strings.Split(item, ",") {
			field := strings.TrimSpace(token)
			if field == "" {
				continue
			}
			if _, ok := seen[field]; ok {
				continue
			}
			seen[field] = struct{}{}
			fields = append(fields, field)
		}
	}
	return fields
}

func (s *Service) resolveContent(content string, contentFile string) (string, error) {
	hasContent := strings.TrimSpace(content) != ""
	hasContentFile := strings.TrimSpace(contentFile) != ""

	if hasContent && hasContentFile {
		return "", fmt.Errorf("content と content-file は同時に指定できません")
	}
	if !hasContent && !hasContentFile {
		return "", fmt.Errorf("content または content-file の指定が必要です")
	}

	if hasContent {
		return content, nil
	}

	data, err := s.fileSystem.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("content-file の読み込みに失敗しました: %w", err)
	}
	fileContent := string(data)
	if strings.TrimSpace(fileContent) == "" {
		return "", fmt.Errorf("content-file が空です: %s", contentFile)
	}
	return fileContent, nil
}

func mergeAttachmentsByName(existing []Attachment, created []Attachment) []Attachment {
	createdByName := make(map[string]struct{}, len(created))
	for _, attachment := range created {
		name := strings.TrimSpace(attachment.Name)
		if name == "" {
			continue
		}
		createdByName[name] = struct{}{}
	}

	merged := make([]Attachment, 0, len(existing)+len(created))
	for _, attachment := range existing {
		name := strings.TrimSpace(attachment.Name)
		if name != "" {
			if _, exists := createdByName[name]; exists {
				continue
			}
		}
		merged = append(merged, attachment)
	}
	merged = append(merged, created...)
	return merged
}

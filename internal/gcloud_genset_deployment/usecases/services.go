package usecases

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Operation はサポートされる処理を表す。
type Operation string

const (
	// OperationListDeployments は Deployment Manager のデプロイメント一覧取得。
	OperationListDeployments Operation = "list-deployments"
)

// Service は gcloud コマンドを生成する。
type Service struct{}

// NewService は Service を生成する。
func NewService() *Service {
	return &Service{}
}

// BuildRequest はコマンド生成に必要な情報をまとめたもの。
type BuildRequest struct {
	Operation       Operation
	ListDeployments *ListDeploymentsOptions
}

// ListDeploymentsOptions は list-deployments 用のオプション。
type ListDeploymentsOptions struct {
	Project string
	Filter  string
	Format  string
	Limit   string
}

// Command は生成されたコマンドを表す。
type Command struct {
	Binary string
	Args   []string
}

// ArgsWithBinary は実行可能を含むコマンド配列を返す。
func (c Command) ArgsWithBinary() []string {
	if c.Binary == "" {
		return append([]string{}, c.Args...)
	}

	result := make([]string, 1+len(c.Args))
	result[0] = c.Binary
	copy(result[1:], c.Args)
	return result
}

// String はコマンドをシェル再現可能な形で連結する。
func (c Command) String() string {
	parts := c.ArgsWithBinary()
	for i, part := range parts {
		if part == "" {
			continue
		}

		if strings.ContainsAny(part, " \t\n\"") {
			parts[i] = strconv.Quote(part)
		}
	}
	return strings.Join(parts, " ")
}

// BuildCommand はオペレーションに応じた gcloud コマンドを構築する。
func (s *Service) BuildCommand(req BuildRequest) (Command, error) {
	switch req.Operation {
	case OperationListDeployments:
		if req.ListDeployments == nil {
			return Command{}, errors.New("list-deployments のオプションが指定されていません")
		}
		return buildListDeploymentsCommand(req.ListDeployments), nil
	default:
		return Command{}, fmt.Errorf("未対応の operation です: %s", req.Operation)
	}
}

func buildListDeploymentsCommand(opts *ListDeploymentsOptions) Command {
	args := []string{"deployment-manager", "deployments", "list"}

	if strings.TrimSpace(opts.Project) != "" {
		args = append(args, fmt.Sprintf("--project=%s", opts.Project))
	}

	if strings.TrimSpace(opts.Filter) != "" {
		args = append(args, fmt.Sprintf("--filter=%s", opts.Filter))
	}

	if strings.TrimSpace(opts.Format) != "" {
		args = append(args, fmt.Sprintf("--format=%s", opts.Format))
	}

	if strings.TrimSpace(opts.Limit) != "" {
		args = append(args, fmt.Sprintf("--limit=%s", opts.Limit))
	}

	return Command{
		Binary: "gcloud",
		Args:   args,
	}
}

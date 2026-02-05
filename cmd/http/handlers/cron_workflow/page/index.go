package page

import (
	"context"
	"io"
	"net/http"
	"strings"

	templ "github.com/a-h/templ"

	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	logging "github.com/landmaster135/devbox/internal/logging"

	body "github.com/landmaster135/devbox/internal/templ_components/core/body"
	button "github.com/landmaster135/devbox/internal/templ_components/core/button"
	code "github.com/landmaster135/devbox/internal/templ_components/core/code"
	div "github.com/landmaster135/devbox/internal/templ_components/core/div"
	form "github.com/landmaster135/devbox/internal/templ_components/core/form"
	head "github.com/landmaster135/devbox/internal/templ_components/core/head"
	headMeta "github.com/landmaster135/devbox/internal/templ_components/core/head_meta"
	headTitle "github.com/landmaster135/devbox/internal/templ_components/core/head_title"
	heading "github.com/landmaster135/devbox/internal/templ_components/core/heading"
	hiddenInput "github.com/landmaster135/devbox/internal/templ_components/core/hidden_input"
	html "github.com/landmaster135/devbox/internal/templ_components/core/html"
	link "github.com/landmaster135/devbox/internal/templ_components/core/link"
	mainComponent "github.com/landmaster135/devbox/internal/templ_components/core/main"
	paragraph "github.com/landmaster135/devbox/internal/templ_components/core/paragraph"
	script "github.com/landmaster135/devbox/internal/templ_components/core/script"
	section "github.com/landmaster135/devbox/internal/templ_components/core/section"
	span "github.com/landmaster135/devbox/internal/templ_components/core/span"
	usecaseArticle "github.com/landmaster135/devbox/internal/templ_components/usecase/article"
	usecaseStyle "github.com/landmaster135/devbox/internal/templ_components/usecase/style"
)

func Serve(
	w http.ResponseWriter,
	r *http.Request,
	logger *logging.StructuredLogger,
	workflows []usecases.Workflow,
	manualRunEndpoint string,
	manualWorkflowFieldName string,
) {
	data, err := buildWorkflowPageData(workflows, manualRunEndpoint, manualWorkflowFieldName)
	if err != nil {
		logger.WithTags("render").Errorf("failed to build workflow page data: %v", err)
		http.Error(w, "failed to render workflow page", http.StatusInternalServerError)
		return
	}

	templ.Handler(createCronWorkflowPage(data)).ServeHTTP(w, r)
}

// createCronWorkflowPage produces the templ component for the dashboard.
func createCronWorkflowPage(data workflowPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		mainChildren := []templ.Component{renderHeroSection(data)}
		if len(data.Categories) == 0 {
			mainChildren = append(mainChildren, renderEmptyState())
		} else {
			categorySections := make([]templ.Component, 0, len(data.Categories))
			for _, cat := range data.Categories {
				categorySections = append(categorySections, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					return renderCategorySection(ctx, w, cat)
				}))
			}
			mainChildren = append(mainChildren, div.Tag("page-sections", categorySections...))
		}
		return html.Document(
			"ja",
			head.Tag(
				headMeta.Base(),
				link.Tag(link.Attributes{
					Rel:   "icon",
					Href:  cronWorkflowFaviconDataURI(),
					Type:  "image/png",
					Sizes: "32x32",
				}),
				link.Tag(link.Attributes{
					Rel:   "shortcut icon",
					Href:  cronWorkflowFaviconDataURI(),
					Type:  "image/png",
					Sizes: "32x32",
				}),
				headTitle.Title(data.Title),
				usecaseStyle.Tag(),
			),
			body.Tag("",
				mainComponent.Tag("page", mainChildren...),
				script.Tag(cronWorkflowPageScript),
			),
		).Render(ctx, w)
	})
}

func renderHeroSection(data workflowPageData) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := paragraph.Text("hero-eyebrow", "Cron Workflow").Render(ctx, w); err != nil {
			return err
		}
		if err := heading.Heading(1, data.Title).Render(ctx, w); err != nil {
			return err
		}
		if data.Description != "" {
			if err := paragraph.Text("hero-description", data.Description).Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
	return section.Section("hero", "", body)
}

func renderEmptyState() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := heading.Heading(2, "No workflows").Render(ctx, w); err != nil {
			return err
		}
		return paragraph.Text("", "workflow.List() did not expose weatherWorkflow or dailyHeadingWorkflow.").Render(ctx, w)
	})
	return section.Section("empty-state", "", body)
}

func renderCategorySection(ctx context.Context, w io.Writer, cat workflowCategoryView) error {
	headerChildren := []templ.Component{heading.Heading(2, cat.Title)}
	if cat.Description != "" {
		headerChildren = append(headerChildren, paragraph.Text("", cat.Description))
	}
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := div.Tag("workflow-category__header", headerChildren...).Render(ctx, w); err != nil {
			return err
		}
		workflowCards := make([]templ.Component, 0, len(cat.Workflows))
		for _, wf := range cat.Workflows {
			workflowCards = append(workflowCards, renderWorkflowCard(wf))
		}
		return div.Tag("workflow-grid", workflowCards...).Render(ctx, w)
	})
	return section.Section("workflow-category", cat.ID, body).Render(ctx, w)
}

func renderWorkflowCard(card workflowCardView) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := div.Tag("workflow-card__header",
			paragraph.Content("workflow-card__process", code.Text("", card.ProcessName)),
			heading.Heading(3, card.Name),
		).Render(ctx, w); err != nil {
			return err
		}
		if card.Summary != "" {
			if err := paragraph.Text("workflow-card__summary", card.Summary).Render(ctx, w); err != nil {
				return err
			}
		}
		if card.CronDefinition != "" {
			if err := div.Tag("workflow-card__cron",
				span.Text("", "CRON"),
				code.Text("", card.CronDefinition),
			).Render(ctx, w); err != nil {
				return err
			}
		}
		if card.ManualAction != "" {
			method := card.ManualMethod
			if method == "" {
				method = http.MethodPost
			}
			upperMethod := strings.ToUpper(method)
			action := card.ManualAction
			body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if card.ManualWorkflowField != "" && card.ManualWorkflowValue != "" {
					if err := hiddenInput.HiddenField(card.ManualWorkflowField, card.ManualWorkflowValue).Render(ctx, w); err != nil {
						return err
					}
				}
				return button.Submit("Manual Run").Render(ctx, w)
			})
			if err := form.Tag("workflow-card__manual", upperMethod, action, action, body).Render(ctx, w); err != nil {
				return err
			}
		}
		return paragraph.Status("workflow-card__status").Render(ctx, w)
	})
	return usecaseArticle.Tag("workflow-card", card.Key, body)
}

const cronWorkflowPageScript = `(() => {
	const forms = document.querySelectorAll(".workflow-card__manual");
	forms.forEach((form) => {
		form.addEventListener("submit", async (event) => {
			event.preventDefault();
			const statusEl = form.parentElement?.querySelector("[data-manual-run-status]");
			if (statusEl) {
				statusEl.textContent = "手動実行中...";
				statusEl.dataset.status = "pending";
			}
			const formData = new FormData(form);
			const payload = new URLSearchParams();
			formData.forEach((value, key) => {
				payload.append(key, value.toString());
			});
			try {
				const response = await fetch(form.action, {
					method: form.method || "POST",
					headers: {
						"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
					},
					body: payload.toString(),
				});
				let data = null;
				try {
					data = await response.json();
				} catch (error) {
					data = null;
				}
				if (statusEl) {
					if (response.ok) {
						statusEl.textContent = (data && (data.message || data.status)) || "完了しました";
						statusEl.dataset.status = "success";
					} else {
						const message = (data && (data.error || data.message)) || ("HTTP " + response.status);
						statusEl.textContent = message;
						statusEl.dataset.status = "error";
					}
				}
			} catch (error) {
				if (statusEl) {
					statusEl.textContent = "リクエストに失敗しました";
					statusEl.dataset.status = "error";
				}
			}
		});
	});
})();`

package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_scheduler/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_scheduler/usecases"
)

func main() {
	conf, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if conf.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()
	command, err := buildCommand(service, conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
}

func buildCommand(service *usecases.Service, conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationCreatePubSubJob:
		return service.BuildCreatePubSubJobCommand(usecases.CreatePubSubJobParams{
			JobName:     conf.JobName,
			ProjectID:   conf.ProjectID,
			Location:    conf.Location,
			Schedule:    conf.Schedule,
			Description: conf.Description,
			TimeZone:    conf.TimeZone,
			PubsubTopic: conf.PubsubTopic,
			MessageBody: conf.MessageBody,
		})
	case cfg.OperationCreateHTTPJob:
		return service.BuildCreateHTTPJobCommand(usecases.CreateHTTPJobParams{
			JobName:                 conf.JobName,
			ProjectID:               conf.ProjectID,
			Location:                conf.Location,
			Schedule:                conf.Schedule,
			Description:             conf.Description,
			TimeZone:                conf.TimeZone,
			HTTPMethod:              conf.HTTPMethod,
			ServiceURL:              conf.ServiceURL,
			MessageBody:             conf.MessageBody,
			Headers:                 conf.Headers,
			OIDCServiceAccountEmail: conf.OIDCServiceAccountEmail,
		})
	case cfg.OperationCreateCloudSQLJob:
		return service.BuildCreateCloudSQLJobCommand(usecases.CreateCloudSQLJobParams{
			JobName:           conf.JobName,
			ProjectID:         conf.ProjectID,
			Location:          conf.Location,
			Schedule:          conf.Schedule,
			Description:       conf.Description,
			TimeZone:          conf.TimeZone,
			PubsubTopic:       conf.PubsubTopic,
			DBInstanceID:      conf.DBInstanceID,
			Action:            conf.Action,
			DiscordWebhookURL: conf.DiscordWebhookURL,
			CloudSQLIconURL:   conf.CloudSQLIconURL,
		})
	case cfg.OperationCreateStartCloudSQLJob:
		return service.BuildCreateStartCloudSQLJobCommand(usecases.CreateStartStopCloudSQLJobParams{
			JobName:           conf.JobName,
			ProjectID:         conf.ProjectID,
			Location:          conf.Location,
			Schedule:          conf.Schedule,
			Description:       conf.Description,
			TimeZone:          conf.TimeZone,
			PubsubTopic:       conf.PubsubTopic,
			DBInstanceID:      conf.DBInstanceID,
			DiscordWebhookURL: conf.DiscordWebhookURL,
			CloudSQLIconURL:   conf.CloudSQLIconURL,
		})
	case cfg.OperationCreateStopCloudSQLJob:
		return service.BuildCreateStopCloudSQLJobCommand(usecases.CreateStartStopCloudSQLJobParams{
			JobName:           conf.JobName,
			ProjectID:         conf.ProjectID,
			Location:          conf.Location,
			Schedule:          conf.Schedule,
			Description:       conf.Description,
			TimeZone:          conf.TimeZone,
			PubsubTopic:       conf.PubsubTopic,
			DBInstanceID:      conf.DBInstanceID,
			DiscordWebhookURL: conf.DiscordWebhookURL,
			CloudSQLIconURL:   conf.CloudSQLIconURL,
		})
	case cfg.OperationListJobs:
		return service.BuildListJobsCommand(usecases.ListJobsParams{
			Location: conf.Location,
			Limit:    conf.Limit,
		})
	case cfg.OperationUpdateHTTPJob:
		return service.BuildUpdateHTTPJobCommand(usecases.UpdateHTTPJobParams{
			JobName:     conf.JobName,
			Schedule:    conf.Schedule,
			MessageBody: conf.MessageBody,
			Headers:     conf.Headers,
		})
	case cfg.OperationUpdatePubSubJob:
		return service.BuildUpdatePubSubJobCommand(usecases.UpdatePubSubJobParams{
			JobName:     conf.JobName,
			Schedule:    conf.Schedule,
			MessageBody: conf.MessageBody,
		})
	case cfg.OperationPauseJob:
		return service.BuildPauseJobCommand(usecases.JobControlParams{
			JobName:  conf.JobName,
			Location: conf.Location,
		})
	case cfg.OperationResumeJob:
		return service.BuildResumeJobCommand(usecases.JobControlParams{
			JobName:  conf.JobName,
			Location: conf.Location,
		})
	case cfg.OperationDeleteJob:
		return service.BuildDeleteJobCommand(usecases.JobControlParams{
			JobName:  conf.JobName,
			Location: conf.Location,
		})
	case cfg.OperationRunJob:
		return service.BuildRunJobCommand(usecases.JobControlParams{
			JobName:  conf.JobName,
			Location: conf.Location,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}

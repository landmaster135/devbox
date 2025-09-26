package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_storage/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_storage/usecases"
)

func main() {
	config, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if config.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()

	command, err := buildCommandFromConfig(service, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
}

func buildCommandFromConfig(service *usecases.Service, conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationUploadFiles:
		return service.BuildUploadFilesCommand(usecases.UploadFilesParams{
			LocalPath: conf.LocalPath,
			BucketURL: conf.BucketURL,
		})
	case cfg.OperationDownloadFiles:
		return service.BuildDownloadFilesCommand(usecases.DownloadFilesParams{
			Sources:     conf.Sources,
			Destination: conf.Destination,
		})
	case cfg.OperationCreateBucket:
		return service.BuildCreateBucketCommand(usecases.CreateBucketParams{
			BucketURL:    conf.BucketURL,
			StorageClass: conf.StorageClass,
			Location:     conf.Location,
		})
	case cfg.OperationListContents:
		return service.BuildListContentsCommand(usecases.ListContentsParams{Target: conf.Target})
	case cfg.OperationShowDetails:
		return service.BuildShowDetailsCommand(usecases.ShowDetailsParams{Target: conf.Target})
	case cfg.OperationListNames:
		return service.BuildListNamesCommand(usecases.ListContentsParams{Target: conf.Target})
	case cfg.OperationDeleteObject:
		return service.BuildDeleteObjectCommand(usecases.DeleteObjectParams{Target: conf.Target})
	case cfg.OperationGetACL:
		return service.BuildGetACLCommand(usecases.ACLParams{Target: conf.Target})
	case cfg.OperationSetACL:
		return service.BuildSetACLCommand(usecases.ACLParams{ACLFile: conf.ACLFile, Target: conf.Target})
	case cfg.OperationGrantReadAll:
		return service.BuildGrantReadAllCommand(usecases.ACLParams{Target: conf.Target})
	case cfg.OperationRemoveReadAll:
		return service.BuildRemoveReadAllCommand(usecases.ACLParams{Target: conf.Target})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}

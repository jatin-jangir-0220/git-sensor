// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/devtron-labs/git-sensor/api"
	"github.com/devtron-labs/git-sensor/internal"
	"github.com/devtron-labs/git-sensor/internal/logger"
	"github.com/devtron-labs/git-sensor/internal/sql"
	"github.com/devtron-labs/git-sensor/pkg"
	"github.com/devtron-labs/git-sensor/pkg/git"
)

// Injectors from wire.go:

func InitializeApp() (*App, error) {
	sugaredLogger := logger.NewSugardLogger()
	config, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	materialRepositoryImpl := sql.NewMaterialRepositoryImpl(db)
	gitUtil := git.NewGitUtil(sugaredLogger)
	repositoryManagerImpl := git.NewRepositoryManagerImpl(sugaredLogger, gitUtil)
	gitProviderRepositoryImpl := sql.NewGitProviderRepositoryImpl(db)
	ciPipelineMaterialRepositoryImpl := sql.NewCiPipelineMaterialRepositoryImpl(db, sugaredLogger)
	repositoryLocker := internal.NewRepositoryLocker(sugaredLogger)
	conn, err := internal.NewNatsConnection()
	if err != nil {
		return nil, err
	}
	gitWatcherImpl, err := git.NewGitWatcherImpl(repositoryManagerImpl, materialRepositoryImpl, sugaredLogger, ciPipelineMaterialRepositoryImpl, repositoryLocker, conn)
	if err != nil {
		return nil, err
	}
	repoManagerImpl := pkg.NewRepoManagerImpl(sugaredLogger, materialRepositoryImpl, repositoryManagerImpl, gitProviderRepositoryImpl, ciPipelineMaterialRepositoryImpl, repositoryLocker, gitWatcherImpl)
	restHandlerImpl := api.NewRestHandlerImpl(repoManagerImpl, sugaredLogger)
	muxRouter := api.NewMuxRouter(sugaredLogger, restHandlerImpl)
	app := NewApp(muxRouter, sugaredLogger, gitWatcherImpl, db, conn)
	return app, nil
}

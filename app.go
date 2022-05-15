package main

import (
	"context"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/multibase-io/multibase/backend/grpc"
)

type OpenProtoFileResult struct {
	ProtoFilePath string `json:"protoFilePath"`
	CurrentDir    string `json:"currentDir"`
}

// App struct
type App struct {
	ctx          context.Context
	grpcProjects map[int]*grpc.Handler // TODO: concurrent-safe, also move to grpc.Module
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		grpcProjects: map[int]*grpc.Handler{},
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) OpenProtoFile() (*OpenProtoFileResult, error) {
	protoFilePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{DisplayName: "Proto Files (*.proto)", Pattern: "*.proto;"},
		},
	})
	if err != nil {
		return nil, err
	}

	return &OpenProtoFileResult{
		ProtoFilePath: protoFilePath,
		CurrentDir:    path.Dir(protoFilePath),
	}, nil
}

func (a *App) OpenImportPath() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}

	return path, nil
}

func (a *App) RefreshProtoDescriptors(projectID int, importPathList, protoFileList []string) ([]*grpc.ProtoTreeNode, error) {
	project, ok := a.grpcProjects[projectID]
	if !ok {
		project = grpc.NewHandler(projectID)
		a.grpcProjects[projectID] = project
	}

	nodes, err := a.grpcProjects[projectID].RefreshProtoDescriptors(importPathList, protoFileList)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (a *App) SelectMethod(projectID int, methodID string) (string, error) {
	project, ok := a.grpcProjects[projectID]
	if !ok {
		project = grpc.NewHandler(projectID)
		a.grpcProjects[projectID] = project
	}

	return a.grpcProjects[projectID].SelectMethod(methodID)
}

func (a *App) SendRequest(projectID int, id int, address, methodID, payload string) (string, error) {
	project, ok := a.grpcProjects[projectID]
	if !ok {
		project = grpc.NewHandler(projectID)
		a.grpcProjects[projectID] = project
	}

	return a.grpcProjects[projectID].SendRequest(id, address, methodID, payload)
}

func (a *App) StopRequest(projectID int, id int) error {
	project, ok := a.grpcProjects[projectID]
	if !ok {
		project = grpc.NewHandler(projectID)
		a.grpcProjects[projectID] = project
	}

	return a.grpcProjects[projectID].StopRequest(id)
}

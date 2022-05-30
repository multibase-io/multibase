package project

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/adrg/xdg"
	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"go.uber.org/multierr"
)

type State struct {
	Projects         map[string]*StateProject `json:"projects"`
	OpenedProjectIDs []string                 `json:"openedProjectIDs"`
	CurrentProjectID string                   `json:"currentProjectID"`
}

type StateProject struct {
	ID   string `json:"-"`
	Type string `json:"type"`
}

type Module struct {
	configFilePath string
	stateMutex     *sync.RWMutex
	state          *State
}

func NewModule() (*Module, error) {
	configFilePath, err := xdg.ConfigFile("multibase/project.json")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project config path: %w", err)
	}

	module := &Module{
		configFilePath: configFilePath,
		stateMutex:     &sync.RWMutex{},
	}

	err = module.readOrInitializeState()
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (m *Module) OpenGRPCProject(newProjectID, grpcProjectID string) (*State, error) {
	if lo.Contains(m.state.OpenedProjectIDs, grpcProjectID) {
		m.state.OpenedProjectIDs = lo.Reject(m.state.OpenedProjectIDs, func(projectID string, _ int) bool {
			return projectID == newProjectID
		})
	} else {
		lo.ReplaceAll(m.state.OpenedProjectIDs, newProjectID, grpcProjectID)
	}

	m.state.CurrentProjectID = grpcProjectID
	delete(m.state.Projects, newProjectID)

	err := m.saveState()
	if err != nil {
		return nil, err
	}
	return m.state, nil
}

func (m *Module) CreateGRPCProject(projectID string) (*State, error) {
	m.state.Projects[projectID] = &StateProject{
		ID:   projectID,
		Type: "grpc",
	}

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) CreateNewProject() (*State, error) {
	projectID := uuid.Must(uuid.NewV4()).String()

	m.state.Projects[projectID] = &StateProject{
		ID:   projectID,
		Type: "new",
	}

	m.state.OpenedProjectIDs = append(m.state.OpenedProjectIDs, projectID)
	m.state.CurrentProjectID = projectID

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) CloseProject(projectID string) (*State, error) {
	if len(m.state.OpenedProjectIDs) <= 1 {
		return m.state, nil
	}

	if m.state.Projects[projectID].Type == "new" {
		delete(m.state.Projects, projectID)
	}

	m.state.OpenedProjectIDs = lo.Reject(m.state.OpenedProjectIDs, func(pID string, _ int) bool {
		return pID == projectID
	})
	m.state.CurrentProjectID = m.state.OpenedProjectIDs[0]

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) State() (*State, error) {
	return m.state, nil
}

func (m *Module) readOrInitializeState() error {
	_, err := os.Stat(m.configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return m.initializeState()
		}

		return err
	}

	return m.readState()
}

func (m *Module) initializeState() (rerr error) {
	m.state = &State{
		Projects: map[string]*StateProject{
			"404f5702-6179-4861-9533-b5ee16161c78": {
				ID:   "404f5702-6179-4861-9533-b5ee16161c78",
				Type: "new",
			},
		},
		OpenedProjectIDs: []string{"404f5702-6179-4861-9533-b5ee16161c78"},
		CurrentProjectID: "404f5702-6179-4861-9533-b5ee16161c78",
	}

	file, err := os.Create(m.configFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = multierr.Combine(rerr, fmt.Errorf("failed to close a config file: %w", err))
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(m.state)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) readState() (rerr error) {
	file, err := os.Open(m.configFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = multierr.Combine(rerr, fmt.Errorf("failed to close a config file: %w", err))
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&m.state)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) saveState() (rerr error) {
	file, err := os.Create(m.configFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = multierr.Combine(rerr, fmt.Errorf("failed to close a config file: %w", err))
		}
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(m.state)
	if err != nil {
		return err
	}

	return nil
}

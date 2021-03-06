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
	Stats            *StateStats              `json:"stats"`
	Projects         map[string]*StateProject `json:"projects"`
	OpenedProjectIDs []string                 `json:"openedProjectIDs"`
	CurrentProjectID string                   `json:"currentProjectID"`
}

type StateProject struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type StateStats struct {
	GRPCProjectCount   int `json:"grpcProjectCount"`
	ThriftProjectCount int `json:"thriftProjectCount"`
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

func (m *Module) OpenProject(newProjectID, projectToOpenID string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	if lo.Contains(m.state.OpenedProjectIDs, projectToOpenID) {
		m.state.OpenedProjectIDs = lo.Reject(m.state.OpenedProjectIDs, func(projectID string, _ int) bool {
			return projectID == newProjectID
		})
	} else {
		m.state.OpenedProjectIDs = lo.ReplaceAll(m.state.OpenedProjectIDs, newProjectID, projectToOpenID)
	}

	m.state.CurrentProjectID = projectToOpenID
	delete(m.state.Projects, newProjectID)

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) CreateGRPCProject(projectID string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.state.Stats.GRPCProjectCount++

	projectName := fmt.Sprintf("gRPC Project %d", m.state.Stats.GRPCProjectCount)

	m.state.Projects[projectID] = &StateProject{
		ID:   projectID,
		Type: "grpc",
		Name: projectName,
	}

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) CreateThriftProject(projectID string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.state.Stats.ThriftProjectCount++

	projectName := fmt.Sprintf("Thrift Project %d", m.state.Stats.ThriftProjectCount)

	m.state.Projects[projectID] = &StateProject{
		ID:   projectID,
		Type: "thrift",
		Name: projectName,
	}

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) DeleteProject(projectID string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	delete(m.state.Projects, projectID)

	m.state.OpenedProjectIDs = lo.Reject(m.state.OpenedProjectIDs, func(pID string, _ int) bool {
		return pID == projectID
	})

	if m.state.CurrentProjectID == projectID {
		m.state.CurrentProjectID = m.state.OpenedProjectIDs[len(m.state.OpenedProjectIDs)-1]
	}

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) RenameProject(projectID, name string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.state.Projects[projectID].Name = name

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) CreateNewProject() (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

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
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	if len(m.state.OpenedProjectIDs) <= 1 {
		return m.state, nil
	}

	if m.state.Projects[projectID].Type == "new" {
		delete(m.state.Projects, projectID)
	}

	m.state.OpenedProjectIDs = lo.Reject(m.state.OpenedProjectIDs, func(pID string, _ int) bool {
		return pID == projectID
	})

	if m.state.CurrentProjectID == projectID {
		m.state.CurrentProjectID = m.state.OpenedProjectIDs[len(m.state.OpenedProjectIDs)-1]
	}

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) SaveCurrentProjectID(projectID string) (*State, error) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	m.state.CurrentProjectID = projectID

	err := m.saveState()
	if err != nil {
		return nil, err
	}

	return m.state, nil
}

func (m *Module) State() (*State, error) {
	m.stateMutex.RLock()
	defer m.stateMutex.RUnlock()

	return m.state, nil
}

func (m *Module) readOrInitializeState() error {
	_, err := os.Stat(m.configFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to describe a project config file: %w", err)
		}

		return m.initializeState()
	}

	return m.readState()
}

func (m *Module) initializeState() (rerr error) {
	m.state = &State{
		Stats: &StateStats{
			GRPCProjectCount:   0,
			ThriftProjectCount: 0,
		},
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
		return fmt.Errorf("failed to create a project config file: %w", err)
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
		return fmt.Errorf("failed to encode a project state: %w", err)
	}

	return nil
}

func (m *Module) readState() (rerr error) {
	file, err := os.Open(m.configFilePath)
	if err != nil {
		return fmt.Errorf("failed to open a project config file: %w", err)
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
		return fmt.Errorf("failed to decode a project state: %w", err)
	}

	return nil
}

func (m *Module) saveState() (rerr error) {
	file, err := os.Create(m.configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create/truncate a project config file: %w", err)
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
		return fmt.Errorf("failed to encode a project state: %w", err)
	}

	return nil
}

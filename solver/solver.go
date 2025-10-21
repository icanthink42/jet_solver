package solver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type SolverType string

const (
	SolverTypeGo     SolverType = "go"
	SolverTypePython SolverType = "python"
)

type InputType string

const (
	InputTypeRawJSON  InputType = "raw_json"
	InputTypeKeyValue InputType = "key_value_pair"
)

type Solver interface {
	Solve(input_data string) (string, error)
	Name() string
	Description() string
	Type() SolverType
	InputType() InputType
	PredefinedKeys() []string
	PredefinedUnits() map[string]string
}

func LoadAllSolvers(dir string) ([]Solver, error) {
	var solvers []Solver

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(dir, entry.Name(), "manifest.toml")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		solver, err := LoadSolverFromFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load solver from %s: %w", manifestPath, err)
		}
		solvers = append(solvers, solver)
	}

	if len(solvers) == 0 {
		return nil, fmt.Errorf("no solvers found in directory: %s", dir)
	}

	return solvers, nil
}

func LoadPythonSolver(name string, description string, manifest map[string]interface{}, solver_path string) (Solver, error) {
	typeSpecific := manifest["type_specific"].(map[string]interface{})
	mainFile := filepath.Join(filepath.Dir(solver_path), typeSpecific["main_file"].(string))
	pythonVersion := typeSpecific["python_version"].(string)

	solverMap := manifest["solver"].(map[string]interface{})
	var input_type InputType = InputTypeRawJSON
	if it, ok := solverMap["input_type"].(string); ok {
		input_type = InputType(it)
	}

	var predefined_keys []string
	var predefined_units = make(map[string]string)

	if keys, ok := solverMap["predefined_keys"].([]interface{}); ok {
		predefined_keys = make([]string, len(keys))
		for i, k := range keys {
			predefined_keys[i] = k.(string)
		}
	}

	if units, ok := solverMap["predefined_units"].(map[string]interface{}); ok {
		for key, unit := range units {
			if unitStr, ok := unit.(string); ok {
				predefined_units[key] = unitStr
			}
		}
	}

	return &PythonSolver{
		solver_name:        name,
		solver_description: description,
		main_file:          mainFile,
		python_version:     pythonVersion,
		input_type:         input_type,
		predefined_keys:    predefined_keys,
		predefined_units:   predefined_units,
	}, nil
}

func LoadSolverFromFile(solver_path string) (Solver, error) {
	var manifest map[string]interface{}
	data, err := os.ReadFile(solver_path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	solverMap := manifest["solver"].(map[string]interface{})
	solver_name := solverMap["name"].(string)
	solver_description := solverMap["description"].(string)
	solver_type := solverMap["type"].(string)
	switch SolverType(solver_type) {
	case SolverTypePython:
		return LoadPythonSolver(solver_name, solver_description, manifest, solver_path)
	default:
		return nil, fmt.Errorf("unsupported solver type: %s", solver_type)
	}
}

package solver

import (
	"bytes"
	"fmt"
	"os/exec"
)

type PythonSolver struct {
	solver_name        string
	solver_description string
	main_file          string
	python_version     string
	input_type         InputType
	predefined_keys    []string
	predefined_units   map[string]string // Maps each key to its unit
}

func (s *PythonSolver) Solve(input_data string) (string, error) {
	cmd := exec.Command("python", "-E", s.main_file, input_data)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%v\nStderr: %s", err, stderr.String())
		}
		return "", err
	}
	return string(output), nil

}

func (s *PythonSolver) Name() string {
	return s.solver_name
}

func (s *PythonSolver) Description() string {
	return s.solver_description
}

func (s *PythonSolver) Type() SolverType {
	return SolverTypePython
}

func (s *PythonSolver) InputType() InputType {
	return s.input_type
}

func (s *PythonSolver) PredefinedKeys() []string {
	return s.predefined_keys
}

func (s *PythonSolver) PredefinedUnits() map[string]string {
	return s.predefined_units
}

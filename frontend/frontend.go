package frontend

import (
	_ "embed"
	"fmt"
	"jet_solver/solver"
	"strings"
)

//go:embed solvers_list/index.css
var solverListCSS string

//go:embed solver_input/index.css
var solverInputCSS string

func SolverOutput(name string, input string) string {
	solvers, err := solver.LoadAllSolvers("solvers")
	if err != nil {
		return fmt.Sprintf("Error: Failed to load solvers: %v", err)
	}

	var selectedSolver solver.Solver
	for _, s := range solvers {
		if s.Name() == name {
			selectedSolver = s
			break
		}
	}

	if selectedSolver == nil {
		return "Error: Solver not found"
	}

	output, err := selectedSolver.Solve(input)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return output
}

func SolverInput(name string) string {
	solvers, err := solver.LoadAllSolvers("solvers")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	var selectedSolver solver.Solver
	for _, s := range solvers {
		if s.Name() == name {
			selectedSolver = s
			break
		}
	}

	if selectedSolver == nil {
		return "Error: Solver not found"
	}

	var html strings.Builder
	html.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<style>%s</style>
</head>
<body>
<div class="container">
	<div class="input-area">
		<form action="/output" method="POST" id="solverForm">
			<input type="hidden" name="name" value="%s">`, solverInputCSS, name))

	if selectedSolver.InputType() == solver.InputTypeKeyValue {
		html.WriteString(`
			<div id="kvPairs">`)

		// Add predefined keys first
		for _, key := range selectedSolver.PredefinedKeys() {
			html.WriteString(fmt.Sprintf(`
				<div class="kv-pair">
					<input type="text" class="key" placeholder="Key" value="%s">
					<input type="text" class="value" placeholder="Value">
				</div>`, key))
		}

		// Add an empty pair if no predefined keys
		if len(selectedSolver.PredefinedKeys()) == 0 {
			html.WriteString(`
				<div class="kv-pair">
					<input type="text" class="key" placeholder="Key">
					<input type="text" class="value" placeholder="Value">
				</div>`)
		}

		html.WriteString(`
			</div>
			<button type="button" onclick="addKVPair()" class="add-button">Add Field</button>
			<input type="hidden" name="json" id="jsonData">
			<script>
				function addKVPair() {
					const div = document.createElement('div');
					div.className = 'kv-pair';
					div.innerHTML = '<input type="text" class="key" placeholder="Key"><input type="text" class="value" placeholder="Value"><button type="button" onclick="this.parentElement.remove()" class="remove-button">Remove</button>';
					document.getElementById('kvPairs').appendChild(div);
				}
				document.getElementById('solverForm').onsubmit = function() {
					const pairs = document.getElementsByClassName('kv-pair');
					const data = {};
					for (const pair of pairs) {
						const key = pair.querySelector('.key').value.trim();
						const value = pair.querySelector('.value').value.trim();
						if (key) data[key] = value;
					}
					document.getElementById('jsonData').value = JSON.stringify(data);
					return true;
				};
			</script>`)
	} else {
		html.WriteString(`<textarea name="json" placeholder="Enter your JSON input here..."></textarea>`)
	}

	html.WriteString(`
			<button type="submit" class="run-button">Run Solver</button>
		</form>
	</div>
</div>
</body>
</html>`)

	return html.String()
}

func SolverList() string {
	solvers, err := solver.LoadAllSolvers("solvers")
	if err != nil {
		return fmt.Sprintf("<div class=\"error\">Failed to load solvers: %v</div>", err)
	}

	var html strings.Builder
	html.WriteString(fmt.Sprintf(`<style>%s</style>`, solverListCSS) + `
<table class="solver-table">
	<thead>
		<tr>
			<th>Name</th>
			<th>Description</th>
			<th>Type</th>
		</tr>
	</thead>
	<tbody>`)

	for _, s := range solvers {
		html.WriteString(fmt.Sprintf(`
		<tr onclick="window.location='/solver?name=%s'">
			<td>%s</td>
			<td>%s</td>
			<td><span class="type-badge">%s</span></td>
		</tr>`, s.Name(), s.Name(), s.Description(), s.Type()))
	}

	html.WriteString(`
	</tbody>
</table>`)

	return html.String()
}

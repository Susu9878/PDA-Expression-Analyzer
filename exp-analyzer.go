package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"CodersSquad/ct-first-pda/stack"

	"gopkg.in/yaml.v2"
)

type PDA struct {
	InitialState  string       `yaml:"S"`
	States        []string     `yaml:"K"`
	Alphabet      []string     `yaml:"E"`
	StackAlphabet []string     `yaml:"R"`
	Transitions   []Transition `yaml:"T"`
	FinalStates   []string     `yaml:"F"`
}

type Transition struct {
	State     string `yaml:"state"`
	Input     string `yaml:"input"`
	StackPop  string `yaml:"stack_pop"`
	ToState   string `yaml:"to_state"`
	StackPush string `yaml:"stack_push"`
}

type Result struct {
	Num        int    `yaml:"num"`
	Expression string `yaml:"expression"`
	Valid      bool   `yaml:"valid"`
}

// loads .yaml returns PDA struct, loads all strings to evaluate returns them as []string
func openInputFiles(pdaFile, inputFile string) (pda PDA, vals []string) {
	yamlFile, err := os.ReadFile(pdaFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &pda)
	var expandedTransitions []Transition //suport all numbers so i dont have to write 10 whole transitions every time i read a number
	for _, t := range pda.Transitions {
		switch t.Input {
		case "num":
			for i := 0; i <= 9; i++ {
				expandedTransitions = append(expandedTransitions, Transition{
					State:     t.State,
					Input:     fmt.Sprintf("%d", i),
					StackPop:  t.StackPop,
					StackPush: t.StackPush,
					ToState:   t.ToState,
				})
			}
		default:
			expandedTransitions = append(expandedTransitions, t)
		}
	}
	pda.Transitions = expandedTransitions

	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	stringsFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer stringsFile.Close()

	scanner := bufio.NewScanner(stringsFile)
	for scanner.Scan() {
		line := scanner.Text()
		vals = append(vals, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func executeTransition(pda PDA, currentState string, inputChar string, stack *stack.Stack) (string, bool) {
	for _, trans := range pda.Transitions {
		if currentState == trans.State && inputChar == trans.Input {
			// Handle stack pop
			if trans.StackPop != "ε" {
				if stack.IsEmpty() {
					return currentState, false
				}
				popped, _ := stack.Pop()
				if popped != trans.StackPop {
					return currentState, false
				}
			}

			if trans.StackPush != "ε" {
				stack.Push(trans.StackPush)
			}

			// Transition executed crrectly
			return trans.ToState, true
		}
	}
	// No valid transition found
	return currentState, false
}

// check if state is a final one
func isFinalState(state string, finalStates []string) bool {
	for _, fs := range finalStates {
		if state == fs {
			return true
		}
	}
	return false
}

func analizeString(pda PDA, expression string, num_expression int) (result Result) {
	stack := stack.New()
	valid := false
	cleanStr := strings.ReplaceAll(expression, " ", "")
	// validate alphabet. if all characters are in alphabet the for loop will set valid to true
	for _, char := range cleanStr {
		valid = false
		for _, charA := range pda.Alphabet {
			if string(char) == charA {
				valid = true
				break
			}
		}
		if !valid {
			break
		}
	}

	//validate with transitions, not even bother if alphabet is invalid
	currentState := pda.InitialState
	if valid {
		//check every char in string
		for _, char := range cleanStr {
			currentState, valid = executeTransition(pda, currentState, string(char), stack)
			if !valid {
				break
			}
		}
	}

	//check final state
	if valid && isFinalState(currentState, pda.FinalStates) && stack.IsEmpty() {
		valid = true
	} else {
		valid = false
	}

	result.Num = num_expression
	result.Expression = expression
	result.Valid = valid
	return result
}

func processData(pda PDA, expressions []string) (results []Result) {
	for i, str := range expressions {
		// log.Println(str)a
		rs := analizeString(pda, str, i)
		results = append(results, rs)
	}
	return
}

func main() {
	var pda = flag.String("pda", "pda.yaml", "The PDA machine definition file")
	var input = flag.String("input", "strings.txt", "Input list of strings file")
	var output = flag.String("output", "results.txt", "Yaml-formatted results file")

	flag.Parse()

	pdaDef, vals := openInputFiles(*pda, *input)

	results := processData(pdaDef, vals)
	//log.Println(results)
	//log.Println(vals)

	var lines []string
	for _, res := range results {
		status := "Invalid"
		if res.Valid {
			status = "Valid"
		}
		line := fmt.Sprintf("Expression %d: %-20s -  %s", res.Num, res.Expression, status)
		lines = append(lines, line)
	}

	outputData := strings.Join(lines, "\n")

	if err := os.WriteFile(*output, []byte(outputData), 0644); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

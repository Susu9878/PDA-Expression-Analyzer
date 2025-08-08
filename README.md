# PDA Expressions Analyzer

This tool reads a Pushdown Automaton (PDA) definition and checks whether input expressions are grammatically valid in the defined language.

---

## PDA Example: Language of the form `0ⁿ 1ᵐ 2ᵐ 3ⁿ`

### Example Accepted Strings
- 0123 → ✅ Valid (n = 1, m = 1)
- 00112233 → ✅ Valid (n = 2, m = 2)
- 000111222333 → ✅ Valid (n = 3, m = 3)

### Example Rejected Strings
- 012 → ❌ Rejected (missing a 3)
- 0011223 → ❌ Rejected (n ≠ number of 3s)
- 0112233 → ❌ Rejected (unbalanced 1s and 2s)

### Notes
- ε (epsilon) represents the empty string or no operation.
- For stack_pop: ε, the stack is not required to pop a value.
- For stack_push: ε, nothing is pushed to the stack.
- The PDA simulates a two-part balancing process: Pushes a 0 for every 0 seen, and a 1 for every 1.
- Pops a 1 for every 2, and a 0 for every 3.
- The input is accepted if the automaton reaches state q3 with matching counts.

### Where:
- The number of `0`s must equal the number of `3`s (`n` times)
- The number of `1`s must equal the number of `2`s (`m` times)
---

### PDA Definition (YAML Format)

```yaml
S: q0  # initial state
K: [q0, q1, q2, q3]  # set of states
E: [0, 1, 2, 3]  # input alphabet
R: [0, 1]  # stack alphabet
F: [q3]  # set of accepting states

T:  # transitions
  - state: q0 
    input: 0
    stack_pop: ε
    stack_push: 0
    to_state: q0

  - state: q0 
    input: 1
    stack_pop: ε
    stack_push: 1
    to_state: q1

  - state: q1
    input: 1
    stack_pop: ε
    stack_push: 1
    to_state: q1

  - state: q1
    input: 2
    stack_pop: 1
    stack_push: ε
    to_state: q2

  - state: q2
    input: 2
    stack_pop: 1
    stack_push: ε
    to_state: q2

  - state: q2
    input: 3
    stack_pop: 0
    stack_push: ε
    to_state: q3

  - state: q3
    input: 3
    stack_pop: 0
    stack_push: ε
    to_state: q3

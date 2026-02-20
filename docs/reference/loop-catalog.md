# Springfield Protocol: Agentic Loop Catalog

A comprehensive catalog of feedback loops, control patterns, and agentic workflows used in the Springfield Protocol framework.

---

## 1. Foundational Loops

### 1.1 Sense-Plan-Act (SPA)
The atomic unit of agentic reasoning and action.

**Flow:**
```
Sense (Perceive) → Plan (Think) → Act (Execute) → [Loop back to Sense]
```

**Characteristics:**
- Observes the current state
- Reasons about the best course of action
- Executes the action
- Feeds results back into the next cycle

**Use Case:** Real-time decision-making, reactive systems, incremental problem-solving.

Phases:
- Design Thinking
- Agile Development

Used By:
- All AI/Agentic Agents

---

### 1.2 ReAct (Reason + Act)
An enhanced version of Sense-Plan-Act that explicitly tracks and communicates reasoning steps.

**Flow:**
```
Thought → Action → Observation → Thought → ... → Final Answer
```

**Characteristics:**
- Verbalizes the reasoning at each step
- Takes concrete tool-using actions
- Observes the results before proceeding
- Maintains a dialogue-like trace for transparency

**Use Case:** Debugging, troubleshooting, clarifying agent thinking, building user trust.

**Example:**
```
Thought: "I need to check the test results"
Action: run_tests()
Observation: "Tests failed with error X"
Thought: "I need to examine the code to fix error X"
Action: read_file("pkg/logger/logger.go")
...
Final Answer: "Fixed the issue by..."
```

Phases:
- Design Thinking
- Agile Development

Used By:
- All AI/Agentic Agents

---

## 2. Refinement & Quality Loops

### 2.1 Generate → Evaluate → Critique → Refine (GECR)
A multi-stage refinement loop for iterative improvement of output quality.

**Flow:**
```
Generate (Create Options)
  ↓
Evaluate (Score/Rank)
  ↓
Critique (Analyze Weaknesses)
  ↓
Refine (Improve Based on Critique)
  ↓
[Loop back to Evaluate if needed]
```

**Characteristics:**
- Generates initial candidates/solutions
- Scores them against criteria
- Identifies weaknesses and failure modes
- Iteratively improves based on feedback
- Continues until quality threshold is reached

**Use Case:** Code review, design validation, output polishing, high-stakes decisions.

Phases:
- Design Thinking

Used By:
- Marge (Design Reviews)


---

### 2.2 Test → Analyze → Learn → Adjust → Retest (TALAR)
The empirical learning loop, rooted in TDD and scientific method.

**Flow:**
```
Test (Run Trial)
  ↓
Analyze (Examine Results)
  ↓
Learn (Extract Insight)
  ↓
Adjust (Change Approach)
  ↓
Retest (Validate Change)
  ↓
[Loop back if improvement < threshold]
```

**Characteristics:**
- Runs experiments systematically
- Analyzes outcomes without bias
- Extracts actionable insights
- Adjusts strategy incrementally
- Validates improvements empirically

**Use Case:** TDD workflows, A/B testing, continuous optimization, performance tuning.

Used By:
- Bart (Breaks & Testing)

---

### 2.3 Observe → Hypothesize → Experiment → Conclude → Iterate (OHECI)
The scientific method applied to agent reasoning and problem-solving.

**Flow:**
```
Observe (Gather Data)
  ↓
Hypothesize (Form Theory)
  ↓
Experiment (Test Theory)
  ↓
Conclude (Draw Results)
  ↓
Iterate (Refine Theory)
  ↓
[Loop back to Observe with refined hypothesis]
```

**Characteristics:**
- Builds knowledge incrementally
- Forms testable hypotheses
- Validates against evidence
- Draws conclusions rigorously
- Refines the mental model

**Use Case:** Root cause analysis, incident investigation, system understanding, research.

Used By:
- Marge (Design Reviews)

---

## 3. Exploration & Decision Loops

### 3.1 Tree of Thoughts (ToT)
Explores multiple reasoning paths in parallel and prunes unpromising branches.

**Flow:**
```
      ┌─────────────────┐
      │  Start Problem  │
      └────────┬────────┘
               │
      ┌────────┴────────┐
      │                 │
   Option A        Option B        Option C
      │                 │                │
   Eval(Good)       Eval(Bad)       Eval(Maybe)
      │                 │                │
    [Continue]       [Stop/Prune]   [Explore]
      │                 │                │
      └────────┬────────┴────────┬──────┘
               │
         ┌─────▼─────┐
         │Final Path │
         └───────────┘
```

**Characteristics:**
- Generates multiple solution paths
- Evaluates each path
- Prunes low-scoring or infeasible options
- Explores promising paths deeper
- Terminates at acceptable solution

**Use Case:** Complex problem-solving, multi-step planning, design decisions, brainstorming.

Phases:
- Design Thinking

Used By:
- Marge (Design Reviews)

---

### 3.2 Generate → Score → Prune → Explore (GSPE)
A variant of ToT that explicitly separates scoring and exploration phases.

**Flow:**
```
Generate (Create Multiple Options)
  ↓
Score (Rank Them)
  ↓
Prune (Remove Low-Scoring)
  ↓
Explore (Dive into Top Options)
  ↓
[Generate refinements and loop]
```

**Characteristics:**
- Generates diverse candidates upfront
- Scores all candidates against criteria
- Prunes implausible or redundant options
- Dedicates exploration to the most promising ones
- Repeats for depth

**Use Case:** Architecture decisions, feature prioritization, variant selection.

Phases:
- Design Thinking

Used By:
- Marge (Design Reviews)

---

### 3.3 Verbalized Sampling (VS)
Generates multiple responses with explicit reasoning, then selects the best.

**Flow:**
```
Generate (Response A + Reasoning A)
Generate (Response B + Reasoning B)
Generate (Response C + Reasoning C)
  ↓
Verbalize & Score (Compare Reasoning Quality)
  ↓
Select (Highest Reasoning Quality)
  ↓
Output (Best Response)
```

**Characteristics:**
- Generates multiple diverse responses
- Asks the model to reason about each
- Evaluates reasoning quality, not just output
- Selects the response with the strongest explanation
- Improves reliability through "meta-reasoning"

**Use Case:** Open-ended questions, creative work, nuanced decisions, reducing hallucinations.

Phases:
- Design Thinking

Used By:
- Marge (Design Reviews)

---

## 4. Planning & Orchestration Loops

### 4.1 Plan-and-Execute (PAE)
Strategic planning followed by guided execution.

**Flow:**
```
User Goal
  ↓
Planner (Creates Task List)
  ↓
[Task List]
  ↓
Executor (Loop)
  ├─ Execute Task N
  ├─ Validate Result
  ├─ Move to Task N+1
  └─ [Loop until all tasks done]
  ↓
Refiner (Polish & Integrate)
  ↓
Final Output
```

**Characteristics:**
- Breaks problems into discrete tasks upfront
- Maintains a task list/plan
- Executes tasks sequentially or in dependency order
- Validates each step before moving to the next
- Refines and integrates at the end

**Use Case:** Large projects, feature implementation, complex workflows, architectural changes.

Phases:
- Design Thinking

Used By:
- Ralph (Code Generation & Refactoring)
- Lisa (Orchestration)

---

### 4.2 The Ralph Wiggum Variant (Stateless Resampling Loop)
The core engine for Agile Agentic Development.

**Flow:**
```
┌────────────────────────────────────────────────┐
│              Control Loop Monitor              │
│  while (exists task in PLAN.json where        │
│         validation == false):                  │
└───────────────────┬────────────────────────────┘
                    │
        ┌───────────┴────────────┐
        │                        │
    Spawn              Ephemeral Context
    Agent              (Clean git worktree)
        │                        │
        └────────┬───────────────┘
                 │
        ┌────────▼──────────┐
        │  Execute Task     │
        │  (TDD Strict)     │
        └────────┬──────────┘
                 │
        ┌────────▼──────────────────┐
        │  Verification Loop        │
        │  (Bart Review)            │
        │  Updates PLAN.json        │
        └────────┬──────────────────┘
                 │
                [Loop back to Control Loop]
```

**Characteristics:**
- Tracks all tasks in `PLAN.json` with validation state
- Control loop monitors and schedules work
- Each iteration uses ephemeral context (clean slate)
- Prevents context rot and hallucination accumulation
- Verification loop feeds results back to plan
- Stateless: no hidden agent memory persists

**Use Case:** High-volume autonomous development, long-running projects, quality-critical work.

Used By:
- Ralph (Code Generation & Refactoring)
- Lisa (Orchestration)

---

## 5. Multi-Agent & Collaboration Loops

### 5.1 Manager-Worker Loop (MWL)
Coordinated collaboration between orchestrator and specialized agents.

**Flow:**
```
                  ┌─────────────┐
                  │  Manager    │
                  │  (Lisa)     │
                  └─────┬───────┘
                        │
                        │ Delegates
                        │
        ┌───────────────┼───────────────┐
        │               │               │
    ┌───▼──┐        ┌───▼──┐       ┌───▼──┐
    │Worker│        │Worker│       │Worker│
    │Ralph │        │Bart  │       │Marge │
    └───┬──┘        └───┬──┘       └───┬──┘
        │               │               │
        └───────────────┼───────────────┘
                        │
                   Feedback & Aggregation
                        │
                  ┌─────▼──────┐
                  │  Manager   │
                  │ Orchestrates
                  │ Next Round │
                  └────────────┘
```

**Characteristics:**
- Central manager (e.g., Lisa) orchestrates work
- Workers specialize in specific tasks (e.g., Ralph on code, Bart on breaks)
- Manager delegates tasks and collects feedback
- Workers execute in parallel or sequence
- Manager aggregates results and decides next steps

**Use Case:** Complex projects requiring multiple specialized roles, parallel execution, asynchronous workflows.

Used By:
- Lisa (Orchestration)
- Ralph (Code Generation & Refactoring)

---

### 5.2 Dialogue Loop (DL)
Back-and-forth communication between agents for iterative refinement.

**Flow:**
```
          Developer          Reviewer
          (Agent)            (Agent)
            │ Message          │
            ├─────────────────>│
            │                  │ Review
            │     Feedback     │
            │<─────────────────┤
            │ Refine           │
            ├─────────────────>│
            │                  │ Approve
            │                  │
            └────────────────┬─┘
                             │
                        Final Output
```

**Characteristics:**
- Two agents exchange messages and feedback
- Developer proposes, reviewer critiques
- Iterative rounds until convergence
- Maintains context across exchanges
- Produces polished final output

**Use Case:** Code review workflows, design collaboration, quality assurance, writing refinement.

Used By:
- Bart (Breaks & Testing)

---

## 6. Knowledge & Learning Loops

### 6.1 Observe → Learn → Apply (OLA)
Captures learning and applies it to future decisions.

**Flow:**
```
Observe (Encounter Situation)
  ↓
Learn (Extract Pattern/Rule)
  ↓
Store (Add to Knowledge Base)
  ↓
Apply (Use in Future Decisions)
  ↓
[New Observation → Loop]
```

**Characteristics:**
- Encounters new situations
- Extracts generalizable knowledge
- Stores in knowledge base (e.g., KEDB, ADRs)
- Reapplies in similar future contexts
- Builds expertise over time

**Use Case:** Error handling (KEDB), architectural patterns (ADRs), best practices, organizational learning.

---

### 6.2 Reflect → Document → Share (RDS)
Captures insights and makes them available to other agents.

**Flow:**
```
Reflect (Process Experience)
  ↓
Summarize (Extract Key Insights)
  ↓
Document (Create Reference)
  ↓
Share (Make Discoverable)
  ↓
[Other Agents Benefit]
```

**Characteristics:**
- Reflects on completed work
- Identifies generalizable insights
- Documents in discoverable format (ADR, runbook, etc.)
- Shares with team/agents
- Accelerates learning across organization

**Use Case:** ADR creation, runbook writing, incident post-mortems, best practice dissemination.

---

## 7. Cross-Cutting Concerns

### 7.1 Backpressure Loop
Monitors system health and slows work when quality degrades.

**Flow:**
```
Monitor Quality Metrics
  ↓
Quality > Threshold? ─YES─> Continue ─┐
  │                                    │
  NO                                   │
  │                                    │
Throttle Work ─> Pause/Investigate ───┤
  │                                    │
  └────────────────────────────────────┘
```

**Characteristics:**
- Monitors test coverage, error rates, queue depth
- Pauses work when quality drops below threshold
- Forces investigation before resuming
- Prevents "rushing and breaking"
- Maintains system health

**Use Case:** CI/CD pipelines, deployment gates, quality enforcement.

---

### 7.2 Feedback Accumulation Loop
Collects and synthesizes feedback over time.

**Flow:**
```
Collect Feedback (From Multiple Sources)
  ↓
Aggregate (Combine & Synthesize)
  ↓
Analyze (Identify Patterns)
  ↓
Act (Update Plan/Process)
  ↓
[Loop back to Collect]
```

**Characteristics:**
- Gathers feedback from various stakeholders
- Aggregates into coherent themes
- Identifies trends and patterns
- Acts on collective feedback
- Adapts based on aggregate learning

**Use Case:** Sprint reviews, user research synthesis, retrospectives, iterative refinement.

---

## 8. Quick Reference Matrix

| Loop | Phase | Complexity | Parallelism | Key Output |
|:-----|:------|:-----------|:------------|:-----------|
| Sense-Plan-Act | Real-time | Low | Sequential | Next Action |
| ReAct | Debugging | Medium | Sequential | Trace + Answer |
| GECR | Quality | High | Sequential | Refined Output |
| TALAR | Learning | Medium | Sequential | Learned Change |
| OHECI | Research | High | Sequential | Validated Hypothesis |
| Tree of Thoughts | Planning | High | Parallel | Best Path |
| GSPE | Decision | Medium | Parallel | Selected Option |
| Verbalized Sampling | Output | Medium | Parallel | Best Reasoning |
| Plan-and-Execute | Implementation | High | Sequential | Complete Delivery |
| Ralph Wiggum | Development | High | Sequential + Parallel | Quality Deliverable |
| Manager-Worker | Orchestration | Very High | Parallel | Integrated Result |
| Dialogue | Collaboration | Medium | Sequential | Refined Agreement |
| OLA | Knowledge | Medium | Sequential | Applied Insight |
| RDS | Learning | Medium | Sequential | Shared Knowledge |
| Backpressure | Governance | Medium | All | System Health |
| Feedback Accumulation | Governance | Medium | Sequential | Adaptive Strategy |

---

## 9. Loop Selection Guide

**Choosing the right loop depends on:**

1. **Problem Type**
   - Vague/Complex → Tree of Thoughts
   - Well-defined → Plan-and-Execute
   - Specific error → ReAct
   - Quality-critical → GECR + Ralph Wiggum

2. **Team Structure**
   - Single agent → Sense-Plan-Act or ReAct
   - Multiple roles → Manager-Worker or Dialogue
   - Distributed → Ralph Wiggum + Plan-and-Execute

3. **Constraints**
   - Time-critical → Tree of Thoughts (parallel pruning)
   - Quality-critical → Ralph Wiggum (iterative verification)
   - Learning-focused → OHECI or OLA (build knowledge)
   - Collaborative → Dialogue or Manager-Worker

4. **Feedback Model**
   - Immediate feedback available → ReAct or Sense-Plan-Act
   - Delayed feedback → TALAR or OHECI
   - No direct feedback → Tree of Thoughts or GSPE

---

## References

- **Sense-Plan-Act**: Classic robotics control pattern
- **ReAct**: Yao et al., "ReAct: Synergizing Reasoning and Acting in Language Models"
- **Tree of Thoughts**: Yao et al., "Tree of Thoughts: Deliberate Problem Solving with LLMs"
- **Verbalized Sampling**: Chain of Thought variants
- **Ralph Wiggum Variant**: Springfield Protocol custom design

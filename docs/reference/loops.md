# Agentic Loops Reference

This catalog defines the feedback patterns and coordination strategies used in the Springfield Protocol.

---

## 1. Foundational Patterns

### Sense-Plan-Act (Observe → Think → Act)
The simplest agent loop. Used for immediate, low-complexity tasks.
- **Use When:** Handling straightforward instructions without significant error risk.

### ReAct (Reason + Act)
Verbalize thinking before taking actions.
- **Use When:** Debugging specific errors where the agent needs to explain its reasoning at each step.

---

## 2. Refinement & Quality

### Ralph Wiggum Loop (The Engine)
Stateless resampling via ephemeral contexts.
- **Use When:** Standard feature implementation and verification.
- **Key:** Fresh context each iteration prevents error accumulation.

### GECR (Generate → Evaluate → Critique → Refine)
Multi-stage refinement for high-quality output.
- **Use When:** Polishing documentation, complex logic, or creative architectural proposals.

### TALAR (Test → Analyze → Learn → Adjust → Retest)
Experiment-driven optimization.
- **Use When:** Performance tuning or resolving complex race conditions.

---

## 3. Exploration & Decisions

### Tree of Thoughts (ToT)
Exploring multiple reasoning paths and pruning low-scoring options.
- **Use When:** Vague or complex problems with many possible solutions.

### OHECI (Observe-Hypothesize-Experiment-Conclude-Integrate)
A scientific method loop for architecture and research spikes.
- **Use When:** Conducting spikes or validating new architectural patterns.

---

## 4. Coordination & Orchestration

### Manager-Worker Loop
One orchestrator (e.g., Planning Agent) coordinates specialized workers (e.g., Build Agents) in parallel.
- **Use When:** Managing large epics or multi-agent workstreams.

### Dialogue Loop
Two agents iterating (propose → critique → refine).
- **Use When:** Peer review between Build and Quality agents.

---

## Loop Selection Guide

| Problem Type | Complexity | Recommended Loop |
| :--- | :--- | :--- |
| Simple task | Low | **Sense-Plan-Act** |
| Specific error | Medium | **ReAct** |
| Vague problem | High | **Tree of Thoughts** |
| Feature build | Medium | **Ralph Wiggum** |
| Optimization | Medium | **TALAR** |
| Peer review | Low | **Dialogue** |

For visual flowcharts of these loops, see [docs/reference/visual-diagrams.md](visual-diagrams.md).

For detailed technical specifications of each loop, see [docs/reference/loop-catalog.md](loop-catalog.md).

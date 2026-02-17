# Springfield Protocol v0.2: The Master Model

This document defines the core logic, philosophy, and architecture of the Springfield Protocol.

---

## 1. Core Principles

The framework is built on nine foundational principles that align Design Thinking and Agile Delivery with the Three Ways of DevOps.

### 1. Systems Thinking (The First Way)
Visualize the whole value stream. Plan before you build. The Discovery Diamond must produce a Feature Brief to ensure implementation is aligned with the overall system goal.

### 2. Flow via Decentralized Coordination
Maintain progress through loose coupling. Agents monitor the state of the 7 Core Documents (the "Blackboard") and "pull" work when triggers are met, eliminating synchronous hand-offs and silos.

### 3. Amplify Feedback Loops (The Second Way)
Shorten and sharpen feedback. Use the Ralph Wiggum Loop to ensure the Quality Agent (Critic) catches defects immediately. Feedback is recorded in `FEEDBACK.md` to prevent errors from moving downstream.

### 4. Continual Experimentation & Learning (The Third Way)
Treat implementation as discovery. Use Spikes (OHECI loop) and "Steer As You Go" tactics to turn every build cycle into an experiment. Capture insights in `CHANGELOG.md` and the KEDB to build organizational mastery.

### 5. Iteration Over Perfection
Embrace stateless resampling. High-quality output emerges through persistent iteration, not one-shot perfection. Starting fresh each time prevents "context rot" where errors compound.

### 6. Explicit Uncertainty
Document what you don't know. Don't pretend to have certainty you lack. Use `Feature.md#unknowns` and ADRs to convert risks into documented decisions.

### 7. Avoid the Distracted Agent
Keep agent context focused and expedient. Split roles by purpose (e.g., Build vs. Quality) to prevent context window overload and "hallucinated" reasoning.
*Reference:* [Distracted Agent Anti-Pattern](https://lexler.github.io/augmented-coding-patterns/anti-patterns/distracted-agent/)

### 8. Leverage Orthogonal Biases
Use personas as cognitive filters during Divergence. Generate options independently using different biases (e.g., Troy for pragmatism vs. Marge for empathy) before pooling. This ensures a wider "Tree of Thoughts."

### 9. Documents Are Shared State
Markdown files are the single source of truth. Agents and skills read from and write to these files to maintain context without direct coupling.

---

## 2. The 5-Agent Team (Single Pizza)

We use specialized agents to keep context windows lean and reasoning sharp.

| Agent | Mindset | Focus | Primary Skills |
| :--- | :--- | :--- | :--- |
| **Product** | Empathetic | What & Why | `discovery`, `triage` |
| **Planning** | Logical | How & Structure | `planning`, `architecture` |
| **Build** | Optimistic | Doing | `implementation`, `testing` |
| **Quality** | Pessimistic | Critiquing | `review`, `verification` |
| **Release** | Ceremonial | Shipping | `release`, `learning` |

---

## 3. The Two Diamonds Flow

The protocol coordinates work across two distinct phases of diverging (exploring) and converging (deciding).

### Discovery Diamond (Design Thinking)
-   **Diverge (Investigate):** Product Agent gathers requirements, conducts Five Whys, and Gemba walks.
-   **Converge (Validate):** Planning Agent selects the best *provisional* option from the pool and creates an ADR. This is a hypothesis to be tested.
-   **Output:** A **Feature Brief** with a Selected Option (Hypothesis).

### Delivery Diamond (Agile)
-   **Diverge (Build):** Build Agent implements the hypothesis (Option A).
-   **Converge (Verify):** Quality Agent verifies the implementation works. **Planning Agent** then confirms the option is proven and locks the ADR as "Verified."
-   **Output:** Verified, production-ready code.

---

## 4. The Ralph Wiggum Loop (The Engine)

The core execution engine uses **stateless resampling** to ensure quality.

1.  **Monitor:** Scheduler checks `PLAN.md` for failed or unstarted tasks.
2.  **Spawn:** Agent is spawned in an **Ephemeral Context** (clean worktree).
3.  **Execute:** Agent exercises skills (Build Agent implements; Quality Agent verifies).
4.  **Update:** Documents are updated. If verification fails, the task is marked for a fresh iteration.

---

## 5. Shared State: The 7 Core Documents

Agents and skills maintain alignment through these documents:

1.  **PLAN.md:** The epic-level roadmap and task status.
2.  **TODO.md:** Immediate executable tasks and implementation learning.
3.  **Feature.md:** The intent (Problem, Requirements, Assumptions, Unknowns).
4.  **ADRs:** Architectural Decision Records (The rationale for the "How").
5.  **BDD Specs:** Gherkin scenarios defining acceptance criteria.
6.  **FEEDBACK.md:** Results of reviews and quality gate checks.
7.  **CHANGELOG.md:** Release history and high-level learning capture.

---

## 6. Architecture & Data Flow

```
Problem → [Discovery Diamond] → Feature Brief → [Delivery Diamond] → Release
             (Product/Planning)                    (Build/Quality)      (Release)
```

For a detailed visual guide to these flows, see [docs/reference/visual-diagrams.md](../reference/visual-diagrams.md).

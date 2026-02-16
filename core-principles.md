# The Springfield Protocol: Core Principles

The Springfield Protocol is a framework for assembling agentic skills inspired by the traits of characters from *The Simpsons*. It leverages the **Ralph Wiggum Loop** concept to create a memorable, modular, and iterative environment for Agile Agentic Development.

## Vision & Principles

The goal is to move away from rigid, linear development towards a system that favors iteration, learning, and rapid feedback.

- **Memorable Personas**: Skills are associated with specific character traits to make their roles intuitive for both users and agents.
- **Iteration over Perfection**: Focus on the "stateless resampling" approach where failures are just prompts for the next iteration.
- **Modular & Lean**: Skills align with the [agentskills.io](https://agentskills.io/specification) specification and are designed to minimize context rot.
- **Dual-Track Agility**: Distinct tracks for **Discovery** (Design Thinking) and **Delivery** (Agile).
- **Plan Before You Build, Steer As You Go**: Gather available data to reduce uncertainty before committing to build, but stay open to learning and discovery during execution. Accept incompleteness and iterate.

---

## The Dual-Track Framework

### 1. The Discovery Track (Design Thinking)
*Focus: Creating the "Intent" and "Compass".*

**Philosophy:** Plan before you build—gather available data to reduce uncertainty and make solid decisions. But leave the door open for discovery that informs product design and strategy as you build.

**Stages:**
1. **Intelligence Gathering** - Investigate root problems through interviews, Gemba walks, and pattern research
2. **Narrative Synthesis** - Frame findings into coherent stories that drive decisions
3. **Architecture Validation** - Ensure proposed solutions fit existing patterns and constraints
4. **Stakeholder Gate** - Confirm user needs are met and roadmap fit is clear

**Key Characters**: **Troy** (Intelligence & Pragmatism), **Frink** (Architecture & Patterns), **Marge** (Empathy & Guardrails)

**Deliverables:**
- Discovery Brief (Troy's intelligence synthesis)
- Feature Brief (validated problem statement)

### 2. The Delivery Track (Agile)
*Focus: Implementation and Iteration.*

**Philosophy:** Steer as you go—stay connected to execution, capture learning, and adjust strategy based on evidence from the build. Watch what works, what doesn't, and feed insights back to refine direction.

**Stages:**
1. **Planning** - Translate feature brief into executable epics and tasks
2. **Execution** - Autonomous TDD-driven implementation
3. **Verification** - Ensure quality, coverage, and adherence to standards
4. **Merge & Release** - Gate decisions and publish with confidence

**Key Characters**: **Lisa** (Strategic Planning), **Ralph** (TDD Implementation), **Bart** (Adversarial Review), **Herb** (Quality Assurance), **Lovejoy** (Release Ceremony)

**Flow Pattern:**
```
Feature Brief (Troy/Frink/Marge)
        ↓
    PLAN.md (Lisa creates epics)
        ↓
    TODO.md (tasks for Ralph)
        ↓
    Ralph executes (TDD loop)
        ↓
    Bart reviews (adversarial check)
        ↓
    Herb validates (95%+ coverage)
        ↓
    Marge gates (user fit check)
        ↓
    Lovejoy releases (ceremony)
```

### 3. The Learning Loop (During Delivery)
*Focus: Feeding execution insights back to discovery.*

**Troy stays engaged during delivery:**
- Watches Ralph's implementation for assumption validation
- Picks up signals: "This is harder than expected" or "Users actually need X, not Y"
- Feeds learning back to Lisa/Marge to adjust strategy mid-sprint
- Adjusts plan if unknowns become known problems

This ensures discovery is **not** a one-time gate, but an ongoing cycle.

---

## The Ralph Wiggum Variant: Core Engine

Inspired by Ralph's unique perspective, this "stateless resampling loop" ensures high-quality output through persistent iteration rather than one-shot perfection.

### How It Works

1. **PLAN.md**: The rolling source of truth. Tracks epics and their current validation state.
2. **Control Loop**: The scheduler. It monitors the plan and launches the next execution for the first failed or unstarted task.
3. **Ephemeral Context**: Each iteration starts from a clean slate (using git worktrees or isolated environments) to prevent "context rot" and hallucination accumulation.
4. **Verification Loop**: A post-process "critic" (e.g., Herb or Bart) that validates results and updates PLAN.md.

### Why This Works

- **No Accumulation of Hallucinations**: Fresh context each iteration means errors don't compound
- **Persistent Iteration**: Failures are just signals for the next iteration, not dead ends
- **Visible Progress**: PLAN.md becomes the artifact that shows steady progress toward completion
- **Graceful Degradation**: If an iteration fails, the next one learns from it without cascading problems

---

## Technical Integration

Skills are designed to be installable into `~/.pi/agent/skills/` and follow a standard directory structure for interoperability across agent harnesses.

```bash
~/.pi/agent/skills/
├── troy/          # Discovery & Intelligence
├── frink/         # Architecture & Patterns
├── lisa/          # Strategic Planning
├── marge/         # Empathy & Guardrails
├── ralph/         # TDD Implementation
├── bart/          # Adversarial Review
├── herb/          # Quality Engineering
├── lovejoy/       # Release Ceremony
└── wiggum/        # Triage & Issue Management
```

---

## Feedback & Loop Patterns

The protocol incorporates various agentic patterns to facilitate feedback and continuous improvement:

- **Sense-Plan-Act**: The foundational agent loop (Observe → Think → Act)
- **ReAct (Reason + Act)**: Thought → Action → Observation cycles for tool usage
- **Generate → Evaluate → Critique → Refine**: A multi-stage refinement loop for high-quality output
- **Tree of Thoughts (ToT)**: Exploring multiple reasoning paths and pruning low-scoring options
- **Manager-Worker Loop**: Multi-agent collaboration where a "Manager" (e.g., Lisa) orchestrates specialized "Workers" (e.g., Ralph)

---

## Next Steps

See the following documents for detailed information:

- **[[character-map.md]]** - All 9 characters, their roles, and relationships
- **[[troy-mcclure.md]]** - Detailed profile of Troy, the Chief Discovery Officer
- **[[discovery-track.md]]** - The discovery workflow in detail
- **[[delivery-track.md]]** - The delivery workflow in detail
- **[[ralph-wiggum-loop.md]]** - Deep dive into the core execution engine

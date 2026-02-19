# Getting Started with Springfield Protocol v0.2

## Quick Start (5 minutes)

### 1. Understand the Model
```bash
# Read the core concepts
cat docs/concepts/model.md          # The complete model
cat docs/concepts/architecture.md   # Architecture overview
cat docs/reference/documents.md     # The 7 documents
```

### 2. Try the CLI
```bash
# List available commands
just

# See what agents are available
just agents

# See what skills are available
just skills

# See the phases
just phases
```

### 3. Impersonate an Agent
```bash
# Load the product agent
just impersonate product-agent

# You're now "in" the product agent
# Try a skill:
#   > just utilize discovery-skill
```

### 4. Flow Through Discovery
```bash
# Navigate the discovery diamond
just flow discovery

# This will show the discovery flow
# and guide you to utilize relevant skills
```

---

## Implementation Steps (Full Setup)

### Step 1: Create Your Justfile

Copy this minimal justfile to your project root:

```justfile
# Springfield Protocol v0.2
# Command interface for agents and skills

impersonate AGENT:
    @cat .github/agents/{{AGENT}}.md
    @echo ""
    @bash

utilize SKILL:
    @cat .github/skills/{{SKILL}}/SKILL.md
    @echo ""

flow PHASE="delivery":
    @echo "Flowing through {{PHASE}} phase..."
    @just --list | grep -i {{PHASE}}

loop:
    @echo "🔄 Ralph Wiggum Loop - monitoring PLAN.md"
    @echo "Find next unstarted task and spawn agent"

agents:
    @echo "🎭 Available Agents:"
    @ls -1 .github/agents/*.md 2>/dev/null | xargs -I {} basename {} .md | sed 's/^/  /'

skills:
    @echo "📚 Available Skills:"
    @ls -1d .github/skills/*/ 2>/dev/null | xargs -I {} basename {} | sed 's/^/  /'

phases:
    @echo "🔷 Phases:"
    @echo "  - discovery    (Design Thinking diamond)"
    @echo "  - delivery     (Agile diamond)"
    @echo "  - devops       (Continuous cycle)"

help:
    @echo "Springfield Protocol v0.2"
    @echo ""
    @echo "just impersonate {agent}    Load an agent"
    @echo "just utilize {skill}        Exercise a skill"
    @echo "just flow {phase}           Navigate a phase"
    @echo "just loop                   Run Ralph Wiggum Loop"
    @echo "just agents                 List agents"
    @echo "just skills                 List skills"
    @echo "just phases                 Show phases"

default:
    @just help
```

### Step 2: Create Directory Structure

```bash
mkdir -p .github/skills
mkdir -p .github/agents
mkdir -p docs/adr
mkdir -p features
touch PLAN.md TODO.md CHANGELOG.md
```

### Step 3: Define Core Agents (5-Agent Team)

Create `.github/agents/product-agent.md`:

```markdown
# Product Agent

**Focus:** Discovery & Triage (The "What & Why")
**Skills:** discovery-skill
**Persona:** Marge

**Responsibilities:**
- Investigate user needs
- Define the problem
- Enforce definition of ready
- Create Feature Briefs

**How to use:**
```bash
just impersonate product-agent
```
```

Create `.github/agents/planning-agent.md`:

```markdown
# Planning Agent

**Focus:** Structure & Architecture (The "How")
**Skills:** planning-skill, architecture-skill
**Persona:** Lisa

**Responsibilities:**
- Break features into tasks
- Validate architecture (ADRs)
- Plan dependencies

**How to use:**
```bash
just impersonate planning-agent
```
```

Create `.github/agents/build-agent.md`:

```markdown
# Build Agent

**Focus:** Implementation (The "Doer")
**Skills:** implementation-skill, testing-skill
**Persona:** Ralph

**Responsibilities:**
- Write code
- Write tests (TDD)
- Configure infrastructure

**How to use:**
```bash
just impersonate build-agent
```
```

Create `.github/agents/quality-agent.md`:

```markdown
# Quality Agent

**Focus:** Verification (The "Critic")
**Skills:** review-skill, verification-skill
**Persona:** Bart

**Responsibilities:**
- Adversarial review
- Security checks
- Verify gates (coverage > 95%)

**How to use:**
```bash
just impersonate quality-agent
```
```

Create `.github/agents/release-agent.md`:

```markdown
# Release Agent

**Focus:** Ceremony (The "Shipper")
**Skills:** release-skill, learning-skill
**Personas:** Lovejoy

**Responsibilities:**
- Manage releases
- Update Changelog
- Capture learning

**How to use:**
```bash
just impersonate release-agent
```
```

### Step 4: Define Discovery Skill

Create `.github/skills/discovery-skill/SKILL.md`:

```markdown
# Discovery Skill

**Purpose:** Investigate problems, gather requirements, understand root causes

**When to exercise:** Before building anything new

**Procedure:**

1. Read the issue/request
2. Conduct interviews with stakeholders
3. Perform Five Whys analysis
4. Conduct Gemba walk (examine docs, code, systems)
5. Synthesize findings into Feature.md
6. Document unknowns (linked to ADRs)
7. List explicit assumptions

**Inputs:**
- GitHub issue or request
- Existing Feature.md (if updating)

**Outputs:**
- Feature.md (problem + requirements + constraints)
- Unknowns list (with links to ADRs to be created)
- Assumptions list

**Tools:**
- Interview template (tools/interview-template.md)
- Five Whys template (tools/five-whys.sh)
- Gemba walk checklist (tools/gemba-walk.md)

**Example:**
See examples/discover-user-needs.md
```

Create `.github/skills/discovery-skill/examples/discover-user-needs.md`:

```markdown
# Example: Discover User Authentication Needs

**Scenario:** Users requested login feature

**Discovery Process:**

1. Interview 5 power users
   - Why can't you log in now? (Feature doesn't exist)
   - Why is login critical? (Need multiple accounts per company)
   - What's the minimum needed? (Email/password + session)

2. Five Whys Analysis
   - Why no authentication? (Not prioritized)
   - Why prioritize now? (Enterprise customer feedback)
   - Why is it critical now? (Account isolation needed)

3. Gemba Walk
   - Check: User profile schema (email exists? ✓)
   - Check: Session handling (framework support? ✓)
   - Check: HTTPS enforcement (all traffic encrypted? ✓)

4. Document Feature.md:
   ```
   Problem: Enterprise customers need account isolation
   Requirements: Email/password login with 24h sessions
   Unknowns: Session storage (Redis vs database?)
   Assumptions: Users have email, HTTPS always on
   ```

5. Link unknowns to ADRs to be created:
   - Unknown: Session storage → ADR-001 to be created
```

### Step 5: Create Root Documents

Create `PLAN.md`:

```markdown
# PLAN.md - Feature Roadmap

## Feature: [Your First Feature]

### Epic 1: [Epic Name]
- [ ] Task 1: [Task description]
  - Status: unstarted
  - Assignee: [agent]
  - ADR: [if any]
  - Feature Brief: Feature.md#[anchor]
  - BDD: scenarios.feature#[anchor]
```

Create `TODO.md`:

```markdown
# TODO.md - Sprint Tasks

## Current Sprint

- [ ] Task 1: [Task description]
  - Assigned to: [agent]
  - Time estimate: [hours]
  - Started: [date]
```

Create `Feature.md`:

```markdown
# Feature.md - [Feature Name]

## Problem
[Root cause analysis]

## Requirements
[User need]

## Acceptance Criteria
See: scenarios.feature

## Constraints
- [Hard limit]

## Unknowns
- [Question] - See ADR-XXX for decision

## Assumptions
- [What we're betting on]

## Scope
✅ [In scope]
❌ [Out of scope]
```

### Step 6: Test the CLI

```bash
# List agents
just agents

# List skills
just skills

# Impersonate product agent
just impersonate product-agent

# (You're now in product agent mode)
# Try:
just skills
just flow discovery
just utilize discovery-skill
```

---

## Common Workflows

### Workflow 1: Discover a Feature
```bash
just impersonate product-agent
just flow discovery
just utilize discovery-skill            # Investigate
# Switch to Planning Agent for architecture
just impersonate planning-agent
just utilize architecture-skill         # Validate
# Feature.md + ADRs created
```

### Workflow 2: Implement a Feature
```bash
just impersonate planning-agent
just flow delivery
just utilize planning-skill             # Plan: creates PLAN.md + TODO.md
just loop                               # Ralph Wiggum Loop:
                                        # - Spawn build-agent
                                        # - Exercise skills
                                        # - Spawn quality-agent
                                        # - Verify results
                                        # - Loop until done
```

### Workflow 3: Release
```bash
just impersonate release-agent
just utilize release-skill              # Create CHANGELOG + tag
# Or via GitHub Actions
# Automated release workflow
```

---

## Next Steps

1. **Create your skills** - Start with discovery-skill (`.github/skills/discovery-skill/SKILL.md`)
2. **Create your agents** - Customize agent definitions for your team
3. **Set up workflows** - Add GitHub Actions for automated testing, coverage, security
4. **Train your team** - Share the docs, run through workflows
5. **Iterate** - Refine skills and agents based on experience

---

## Resources

- **Model Overview:** [docs/concepts/model.md](../concepts/model.md)
- **Architecture:** [docs/concepts/architecture.md](../concepts/architecture.md)
- **Documents Reference:** [docs/reference/documents.md](../reference/documents.md)
- **All Skills:** See `.github/skills/*/SKILL.md`
- **All Agents:** See `.github/agents/*.md`

---

## Troubleshooting

**"just command not found"**
- Install justfile: https://github.com/casey/just#installation

**"SKILL.md not found"**
- Ensure `.github/skills/{skill-name}/SKILL.md` exists
- Check skill name matches exactly

**"Agent not found"**
- Ensure `.github/agents/{agent-name}.md` exists

**"Not sure which skill to use"**
- Run `just skills` to see all available
- Run `just flow {phase}` to see recommended skills for that phase

This is your entry point to Springfield Protocol v0.2.

# Springfield Protocol: Character Skills Reference

A detailed guide to each Simpson character skill, their roles in the Springfield Protocol, and how to invoke them.

---

## Discovery Track Characters

### Lisa Simpson - Strategic Planner (flow-planner)
**Archetype:** The Visionary Organizer

**Core Responsibility:**
Translates high-level intent and Feature Briefs into executable plans and task breakdowns.

**Primary Loop:** Plan-and-Execute (PAE)

**Key Skills:**
- Synthesizing complex requirements into coherent plans
- Breaking down epics into manageable tasks with clear dependencies
- Identifying impediments and bottlenecks
- Prioritizing work using Flow Scoring
- Managing scope and preventing scope creep

**Outputs:**
- `PLAN.md` - Rolling implementation plan
- `TODO.md` - Atomic task list with dependencies
- Impediment reports

**When to Use:**
- After a Feature Brief is validated
- When starting a new epic
- When priorities need to shift
- When scope needs tightening

**Interface:**
```
Input: Feature Brief, ADR context, team capacity
Output: PLAN.md with flow scores, TODO.md with tasks
```

---

### Marge Simpson - Empathy & Guardrails (marge)
**Archetype:** The Advocate

**Core Responsibility:**
Ensures that solutions meet real user needs and maintains quality guardrails through the PR-to-Merge flow.

**Primary Loops:** Dialogue, Feedback Accumulation

**Key Skills:**
- User empathy and need validation
- Scope refinement (cutting unnecessary features)
- PR review with a user-centric lens
- Backpressure enforcement (slowing down when needed)
- Updating PLAN.md with field feedback

**Outputs:**
- User feedback synthesis
- PR review comments
- Refined scope for next iteration
- Updated PLAN.md with learnings

**When to Use:**
- Before implementation begins (validate assumptions)
- During PR review (ensure user alignment)
- After implementing a feature (gather feedback)
- When backpressure signals trigger

**Interface:**
```
Input: Feature Brief, User interviews, PR code
Output: Validation report, Review feedback, Updated PLAN.md
```

---

### Frink Simpson - Architect (frink)
**Archetype:** The Technologist

**Core Responsibility:**
Specializes in ADRs, pattern specifications, and ensuring technical coherence and reuse.

**Primary Loops:** OHECI (Observe → Hypothesize → Experiment → Conclude → Iterate), OLA (Knowledge)

**Key Skills:**
- Architecture Decision Records (ADRs)
- Pattern specification and design
- Technical guardrails definition
- Reusability analysis
- Cross-system impact assessment

**Outputs:**
- ADRs documenting key decisions
- Pattern specifications for reuse
- Technical guardrails documentation
- Coherence audits

**When to Use:**
- Before major architectural decisions
- When defining reusable patterns
- When establishing technical standards
- During architecture review phase

**Interface:**
```
Input: Problem statement, Constraints, Existing patterns
Output: ADR, Pattern spec, Guardrails document
```

---

### Product Discovery Agent (product-discovery)
**Archetype:** The Investigator

**Core Responsibility:**
Elicit user requirements, conduct "Gemba walks" through repository documentation, and draft Feature Briefs.

**Primary Loops:** Tree of Thoughts (exploration), Dialogue (user interviews)

**Key Skills:**
- User interview facilitation
- Requirements elicitation
- Competitive analysis
- Repository analysis (docs, ADRs, BDDs)
- Feature Brief authoring

**Outputs:**
- Validated Feature Brief
- User personas
- Market/competitive analysis
- Discovered constraints and opportunities

**When to Use:**
- At the start of a new initiative
- When requirements are unclear
- When validating product-market fit
- Before committing significant resources

**Interface:**
```
Input: User questions, Repo documentation, Market context
Output: Feature Brief (validated), User personas, Analysis
```

---

## Delivery Track Characters

### Ralph Wiggum - TDD Executor (ralph)
**Archetype:** The Tireless Implementer

**Core Responsibility:**
Executes tasks from TODO.md using strict TDD practices and project standards. The workhorse of implementation.

**Primary Loop:** Ralph Wiggum Variant (Stateless Resampling Loop)

**Key Skills:**
- Test-Driven Development (Red-Green-Refactor)
- Code implementation following standards
- Git workflow and atomic commits
- Following project conventions
- Incremental, verifiable progress

**Outputs:**
- Passing test suites
- Implemented features
- Atomic git commits
- PR-ready code

**When to Use:**
- When you have a clear TODO.md
- When you need reliable, tested code
- When quality and coverage matter
- When you want to avoid tech debt

**Interface:**
```
Input: Task from TODO.md, Project standards, Test template
Output: Passing tests + Implementation + PR ready code
Validation: 95%+ coverage, All tests green
```

---

### Bart Simpson - Adversarial Reviewer (bart)
**Archetype:** The Breaker

**Core Responsibility:**
Tries to "break" implementations before formal review. Looks for shortcuts, security holes, lazy code, and missing edge cases.

**Primary Loops:** ReAct (debugging/breaking), GECR (critique phase)

**Key Skills:**
- Security vulnerability detection
- Edge case identification
- Code smell detection
- Testing lazy code finding
- Breaking assumptions

**Outputs:**
- Issues and break reports
- Security recommendations
- Edge case test suggestions
- Code quality feedback

**When to Use:**
- After Ralph completes work (adversarial check)
- Before PR merge
- For security-critical code
- When you need stress-testing

**Interface:**
```
Input: PR code, Test suite, Security context
Output: Break report with issues and suggestions
Pass Criteria: No critical issues, edge cases covered
```

---

### Herb Powell - Quality Engineer (herb)
**Archetype:** The Verifier

**Core Responsibility:**
Enforces 95%+ test coverage, mock-first testing, and zero-change import goals. Ensures quality standards are met.

**Primary Loops:** TALAR (Test → Analyze → Learn → Adjust → Retest), Ralph Wiggum Verification

**Key Skills:**
- Test coverage analysis
- Mock-first methodology
- Import validation (zero-change)
- Quality metrics enforcement
- Coverage gap identification

**Outputs:**
- Coverage reports
- Mock validation reports
- Import compatibility reports
- Quality sign-off

**When to Use:**
- As the final verification step (Ralph Wiggum Verification Loop)
- When coverage requirements are strict
- For import validation
- Before release approval

**Interface:**
```
Input: Code, Test suite, Coverage requirements
Output: Coverage report, Quality certification
Pass Criteria: 95%+ coverage, All mocks valid, Zero-change imports
```

---

### Chief Wiggum - Triage Officer (wiggum)
**Archetype:** The Gatekeeper

**Core Responsibility:**
Manages the "Issue → TODO.md" bridge. Enforces the "Definition of Ready" and performs ecosystem investigation.

**Primary Loops:** OHECI (investigation), Sense-Plan-Act (triage)

**Key Skills:**
- Issue analysis and categorization
- Ecosystem investigation
- Definition of Ready enforcement
- Dependency discovery
- Feasibility assessment

**Outputs:**
- Triaged issue classification
- Ecosystem map
- Definition of Ready checklist
- Issues ready for TODO.md

**When to Use:**
- When an issue arrives
- Before adding to PLAN.md
- When understanding impact scope
- When dependencies are unclear

**Interface:**
```
Input: GitHub issue, Project context
Output: Triaged issue, Ecosystem map, DoR checklist
Pass Criteria: DoR met, Dependencies mapped, Ready for implementation
```

---

## Brownfield & Infrastructure Characters

### Homer Simpson - Brownfield Specialist (homer)
**Archetype:** The Safe Importer

**Core Responsibility:**
Handles safe, zero-change brownfield imports of AWS resources into Takeda Building Blocks.

**Primary Loops:** Plan-and-Execute, Zero-change validation

**Key Skills:**
- AWS resource discovery and analysis
- Terraform module authoring
- State migration planning
- Zero-change import validation
- Rollback planning

**Outputs:**
- Terraform modules for existing resources
- Import plan
- Zero-change validation report
- Runbooks for rollback

**When to Use:**
- When migrating existing AWS infrastructure to Terraform
- When preserving existing resources
- When you need a safety net
- When infrastructure is complex

**Interface:**
```
Input: AWS account, Resource list, Building Block template
Output: Terraform modules, Import plan, Validation report
Validation: Zero-change, Resource mapping, Rollback ready
```

---

## Release & Ceremony Characters

### Reverend Lovejoy - Release Master (lovejoy)
**Archetype:** The Ceremonializer

**Core Responsibility:**
Manages semantic versioning, changelogs, and the "ceremony" of publishing new versions to the TFC Registry.

**Primary Loops:** RDS (Reflect → Document → Share)

**Key Skills:**
- Semantic versioning (semver)
- Changelog generation and curation
- Release notes authoring
- Registry publication
- Version tag management

**Outputs:**
- Updated CHANGELOG.md
- Semantic version tag
- Release notes
- Published artifact in registry

**When to Use:**
- When completing a release-ready feature
- Before publishing to registry
- When managing version history
- For ceremony and communication

**Interface:**
```
Input: Commits since last version, ADRs, Breaking changes
Output: Version tag, CHANGELOG.md update, Release notes
Validation: Semver correct, CHANGELOG accurate, Registry updated
```

---

## Utility & Data Characters

### Document Processing (docx, pdf, pptx, xlsx)
**Archetype:** The Translator

**Core Responsibility:**
Handle creation, reading, and manipulation of various document formats for artifact generation and reporting.

**Key Skills:**
- Word document creation and editing (.docx)
- PDF generation and manipulation
- PowerPoint slide deck creation
- Spreadsheet data handling (.xlsx)
- Format conversion

**Outputs:**
- Professional documents with formatting
- Reports and dashboards
- Presentations
- Tabular data exports

**When to Use:**
- When producing reports or deliverables
- For stakeholder communication
- For data visualization
- For artifact archival

---

### Takeda Building Blocks (takeda-building-blocks)
**Archetype:** The Infrastructure Librarian

**Core Responsibility:**
Discover, explore, understand, and implement Takeda Terraform Building Blocks (standardized AWS infrastructure modules).

**Key Skills:**
- Building Block discovery
- Pattern matching to existing blocks
- Module composition
- Reuse validation
- Custom extension when needed

**Outputs:**
- Building Block recommendations
- Module implementations
- Composition patterns
- Extension suggestions

**When to Use:**
- When designing infrastructure
- To maximize reuse
- To follow Takeda standards
- To avoid custom implementations

---

### Terraform Cloud API (tfc-api)
**Archetype:** The State Inspector

**Core Responsibility:**
Query Terraform Cloud workspaces, runs, plans, and logs for inspection and validation.

**Key Skills:**
- Workspace state inspection
- Run history analysis
- Plan diff reading
- Log analysis
- Cost estimation

**Outputs:**
- State snapshots
- Run reports
- Plan analyses
- Logs and diagnostics

**When to Use:**
- When debugging infrastructure issues
- For state inspection
- When validating TFC runs
- For cost and impact analysis

---

### Known Error Database (kedb)
**Archetype:** The Knowledge Keeper

**Core Responsibility:**
Search, create, and manage Takeda's Known Error Database (KEDB) - a crowdsourced troubleshooting knowledge base.

**Key Skills:**
- Error pattern matching
- Solution discovery
- Error report creation
- Resolution documentation
- Knowledge base maintenance

**Outputs:**
- Error solutions
- Issue reports
- Resolution documentation
- Updated KEDB entries

**When to Use:**
- When encountering Takeda-specific errors
- To discover known solutions
- To document new errors
- To contribute to team knowledge

---

## Character Selection Matrix

| Challenge | Primary Character | Secondary Support | Loop Used |
|:----------|:------------------|:------------------|:----------|
| Understand user need | Product Discovery | Marge | Dialogue, Tree of Thoughts |
| Design architecture | Frink | Lisa | OHECI, OLA |
| Plan implementation | Lisa | Frink | Plan-and-Execute |
| Implement feature | Ralph | - | Ralph Wiggum Loop |
| Find bugs/security | Bart | Herb | ReAct, GECR |
| Verify quality | Herb | Bart | TALAR, Ralph Verification |
| Triage issue | Wiggum | Frink | OHECI, Sense-Plan-Act |
| Import infrastructure | Homer | Herb | Plan-and-Execute |
| Release version | Lovejoy | - | RDS |
| Gather feedback | Marge | Lisa | Dialogue, Feedback Accumulation |
| Find solution in KEDB | - | kedb skill | OLA |

---

## Cross-Functional Workflows

### Feature Implementation End-to-End

```
Product Discovery (Validate need)
    ↓
Lisa (Plan & break down)
    ↓
Frink (Design & ADR)
    ↓
Ralph (Implement with TDD)
    ↓
Bart (Adversarial review)
    ↓
Herb (Verify coverage & quality)
    ↓
Marge (User alignment check)
    ↓
Lovejoy (Release & publish)
```

### Infrastructure Implementation

```
Wiggum (Triage & investigate)
    ↓
Frink (Architecture & ADR)
    ↓
Lisa (Plan implementation)
    ↓
Homer (Brownfield import or new build)
    ↓
Herb (Validate no side effects)
    ↓
tfc-api (Verify state & runs)
    ↓
Lovejoy (Release version)
```

### Issue to Resolution

```
Wiggum (Receive & triage)
    ↓
kedb (Search known solutions)
    ↓
IF found → Use & document
    ↓
IF not found → Product Discovery (analyze)
    ↓
[Normal feature workflow]
    ↓
kedb (Document new solution)
```

---

## Skill Installation & Invocation

All skills are installed in `~/.pi/agent/skills/` and follow the structure:

```
~/.pi/agent/skills/
├── lisa/
│   ├── SKILL.md
│   ├── instructions.md
│   └── templates/
├── ralph/
│   ├── SKILL.md
│   ├── instructions.md
│   └── templates/
└── ...
```

**To invoke a skill in pi:**

```bash
@lisa "Break down this epic into tasks"
@ralph "Implement the feature from TODO.md"
@bart "Review this PR for vulnerabilities"
```

**In other agent harnesses**, load the skill's instructions and persona from the SKILL.md file.

# How to Design Architecture

Step-by-step guide for making and documenting architectural decisions.

## The Workflow

```
Problem/Decision needed
    ↓
@frink: Propose architecture
    ↓
@bart: Adversarial review (poke holes)
    ↓
@frink: Refine based on feedback
    ↓
Create ADR (Architecture Decision Record)
    ↓
Document pattern for reuse
    ↓
Implement with guardrails
```

## Step 1: Propose (@frink)

**Who:** Frink (Architect)

**What to do:**
- Understand the problem deeply
- Research possible approaches
- Sketch 2-3 candidate solutions
- Evaluate against:
  - Existing ADRs and patterns
  - Performance requirements
  - Team capabilities
  - Maintenance burden
- Propose the strongest option
- Document rationale

**Output:**
- Architecture proposal document
- Comparison of alternatives
- Recommended approach

**When you're done:** Proposal is ready for challenge.

---

## Step 2: Adversarial Review (@bart)

**Who:** Bart (Adversarial Reviewer)

**What to do:**
- Challenge the assumptions
- Ask tough questions:
  - What could go wrong?
  - How will this scale?
  - What's the failure mode?
  - Is this consistent with other ADRs?
  - Can a new team member understand this?
- Identify risks
- Suggest improvements
- Play devil's advocate

**Output:**
- List of challenges & concerns
- Questions for Frink

**When you're done:** Frink has things to think about.

---

## Step 3: Refine (@frink)

**Who:** Frink (Architect)

**What to do:**
- Address Bart's concerns
- Update proposal with:
  - Risk mitigations
  - Clarified assumptions
  - Trade-off explanations
  - Implementation guardrails
- Document decisions that changed
- Re-propose if significantly revised

**Output:**
- Refined architecture proposal
- Risk mitigation strategies
- Implementation guidelines

**When you're done:** Proposal is robust.

---

## Step 4: Create ADR

**Who:** Frink (Architect)

**What to do:**
- Create ADR file: `docs/adr/ADR-XXX-title.md`
- Structure:
  ```
  # ADR-XXX: Title
  
  Date: YYYY-MM-DD
  Status: Proposed | Accepted | Superseded
  
  ## Problem
  Why are we making this decision?
  
  ## Decision
  What architecture are we adopting?
  
  ## Rationale
  Why this over alternatives?
  
  ## Consequences
  What are the trade-offs?
  - Benefits
  - Risks & mitigations
  - Maintenance burden
  
  ## Related ADRs
  Links to related decisions
  ```
- Get team approval
- Mark as "Accepted"

**Output:**
- ADR document
- Decision is recorded

**When you're done:** Decision is documented for posterity.

---

## Step 5: Document Pattern

**Who:** Frink (Architect)

**What to do:**
- Extract reusable pattern
- Create pattern doc: `docs/patterns/PATTERN-name.md`
- Include:
  - When to use this pattern
  - How to implement it
  - Code examples
  - Common mistakes
  - Links to relevant ADRs

**Output:**
- Pattern documentation
- Team can reuse this design

**When you're done:** Future developers can apply this pattern.

---

## Step 6: Implement with Guardrails

**Who:** Ralph (TDD Executor) + Frink (oversight)

**What to do:**
- Ralph implements following the ADR
- Frink reviews for:
  - Adherence to ADR
  - Consistency with pattern
  - No shortcuts
- Use tests to encode guardrails
- Document any deviations

**Output:**
- Implementation following architecture
- Tests that prevent violations
- Documentation of decisions

**When you're done:** Architecture is implemented and guarded.

---

## Example: Event Sourcing Decision

### Step 1: Propose (Frink)
```
Problem: We need to track all changes to critical data
Proposal: Use Event Sourcing pattern
Alternatives considered:
- Change logs in database (simpler but harder to query)
- Audit tables (works but duplicates data)
- Event Sourcing (immutable, queryable, complex to learn)
Recommendation: Event Sourcing due to audit trail requirements
```

### Step 2: Adversarial Review (Bart)
```
Questions:
- How do we handle event upcasting when schema changes?
- What's the performance impact at scale?
- How do we query current state efficiently?
- Can a junior dev understand this?
```

### Step 3: Refine (Frink)
```
Mitigation for schema changes: Document event versioning strategy
Performance: Add read model / materialized view
Current state: Define aggregate snapshots at 100 events
Learnability: Include template & examples in pattern doc
```

### Step 4: Create ADR
```
ADR-005: Event Sourcing for Audit Trail

Problem: Need immutable audit trail for critical data
Decision: Implement Event Sourcing pattern
Rationale: Enables compliance, debugging, time-travel queries
Consequences:
  + Complete history preserved
  + Query flexibility
  - Higher complexity, learning curve
  - Need read models for performance
```

### Step 5: Document Pattern
```
PATTERN-event-sourcing.md

When to use: Critical audit trails, temporal queries, debugging
How to implement:
1. Define events (immutable)
2. Store in event log
3. Rebuild state from events
4. Add read model for queries
5. Snapshots for performance

Common mistakes:
- Mutable events (don't do this!)
- No schema versioning (plan for it!)
```

---

## Troubleshooting

**Bart finds fundamental flaws**
- Don't be defensive
- Loop back to Step 1
- Explore alternatives seriously
- It's better to catch it now

**Team doesn't understand the proposal**
- Add diagrams
- Create examples
- Simplify language
- Present in person

**Implementation reveals problems**
- Document as ADR amendment
- May need new ADR if significant
- Don't silently deviate

---

## Files Involved

- `docs/adr/ADR-*.md` - Architecture Decision Records
- `docs/patterns/PATTERN-*.md` - Reusable patterns
- Implementation code
- Tests that guard the architecture

## See Also

- **frink.md** (`.github/agents/`) - Architect profile
- **CHARACTER_SKILLS.md** (root) - Frink's detailed responsibilities
- **VISUAL_REFERENCE.md** (root) - Architecture diagrams

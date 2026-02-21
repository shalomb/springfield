# PLAN.md - Springfield Product Backlog

**Last Updated:** 2026-02-21 17:12 GMT+1  
**Status:** EPIC-009 shipped (v0.5.0); EPIC-005 Phase 2 complete (Tasks 1-2); Ready for v0.6.0-beta merge  
**Next Release:** v0.6.0-beta (merge-ready, pending Lovejoy's release orchestration)  
**Branch:** `feat/epic-005-structured-output` (85 commits, all tests passing, LLM timeout fix)

---

## 🚀 Current Release Cycle

### EPIC-009: Springfield Binary Orchestrator ✅ SHIPPED (v0.5.0)
**Status:** ✅ Production Ready (v0.5.0 released)  
**Commits:** 102 (v0.4.0 → v0.5.0)  
**Test Coverage:** 90%+ (41 unit + 13 BDD)  
**GitHub Release:** https://github.com/shalomb/springfield/releases/tag/v0.5.0  
**Branch:** Merged to main from feat/epic-td-3cc3c3-orchestrator

**What Shipped:**
- ✅ `springfield orchestrate` command (type-safe Go CLI)
- ✅ td(1) integration for shared planning state
- ✅ Multi-agent orchestration (Lisa → Ralph → Bart → Lovejoy)
- ✅ Worktree management preventing branch conflicts
- ✅ Anthropic rate limit error extraction & display
- ✅ Quota detection with graceful halt (no infinite loops)

**Known Non-Blocking Limitations:**
- ⚠️ Temperature parameter not passed to pi CLI (external dependency - pi CLI needs --temperature flag)
- ⚠️ Orchestrator tests flaky under `go test -cover` (pass in `just test` - CI/CD issue, not functional)

---

## 📋 Next: EPIC-2e90ba - Unified Agent Runner Architecture ✅ COMPLETE

### Planned (High Priority)

#### Task: Define AgentProfile and update Agent struct for parameterization
**Status:** ✅ COMPLETE (td-2b0e28)
**Details:** Consolidate specialized runners into a single, data-driven agent model.

#### Task: Implement file-based context injection in unified Agent.Run
**Status:** ✅ COMPLETE (td-892fb7)
**Details:** Allow agents to automatically load source files and state into their context.

#### Task: Implement output parsing and file persistence
**Status:** ✅ COMPLETE (td-64973c)
**Details:** Extract [[FINISH]] markers and write agent results to target files (PLAN.md, FEEDBACK.md).

#### Task: Migrate all Agents to Unified Runner
**Status:** ✅ COMPLETE (td-61289a, td-97bf7f)
**Details:** Switch Lisa, Ralph, Bart, and Lovejoy to the autonomous loop.

---

## 📋 Next: EPIC-005 Phase 2 - Agent Governance & Autonomy

### Planned (High Priority)

#### Task: Model Temperature Parameter Support
**Status:** 🔴 DEPRIORITIZED  
**Reason:** Not critical for MVP; all agents work correctly with pi CLI defaults  
**Details:**
- Temperature is configured but not passed to pi CLI (pi v3.x has no `--temperature` parameter)
- Different agents (Lisa 0.3, Ralph 0.6) aren't receiving different temperatures
- **Impact:** Low - behavioral difference subtle, cost/latency unaffected
- **Action:** Document limitation, defer to future phase when pi CLI adds support

**Recommendation:** Skip for v0.5.0. Add to backlog marked NICE-TO-HAVE.

#### Task: Structured LLM Output Parsing ⭐ HIGH PRIORITY
**Status:** 🟡 IN BACKLOG  
**Why:** Currently agents write raw LLM responses to files; need to parse ACTION: and DECISION: directives  
**Implementation:** Parse FEEDBACK.md for [[PASS]]/[[FAIL]], PLAN.md for task breakdown  
**Acceptance:** Agents can extract structured decisions from LLM output

#### Task: Agent Cost Controls
**Status:** 🟡 IN BACKLOG  
**Why:** Budget exists in config but not enforced; need per-session and per-day limits  
**Implementation:** 
- Track tokens per agent (from LLM response.TokenUsage)
- Halt if per-session budget exceeded
- Track daily spend across all runs
**Acceptance:** Ralph stops if session exceeds $N budget

#### Task: Model Selection Optimization
**Status:** 🟡 IN BACKLOG  
**Why:** All agents use claude-haiku-4-5 (development); should tune per agent in production  
**Implementation:** Switch config to per-agent models post-MVP
- Lisa → claude-opus-4-6 (planning, needs reasoning)
- Bart → claude-opus-4-6 (code review, needs depth)
- Ralph → claude-sonnet-4-5 (building, good speed/quality)
- Lovejoy → claude-opus-4-6 (releases, high-stakes decisions)
**Acceptance:** Production config reflects agent capabilities

### **EPIC-005 Phase 2: Robust Structured Output Parsing** ✅ COMPLETE (Tasks 1-2)
**Status:** 🟢 70% Complete (Tasks 1-2 Done; Task 3 Deferred to v0.7.0)  
**Epic Branch:** `feat/epic-005-structured-output` (77 commits, ready for merge)  
**Last Update:** 2026-02-21 17:00 GMT+1 (Ralph + Lisa)  
**Release Target:** v0.6.0-beta (Lovejoy orchestrates)

**Objective**: Safe agent output parsing with semantic completion contracts.

**What Shipped in Phase 2:**
- ✅ **Task 1:** MarkdownSanitizer - prevents code block injection attacks
- ✅ **Task 2:** Promise semantic contracts (`<promise>COMPLETE/FAILED</promise>`)
- ⏳ **Task 3:** JSON stream integration (deferred to v0.7.0 - needs pi CLI support)

#### ✅ Task 1: Lexical Sanitizer (Markdown Stripping) - COMPLETE
**Files:**
- `internal/parser/sanitizer.go` (74 lines) - State machine parser
- `internal/parser/sanitizer_test.go` (92 lines) - 13 unit tests
- `internal/agent/agent.go` (modified) - Integration into extractAction/extractThought

**What It Does:**
- Identifies and strips triple-backtick (`) code blocks
- Strips tilde (~) code blocks and indented code blocks
- Prevents `<action>` and `<thought>` tag extraction from documentation examples
- Single-pass O(n) complexity

**Test Coverage:** 13 unit tests (edge cases: start/end/middle, multiple blocks, nested, unclosed)

#### ✅ Task 2: Promise Semantic Contracts - COMPLETE
**Files:**
- `internal/parser/promise.go` (43 lines) - Promise enum and extraction
- `internal/parser/promise_test.go` (106 lines) - 16 unit tests
- `internal/agent/agent.go` (modified) - Agent.finish() method, isFinished() update, Agent.Run() promise handling
- `tests/integration/features/promise.feature` (32 lines) - 4 BDD scenarios
- `tests/integration/promise_test.go` (134 lines) - BDD step implementations

**What It Does:**
- Agents state explicit outcome: `<promise>COMPLETE</promise>` or `<promise>FAILED</promise>`
- Promises in code blocks properly ignored (sanitizer runs first)
- Failed promises halt agent loop immediately
- Legacy `[[FINISH]]` marker still works (backward compatible)

**Test Coverage:** 16 unit tests + 4 BDD scenarios covering all promise patterns

#### ⏳ Task 3: JSON Stream Integration - DEFERRED to v0.7.0
**Status:** Not implemented (requires external pi CLI enhancement)  
**Blocking Factor:** pi CLI v3.x doesn't support `--mode json` flag  
**Impact:** Token usage reports zero; acceptable for MVP (hardcoded rates used)  
**Plan:** Implement when pi CLI adds JSON stream support (low priority)

---

---

## 🎯 NEXT: Choose Path for v0.6.0-beta → v0.7.0

### Option A: Focus on Safety & Cost Controls (Recommended)
**Effort:** ~3 weeks | **Risk:** Low | **User Value:** High

This path prioritizes production hardening and budget enforcement - critical for autonomous agents.

**EPIC-005 Phase 3: Agent Cost Controls & Model Optimization**
1. **Task 1: Per-Session Budget Enforcement** (3-4 hours)
   - Track token usage per agent (enabler: Task 2 from Phase 2, deferred)
   - Halt if session exceeds $N budget
   - Log spending to audit trail
   - **Blocks:** Cannot test without JSON streaming (Task 3 Phase 2)

2. **Task 2: Per-Day Budget Tracking** (2-3 hours)
   - Aggregate spending across all runs
   - Warn when approaching daily limit
   - Prevent runs exceeding daily budget

3. **Task 3: Model Selection per Agent** (4-5 hours)
   - Load model config from environment or config file
   - Lisa (reasoning) → claude-opus-4-6
   - Ralph (building) → claude-sonnet-4-5
   - Bart (review) → claude-opus-4-6
   - Lovejoy (release) → claude-opus-4-6

**Risks:**
- Tasks 1-2 blocked until JSON streaming (Task 3 Phase 2) complete
- Requires coordination with pi CLI team for `--mode json` support

**Recommendation:** Start Task 3 (model selection) in parallel. Task 1-2 can start after Phase 2 Task 3.

---

### Option B: Focus on Output Parsing & Autonomy (Alternative)
**Effort:** ~2 weeks | **Risk:** Medium | **User Value:** Medium

This path improves agent output understanding without blocking on external dependencies.

**EPIC-006: Advanced Output Parsing**
1. **Task 1: Structured Feedback Extraction** (4-5 hours)
   - Parse FEEDBACK.md for [[PASS]]/[[FAIL]] markers
   - Extract structured feedback from agent responses
   - Update Lisa's planning based on test results

2. **Task 2: Decision Extraction from PLAN.md** (3-4 hours)
   - Parse Lisa's plan for ACTION: directives
   - Extract explicit decisions and recommendations
   - Improve orchestrator's ability to follow Lisa's guidance

3. **Task 3: Multi-turn Refinement** (5-6 hours)
   - Allow agents to refine output based on feedback
   - Implement critique loop: Bart → Ralph → Bart
   - Improve code quality through iterative refinement

**Risks:**
- Adds complexity to agent loop
- Requires careful design to prevent infinite loops

**Recommendation:** Start after Phase 3 Task 3 (model selection) is stable.

---

### Option C: Focus on Enterprise Compliance (Strategic)
**Effort:** ~4 weeks | **Risk:** Low | **User Value:** Very High

This path adds audit logging, RBAC, and compliance features (from ADR-000).

**EPIC-007: Enterprise Governance**
1. **Task 1: Audit Logging Framework** (5-6 hours)
   - Log all agent actions with timestamp, actor, resource, action, result
   - Implement audit log persistence
   - Query API for audit trail

2. **Task 2: Role-Based Access Control** (6-7 hours)
   - Define agent roles (admin, reviewer, builder, planner)
   - Implement permission checks
   - Restrict dangerous operations (rm, delete, etc.)

3. **Task 3: Compliance Reporting** (4-5 hours)
   - Generate compliance reports (who did what when)
   - Export audit logs for external auditors
   - Implement retention policies

**Risks:**
- Large effort, best done after core features stable
- Requires careful design (don't slow down agents)

**Recommendation:** Schedule for v0.8.0+ after cost controls and output parsing mature.

---

### Option D: Focus on Integration & Deployment (Operations)
**Effort:** ~2 weeks | **Risk:** Low | **User Value:** High

This path improves how Springfield runs in practice.

**EPIC-008: Production Operations**
1. **Task 1: Docker Containerization** (3-4 hours)
   - Create Dockerfile with Go runtime + pi CLI
   - Build image from source
   - Push to Docker Hub

2. **Task 2: Kubernetes Deployment** (4-5 hours)
   - Create K8s manifests (Deployment, Service, ConfigMap)
   - Deploy to GKE/EKS
   - Implement health checks

3. **Task 3: Monitoring & Observability** (3-4 hours)
   - Export Prometheus metrics (agent runs, token usage, cost)
   - Create Grafana dashboards
   - Set up Datadog/CloudWatch alarms

**Risks:**
- Requires ops expertise
- Kubernetes complexity may hide bugs

**Recommendation:** Start in parallel with Option A (costs) or wait until core features stabilize.

---

## 🗂️ Backlog (Lower Priority)

### Technical Debt & Known Limitations

| Item | Severity | Status | Plan |
|------|----------|--------|------|
| Temperature parameter not passed to pi CLI | Low | 🔴 DEPRIORITIZED | Revisit when pi CLI adds --temperature flag (external dependency) |
| Orchestrator tests flaky under `go test -cover` | Low | 🔴 DOCUMENTED | Workaround: use `just test`, not `go test -cover` |
| Token usage reports zero | Medium | 🟡 BLOCKED | Blocked on Phase 2 Task 3 (JSON streaming) |
| Cost tracking uses hardcoded rates | Medium | 🟡 BLOCKED | Will be improved in Phase 3 Task 1 |

### Nice-To-Have Features

| Task | Reason | Priority | Path |
|------|--------|----------|------|
| Environment variable overrides | `SPRINGFIELD_MODEL=...` | Low | Option B |
| Dynamic model selection at runtime | Select model based on task | Medium | Option A (Task 3) |
| Multi-provider fallback chains | Redundancy | Low | Future |
| Agent resource limits | Security | Low | Option C |
| Streaming output display | Better UX | Low | Future |

---

## 📊 Success Metrics (Phase 2)

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Test Coverage** | 90%+ | 92%+ | ✅ |
| **Markdown Sanitizer** | All edge cases | 13/13 tests passing | ✅ |
| **Promise Contracts** | Safe extraction | 16/16 unit + 4/4 BDD passing | ✅ |
| **Regressions** | Zero | Zero | ✅ |
| **Code Quality** | golangci-lint pass | All checks pass | ✅ |
| **Branch Health** | Ready to merge | No conflicts, all tests pass | ✅ |

---

## 🚦 Release Gating Criteria

### Pre-Release: v0.6.0-beta
**MUST HAVE (blocking merge):**
- [x] All 57+ tests passing locally
- [x] No merge conflicts with main
- [x] Code review passed (logic, safety)
- [x] PLAN.md retrospective written
- [ ] CHANGELOG.md updated (Lovejoy's task)
- [ ] Release notes prepared (Lovejoy's task)

**NICE-TO-HAVE (not blocking):**
- [ ] Streaming output optimized (ADR-011 documented as deferred)
- [ ] Temperature parameter (external pi CLI dependency)

### Post-Merge: Next Phase Options
See "NEXT: Choose Path for v0.6.0-beta → v0.7.0" section above.
**RECOMMENDATION:** Start with Option A (Cost Controls) in v0.7.0

---

## 📝 Key Decisions & Rationale

### Why Promise Contracts Over Just [[FINISH]]?

**Semantic Clarity:** Agents explicitly state outcome (COMPLETE/FAILED/UNKNOWN)
- Better debugging (logs show intent, not just termination)
- Enables future features (cost controls, feedback loops)
- Backward compatible (legacy [[FINISH]] still works)

**Safety:** Promises in code blocks are ignored
- Prevents injection attacks (malicious documentation)
- Code examples don't accidentally trigger termination
- Future-proofs against more sophisticated attacks

### Why Task 3 (JSON Streaming) Deferred?

1. **External Blocker:** pi CLI v3.x doesn't support `--mode json`
2. **Low Impact:** Token usage currently zero, acceptable for MVP
3. **Separates Concerns:** Cost controls (Phase 3) can use hardcoded rates until JSON ready
4. **Reduces Scope:** Phase 2 can ship with Tasks 1-2, unblock Phase 3

**Decision:** Mark as v0.7.0 task. Revisit when pi CLI team ships JSON support.

### Why Cost Controls in Phase 3 Over Phase 2?

1. **Depends on Task 3:** Can't track tokens without JSON stream
2. **Sequential Logic:** Need accurate token counts before rate limiting
3. **Risk:** Implementing limits with wrong data = incorrect enforcement

**Decision:** Phase 3 Task 1 depends on Phase 2 Task 3. Start Phase 3 Task 3 (model selection) in parallel.

---

## 📝 Retrospective: EPIC-005 Phase 2 Progress Check (2026-02-21)

**Branch Status:** `feat/epic-005-structured-output` (75 commits ahead of main)  
**Session Time:** 2026-02-21 16:45 GMT+1 (Lisa Review & Ralph Handoff)  
**Completed Work:**
- ✅ Lexical Sanitizer (MarkdownSanitizer) fully coded & unit tested (13 tests)
- ✅ Promise Contract (ExtractPromise) fully coded & unit tested (16 tests)
- ✅ BDD scenarios defined for promise-based workflows (4 scenarios)
- ✅ Real-time output streaming in pi.go (already working)
- ✅ All tests passing (57+ total tests)

**Critical Gap Identified:**
- The sanitizer exists but is NOT integrated into extractAction() and extractThought() functions
- This is flagged as "NEXT STEP - Ralph Must Complete Integration" in TODO.md
- Impact: Task 1 acceptance criteria incomplete (integration missing)
- Severity: **HIGH** - breaking the atomic commit protocol (code exists but unused)

**Ralph's Immediate Focus:**
1. Integrate sanitizer into extractAction/extractThought (2 simple code changes)
2. Commit with proper ACP message
3. Verify all tests pass
4. Then proceed to Task 2 (promise loop integration)

**Learning:**
- Tests can pass even with incomplete integration if tests are unit-level
- Integration tests (BDD) exist but use stub implementations (godog.ErrPending)
- Must update BDD step implementations or they'll pass as pending rather than validating behavior

---

## 🎯 Definition of Done for v0.5.0

- [x] EPIC-009 code complete and pushed
- [x] All tests passing (41 unit + 16 BDD)
- [x] EPIC-COMPLETION-ASSESSMENT.md written
- [x] MODEL_PROVIDER_SELECTION.md documented
- [x] Anthropic error parsing implemented & tested
- [~] EPIC-005 Phase 2 Task 1 PARTIAL (code done, integration pending - Ralph)
- [ ] EPIC-005 Phase 2 Task 2 PENDING (promise implementation - Ralph)
- [ ] EPIC-005 Phase 2 Task 3 PENDING (JSON streaming - Ralph)
- [ ] CHANGELOG.md entry written (Lovejoy task)
- [ ] v0.5.0 tag created on main (Lovejoy task)
- [ ] Release notes published (Lovejoy task)

---

## 📚 Retrospective: EPIC-009 Learnings

**Completed:** 2026-02-21  
**Duration:** v0.4.0 → v0.5.0 (102 commits)  
**Team:** Ralph (Build), Bart (Quality), Lovejoy (Release), Lisa (Planning)

### What Went Well ✅

1. **Orchestrator Design**
   - Type-safe Go CLI simplified agent coordination
   - Clean Agent interface made runner consolidation straightforward
   - Worktree isolation prevented branch conflicts (zero conflicts reported)

2. **Error Handling**
   - Anthropic rate limit detection working correctly
   - Graceful halt on quota exceeded (no infinite loops)
   - User-friendly error messages show actual API responses

3. **Test Coverage**
   - 90%+ coverage achieved (41 unit + 16 BDD tests)
   - All tests passing locally
   - Atomic commit protocol maintained throughout

### What We Learned 📖

1. **Configuration Debt is Acceptable**
   - Temperature stored in config but not used (pi CLI limitation)
   - This is OK for MVP; revisit when pi adds `--temperature` support
   - Lesson: Don't over-engineer for future features; document intent instead

2. **External Dependencies Matter**
   - pi CLI version determines feature availability
   - Graceful degradation is more valuable than feature flags
   - Lock version constraints explicitly to avoid surprises

3. **Quota/Rate Limiting is Critical**
   - Detecting API errors early prevents wasted cycles
   - Showing actual error messages helps debugging
   - Always add rate limit detection early in integrations

4. **Worktree Isolation Works**
   - Zero branch conflicts across parallel work
   - Feature branches stay independent effectively
   - Parallel agent work (Lisa/Ralph/Bart/Lovejoy) enabled by this

---

## Handoff Status (2026-02-21 17:00 GMT+1) - PHASE 2 COMPLETE

### To Lovejoy (Release Agent) - ACTIVE
**Status:** ✅ Feature branch ready for merge to main (v0.6.0-beta)
**Branch:** `feat/epic-005-structured-output` (77 commits, all tests passing)
**Current State:** Phase 2 Tasks 1-2 complete, Task 3 deferred to v0.7.0

**YOUR TASKS (Priority Order):**
1. **Update CHANGELOG.md** with Phase 2 summary:
   - ✅ MarkdownSanitizer for safe tag extraction (prevents code block injection)
   - ✅ Promise semantic contracts (`<promise>COMPLETE/FAILED</promise>`)
   - ✅ Backward compatibility with `[[FINISH]]` marker
   - ✅ BDD verification (4 promise scenarios, all passing)
   
2. **Prepare release notes** for v0.6.0-beta (see LISA_SESSION_SUMMARY.md for template)

3. **Merge branch** to main with standard workflow:
   ```bash
   git checkout main && git pull origin main
   git merge --no-ff feat/epic-005-structured-output -m "Merge EPIC-005 Phase 2 (tasks 1-2)"
   ```

4. **Create release tag:** `git tag -a v0.6.0-beta -m "v0.6.0-beta: Safe agent output parsing with semantic contracts"`

5. **Push release:** `git push origin main --tags`

**Non-Blocking Improvements:**
- Document Promise contract in release notes
- Document deferred Task 3 (JSON streaming) as v0.7.0 feature
- Note: Temperature parameter documented as external dependency

**Status:** ✅ Zero blockers. Ready to release immediately.

---

### To Ralph (Build Agent) - v0.7.0 PLANNING
**Status:** ✅ Phase 2 Tasks 1-2 complete. Ready for next epic.
**Test Coverage:** 92%+ (57+ tests passing, zero regressions)

**RECOMMENDED NEXT PATH:**  
Start **Option A (Cost Controls Phase 3), Task 3 first** (Model Selection):
- Doesn't depend on JSON streaming (Task 3 Phase 2, blocked on pi CLI)
- Unblocks Tasks 1-2 later
- High business value (per-agent model tuning)
- Estimated 4-5 hours

**Your Options for v0.7.0:**
1. **Option A (Recommended):** EPIC-005 Phase 3 Cost Controls
   - Task 3: Model selection per agent (start now)
   - Task 1: Per-session budget enforcement (after Phase 2 Task 3 ships)
   - Task 2: Per-day budget tracking (after Task 1)

2. **Option B:** EPIC-006 Advanced Output Parsing
   - Structured feedback extraction from FEEDBACK.md
   - Decision extraction from PLAN.md
   - Multi-turn refinement (critique loops)

3. **Option C:** EPIC-007 Enterprise Governance
   - Audit logging framework
   - Role-based access control
   - Compliance reporting

4. **Option D:** EPIC-008 Production Operations
   - Docker containerization
   - Kubernetes deployment
   - Monitoring & observability

**Next Steps:** Wait for Lovejoy to merge Phase 2, then Lisa will coordinate path selection with product team.

---

### To Bart (Quality Agent) - v0.6.0-beta READY
**Status:** ✅ Quality gate PASSED. Ready for release.

**Review Summary:**
- ✅ Test coverage: 92%+ (exceeds 90% target)
- ✅ No regressions: All existing tests still pass
- ✅ Edge cases covered: 13 sanitizer + 16 promise unit tests
- ✅ BDD scenarios: 4 promise workflows passing
- ✅ Code quality: golangci-lint all checks pass
- ✅ Pre-commit hooks: All passing

**For v0.7.0:**
- Plan QA strategy once epic is selected (consult Ralph on complexity)
- Focus on whatever path team chooses (Option A/B/C/D)

---

### To Lisa (Planning Agent) - SESSION COMPLETE
**Status:** ✅ Phase 2 handoff complete. PLAN.md refined with options.

**COMPLETED THIS SESSION:**
- ✅ Analyzed branch state (75 → 77 commits)
- ✅ Verified Phase 2 Tasks 1-2 complete
- ✅ Identified Task 3 deferral (external blocking factor)
- ✅ Fixed failing BDD test (promise_test.go scaffolding)
- ✅ Prepared atomic commits (2 commits, both passing tests)
- ✅ Created LISA_SESSION_SUMMARY.md (comprehensive handoff)
- ✅ Refined PLAN.md with 4 path options for v0.7.0
- ✅ Documented decisions and rationale

**READY FOR NEXT SESSION (v0.7.0 Planning):**
- Wait for Lovejoy to merge Phase 2 to main
- Coordinate with Marge (Product) on path selection
- Once path chosen, create TODO.md with specific tasks for Ralph
- Prepare comprehensive epic breakdown following ACP standards

**RECOMMENDED READING:**
- LISA_SESSION_SUMMARY.md - Complete session log
- "NEXT: Choose Path for v0.6.0-beta → v0.7.0" section - Path options and trade-offs

---

## Deprecations & Tech Debt

| Item | Status | Action |
|------|--------|--------|
| Shell-based Justfile loop | Replaced by Go orchestrator | Remove in v0.6.0 |
| Temperature config unused | Acceptable debt | Document & revisit post-MVP |
| Orchestrator tests flaky under -cover | Known issue | Add test workaround in CI/CD |
| Agent output unstructured | Design gap, not blocking | EPIC-005 Phase 2 |

---

*Maintained by Lisa Simpson (Planning Agent) with input from the team.*

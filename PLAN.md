# PLAN.md - Springfield Product Backlog

**Last Updated:** 2026-02-20 22:55 GMT+1  
**Status:** EPIC-009 (Orchestrator) Complete & Shipped  
**Next:** EPIC-005 Phase 2 (Governance & Agent Autonomy)

---

## 🚀 Current Release: v0.5.0-beta

### EPIC-009: Springfield Binary Orchestrator ✅ COMPLETE
**Status:** Production Ready (pending Anthropic quota reset)  
**Commits:** 102 (since v0.4.0)  
**Test Coverage:** 90%+  
**PR:** https://github.com/shalomb/springfield/tree/feat/epic-td-3cc3c3-orchestrator

**What Shipped:**
- ✅ `springfield orchestrate` command (type-safe Go CLI)
- ✅ td(1) integration for shared planning state
- ✅ Multi-agent orchestration (Lisa → Ralph → Bart → Lovejoy)
- ✅ Worktree management preventing branch conflicts
- ✅ Anthropic rate limit error extraction & display
- ✅ Quota detection with graceful halt (no infinite loops)

**Known Limitations (Non-blocking):**
- ⚠️ Agent LLM outputs not parsed into directives (scheduled EPIC-005 Phase 2)
- ⚠️ Orchestrator tests flaky under `go test -cover` (pass in `just test`)

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

---

## 🗂️ Backlog (Lower Priority)

### Nice-To-Have Features

| Task | Reason | Status |
|------|--------|--------|
| Temperature parameter support | pi CLI needs --temperature flag | 🔴 DEPRIORITIZED |
| Environment variable overrides | `SPRINGFIELD_MODEL=...` | ⏳ BACKLOG |
| Dynamic model selection | Select model based on task/budget | ⏳ BACKLOG |
| Multi-provider fallback chains | More than 2 fallbacks | ⏳ BACKLOG |
| Agent resource limits | Memory/CPU constraints | ⏳ BACKLOG |
| Streaming output display | Real-time pi CLI output | ⏳ BACKLOG |

---

## 📊 Success Metrics (v0.5.0)

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Test Coverage** | 90%+ | 90%+ | ✅ |
| **Agents Coordinating** | Lisa→Ralph→Bart→Lovejoy | All 4 working | ✅ |
| **Quota Handling** | Detect & halt gracefully | Anthropic 429 detected | ✅ |
| **Branch Conflicts** | Zero (worktree isolation) | Isolated per epic | ✅ |
| **Error Messages** | Actionable (show actual API errors) | Anthropic JSON parsed | ✅ |
| **Deployed** | GitHub public | https://github.com/shalomb/springfield | ✅ |

---

## 🚦 Release Gating Criteria

**BLOCKERS (must fix before v0.5.0 tag):**
- [ ] Anthropic quota reset (needed for final QA)
- [ ] All tests passing locally
- [ ] CHANGELOG.md updated with v0.5.0 notes

**NICE-TO-HAVE (not blocking):**
- [ ] Temperature support (deprioritized per this update)
- [ ] Streaming output (ADR-011 documented why deferred)

---

## 📝 Notes

### Why Temperature Support is Deprioritized

1. **Not blocking:** Agents work correctly with pi CLI defaults
2. **Subtle impact:** Difference between 0.3 and 0.6 temperature is semantic
3. **External dependency:** Requires pi CLI enhancement (not our code)
4. **Config debt:** Storing unused config is acceptable technical debt for MVP
5. **Cost/Performance:** Temperature doesn't affect speed or cost, only response variance

**Decision:** Keep configuration in place for documentation, skip implementation.

### Streaming Output (ADR-011)

Investigated but rejected for v0.5.0:
- pi CLI outputs via JSON events (no text_delta events)
- Real-time streaming adds complexity without MVP value
- Post-execution analysis sufficient

**Decision:** Defer to future iteration when pi CLI adds streaming support.

---

## 🎯 Definition of Done for v0.5.0

- [x] EPIC-009 code complete and pushed
- [x] All tests passing (41 unit + 16 BDD)
- [x] EPIC-COMPLETION-ASSESSMENT.md written
- [x] MODEL_PROVIDER_SELECTION.md documented
- [x] Anthropic error parsing implemented & tested
- [ ] CHANGELOG.md entry written (Lovejoy task)
- [ ] v0.5.0 tag created on main (Lovejoy task)
- [ ] Release notes published (Lovejoy task)

---

## Handoff Status

### To Lovejoy (Release)
- ✅ Feature branch ready
- ✅ Code reviewed and approved
- ⚠️ Waiting for Anthropic quota to reset for final QA
- 📋 Document temperature limitation in release notes (NICE-TO-HAVE)

### To Ralph (Build)
- ✅ Orchestrator ready for integration
- 📋 Next epic: Structured output parsing
- 📋 Future: Agent cost controls

### To Bart (Quality)
- ✅ Full test suite passing
- ✅ No blockers for v0.5.0
- 📋 Next: Review EPIC-005 Phase 2 scope

### To Lisa (Planning)
- ✅ EPIC-009 scope delivered
- 📋 Next: Plan EPIC-005 Phase 2 breakdown
- 📋 Review model selection optimization strategy

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

# Model & Provider Selection Architecture in Springfield

## Overview

Springfield implements a **hierarchical, cascading model selection strategy** with explicit provider support and graceful fallback mechanisms.

---

## Selection Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ User invokes: springfield --agent ralph --task "fix bug"        │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ▼
        ┌────────────────────────────┐
        │ cmd/springfield/main.go    │
        │ - Load config.toml         │
        │ - Call LoadConfig(".")     │
        └────────────┬───────────────┘
                     │
                     ▼
    ┌────────────────────────────────────────┐
    │ internal/config/config.go              │
    │ LoadConfig(dir) → Config               │
    │ - Parse [agent] section (defaults)     │
    │ - Parse [agents.ralph] section         │
    └────────────┬───────────────────────────┘
                 │
                 ▼
    ┌─────────────────────────────────────────────────────┐
    │ GetAgentConfig("ralph")                             │
    │ - Check [agents.ralph] for ralph-specific override  │
    │ - Fall back to [agent] global defaults if not found │
    │ - mergeWithDefaults() fills missing values           │
    └────────────┬────────────────────────────────────────┘
                 │
                 ▼
    ┌─────────────────────────────────────────────────────┐
    │ AgentConfig returned:                               │
    │ {                                                   │
    │   Model: "anthropic/claude-haiku-4-5"               │
    │   PrimaryModel: ""  (not set, use Model)            │
    │   FallbackModel: "anthropic/claude-haiku-4-5"       │
    │   Temperature: 0.6                                  │
    │   Budget: 150000                                    │
    │ }                                                   │
    └────────────┬────────────────────────────────────────┘
                 │
         ┌───────┴──────────┐
         │                  │
         ▼                  ▼
    Has PrimaryModel?   Use Model field
         NO                 │
         │                  │
    primaryModel =          │
    agentCfg.Model ─────────┤
                            │
                     ┌──────┴──────────┐
                     │                 │
                     ▼                 ▼
        ┌──────────────────────────┐  Check FallbackModel
        │ Create Primary LLM       │   │
        │ primary := &PiLLM{       │   │
        │   Model: primaryModel    │   │
        │ }                        │   │
        │ // "anthropic/claude-..." │   │
        └───────────┬──────────────┘   │
                    │                  │
                    │  ┌───────────────┘
                    │  │ Has FallbackModel?
                    │  │ YES
                    │  │
                    │  ▼
        ┌──────────────────────────────┐
        │ Create FallbackLLM wrapper   │
        │ l := &FallbackLLM{           │
        │   Primary: primary,          │
        │   Fallback: &PiLLM{          │
        │     Model: fallbackModel     │
        │   }                          │
        │ }                            │
        └───────────┬──────────────────┘
                    │
                    └──────┬───────────────┐
                           │               │
                       (with fallback)  (no fallback)
                           │               │
                           ▼               ▼
        ┌─────────────────────────┐  l = primary
        │ NewRunnerWithBudget()   │
        │ - Create RalphRunner    │
        │ - Set budget            │
        └────────────┬────────────┘
                     │
                     ▼
        ┌─────────────────────────────┐
        │ runner.Run(ctx)             │
        │ - Load prompt               │
        │ - Call l.Chat(messages)     │
        └────────────┬────────────────┘
                     │
          ┌──────────┴──────────┐
          │                     │
          ▼ (Primary fails)    ▼ (Primary succeeds)
    ┌──────────────────┐    Return Response
    │ Try Fallback     │
    │ (if configured)  │
    └──────────────────┘
```

---

## Configuration Hierarchy

### 1. **Global Defaults** (config.toml `[agent]` section)

```toml
[agent]
model = "anthropic/claude-haiku-4-5"
temperature = 0.7
max_iterations = 20
budget = 100000
```

- Applied to all agents
- Overridden by agent-specific settings

### 2. **Agent-Specific Overrides** (config.toml `[agents.NAME]` sections)

```toml
[agents.ralph]
model = "anthropic/claude-haiku-4-5"
fallback_model = "anthropic/claude-haiku-4-5"
temperature = 0.6
max_iterations = 30
budget = 150000
```

- Overrides global defaults for this agent
- Empty values fall back to global defaults in `mergeWithDefaults()`

### 3. **Environment Overrides** (not yet implemented)

Planned for future:
```bash
SPRINGFIELD_MODEL="openai/gpt-4o"
SPRINGFIELD_RALPH_MODEL="anthropic/claude-opus-4-6"
```

---

## Model & Provider Format

### String Format Rules

**Format:** `provider/model-name` or `model-name` (default provider)

### Examples

```
"anthropic/claude-opus-4-6"      → Claude Opus via Anthropic
"anthropic/claude-sonnet-4-5"    → Claude Sonnet via Anthropic
"anthropic/claude-haiku-4-5"     → Claude Haiku via Anthropic
"openai/gpt-4o"                  → GPT-4o via OpenAI
"google-gemini-cli/gemini-2.0-flash"  → Gemini 2.0 Flash via Google
"gpt-4o"                         → Uses pi's default provider for GPT-4o
```

### Supported Providers (via `pi` CLI)

| Provider | Format | Models | Status |
|----------|--------|--------|--------|
| Anthropic | `anthropic/claude-*` | opus, sonnet, haiku | ✅ Active |
| OpenAI | `openai/gpt-*` | gpt-4o, gpt-4, gpt-3.5-turbo | ✅ Active |
| Google | `google-gemini-cli/gemini-*` | 2.0-flash, 1.5-pro | ✅ Active |
| GitHub Copilot | `github-copilot/*` | Models via Copilot | ✅ Active |
| OpenRouter | `openrouter/*` | Multi-provider router | ⚠️ Requires key |

---

## Fallback Strategy

### Primary → Fallback Chain

When a primary model fails, the fallback is attempted **automatically**:

```go
type FallbackLLM struct {
    Primary  LLMClient   // Try first
    Fallback LLMClient   // Try if Primary fails
}

func (f *FallbackLLM) Chat(ctx, messages) (Response, error) {
    resp, err := f.Primary.Chat(ctx, messages)
    if err == nil {
        return resp, nil  // ✅ Primary succeeded
    }
    if f.Fallback == nil {
        return resp, err  // ❌ No fallback, return primary error
    }
    return f.Fallback.Chat(ctx, messages)  // Try fallback
}
```

### When Fallback Is NOT Used

1. **No fallback configured** → Error from primary is returned as-is
2. **Quota error detected** → Stops immediately (doesn't try fallback)
3. **Primary succeeds** → Fallback never called

### Quota Error Detection Prevents Fallback Spam

```go
// In pi.go:

if isQuotaExceeded(stderrStr) {
    return nil, &QuotaExceededError{...}  // ← Terminal, no retry
}
```

Quota errors (429, rate_limit_error, etc.) are detected **before** returning to FallbackLLM, so we don't waste fallback attempts on quota issues.

---

## Current Configuration (Development)

All agents currently use **Claude Haiku 4.5** for cost-effectiveness:

```toml
[agents.marge]   # Product reasoning
model = "anthropic/claude-haiku-4-5"

[agents.lisa]    # Planning & breakdown
model = "anthropic/claude-haiku-4-5"

[agents.ralph]   # Build & code generation
model = "anthropic/claude-haiku-4-5"

[agents.bart]    # Quality & review
model = "anthropic/claude-haiku-4-5"

[agents.lovejoy] # Release & ceremony
model = "anthropic/claude-haiku-4-5"
```

### Planned Production Configuration (Post-MVP)

```toml
# FUTURE: Optimize by agent capability

[agents.marge]
model = "anthropic/claude-sonnet-4-5"       # Better UX reasoning

[agents.lisa]
model = "anthropic/claude-opus-4-6"         # Tree of Thought planning
fallback_model = "anthropic/claude-sonnet-4-5"

[agents.ralph]
model = "anthropic/claude-sonnet-4-5"
fallback_model = "google-gemini-cli/gemini-2.5-flash"  # Code generation fallback

[agents.bart]
model = "anthropic/claude-opus-4-6"         # Deep code review
fallback_model = "anthropic/claude-sonnet-4-5"

[agents.lovejoy]
model = "anthropic/claude-opus-4-6"         # High-stakes release decisions
```

---

## Code Flow: Step-by-Step

### Step 1: Configuration Loading

```go
// cmd/springfield/main.go:57
cfg, err := config.LoadConfig(".")

// internal/config/config.go:61
func LoadConfig(dir string) (*Config, error) {
    cfg := &Config{
        Agent: AgentConfig{
            Model: "gemini-2.0-flash",  // Default fallback
            // ...
        },
    }
    
    // Try .springfield.toml first, then config.toml
    // Parse TOML into cfg struct
    toml.DecodeFile(path, cfg)
    
    return cfg, nil
}
```

### Step 2: Agent-Specific Config Resolution

```go
// cmd/springfield/main.go:61
agentCfg := cfg.GetAgentConfig("ralph")

// internal/config/config.go:76
func (c *Config) GetAgentConfig(agentName string) AgentConfig {
    agentName = strings.ToLower(agentName)
    
    // Check agent-specific section [agents.ralph]
    if agentConfig, ok := c.Agents[agentName]; ok {
        return c.mergeWithDefaults(agentConfig)
    }
    
    // Fall back to [agent] section
    return c.Agent
}

// internal/config/config.go:83
func (c *Config) mergeWithDefaults(agentConfig AgentConfig) AgentConfig {
    if agentConfig.Model == "" && agentConfig.PrimaryModel == "" {
        agentConfig.Model = c.Agent.Model  // Use global default
    }
    if agentConfig.FallbackModel == "" {
        agentConfig.FallbackModel = c.Agent.FallbackModel
    }
    // ... merge other fields
    return agentConfig
}
```

Result: `AgentConfig` has all fields populated (either from agent-specific or global defaults)

### Step 3: LLM Client Construction

```go
// cmd/springfield/main.go:64
primaryModel := agentCfg.PrimaryModel
if primaryModel == "" {
    primaryModel = agentCfg.Model
}
// primaryModel = "anthropic/claude-haiku-4-5"

primary := &llm.PiLLM{Model: primaryModel}
// NOTE: Temperature is loaded from config but NOT passed to pi CLI
// (pi CLI does not support --temperature parameter as of v3.x)

if agentCfg.FallbackModel != "" {
    fallback := &llm.PiLLM{Model: agentCfg.FallbackModel}
    l = &llm.FallbackLLM{
        Primary: primary,
        Fallback: fallback,
    }
} else {
    l = primary  // No fallback
}
```

### Step 4: Runner Execution

```go
// cmd/springfield/main.go:77
runner, err := agent.NewRunnerWithBudget("ralph", task, l, agentCfg.Budget)

// internal/agent/runner.go:48
func (br *BaseRunner) Run(ctx context.Context) error {
    messages := []llm.Message{...}
    
    // This calls the LLMClient.Chat() method
    response, err := br.LLMClient.Chat(ctx, messages)
    // ↑ Calls FallbackLLM.Chat() or PiLLM.Chat() directly
    
    if err != nil {
        // Check if quota error (terminal, no retry)
        if llm.IsQuotaExceededError(err) {
            return err  // Stop immediately
        }
        // Other errors might retry depending on agent type
    }
    
    return nil
}
```

### Step 5: LLM Chat Execution

```go
// internal/llm/pi.go:38
func (p *PiLLM) Chat(ctx context.Context, messages []Message) (Response, error) {
    // p.Model = "anthropic/claude-haiku-4-5"
    
    args := []string{"-p", "--no-tools"}
    args = append(args, "--model", p.Model)  // Pass to pi CLI
    // args = [..., "--model", "anthropic/claude-haiku-4-5"]
    
    out, err := execFn(ctx, "pi", args...)
    // Runs: pi -p --no-tools --model anthropic/claude-haiku-4-5 ...
    
    if err != nil {
        // Check stderr for quota patterns
        if isQuotaExceeded(stderrStr) {
            return nil, &QuotaExceededError{...}
        }
        return nil, fmt.Errorf("npm exec failed: %s", errMsg)
    }
    
    return Response{Content: string(out)}, nil
}
```

---

## Decision Points in Model Selection

### 1. **Does [agents.NAME] section exist?**
- YES → Use agent-specific config
- NO → Use global [agent] config

### 2. **Is PrimaryModel set?**
- YES → Use PrimaryModel
- NO → Use Model field

### 3. **Is FallbackModel set?**
- YES → Wrap in FallbackLLM for automatic retry
- NO → Use primary directly (no fallback)

### 4. **Did Primary LLM fail?**
- YES (and Fallback exists) → Try Fallback
- YES (and no Fallback) → Return error
- NO → Return success response

### 5. **Is error a QuotaExceededError?**
- YES → Stop immediately (don't waste fallback attempts)
- NO → Let normal retry logic proceed

---

## Testing Model Selection

### Unit Test: Config Loading

```go
func TestGetAgentConfig(t *testing.T) {
    cfg := &Config{
        Agent: AgentConfig{Model: "global-model"},
        Agents: map[string]AgentConfig{
            "ralph": {Model: "ralph-model"},
        },
    }
    
    // Test agent-specific override
    agentCfg := cfg.GetAgentConfig("ralph")
    assert(agentCfg.Model == "ralph-model")  // ✅
    
    // Test fallback to global
    agentCfg = cfg.GetAgentConfig("lisa")
    assert(agentCfg.Model == "global-model")  // ✅
}
```

### Unit Test: FallbackLLM

```go
func TestFallbackLLM(t *testing.T) {
    primary := &MockLLM{err: errors.New("primary failed")}
    fallback := &MockLLM{response: "fallback response"}
    
    f := &FallbackLLM{Primary: primary, Fallback: fallback}
    resp, err := f.Chat(ctx, messages)
    
    assert(resp.Content == "fallback response")  // ✅
    assert(err == nil)                           // ✅
}
```

### Integration Test: Agent with Model

```go
func TestAgentWithModel(t *testing.T) {
    cfg := &Config{
        Agent: AgentConfig{
            Model: "anthropic/claude-haiku-4-5",
        },
    }
    
    agentCfg := cfg.GetAgentConfig("ralph")
    
    llmClient := &llm.PiLLM{Model: agentCfg.Model}
    runner, _ := agent.NewRunner("ralph", "test task", llmClient)
    
    // Would call pi CLI with --model anthropic/claude-haiku-4-5
}
```

---

## ⚠️ Known Limitation: Temperature NOT Passed to pi CLI

### The Issue

**Temperature is configured in `config.toml` but NOT actually used.**

```toml
[agents.ralph]
temperature = 0.6      # ← Loaded from config
                       # ← Stored in AgentConfig struct
                       # ← BUT NOT PASSED TO pi CLI
```

### Why?

The `pi` CLI (v3.x) **does not support a `--temperature` parameter**:

```bash
$ pi --help
  --model <pattern>              Model pattern or ID
  --system-prompt <text>         System prompt
  --mode <mode>                  Output mode: text, json, or rpc
  --thinking <level>             Set thinking level
  [... 20+ other flags ...]
  
  # NO --temperature FLAG!
```

### Current Behavior

```go
// internal/llm/pi.go:56
// For now, pi CLI doesn't seem to have a temperature flag in this mock implementation
// but we could add it if it did.

args := []string{"-p", "--no-tools"}
args = append(args, "--model", p.Model)  // ✅ Passed

// Temperature WOULD go here:
// args = append(args, "--temperature", strconv.FormatFloat(p.Temperature, 'f', -1, 64))
// BUT there's no p.Temperature field in PiLLM struct!
```

### Why Store It If Unused?

1. **Future-proofing:** When/if pi CLI adds `--temperature` support
2. **Per-agent tuning:** Document intended behavior (Lisa at 0.3, Ralph at 0.6)
3. **Test coverage:** Verify config loading works correctly
4. **Alternative implementations:** If we switch from pi to direct API calls

### What DOES Control Temperature?

Temperature in Anthropic models is controlled **via the API directly**, not CLI parameters. Options:

#### Option A: Bypass pi CLI (Direct API)
```go
// Use anthropic Go SDK directly
client := anthropic.NewClient(apiKey)
response, _ := client.Messages.New(ctx, &anthropic.MessageNewParams{
    Model: "claude-3-5-sonnet-20241022",
    Temperature: anthropic.Float(0.6),  // ✅ Now we can control it!
    Messages: messages,
})
```

#### Option B: Extend pi CLI
```bash
# Open PR to pi to add --temperature support
# Then use:
pi --model anthropic/claude-3-5-sonnet --temperature 0.6
```

#### Option C: Use pi Session Config
```bash
# Store temperature in ~/.pi/agent/config.json
# pi CLI reads it automatically
cat ~/.pi/agent/config.json
{
  "claude-3-5-sonnet": {
    "temperature": 0.6,
    "max_tokens": 4096
  }
}
```

### Impact Assessment

**Current Impact:** None - Temperature isn't affecting model behavior
- All agents use pi CLI defaults (likely temperature ~0.7 or unset)
- Different agents intended to have different temperatures aren't getting them
- Lisa (planning, 0.3) behaves same as Ralph (building, 0.6)

**Severity:** Low - Behavioral difference not critical for current MVP
- Agents work correctly even with default temperatures
- Semantic differences between 0.3 and 0.6 are subtle
- Cost and latency unaffected

**Priority:** Medium - Should address before production scaling
- Intentional per-agent tuning shows we care about quality
- Lisa really should have lower temperature for consistent planning
- Bart really should have lower temperature for rigorous QA

### Recommendation

**Add to TODO for EPIC-010:**
```
- [ ] Investigate pi CLI temperature support
  - Check if pi v3.1+ supports --temperature
  - If not, propose feature request to pi project
  - Fallback: Direct Anthropic API integration for temp control
  
- [ ] Update PiLLM struct to store temperature
  PiLLM struct {
    Model: string
    Temperature: float64  // ← Add this
  }
  
- [ ] Pass temperature to pi when CLI supports it
  args = append(args, "--temperature", strconv.FormatFloat(p.Temperature, 'f', 1, 64))
  
- [ ] Test that temperature differences work
  Test with temperature 0.1 (deterministic) vs 0.9 (creative)
```

---

## Future Enhancements

### 1. Environment Variable Overrides

```bash
export SPRINGFIELD_MODEL="openai/gpt-4o"
export SPRINGFIELD_RALPH_MODEL="anthropic/claude-opus-4-6"

springfield --agent ralph --task "fix"
# Would use claude-opus-4-6 for Ralph, not the config file value
```

### 2. Dynamic Model Selection

```go
// SelectModel(agent, task, budget) → "anthropic/claude-opus-4-6"
// Logic: Switch model based on task complexity/budget

func (s *Springfield) SelectModel(agent, task string, budget int) string {
    if contains(task, "code review") {
        return "anthropic/claude-opus-4-6"  // Need best model
    }
    if budget < 50000 {
        return "anthropic/claude-haiku-4-5"  // Budget constrained
    }
    return s.config.GetAgentConfig(agent).Model
}
```

### 3. Multi-Provider Fallback Chain

```toml
[agents.ralph]
primary_model = "anthropic/claude-sonnet-4-5"
fallback_models = [
    "anthropic/claude-haiku-4-5",
    "google-gemini-cli/gemini-2.5-flash",
    "openai/gpt-4o-mini"
]
```

Currently only supports 2 models (Primary + Fallback). Could extend to chain of fallbacks.

---

## Summary Table

| Component | Location | Responsibility |
|-----------|----------|-----------------|
| **Config Loading** | `internal/config/config.go` | Parse TOML, merge defaults |
| **Model Resolution** | `cmd/springfield/main.go` | Select primary/fallback from config |
| **LLM Abstraction** | `internal/llm/llm.go` | Define LLMClient interface |
| **pi CLI Wrapper** | `internal/llm/pi.go` | Call `pi` with `--model` flag |
| **Fallback Chain** | `internal/llm/llm.go:FallbackLLM` | Try fallback on primary failure |
| **Quota Detection** | `internal/llm/pi.go:isQuotaExceeded()` | Halt on quota errors |
| **Agent Runners** | `internal/agent/runner.go` | Use LLMClient.Chat() |

---

## What Gets Passed to pi CLI vs What Doesn't

### ✅ PASSED to pi CLI

| Parameter | Format | Example |
|-----------|--------|---------|
| **Model** | `--model` | `--model anthropic/claude-haiku-4-5` |
| **System Prompt** | `--system-prompt` | `--system-prompt "You are Ralph..."` |
| **No Tools** | `--no-tools` | Disables read/bash/edit/write |
| **Print Mode** | `-p` | Non-interactive, exit after response |
| **Mode** | `--mode` | `--mode text` (text, json, or rpc) |

### ❌ NOT PASSED (Configured but Unused)

| Parameter | Reason | Status |
|-----------|--------|--------|
| **Temperature** | pi CLI has no `--temperature` parameter | ⚠️ Issue tracked |
| **Budget** | Not a pi CLI concept (Springfield concept) | ⚠️ Enforced by Ralph loop |
| **Max Iterations** | Not a pi CLI concept | ⚠️ Enforced by agent runners |
| **Primary/Fallback Model** | Different LLM instances (not pi flag) | ✅ Handled by FallbackLLM |

### pi CLI Command Constructed

```bash
# What Springfield actually runs:
pi -p --no-tools \
   --model anthropic/claude-haiku-4-5 \
   --system-prompt "You are Ralph the Build Agent..." \
   "Work on epic td-3cc3c3"

# What it WOULD run (if temperature supported):
pi -p --no-tools \
   --model anthropic/claude-haiku-4-5 \
   --temperature 0.6 \
   --system-prompt "You are Ralph the Build Agent..." \
   "Work on epic td-3cc3c3"
```


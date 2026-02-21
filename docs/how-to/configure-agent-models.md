# Configuring Agent Models

Springfield allows you to configure which LLM model and provider each agent uses via `config.toml`. This enables fine-tuned control over model selection, cost optimization, and performance tuning per agent role.

## Quick Start

Edit `config.toml` to set per-agent models with explicit providers:

```toml
[agent]
model = "google-gemini-cli/gemini-2.0-flash"  # Default: provider/model

[agents.lisa]
model = "anthropic/claude-opus-4-1"    # Lisa (planner) gets the strongest model

[agents.ralph]
model = "openai/gpt-4o-mini"           # Ralph (builder) uses a cost-effective model
fallback_model = "google-gemini-cli/gemini-2.0-flash"  # Falls back if primary fails

[agents.bart]
model = "anthropic/claude-3.5-sonnet-20241022"  # Bart (quality) uses a thorough reviewer
```

## Model Specification Format

Models can be specified in two ways:

1. **Model only** (uses default provider): `"claude-opus-4-1"`
2. **Provider and model** (explicit): `"anthropic/claude-opus-4-1"`

The explicit `provider/model` format is **recommended** because many models are available from multiple providers with different capabilities and costs.

### Supported Providers

List available providers and models:
```bash
npm exec @mariozechner/pi-coding-agent -- --list-models
```

**Common providers:**
- `anthropic` — Anthropic Claude models (via API)
- `openai` — OpenAI GPT models (via API)
- `google-gemini-cli` — Google Gemini (via Cloud Code Assist)
- `github-copilot` — GitHub Copilot (subscription)
- `openrouter` — Multi-provider router (single API key for many models)
- `openai-azure` — Azure OpenAI deployment
- `mistral` — Mistral AI models
- `groq` — Groq cloud inference
- And more...

## Configuration Structure

### Global Defaults

```toml
[agent]
model = "google-gemini-cli/gemini-2.0-flash"  # provider/model format
temperature = 0.7                             # Sampling temperature (0.0-1.0)
max_iterations = 20                           # Max loop iterations per task
budget = 100000                               # Token budget per session
```

All agents inherit these defaults unless overridden.

### Per-Agent Overrides

```toml
[agents.{agent_name}]
model = "anthropic/claude-opus-4-1"           # Override model (with provider)
fallback_model = "google-gemini-cli/gemini-2.0-flash"  # Optional fallback
temperature = 0.3                             # Agent-specific temperature
max_iterations = 15                           # Agent-specific iteration limit
budget = 50000                                # Agent-specific budget
```

Replace `{agent_name}` with one of:
- `marge` — Product Agent
- `lisa` — Planning Agent
- `ralph` — Build Agent
- `bart` — Quality Agent
- `lovejoy` — Release Agent

## Recommendations by Role

### **Marge (Product Agent)**
- **Primary:** `anthropic/claude-3.5-sonnet-20241022` or `openai/gpt-4o`
- **Why:** Needs strong reasoning for user empathy and product definition
- **Temperature:** 0.5 (balanced)

```toml
[agents.marge]
model = "anthropic/claude-3.5-sonnet-20241022"
temperature = 0.5
```

### **Lisa (Planning Agent)**
- **Primary:** `anthropic/claude-opus-4-1` (most capable)
- **Why:** Tree of Thought option exploration and Self-Consistency validation are expensive; requires the strongest model
- **Temperature:** 0.3 (focused, deterministic)

```toml
[agents.lisa]
model = "anthropic/claude-opus-4-1"
temperature = 0.3
max_iterations = 10  # Lisa spends more tokens per iteration
budget = 200000
```

### **Ralph (Build Agent)**
- **Primary:** `openai/gpt-4o-mini` or `openrouter/openai/gpt-4o-mini` (cost-effective, fast)
- **Fallback:** `google-gemini-cli/gemini-2.0-flash` (reliable backup)
- **Why:** Fast iteration on TDD; fallback ensures robustness
- **Temperature:** 0.6 (moderate randomness for code variety)

```toml
[agents.ralph]
model = "openai/gpt-4o-mini"
fallback_model = "google-gemini-cli/gemini-2.0-flash"
temperature = 0.6
max_iterations = 30  # Ralph iterates more than others
budget = 150000
```

**Note:** Use `openrouter/openai/gpt-4o-mini` if you don't have direct OpenAI API access but have OpenRouter.

### **Bart (Quality Agent)**
- **Primary:** `anthropic/claude-3.5-sonnet-20241022`
- **Why:** Needs thorough reasoning for code review, test coverage analysis, adversarial testing
- **Temperature:** 0.3 (strict, deterministic)

```toml
[agents.bart]
model = "anthropic/claude-3.5-sonnet-20241022"
temperature = 0.3
budget = 80000
```

### **Lovejoy (Release Agent)**
- **Primary:** `anthropic/claude-opus-4-1` (highly reliable)
- **Why:** Release decisions are high-stakes; needs careful reasoning
- **Temperature:** 0.2 (very conservative)

```toml
[agents.lovejoy]
model = "anthropic/claude-opus-4-1"
temperature = 0.2
max_iterations = 5  # Release ceremony is brief and decisive
budget = 40000
```

## Model Selection Strategy

### Cost Optimization
Use faster, cheaper models for agents that iterate frequently (Ralph). Reserve expensive models for agents that reason deeply once (Lisa, Lovejoy).

```toml
# Fast iteration → cheaper model with provider
[agents.ralph]
model = "openai/gpt-4o-mini"

# Deep reasoning → expensive model with provider
[agents.lisa]
model = "anthropic/claude-opus-4-1"
```

### Provider Selection
Different providers offer the same model with different trade-offs:

| Model | Anthropic | OpenAI | GitHub Copilot | OpenRouter |
|:---|:---|:---|:---|:---|
| Claude Opus | ✅ Native | ❌ No | ✅ Via Copilot | ✅ Routed |
| GPT-4o | ❌ No | ✅ Native | ✅ Via Copilot | ✅ Routed |
| Gemini | ❌ No | ❌ No | ✅ Via Copilot | ✅ Routed |

**Recommendation:** Use OpenRouter (`openrouter/*`) if you want unified access to multiple models with a single API key.

### Reliability and Fallback
For critical agents (Ralph, Lovejoy), configure a fallback model:

```toml
[agents.ralph]
model = "openai/gpt-4o-mini"
fallback_model = "google-gemini-cli/gemini-2.0-flash"  # If GPT fails, try Gemini
```

### Temperature Control
Lower temperature (0.0-0.3) for deterministic tasks (planning, quality review).
Higher temperature (0.5-0.9) for creative tasks (product discovery, code generation).

## Finding Available Models

List all available models on your system:
```bash
npm exec @mariozechner/pi-coding-agent -- --list-models
```

Search for a specific model:
```bash
npm exec @mariozechner/pi-coding-agent -- --list-models sonnet
npm exec @mariozechner/pi-coding-agent -- --list-models gpt-4o
npm exec @mariozechner/pi-coding-agent -- --list-models anthropic/*
```

## Supported Models (Sample)

Pi supports 50+ models across multiple providers. See `npm exec @mariozechner/pi-coding-agent -- --list-models` for the full list.

**Anthropic models:**
- `claude-opus-4-1` (most capable)
- `claude-opus-4-5` (latest)
- `claude-3.5-sonnet-20241022` (strong and fast)
- `claude-3-haiku-20240307` (fast and cheap)

**OpenAI models:**
- `gpt-4o` (most capable)
- `gpt-4o-mini` (cost-effective)
- `gpt-4-turbo` (legacy, still available)

**Google models:**
- `gemini-2.5-pro` (most capable)
- `gemini-2.0-flash` (fast and reliable)
- `gemini-2.5-flash` (latest fast model)

**GitHub Copilot models:**
- Access Claude, GPT, and Gemini via Copilot subscription
- See `github-copilot/*` in model list

## Config File Resolution

Springfield looks for configuration in this order:

1. `.springfield.toml` (if it exists)
2. `config.toml` (fallback)
3. Built-in defaults (if neither exists)

## Example: Full Config

```toml
# Global defaults
[agent]
model = "google-gemini-cli/gemini-2.0-flash"
temperature = 0.7
max_iterations = 20
budget = 100000

# Per-agent configuration
[agents.marge]
model = "anthropic/claude-3.5-sonnet-20241022"
temperature = 0.5
max_iterations = 15
budget = 50000

[agents.lisa]
model = "anthropic/claude-opus-4-1"
temperature = 0.3
max_iterations = 10
budget = 200000

[agents.ralph]
model = "openai/gpt-4o-mini"
fallback_model = "google-gemini-cli/gemini-2.0-flash"
temperature = 0.6
max_iterations = 30
budget = 150000

[agents.bart]
model = "anthropic/claude-3.5-sonnet-20241022"
temperature = 0.3
max_iterations = 20
budget = 80000

[agents.lovejoy]
model = "anthropic/claude-opus-4-1"
temperature = 0.2
max_iterations = 5
budget = 40000

# Sandbox configuration
[sandbox]
image = "docker.io/library/debian:trixie-slim"
image_builder = "podman"
```

## Programmatic Access

In Go code, retrieve agent-specific configuration:

```go
cfg, _ := config.LoadConfig(".")
ralphCfg := cfg.GetAgentConfig("ralph")
fmt.Println(ralphCfg.Model)  // "openai/gpt-4o-mini"
```

The `GetAgentConfig()` method:
- Returns agent-specific settings if configured
- Falls back to defaults for missing values
- Is case-insensitive on agent names
- Returns models in `provider/model-name` format (or just `model-name` if no provider specified)

## Integration with Pi CLI

When passing configuration to the pi CLI, use the model string directly:

```bash
# From config: model = "anthropic/claude-opus-4-1"
# Pass to pi:
npm exec @mariozechner/pi-coding-agent -- --model anthropic/claude-opus-4-1 -p "Your prompt"

# Or use the shorthand (pi infers the provider):
npm exec @mariozechner/pi-coding-agent -- --model claude-opus-4-1 -p "Your prompt"
```

When the Springfield binary invokes pi agents, it reads the model from configuration and passes it via the `--model` flag.

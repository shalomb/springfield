# ADR-011: Streaming Output - Investigation & Conclusion

**Status:** RESOLVED (Rejected Approach, Documented Lessons)  
**Date:** 2026-02-20  
**Deciders:** Ralph (Build), Bart (Quality)  
**Context:** Need real-time transparency for agent execution (streaming output from pi CLI process).  

---

## Problem

When agents run via the pi CLI, there's a perceived "hang" or lack of visibility. Users requested streaming output to help diagnose what's happening during LLM execution.

## Investigation

### Hypothesis (Rejected)
We hypothesized that `pi --mode json` with streaming support would allow us to:
1. Parse JSONL events from pi's stdout
2. Extract `text_delta` events for real-time character accumulation
3. Display a spinner + status line during execution
4. Adapt display based on environment (TTY vs piped)

### Implementation Attempted
- Created `internal/llm/stream.go` with JSONL parser
- Created `internal/llm/display.go` with TTY detection and spinner animation
- Implemented fallback from `pi` binary to `npm exec`
- Added unit tests for stream parsing

### What We Discovered

When running `pi --mode json`, the actual output structure is:
```json
{"type":"session",...}
{"type":"agent_start"}
{"type":"turn_start"}
{"type":"message_start",...}
{"type":"message_update","assistantMessageEvent":{"type":"toolcall_start",...}}
{"type":"message_update","assistantMessageEvent":{"type":"toolcall_end",...}}
{"type":"message_end",...}
{"type":"turn_end"}
```

**Key Finding:** `assistantMessageEvent` contains **tool call events**, NOT **text_delta events**.

The `--mode json` output is designed to capture:
- Tool interactions (read files, execute commands)
- Token usage metadata
- Response structure

It is NOT designed for streaming text character-by-character.

### Blockers

1. **No text_delta events**: The stream never produces the `text_delta` events our parser expected
2. **Tool-centric design**: JSON mode prioritizes tool execution over text streaming
3. **Incomplete buffering**: Even with stderr redirection, the parser hangs waiting for non-existent events
4. **No streaming API**: The pi CLI doesn't expose a dedicated streaming endpoint like `openai/streaming` or Claude's `stream: true`

## Decision

**Reject streaming output feature** for now. Instead:

1. **Use standard text output** (`pi` without `--mode json`)
   - Works reliably
   - Provides real feedback  
   - No parsing complexity
   
2. **Improve observability differently**:
   - Add debug logging with `DEBUG=1` (already implemented in LLM layer)
   - Show LLM call details in verbose mode
   - Log token counts and timing to FEEDBACK.md
   - Document what's happening in agent output

3. **Keep npm exec fallback**
   - Works with both `pi` binary and fallback
   - No changes needed to executor logic

## Lessons Learned

### Technical
- **Real-time streaming requires a dedicated protocol**, not JSON event extraction
- **Tool-centric APIs** (like `--mode json`) aren't suitable for text streaming
- **Assumptions about output format matter**: We assumed `text_delta` without verifying

### Process
- **Verify interfaces before building**: Check actual CLI output before designing a parser
- **Prefer simplicity over features**: Plain text works, JSON complexity wasn't worth the gain
- **Debug logging beats custom parsing**: Structured logs are more valuable than streaming UI

## Alternative Approaches Considered

### Option A: Real-time Streaming (REJECTED - This ADR)
- Status: Blocked by pi CLI design
- Cost: 3 full implementations (stream.go, display.go, parser)
- Benefit: Real-time text visibility
- Conclusion: pi CLI doesn't support it

### Option B: Post-Execution Log Analysis
- Parse stdout after execution completes
- Extract key information (tokens, actions, errors)
- Add to FEEDBACK.md
- **Status:** RECOMMENDED for future work
- Cost: Medium
- Benefit: Simple, testable, no hanging

### Option C: Custom Streaming Wrapper
- Build our own streaming abstraction around pi
- Use goroutines to read lines and send updates
- Display spinner while waiting
- **Status:** VIABLE but complex
- Cost: High (requires concurrency primitives)
- Benefit: Works with any CLI
- Recommendation: Defer to EPIC-005 (Governance) when needed for detailed observability

## Action Items

- [x] Clean up failed streaming code
- [x] Document learnings in this ADR
- [ ] Implement Option B (post-execution log analysis) in EPIC-005 Phase 2
- [ ] Consider Option C (custom streaming wrapper) for future release cadence visibility

## Files Modified

**Reverted:**
- `internal/llm/stream.go` - JSONL parser (deleted)
- `internal/llm/display.go` - TTY detection & spinner (deleted)
- `internal/llm/stream_test.go` - Parser tests (deleted)
- Streaming-related changes to `internal/llm/pi.go`

**Kept:**
- Basic pi CLI integration (no `--mode json`)
- npm exec fallback
- Error handling for "not found"

## Related Issues

- **EPIC-005**: Agent Governance (planned Phase 2 work on observability)
- **TODO**: Add structured logging to LLM calls for better debugging
- **FEEDBACK**: Provide template for agent output analysis

---

**Outcome:** The system works reliably without streaming. Real-time output adds complexity without sufficient benefit given pi's output design. Future observability improvements will come through debug logging and post-execution analysis.

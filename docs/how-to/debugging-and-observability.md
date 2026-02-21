# Debugging and Observability

This guide explains how to troubleshoot Springfield and understand what's happening during agent execution.

## Real-Time Debug Logging

Springfield uses [sirupsen/logrus](https://github.com/sirupsen/logrus) for structured logging. Enable debug output with the `DEBUG` environment variable.

### Enable Debug Logging

```bash
# Run with debug output
DEBUG=1 ./bin/springfield --agent ralph

# Or via Justfile
DEBUG=1 just ralph
```

### Log Output Format

```
time="21:53:38.714" level=debug msg="Starting LLM call with 2 messages" ctx=PiLLM.Chat
time="21:53:38.714" level=debug msg="  Message 0 (role=system): 817 chars" ctx=PiLLM.Chat
time="21:53:38.714" level=debug msg="Executing pi CLI..." ctx=PiLLM.Chat
time="21:53:38.714" level=debug msg="Attempting to execute: pi with 5 arguments" ctx=executorWithFallback
```

Fields:
- **time**: Timestamp with millisecond precision
- **level**: Log level (debug, info, warning, error)
- **msg**: Human-readable message
- **ctx**: Context/module name (function or package)

## What Gets Logged

### LLM Execution (internal/llm/pi.go)

When you enable `DEBUG=1`, you'll see:

1. **LLM Call Start**
   ```
   msg="Starting LLM call with 2 messages" ctx=PiLLM.Chat
   msg="  Message 0 (role=system): 817 chars" ctx=PiLLM.Chat
   msg="  Message 1 (role=user): 4 chars" ctx=PiLLM.Chat
   ```
   - Shows message count and content sizes
   - Helps identify if system prompts are being passed correctly

2. **Executor Selection & Fallback**
   ```
   msg="Attempting to execute: pi with 5 arguments" ctx=executorWithFallback
   msg="pi not found in PATH, falling back to npm exec" ctx=executorWithFallback
   msg="Executing: npm exec @mariozechner/pi-coding-agent with 5 arguments" ctx=executorWithFallback
   ```
   - Shows whether `pi` binary is available
   - Confirms fallback to npm is working
   - Indicates argument count for troubleshooting

3. **Execution Results**
   ```
   msg="Successfully executed pi, got 4821 bytes" ctx=executorWithFallback
   msg="npm exec succeeded, got 5034 bytes" ctx=executorWithFallback
   msg="After filtering npm output: 4891 bytes" ctx=executorWithFallback
   ```
   - Shows output size before and after npm warning filtering
   - Helps diagnose if output is being truncated

4. **Errors**
   ```
   level=error msg="pi CLI execution failed" ctx=PiLLM.Chat err="exit status 127"
   level=error msg="npm exec failed" ctx=executorWithFallback err="not found"
   ```
   - Shows error type and context
   - Exit status helps identify the failure

## Troubleshooting Common Issues

### Hang During Execution

**Symptom:** Command starts but never returns

**Diagnosis:**
```bash
DEBUG=1 timeout 10 ./bin/springfield --agent ralph
```

**Expected sequence:**
```
[21:53:38.714] msg="Attempting to execute: pi..." 
[21:53:38.714] msg="Successfully executed pi, got 1234 bytes"  # or npm exec
[21:53:38.715] msg="LLM call completed. Response: 1200 chars"
```

**If hang occurs after "Attempting to execute":**
- Check if `pi` binary is in PATH: `which pi`
- Check if `npm` is available: `which npm`
- Verify pi/npm credentials are configured
- Check disk space and memory

### Empty or No Output

**Symptom:** Agent runs but produces no response

**Diagnosis:**
```bash
DEBUG=1 ./bin/springfield --agent bart --task "test" 2>&1 | grep "Response:"
```

**Look for:**
```
msg="LLM call completed. Response: 0 chars"  # Response is empty
```

**Solutions:**
- Verify pi CLI is working: `pi -p "hello"` or `npm exec @mariozechner/pi-coding-agent -- -p "hello"`
- Check system prompt is being passed
- Verify LLM API credentials (ANTHROPIC_API_KEY, OPENAI_API_KEY, etc.)

### Fallback Not Working

**Symptom:** Error says "pi not found" but npm should work

**Diagnosis:**
```bash
DEBUG=1 ./bin/springfield --agent ralph 2>&1 | grep -E "fallback|npm"
```

**Expected:**
```
msg="pi not found in PATH, falling back to npm exec"
msg="npm exec succeeded, got 5034 bytes"
```

**If no fallback happens:**
- Check npm is installed: `npm --version`
- Check pi package is available: `npm list -g @mariozechner/pi-coding-agent`
- Verify PATH includes npm: `echo $PATH`

## Integration with Log Aggregation

The logrus output format can be easily parsed by log aggregation systems:

### JSON Format

To output JSON (for ELK, Loki, Datadog, etc.), modify `internal/llm/debug.go`:

```go
func init() {
    log.SetFormatter(&log.JSONFormatter{
        TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
    })
    // ... rest of config
}
```

Then use with aggregation tools:
```bash
DEBUG=1 ./bin/springfield --agent ralph 2>&1 | jq .
```

### Filtering Logs

```bash
# Only debug logs
DEBUG=1 ./bin/springfield ... 2>&1 | grep 'level=debug'

# Only errors
DEBUG=1 ./bin/springfield ... 2>&1 | grep 'level=error'

# Specific module
DEBUG=1 ./bin/springfield ... 2>&1 | grep 'ctx=executorWithFallback'

# Time range
DEBUG=1 ./bin/springfield ... 2>&1 | grep '21:5[3-9]'
```

## What's NOT Logged Yet

These could be added in the future (EPIC-005 Phase 2):

- **Token usage**: The pi CLI doesn't expose token counts currently
- **Timing**: Duration of LLM calls
- **Memory usage**: Process resource consumption
- **Agent decisions**: Ralph's iteration loop, Lisa's planning decisions
- **File I/O**: What files were read/written during agent execution

See [ADR-011: Streaming Output Investigation](../adr/ADR-011-streaming-output-discovery.md) for more on observability design decisions.

## Adding Custom Logging

To add logging to new code:

```go
package mypackage

import log "github.com/sirupsen/logrus"

func MyFunction() {
    logger := GetLogger("MyFunction")
    logger.Debugf("Starting with %d items", len(items))
    
    if err != nil {
        logger.WithError(err).Errorf("Failed to process")
        return
    }
    
    logger.Debugf("Success!")
}
```

Remember:
- Import logrus in `internal/llm/debug.go` (already done)
- Use `GetLogger("FunctionName")` to get a logger
- Use `Debugf()` for debug messages
- Use `WithError(err).Errorf()` for errors
- The `ctx` field will be auto-populated with your function name

## Related Documentation

- [ADR-011: Streaming Output Investigation](../adr/ADR-011-streaming-output-discovery.md)
- [PLAN.md - EPIC-005 Phase 2](../../PLAN.md) - Planned observability improvements
- [sirupsen/logrus Documentation](https://github.com/sirupsen/logrus)

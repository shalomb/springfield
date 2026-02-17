# Agentic Loops Reference

A quick lookup guide for the 16+ feedback loops used in Springfield Protocol.

## Loop Selection

| Loop | Use When | Complexity |
|------|----------|-----------|
| **Sense-Plan-Act** | Real-time decisions | Low |
| **ReAct** | Debugging with reasoning | Medium |
| **Tree of Thoughts** | Complex multi-path decisions | High |
| **Plan-and-Execute** | Clear specs, sequential work | Medium |
| **Ralph Wiggum** | Feature delivery with quality gates | Medium |
| **GECR** | Polishing output | Medium |
| **TALAR** | Test-driven optimization | Medium |
| **Manager-Worker** | Multi-agent coordination | Medium |
| **Dialogue** | Two-agent iteration | Low |

## Decision Tree

```
What's your problem?

├─ Vague/complex problem
│  └─ Tree of Thoughts
│
├─ Specific error with reasoning
│  └─ ReAct
│
├─ Clear spec, sequential work
│  └─ Plan-and-Execute
│
├─ Feature delivery with gates
│  └─ Ralph Wiggum
│
├─ Polish & refine output
│  └─ GECR
│
├─ Optimize via testing
│  └─ TALAR
│
├─ Multiple agents in parallel
│  └─ Manager-Worker
│
└─ Two agents iterating
   └─ Dialogue
```

## Full Specifications

See [`../LOOP_CATALOG.md`](../LOOP_CATALOG.md) in the root for complete technical specifications with diagrams.

## Related

- **QUICK_START.md** (root) - Loop quick reference
- **concepts/ralph-wiggum-loop.md** - Conceptual explanation
- **how-to/** - Workflows using these loops

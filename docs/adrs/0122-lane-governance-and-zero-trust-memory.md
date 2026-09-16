# ADR 0122: Lane Governance & Zero-Trust Memory

## Context
As CognitiveOS manages local AI workflows across multiple projects, there is a severe risk of "Context Bleeding". An AI assisting with a corporate codebase (`lane: corp-xyz`) might inadvertently leak proprietary code or API keys when the user switches to a personal project (`lane: personal`), because standard AI tools merge context globally. Conversely, isolating projects completely prevents the AI from carrying over learned, generalized skills (like design patterns).

## Decision
We will implement a **Lane Governance** system at Layer 1 (Substrate) functioning as a "Context Firewall". 
1. **Physical Isolation:** Each workspace maintains its own `.cog/` local ledger, while a shared `~/.cog_global/` acts as the `commons` lane.
2. **Pre-scoring Filter:** The Context Engine will strictly filter out any `CogDoc` that does not match the active `lane_id` or `commons` *before* the Mamba SSM salience scoring runs.
3. **Abstraction Bridge:** To share knowledge safely, the Lexis Engine (Layer 3) will be responsible for extracting abstract patterns from private lanes, anonymizing them, and proposing them to the `commons` lane (subject to user approval).

## Consequences
**Pros:**
- **Zero Data Leakage:** Cryptographically and physically bounds the context window.
- **Model Hot-swapping:** The OS holds the state, allowing seamless switching between Claude, Gemini, or local Ollama models without losing workspace context.
- **Cognitive Evolution:** Skills and patterns can be transferred safely via the `commons` lane without carrying the proprietary business logic.

**Cons:**
- **Complexity:** Managing multi-lane ledgers and global vs. local mounts adds overhead to the daemon's boot and reconciliation loops.
- **Quarantine Overhead:** Anonymization requires a robust "Human-in-the-loop" approval process to prevent accidental leaks during promotion to `commons`.

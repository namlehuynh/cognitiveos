# CognitiveOS Agent Proxy Rules

> This is a Proxy file to inform AI Agents (Cursor, Antigravity, Github Copilot) about the context and rules of the project.

## 🚨 To AI Agents:
You are working in the **CognitiveOS** project. The context, day/night shift workflows, and roadmaps are **not** located in the root directory to keep the codebase clean.

Instead, the local "brain" for AI is located in the hidden `.cdlc/` directory.

### Please read the following files before starting any task:
1. Read `.cdlc/ROADMAP.md` to know which features are assigned to you.
2. Read `.cdlc/REVIEW_QUEUE.md` to check for any pending issues/blocks.
3. Read `.cdlc/project-rules.md` for project-specific rules.
4. If you are assigned to a Night Shift (`Window: REM`), ensure all error logs are updated in `REVIEW_QUEUE.md` instead of popping up questions to the user.
5. Read `agents.md` in the parent `notes-workspace` directory to master the **CDLC** and **Atomic Commits** rules.

- alt+del will delete the word or the surrounding braces if any
- ai integration
- ctrl-space -> cmd panel
- convinient shortcut to delete between scopes? [cursor](this wil be deleted)

- add an advanced mode to override core configurations
	- if its active show the info in the statusbar, A all the way to the right
	- if you click on it, it shows you which core functionality has been changed

- plugins:
	execute external go/c/whatever programs with context


# QOL things
- smear cursor by default
	https://github.com/sphamba/smear-cursor.nvim/tree/main/lua/smear_cursor
	https://github.com/karb94/neoscroll.nvim



## Multiplayer Editing via SSH — Concept

Terminal-native, multiplayer-aware editing in kitsune. No cloud, no GUI — just SSH.

- When users SSH into the same machine and launch kit, they appear in a shared session list.
- Each open file has a "master" user — only they can edit.
- Others can view the file live, chat inline, or request control.
- The master can grant/reject control instantly.
- Chat is scoped per file or global.
- All edits are versioned per user and can be reviewed, merged, or replayed.
- Each user has their own config (.kit.lua) — color, name, macros.
- Optional: show presence (ghost cursors, file activity), per-user history, and command audit.

TL;DR:
Collaborative terminal editing with scoped control, version trace, and real-time presence — powered by SSH. Like tmux + vim + git + Slack, minus the bloat.

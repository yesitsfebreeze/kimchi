# Kitsune

<img src=".doc/logo.svg" alt="Logo" height="60">

**Kitsune** is a terminal-native, modalless code editor written in Go — designed for speed, scriptability, and minimalism.

It’s *insert-first*, with fast Ctrl-based commands and advanced leader key strokes for power tools. Designed to be instantly usable but deeply configurable.

## Philosophy

- Insert-first: typing is immediate
- Commands via Ctrl binds (e.g. Ctrl+S to save)
- Advanced actions via leader strokes (e.g. `ff` for fuzzy find)
- Zero-modality: no `insert/normal` mode toggling
- Lua-configurable with project awareness

## Example Config (`kitsune.lua`)

```lua
cfg('indent.saved.style', 'spaces')
cfg('indent.visual.style', 'tabs')
cfg('statusbar.position', 'top')

-- some binds
bind('SaveAll', 'ctrl-shift-s')
bind('Prompt', 'ctrl-space')

-- ctrl-space then tap ff -> FuzzyFind
prompt('FuzzyFind', 'ff')

style('text', '#ffffff', '#000000')
```

Usage:
```shell
kit /etc/hosts
```

#### Kitsune:
Fox spirits that grow additional tails (up to nine) as they age and gain power. With time, they can become celestial beings or take permanent human form.

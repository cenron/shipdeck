---
name: systematic-debugging
description: Use when a bug, failing test, or unexpected behavior needs root-cause analysis before proposing a fix
---

# Systematic Debugging

## Overview

Investigate the actual failure before changing code. Prefer evidence, reproduction, and narrowing the problem over guessing.

## Rules

- Reproduce the issue first.
- Read logs and failing output carefully.
- Narrow the scope with the smallest useful test or command.
- Check the most likely boundary failures before changing internals.
- Verify the fix against the original failure.

## Common Mistakes

- Jumping to a fix before reproducing.
- Changing several things at once.
- Trusting memory over actual output.

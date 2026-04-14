# authkit

Team and organization management.

## Overview

The `authkit` package provides team and organization management primitives
inspired by Upstream AuthKit.

**Module:** `github.com/gocanto/bedrock/packages/authkit`

```bash
go get github.com/gocanto/bedrock/packages/authkit@latest
```

## Features

- Create and delete teams
- Invite users to a team via email with invitation tokens
- Assign per-team roles to members
- Transfer team ownership
- Multi-tenancy helpers that scope queries to the currently active team

## Concepts

| Concept      | Description                                                |
|--------------|------------------------------------------------------------|
| `Team`       | An organizational unit with an owner and members          |
| `Membership` | A user–team association with an assigned role             |
| `Invitation` | A pending email invite to join a team                     |
| `Role`       | A named set of permissions scoped to a team               |

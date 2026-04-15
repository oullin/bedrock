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

| Concept      | Description                                      |
| ------------ | ------------------------------------------------ |
| `Team`       | An organizational unit with an owner and members |
| `Membership` | A user–team association with an assigned role    |
| `Invitation` | A pending email invite to join a team            |
| `Role`       | A named set of permissions scoped to a team      |

## Usage

```go
actions := authkit.NewActions(teamRepo, memberRepo, inviteRepo, mailer)

// Create a team
team, err := actions.CreateTeam(ctx, owner, "Acme Corp")

// Invite a member
invite, err := actions.InviteTeamMember(ctx, team, "alice@example.com", "editor")

// Accept an invitation (called after user clicks the email link)
err = actions.AcceptInvitation(ctx, invite)

// Remove a member
err = actions.RemoveTeamMember(ctx, team, member)

// Update a member's role
err = actions.UpdateTeamMember(ctx, team, member, "admin")

// Transfer ownership
err = actions.TransferOwnership(ctx, team, newOwner)

// Delete a team
err = actions.DeleteTeam(ctx, team)
```

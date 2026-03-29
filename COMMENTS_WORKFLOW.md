# Comments Workflow

This document explains how comments work across routes, handlers, storage, SQL queries, and templates.

## Overview

The comments system supports:
- Top-level comments on show pages.
- Replies to existing comments.
- Thread display using depth values returned from a recursive SQL query.

Comments are currently used on:
- Writing show pages
- Micro show pages
- Books show pages

## End-to-End Flow

1. User opens a show page.
2. Show handler loads the entry and attaches comments.
3. Template renders existing comments plus a form for new comments.
4. User submits a top-level comment or reply form.
5. POST handler validates input and parent comment relationship.
6. Store writes the comment to SQLite.
7. User is redirected back to the page anchor at #comments.
8. On reload, comments are queried again and rendered in thread order.

## Routing

Each section that supports comments has a POST route:
- POST /writing/{slug}/comments
- POST /micro/{slug}/comments
- POST /books/{slug}/comments
- POST /projects/{slug}/comments

All these routes call the same handler factory with a section argument.

## Handler Responsibilities

The sectionAddComment handler does the following:

1. Reads slug from the route path.
2. Parses form values.
3. Reads name and content.
4. Reads optional parent_id.
5. If parent_id is present:
- Parses it as int64.
- Rejects non-positive or invalid values.
- Verifies that the parent comment exists on the same section and slug.
6. Calls store.AddComment.
7. Redirects to /{section}/{slug}#comments.

## Store Responsibilities

### Reading comments for a page

commentsForEntry builds a canonical slug path in this format:
- /{section}/{slug}

Then it queries comments with that canonical slug. If no rows are found, it tries a fallback using just slug to support legacy rows created before canonical normalization.

### Adding a comment

AddComment currently enforces:
- Datastore must be available.
- Name is trimmed and defaults to Anonymous if empty.
- Content is trimmed and must be non-empty.
- Slug is normalized to /{section}/{slug} before insert.
- parent_id is optional; when present it is stored as a valid nullable int64.

### Parent safety check

CommentExists checks whether the provided parent comment ID appears in the comment set for the same page. This prevents cross-page reply linkage.

## Database Model

SQLite table:
- id INTEGER PRIMARY KEY AUTOINCREMENT
- name TEXT NOT NULL
- content TEXT NOT NULL
- slug TEXT NOT NULL
- created_at DATETIME DEFAULT CURRENT_TIMESTAMP
- parent_id INTEGER REFERENCES comments(id)

parent_id is NULL for top-level comments and set for replies.

## SQL Query for Threading

GetCommentsBySlug uses a recursive CTE:

1. Base case selects root comments where parent_id IS NULL and slug matches.
2. Recursive case joins children on parent_id = parent row id.
3. depth increments by 1 per level.
4. path is constructed with zero-padded IDs via printf('%020d', id).
5. Final ordering is ORDER BY path, created_at.

This yields deterministic parent-first ordering with stable child grouping.

## How Flat SQL Rows Become a Threaded UI

If you come from React, it is natural to expect a nested JSON shape like:
- comment
- children: []comment

This implementation intentionally does not build a nested tree in Go.

Instead, SQLite returns a flat list that already contains enough metadata to render a thread:
- id: unique comment id
- parent_id: null for root, set for replies
- depth: visual nesting level computed by the recursive CTE
- path: hierarchical sort key so parents and descendants stay adjacent

Think of it as a pre-ordered traversal stream. The list is flat, but order + depth carry the tree semantics.

### Why this works for server-rendered templates

The Go template loops once over .Data.Entry.Comments and for each row:
- Uses depth to pick margin classes (ml-8, ml-16, etc.)
- Draws a guide line when depth > 0
- Uses parent_id only for the "Replying to comment #N" anchor

Because rows arrive in parent-first order, children naturally appear below their parent in the final HTML.

### Concrete example

Imagine SQL returns rows in this order:

1. id=1, parent_id=null, depth=0, path=...0001
2. id=2, parent_id=1, depth=1, path=...0001....0002
3. id=4, parent_id=2, depth=2, path=...0001....0002....0004
4. id=3, parent_id=1, depth=1, path=...0001....0003

Single-pass template output will visually render:
- Comment 1 (root)
- Comment 2 (indented)
- Comment 4 (more indented)
- Comment 3 (back to previous indentation)

No in-memory tree construction is required.

### Trade-off vs nested JSON

Flat threaded rows (current approach):
- Simpler backend transformation layer
- Efficient for server-side rendering
- Easy to stream in one pass

Nested JSON tree (typical SPA approach):
- Convenient for recursive client components
- Usually requires extra backend transform or client-side tree building
- More moving parts if you do not need rich client interactivity

For this codebase (Go templates, server-rendered pages), flat rows + depth is the simpler and cleaner fit.

## Template Rendering

The comments partial handles three concerns:

1. New top-level comment form.
2. List rendering for existing comments.
3. Per-comment reply form with hidden parent_id.

For visual nesting, each comment applies indentation classes based on depth:
- depth 0: no margin class
- depth 1: ml-8
- depth 2: ml-16
- depth 3: ml-24
- depth 4: ml-32
- depth 5+: ml-40

When depth is greater than zero, a vertical guide line is rendered on the left.

## Included Data Shape

Entry includes comments as a list of sqlc rows. Each row includes at least:
- ID
- Name
- Content
- Slug
- CreatedAt
- ParentID
- Depth
- Path

## Operational Notes

- Canonical slugs should always be stored as /section/slug.
- Legacy rows with bare slug format may still display due to fallback lookup.
- If Tailwind classes change in the template, regenerate CSS with npm run css:build.

## Manual Test Checklist

1. Post a top-level comment with empty name and verify it displays as Anonymous.
2. Post a reply and confirm:
- It appears indented under the parent group.
- Replying to comment #N link points to the parent anchor.
3. Refresh page and confirm ordering remains stable.
4. Try invalid parent_id manually and confirm handler returns bad request behavior.

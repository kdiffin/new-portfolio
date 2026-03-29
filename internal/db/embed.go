package dbschema

import _ "embed"

// SchemaSQL contains the embedded SQLite schema used for bootstrapping.
//
//go:embed schema.sql
var SchemaSQL string

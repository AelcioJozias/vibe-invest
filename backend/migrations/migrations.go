package migrations

import "embed"

// FS holds all SQL migration files embedded at compile time.
// golang-migrate reads from this FS via the iofs source driver,
// similar to how Flyway scans the classpath for migration scripts in Spring.
//
//go:embed *.sql
var FS embed.FS

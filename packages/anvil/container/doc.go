// Package container provides a Laravel-inspired IoC service container.
//
// The container manages service bindings (transient and singleton) and
// resolves them on demand. It supports aliases, tagging, contextual
// bindings, and lifecycle callbacks. It is safe for concurrent use.
//
//	c := container.New()
//	c.Singleton("db", func(c *container.Container) (any, error) {
//	    dsn, _ := c.Make("db.dsn")
//	    return sql.Open("postgres", dsn.(string))
//	})
//
//	db := container.MustMake[*sql.DB](c, "db")
package container

// nolint:revive
package design

import . "goa.design/goa/v3/dsl"

var CacheGetRequest = Type("CacheGetRequest", func() {
	Field(1, "key", String)
	Field(2, "namespace", String)
	Field(3, "scope", String)
	Field(4, "strategy", String)
	Required("key")
})

var CacheSetRequest = Type("CacheSetRequest", func() {
	Field(1, "data", Any)
	Field(2, "key", String)
	Field(3, "namespace", String)
	Field(4, "scope", String) // Initial implementation with a single scope
	Field(5, "ttl", Int)
	Required("data", "key")
})

var HealthResponse = Type("HealthResponse", func() {
	Field(1, "service", String, "Service name.")
	Field(2, "status", String, "Status message.")
	Field(3, "version", String, "Service runtime version.")
	Required("service", "status", "version")
})

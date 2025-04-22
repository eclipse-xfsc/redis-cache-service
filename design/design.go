// nolint:revive
package design

import . "goa.design/goa/v3/dsl"

var _ = API("cache", func() {
	Title("Cache Service")
	Description("The cache service exposes interface for working with Redis.")
	Server("cache", func() {
		Description("Cache Server")
		Host("development", func() {
			Description("Local development server")
			URI("http://localhost:8083")
		})
	})
})

var _ = Service("health", func() {
	Description("Health service provides health check endpoints.")

	Method("Liveness", func() {
		Payload(Empty)
		Result(HealthResponse)
		HTTP(func() {
			GET("/liveness")
			Response(StatusOK)
		})
	})

	Method("Readiness", func() {
		Payload(Empty)
		Result(HealthResponse)
		HTTP(func() {
			GET("/readiness")
			Response(StatusOK)
		})
	})
})

var _ = Service("cache", func() {
	Description("Cache service allows storing and retrieving data from distributed cache.")

	Method("Get", func() {
		Description("Get JSON value from the cache.")

		Payload(CacheGetRequest)
		Result(Any)

		HTTP(func() {
			GET("/v1/cache")

			Header("key:x-cache-key", String, "Cache entry key", func() {
				Example("did:web:example.com")
			})
			Header("namespace:x-cache-namespace", String, "Cache entry namespace", func() {
				Example("Login")
			})
			Header("scope:x-cache-scope", String, "Cache entry scope", func() {
				Example("multiple scopes", "administration,user")
				Example("default", "administration")
			})
			Header("strategy:x-cache-flatten-strategy", String, "Flatten strategy.", func() {
				Example("first key value only", "first")
				Example("last key value only", "last")
				Example("default", "merge")
			})

			Response(StatusOK, func() {
				ContentType("application/json")
			})
		})
	})

	Method("Set", func() {
		Description("Set a JSON value in the cache.")

		Payload(CacheSetRequest)
		Result(Empty)

		HTTP(func() {
			POST("/v1/cache")

			Header("key:x-cache-key", String, "Cache entry key", func() {
				Example("did:web:example.com")
			})
			Header("namespace:x-cache-namespace", String, "Cache entry namespace", func() {
				Example("Login")
			})
			Header("scope:x-cache-scope", String, "Cache entry scope", func() {
				Example("administration")
			})
			Header("ttl:x-cache-ttl", Int, "Cache entry TTL in seconds", func() {
				Example(60)
			})
			Body("data")

			Response(StatusCreated)
		})
	})

	Method("SetExternal", func() {
		Description("Set an external JSON value in the cache and provide an event for the input.")

		Payload(CacheSetRequest)
		Result(Empty)

		HTTP(func() {
			POST("/v1/external/cache")

			Header("key:x-cache-key", String, "Cache entry key", func() {
				Example("did:web:example.com")
			})
			Header("namespace:x-cache-namespace", String, "Cache entry namespace", func() {
				Example("Login")
			})
			Header("scope:x-cache-scope", String, "Cache entry scope", func() {
				Example("administration")
			})
			Header("ttl:x-cache-ttl", Int, "Cache entry TTL in seconds", func() {
				Example(60)
			})
			Body("data")

			Response(StatusOK)
		})
	})
})

var _ = Service("openapi", func() {
	Description("The openapi service serves the OpenAPI(v3) definition.")
	Meta("swagger:generate", "false")
	HTTP(func() {
		Path("/swagger-ui")
	})
	Files("/openapi.json", "./gen/http/openapi3.json", func() {
		Description("JSON document containing the OpenAPI(v3) service definition")
	})
	Files("/{*filepath}", "./swagger/")
})

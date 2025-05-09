// Code generated by goa v3.20.1, DO NOT EDIT.
//
// cache HTTP client CLI support package
//
// Command:
// $ goa gen github.com/eclipse-xfsc/redis-cache-service/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	cachec "github.com/eclipse-xfsc/redis-cache-service/gen/http/cache/client"
	healthc "github.com/eclipse-xfsc/redis-cache-service/gen/http/health/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `cache (get|set|set-external)
health (liveness|readiness)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` cache get --key "Iusto consequatur voluptatem eligendi et eligendi." --namespace "Optio natus." --scope "Ratione quasi perspiciatis qui." --strategy "Animi non alias occaecati esse."` + "\n" +
		os.Args[0] + ` health liveness` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, any, error) {
	var (
		cacheFlags = flag.NewFlagSet("cache", flag.ContinueOnError)

		cacheGetFlags         = flag.NewFlagSet("get", flag.ExitOnError)
		cacheGetKeyFlag       = cacheGetFlags.String("key", "REQUIRED", "")
		cacheGetNamespaceFlag = cacheGetFlags.String("namespace", "", "")
		cacheGetScopeFlag     = cacheGetFlags.String("scope", "", "")
		cacheGetStrategyFlag  = cacheGetFlags.String("strategy", "", "")

		cacheSetFlags         = flag.NewFlagSet("set", flag.ExitOnError)
		cacheSetBodyFlag      = cacheSetFlags.String("body", "REQUIRED", "")
		cacheSetKeyFlag       = cacheSetFlags.String("key", "REQUIRED", "")
		cacheSetNamespaceFlag = cacheSetFlags.String("namespace", "", "")
		cacheSetScopeFlag     = cacheSetFlags.String("scope", "", "")
		cacheSetTTLFlag       = cacheSetFlags.String("ttl", "", "")

		cacheSetExternalFlags         = flag.NewFlagSet("set-external", flag.ExitOnError)
		cacheSetExternalBodyFlag      = cacheSetExternalFlags.String("body", "REQUIRED", "")
		cacheSetExternalKeyFlag       = cacheSetExternalFlags.String("key", "REQUIRED", "")
		cacheSetExternalNamespaceFlag = cacheSetExternalFlags.String("namespace", "", "")
		cacheSetExternalScopeFlag     = cacheSetExternalFlags.String("scope", "", "")
		cacheSetExternalTTLFlag       = cacheSetExternalFlags.String("ttl", "", "")

		healthFlags = flag.NewFlagSet("health", flag.ContinueOnError)

		healthLivenessFlags = flag.NewFlagSet("liveness", flag.ExitOnError)

		healthReadinessFlags = flag.NewFlagSet("readiness", flag.ExitOnError)
	)
	cacheFlags.Usage = cacheUsage
	cacheGetFlags.Usage = cacheGetUsage
	cacheSetFlags.Usage = cacheSetUsage
	cacheSetExternalFlags.Usage = cacheSetExternalUsage

	healthFlags.Usage = healthUsage
	healthLivenessFlags.Usage = healthLivenessUsage
	healthReadinessFlags.Usage = healthReadinessUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "cache":
			svcf = cacheFlags
		case "health":
			svcf = healthFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "cache":
			switch epn {
			case "get":
				epf = cacheGetFlags

			case "set":
				epf = cacheSetFlags

			case "set-external":
				epf = cacheSetExternalFlags

			}

		case "health":
			switch epn {
			case "liveness":
				epf = healthLivenessFlags

			case "readiness":
				epf = healthReadinessFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "cache":
			c := cachec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get":
				endpoint = c.Get()
				data, err = cachec.BuildGetPayload(*cacheGetKeyFlag, *cacheGetNamespaceFlag, *cacheGetScopeFlag, *cacheGetStrategyFlag)
			case "set":
				endpoint = c.Set()
				data, err = cachec.BuildSetPayload(*cacheSetBodyFlag, *cacheSetKeyFlag, *cacheSetNamespaceFlag, *cacheSetScopeFlag, *cacheSetTTLFlag)
			case "set-external":
				endpoint = c.SetExternal()
				data, err = cachec.BuildSetExternalPayload(*cacheSetExternalBodyFlag, *cacheSetExternalKeyFlag, *cacheSetExternalNamespaceFlag, *cacheSetExternalScopeFlag, *cacheSetExternalTTLFlag)
			}
		case "health":
			c := healthc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "liveness":
				endpoint = c.Liveness()
			case "readiness":
				endpoint = c.Readiness()
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// cacheUsage displays the usage of the cache command and its subcommands.
func cacheUsage() {
	fmt.Fprintf(os.Stderr, `Cache service allows storing and retrieving data from distributed cache.
Usage:
    %[1]s [globalflags] cache COMMAND [flags]

COMMAND:
    get: Get JSON value from the cache.
    set: Set a JSON value in the cache.
    set-external: Set an external JSON value in the cache and provide an event for the input.

Additional help:
    %[1]s cache COMMAND --help
`, os.Args[0])
}
func cacheGetUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] cache get -key STRING -namespace STRING -scope STRING -strategy STRING

Get JSON value from the cache.
    -key STRING: 
    -namespace STRING: 
    -scope STRING: 
    -strategy STRING: 

Example:
    %[1]s cache get --key "Iusto consequatur voluptatem eligendi et eligendi." --namespace "Optio natus." --scope "Ratione quasi perspiciatis qui." --strategy "Animi non alias occaecati esse."
`, os.Args[0])
}

func cacheSetUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] cache set -body JSON -key STRING -namespace STRING -scope STRING -ttl INT

Set a JSON value in the cache.
    -body JSON: 
    -key STRING: 
    -namespace STRING: 
    -scope STRING: 
    -ttl INT: 

Example:
    %[1]s cache set --body "Enim vel." --key "Ut in." --namespace "Ab dolores distinctio quis." --scope "Optio aliquam error nam." --ttl 2227603043401673122
`, os.Args[0])
}

func cacheSetExternalUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] cache set-external -body JSON -key STRING -namespace STRING -scope STRING -ttl INT

Set an external JSON value in the cache and provide an event for the input.
    -body JSON: 
    -key STRING: 
    -namespace STRING: 
    -scope STRING: 
    -ttl INT: 

Example:
    %[1]s cache set-external --body "Sint ipsa fugiat et id rem." --key "Molestiae minima." --namespace "Quia dolores rem." --scope "Est illum." --ttl 6207033275224297400
`, os.Args[0])
}

// healthUsage displays the usage of the health command and its subcommands.
func healthUsage() {
	fmt.Fprintf(os.Stderr, `Health service provides health check endpoints.
Usage:
    %[1]s [globalflags] health COMMAND [flags]

COMMAND:
    liveness: Liveness implements Liveness.
    readiness: Readiness implements Readiness.

Additional help:
    %[1]s health COMMAND --help
`, os.Args[0])
}
func healthLivenessUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] health liveness

Liveness implements Liveness.

Example:
    %[1]s health liveness
`, os.Args[0])
}

func healthReadinessUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] health readiness

Readiness implements Readiness.

Example:
    %[1]s health readiness
`, os.Args[0])
}

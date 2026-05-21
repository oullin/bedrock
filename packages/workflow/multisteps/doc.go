// sync/async jobs with dependency resolution, conditional execution, and retry
// policies.
//
// Usage:
//
//	wf := multisteps.Workflow("signup",
//	    multisteps.Sync("create", createUserAccount, multisteps.Args(multisteps.A{
//	        "name": multisteps.Variable("name"),
//	    })),
//	    multisteps.Async("email", sendWelcomeEmail,
//	        multisteps.Args(multisteps.A{"userId": multisteps.Response("create", "id")}),
//	        multisteps.WithRetry(3, time.Second, 30*time.Second),
//	    ),
//	    multisteps.Async("notify", notifyTeam,
//	        multisteps.Args(multisteps.A{"userId": multisteps.Response("create", "id")}),
//	        multisteps.WithRunIf(func(in multisteps.JobInput) bool {
//	            return in.Vars["env"] == "prod"
//	        }),
//	    ),
//	)
//
//	eng := multisteps.NewEngine(multisteps.WithDriver(driver))
//	res, err := eng.Run(ctx, wf, map[string]any{"name": "Jane", "env": "prod"})
//
// Async siblings fan out via the injected concurrency.Driver; failures cancel
// the sibling group by default (toggle via WithContinueOnError).
package multisteps

// 8. 0's named AWS
// credential providers for the SQS queue driver. Instead of relying on
// whatever ambient credentials the process happens to inherit, callers
// register a set of named providers (a profile, an SSO session, an EC2
// role, a static credential pair, a chain) and reference them by name
// from connection configuration.
//
// Typical wiring at service-bootstrap time:
//
//	reg := credentials.NewRegistry()
//	reg.Register("prod-profile",   credentials.Profile("prod"))
//	reg.Register("staging-static", credentials.Static("AKIA…", "secret", ""))
//	reg.Register("ec2",            credentials.EC2InstanceRole())
//
// Then, when constructing the SQS client at connection-resolve time:
//
//	provider, err := reg.Resolve(ctx, cfg["credentials"].(string))
//	if err != nil { return nil, err }
//
//	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider))
//	// hand awsCfg to your sqs.Client constructor and adapt it to drivers.SQSClient.
//
// Why client construction lives outside the registry: each deployment
// tweaks the SQS client differently (custom endpoints, retry middleware,
// regional fallbacks, observability hooks). Putting the wiring in
// service code keeps this subpackage focused on the credential
// resolution problem and avoids pulling the entire AWS SQS service
// module into the queue package's dependency closure.
package credentials

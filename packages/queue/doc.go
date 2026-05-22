// Package queue provides job queue management.
// It defines Queue, Job, and Connector interfaces with multiple driver
// implementations (sync, database, redis, beanstalkd, sqs, null, background,
// deferred, failover) and a Worker for processing jobs.
package queue

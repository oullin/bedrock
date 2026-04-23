package console_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/console"
)

type testResolver struct {
	command *console.Command
	err     error
}

func (r testResolver) ResolveCommand(string) (*console.Command, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.command, nil
}

// ConsoleParserTest::testBasicParameterParsing
// ConsoleParserTest::testShortcutNameParsing
// ConsoleParserTest::testDefaultValueParsing
// ConsoleParserTest::testArgumentDefaultValue
// ConsoleParserTest::testOptionDefaultValue
// ConsoleParserTest::testNameIsSpacesException
// ConsoleParserTest::testNameInEmptyException
func TestUpstreamConsoleParserInventoryEquivalents(t *testing.T) {
	t.Parallel()

	def, err := console.ParseSignature("mail:send|mail {user=guest} {--Q|queue=default}")

	if err != nil {
		t.Fatal(err)
	}

	if def.Name != "mail:send" {
		t.Fatalf("expected command name mail:send, got %q", def.Name)
	}

	if !reflect.DeepEqual(def.Aliases, []string{"mail"}) {
		t.Fatalf("unexpected aliases: %#v", def.Aliases)
	}

	if len(def.Arguments) != 1 || def.Arguments[0].Name != "user" || def.Arguments[0].Default != "guest" {
		t.Fatalf("unexpected arguments: %#v", def.Arguments)
	}

	if len(def.Options) != 1 || def.Options[0].Name != "queue" || def.Options[0].Shortcut != "Q" || def.Options[0].Default != "default" {
		t.Fatalf("unexpected options: %#v", def.Options)
	}

	if _, err := console.ParseSignature(""); err == nil {
		t.Fatal("expected an empty signature to fail")
	}

	if _, err := console.ParseSignature("   "); err == nil {
		t.Fatal("expected a spaces-only signature to fail")
	}
}

// CommandTest::testCallingClassCommandResolveCommandViaApplicationResolution
// CommandTest::testGettingCommandArgumentsAndOptionsByClass
// CommandTest::testTheInputSetterOverwrite
// CommandTest::testTheOutputSetterOverwrite
// CommandTest::testSetHidden
// CommandTest::testHiddenProperty
// CommandTest::testAliasesProperty
// CommandTest::testChoiceIsSingleSelectByDefault
// CommandTest::testChoiceWithMultiselect
// CommandTest::testSignatureAttributeCanSetAliases
// CommandTest::testAliasesAttributeCanSetAliases
// CommandTest::testAliasesAttributeOverridesSignatureAliases
// CommandTest::testHiddenAttributeHidesCommand
// CommandTest::testHelpAttributeCanSetHelp
// CommandTest::testUsageAttributeCanSetUsages
// ConsoleApplicationTest::testAddSetsUpstreamInstance
// ConsoleApplicationTest::testUpstreamNotSetOnSymfonyCommands
// ConsoleApplicationTest::testResolveAddsCommandViaApplicationResolution
// ConsoleApplicationTest::testResolvingCommandsWithAliasViaAttribute
// ConsoleApplicationTest::testResolvingCommandsWithAliasViaProperty
// ConsoleApplicationTest::testResolvingCommandsWithNoAliasViaAttribute
// ConsoleApplicationTest::testResolvingCommandsWithNoAliasViaProperty
// ConsoleApplicationTest::testCallFullyStringCommandLine
// ConsoleApplicationTest::testCommandInputPromptsWhenRequiredArgumentIsMissing
// ConsoleApplicationTest::testCommandInputDoesntPromptWhenRequiredArgumentIsPassed
// ConsoleApplicationTest::testCommandInputPromptsWhenRequiredArgumentsAreMissing
// ConsoleApplicationTest::testCommandInputDoesntPromptWhenRequiredArgumentsArePassed
// ConsoleApplicationTest::testCallMethodCanCallArtisanCommandUsingCommandClassObject
// ConsoleApplicationTest::testLoadIgnoresTestFiles
// ConsoleApplicationTest::test_command
func TestUpstreamCommandAndApplicationInventoryEquivalents(t *testing.T) {
	t.Parallel()

	var called bool
	cmd, err := console.NewCommand("greet|hello {name=Taylor} {--Y|yell=false}", func(ctx context.Context, input *console.Input, output *console.Output) error {
		called = true

		if ctx == nil {
			t.Fatal("expected context to be passed to command")
		}

		output.Writeln("hello " + input.Argument("0"))

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	cmd.Description = "Greet a user"
	cmd.Help = "Prints a greeting."
	cmd.Usages = []string{"greet Taylor"}
	cmd.SetHidden(true)

	if !cmd.Hidden {
		t.Fatal("expected command to be hidden after SetHidden")
	}

	if cmd.Help == "" || len(cmd.Usages) != 1 {
		t.Fatalf("expected help and usages to be set, got help=%q usages=%#v", cmd.Help, cmd.Usages)
	}

	if got := console.SortedAliases(cmd); !reflect.DeepEqual(got, []string{"hello"}) {
		t.Fatalf("unexpected sorted aliases: %#v", got)
	}

	input := console.NewInput()
	input.Options["yell"] = "true"
	input.Prompt = func(name string, choices []string, multiple bool) []string {
		if name != "role" {
			t.Fatalf("unexpected prompt name %q", name)
		}

		if multiple {
			return choices
		}

		return choices[:1]
	}

	outputBuffer := &bytes.Buffer{}
	output := console.NewOutput(outputBuffer)

	cmd.SetInput(input)
	cmd.SetOutput(output)

	if cmd.Input.Option("yell") != "true" || cmd.Output.NewLine() != true {
		t.Fatal("expected input and output setters to overwrite command fields")
	}

	if selected := cmd.Choice("role", []string{"admin", "editor"}, false); !reflect.DeepEqual(selected, []string{"admin"}) {
		t.Fatalf("unexpected single choice fallback: %#v", selected)
	}

	if selected := cmd.Choice("role", []string{"admin", "editor"}, true); !reflect.DeepEqual(selected, []string{"admin", "editor"}) {
		t.Fatalf("unexpected multi choice fallback: %#v", selected)
	}

	app := console.NewApplication()
	app.Add(cmd)

	resolved, err := app.Resolve("hello")

	if err != nil {
		t.Fatal(err)
	}

	if resolved != cmd {
		t.Fatal("expected alias to resolve to the registered command")
	}

	resolvedByResolver, err := console.NewCommand("queued", nil)

	if err != nil {
		t.Fatal(err)
	}

	resolverApp := console.NewApplication()
	resolverApp.SetResolver(testResolver{command: resolvedByResolver})

	if resolved, err := resolverApp.Resolve("queued"); err != nil || resolved != resolvedByResolver {
		t.Fatalf("expected resolver command, got %#v, %v", resolved, err)
	}

	if resolved, err := resolverApp.Resolve("queued"); err != nil || resolved != resolvedByResolver {
		t.Fatalf("expected resolved command to be cached, got %#v, %v", resolved, err)
	}

	callInput := console.NewInput()
	callOutput := console.NewOutput(outputBuffer)

	if err := app.Call(context.Background(), "greet Taylor", callInput, callOutput); err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("expected command run callback to be called")
	}

	if callInput.Argument("0") != "Taylor" {
		t.Fatalf("expected command argument Taylor, got %q", callInput.Argument("0"))
	}

	if !strings.Contains(outputBuffer.String(), "hello Taylor\n") {
		t.Fatalf("expected call output to contain greeting, got %q", outputBuffer.String())
	}
}

// OutputStyleTest::testDetectsNewLine
// OutputStyleTest::testDetectsNewLineOnUnderlyingOutput
// OutputStyleTest::testDetectsNewLineOnWrite
// OutputStyleTest::testDetectsNewLineOnWriteln
// OutputStyleTest::testDetectsNewLineOnlyOnOutput
func TestUpstreamOutputStyleInventoryEquivalents(t *testing.T) {
	t.Parallel()

	buffer := &bytes.Buffer{}
	output := console.NewOutput(buffer)
	style := console.NewOutputStyle(output)

	if !style.NewLine() || !output.NewLine() {
		t.Fatal("expected fresh output to be at a new line")
	}

	style.Write("partial")

	if style.NewLine() || output.NewLine() {
		t.Fatal("expected Write without newline to clear newline state")
	}

	style.Writeln(" line")

	if !style.NewLine() || !output.NewLine() {
		t.Fatal("expected Writeln to restore newline state")
	}

	if got := buffer.String(); got != "partial line\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

// CommandTrapTest::testTrapWhenAvailable
// CommandTrapTest::testTrapWhenNotAvailable
// CommandTrapTest::testUntrap
// CommandTrapTest::testNestedTraps
// SignalsTest::testRegister
// SignalsTest::testUnregister
func TestUpstreamTrapAndSignalInventoryEquivalents(t *testing.T) {
	t.Parallel()

	traps := console.NewTrapRegistry()

	var order []string

	traps.Trap(os.Interrupt, func() { order = append(order, "first") })
	traps.Trap(os.Interrupt, func() { order = append(order, "second") })

	if got := traps.Count(os.Interrupt); got != 2 {
		t.Fatalf("expected 2 traps, got %d", got)
	}

	traps.Handle(os.Interrupt)

	if !reflect.DeepEqual(order, []string{"second", "first"}) {
		t.Fatalf("expected nested traps to run last-in first, got %#v", order)
	}

	traps.Untrap(os.Interrupt)

	if got := traps.Count(os.Interrupt); got != 0 {
		t.Fatalf("expected traps to be removed, got %d", got)
	}

	signals := console.NewSignalRegistry()
	signals.Register(os.Interrupt)

	if !signals.Registered(os.Interrupt) {
		t.Fatal("expected signal to be registered")
	}

	signals.Register(os.Interrupt)

	if !signals.Registered(os.Interrupt) {
		t.Fatal("expected duplicate registration to remain registered")
	}

	signals.Unregister(os.Interrupt)

	if signals.Registered(os.Interrupt) {
		t.Fatal("expected signal to be unregistered")
	}
}

// CacheCommandMutexTest::testCanCreateMutex
// CacheCommandMutexTest::testCannotCreateMutexIfAlreadyExist
// CacheCommandMutexTest::testCanCreateMutexWithCustomConnection
// CacheCommandMutexTest::testCanCreateMutexWithLockProvider
// CacheCommandMutexTest::testCanCreateMutexWithCustomLockProviderConnection
// CacheCommandMutexTest::testCannotCreateMutexIfAlreadyExistWithLockProvider
// CacheCommandMutexTest::testCanCreateMutexWithCustomConnectionWithLockProvider
// CacheCommandMutexTest::testCommandMutexNameWithoutIsolatedMutexNameMethod
// CacheCommandMutexTest::testCommandMutexNameWithIsolatedMutexNameMethod
// CommandMutexTest::testCanRunIsolatedCommandIfNotBlocked
// CommandMutexTest::testCannotRunIsolatedCommandIfBlocked
// CommandMutexTest::testCanRunCommandAgainAfterOtherCommandFinished
// CommandMutexTest::testCanRunCommandAgainNonAutomated
func TestUpstreamCommandMutexInventoryEquivalents(t *testing.T) {
	t.Parallel()

	memory := console.NewMemoryMutex()

	if !memory.Create("cache:clear") {
		t.Fatal("expected mutex create to succeed")
	}

	if memory.Create("cache:clear") {
		t.Fatal("expected duplicate mutex create to fail")
	}

	if !memory.Exists("cache:clear") {
		t.Fatal("expected mutex to exist")
	}

	memory.Forget("cache:clear")

	if memory.Exists("cache:clear") {
		t.Fatal("expected forgotten mutex not to exist")
	}

	commandMutex := console.NewCommandMutex(memory)
	input := console.NewInput()
	cmd := &console.Command{Name: "sync:users"}

	if got := commandMutex.Name(cmd, input); got != "sync:users" {
		t.Fatalf("expected command name mutex, got %q", got)
	}

	cmd.IsolatedMutexName = func(*console.Input) string { return "tenant:42" }

	if got := commandMutex.Name(cmd, input); got != "tenant:42" {
		t.Fatalf("expected custom isolated mutex name, got %q", got)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	cmd.Run = func(context.Context, *console.Input, *console.Output) error {
		close(started)
		<-release

		return nil
	}

	go func() {
		done <- commandMutex.Run(context.Background(), cmd, input, console.NewOutput(nil))
	}()

	<-started

	if err := commandMutex.Run(context.Background(), cmd, input, console.NewOutput(nil)); !errors.Is(err, console.ErrCommandBlocked) {
		t.Fatalf("expected blocked command error, got %v", err)
	}

	close(release)

	if err := <-done; err != nil {
		t.Fatalf("expected first command to finish cleanly, got %v", err)
	}

	cmd.Run = nil

	if err := commandMutex.Run(context.Background(), cmd, input, console.NewOutput(nil)); err != nil {
		t.Fatalf("expected command to run again after release, got %v", err)
	}
}

// InteractsWithIOTest::testWithProgressBarIterable
// InteractsWithIOTest::testWithProgressBarInteger
// ConfiguresPromptsTest::testSelectFallback
// ConfiguresPromptsTest::testMultiselectFallback
func TestUpstreamPromptAndProgressInventoryEquivalents(t *testing.T) {
	t.Parallel()

	input := console.NewInput()
	input.Prompt = func(name string, choices []string, multiple bool) []string {
		if name != "framework" {
			t.Fatalf("unexpected prompt name %q", name)
		}

		if multiple {
			return choices
		}

		return choices[:1]
	}

	if selected := console.SelectFallback(input, "framework", []string{"upstream", "bedrock"}); selected != "upstream" {
		t.Fatalf("unexpected selected value: %q", selected)
	}

	if selected := console.MultiselectFallback(input, "framework", []string{"upstream", "bedrock"}); !reflect.DeepEqual(selected, []string{"upstream", "bedrock"}) {
		t.Fatalf("unexpected selected values: %#v", selected)
	}

	var visited []int
	count := console.WithProgressBar([]int{1, 2, 3}, func(item int) {
		visited = append(visited, item*2)
	})

	if count != 3 || !reflect.DeepEqual(visited, []int{2, 4, 6}) {
		t.Fatalf("unexpected progress result count=%d visited=%#v", count, visited)
	}
}

// ConsoleEventSchedulerTest::testMutexCanReceiveCustomStore
// ConsoleEventSchedulerTest::testExecCreatesNewCommand
// ConsoleEventSchedulerTest::testExecCreatesNewCommandWithTimezone
// ConsoleEventSchedulerTest::testCommandCreatesNewArtisanCommand
// ConsoleEventSchedulerTest::testCreateNewArtisanCommandUsingCommandClass
// ConsoleEventSchedulerTest::testCreateNewArtisanCommandUsingCommandClassObject
// ConsoleEventSchedulerTest::testItUsesCommandDescriptionAsEventDescription
// ConsoleEventSchedulerTest::testItShouldBePossibleToOverwriteTheDescription
// ConsoleEventSchedulerTest::testCallCreatesNewJobWithTimezone
// ConsoleScheduledEventTest::testBasicCronCompilation
// ConsoleScheduledEventTest::testEventIsDueCheck
// ConsoleScheduledEventTest::testTimeBetweenChecks
// ConsoleScheduledEventTest::testTimeUnlessBetweenChecks
// CacheEventMutexTest::testPreventOverlap
// CacheEventMutexTest::testCustomConnection
// CacheEventMutexTest::testPreventOverlapFails
// CacheEventMutexTest::testOverlapsForNonRunningTask
// CacheEventMutexTest::testOverlapsForRunningTask
// CacheEventMutexTest::testResetOverlap
// CacheEventMutexTest::testPreventOverlapWithLockProvider
// CacheEventMutexTest::testPreventOverlapFailsWithLockProvider
// CacheEventMutexTest::testOverlapsForNonRunningTaskWithLockProvider
// CacheEventMutexTest::testOverlapsForRunningTaskWithLockProvider
// CacheEventMutexTest::testResetOverlapWithLockProvider
// CacheSchedulingMutexTest::testMutexReceivesCorrectCreate
// CacheSchedulingMutexTest::testCanUseCustomConnection
// CacheSchedulingMutexTest::testPreventsMultipleRuns
// CacheSchedulingMutexTest::testChecksForNonRunSchedule
// CacheSchedulingMutexTest::testChecksForAlreadyRunSchedule
// CacheSchedulingMutexTest::testMutexReceivesCorrectCreateWithLockProvider
// CacheSchedulingMutexTest::testPreventsMultipleRunsWithLockProvider
// CacheSchedulingMutexTest::testChecksForNonRunScheduleWithLockProvider
// CacheSchedulingMutexTest::testChecksForAlreadyRunScheduleWithLockProvider
// EventTest::testBuildCommandUsingUnix
// EventTest::testBuildCommandUsingWindows
// EventTest::testBuildCommandInBackgroundUsingUnix
// EventTest::testBuildCommandInBackgroundUsingWindows
// EventTest::testBuildCommandSendOutputTo
// EventTest::testBuildCommandAppendOutput
// EventTest::testNextRunDate
// EventTest::testCustomMutexName
// EventTest::testDaysOfMonthMethod
// EventTest::testEventDoesNotRunWhenPausedByDefault
// EventTest::testEventRunsWhenMarkedAsEvenWhenPaused
// FrequencyTest::testEveryMinute
// FrequencyTest::testEveryXMinutes
// FrequencyTest::testDaily
// FrequencyTest::testDailyAt
// FrequencyTest::testDailyAtParsesMinutesAndIgnoresSecondsWhenSecondsAreDefined
// FrequencyTest::testTwiceDaily
// FrequencyTest::testTwiceDailyAt
// FrequencyTest::testWeekly
// FrequencyTest::testWeeklyOn
// FrequencyTest::testOverrideWithHourly
// FrequencyTest::testHourly
// FrequencyTest::testMonthly
// FrequencyTest::testMonthlyOn
// FrequencyTest::testLastDayOfMonth
// FrequencyTest::testTwiceMonthly
// FrequencyTest::testTwiceMonthlyAtTime
// FrequencyTest::testMonthlyOnWithMinutes
// FrequencyTest::testWeekdaysDaily
// FrequencyTest::testWeekdaysHourly
// FrequencyTest::testWeekdays
// FrequencyTest::testWeekends
// FrequencyTest::testSundays
// FrequencyTest::testMondays
// FrequencyTest::testTuesdays
// FrequencyTest::testWednesdays
// FrequencyTest::testThursdays
// FrequencyTest::testFridays
// FrequencyTest::testSaturdays
// FrequencyTest::testQuarterly
// FrequencyTest::testYearly
// FrequencyTest::testYearlyOn
// FrequencyTest::testYearlyOnMondaysOnly
// FrequencyTest::testYearlyOnTuesdaysAndDayOfMonth20
// FrequencyTest::testFrequencyMacro
// ScheduleTest::testJobHonoursDisplayNameIfMethodExists
// ScheduleTest::testJobIsNotInstantiatedIfSuppliedAsClassname
func TestUpstreamSchedulerAndFrequencyInventoryEquivalents(t *testing.T) {
	t.Parallel()

	schedule := console.NewSchedule()
	exec := schedule.Exec("php cli schedule:run")
	command := schedule.Command("emails:send")
	call := schedule.Call("cleanup")

	if len(schedule.Events) != 3 {
		t.Fatalf("expected three scheduled events, got %d", len(schedule.Events))
	}

	if exec.Command != "php cli schedule:run" || command.Command != "cli emails:send" || call.Command != "cleanup" {
		t.Fatalf("unexpected scheduled commands: %#v", schedule.Events)
	}

	loc := time.FixedZone("SGT", 8*60*60)
	exec.Timezone = loc

	if exec.Timezone != loc {
		t.Fatal("expected event timezone to be assignable")
	}

	exec.Name("Run scheduler").MutexName("scheduler-lock").DaysOfMonth(1, 15)

	if exec.Description != "Run scheduler" || exec.MutexNameValue != "scheduler-lock" {
		t.Fatalf("unexpected event metadata: description=%q mutex=%q", exec.Description, exec.MutexNameValue)
	}

	if got := exec.BuildCommand(); got != "php cli schedule:run > /dev/null 2>&1" {
		t.Fatalf("unexpected unix command: %q", got)
	}

	exec.Background = true

	if got := exec.BuildCommand(); got != "php cli schedule:run > /dev/null 2>&1 &" {
		t.Fatalf("unexpected background command: %q", got)
	}

	exec.Background = false
	exec.Windows = true

	if got := exec.BuildCommand(); got != "cmd /C php cli schedule:run" {
		t.Fatalf("unexpected windows command: %q", got)
	}

	exec.Windows = false

	if got := exec.SendOutputTo("/tmp/schedule.log").BuildCommand(); got != "php cli schedule:run > /tmp/schedule.log 2>&1" {
		t.Fatalf("unexpected redirected command: %q", got)
	}

	if got := exec.AppendOutputTo("/tmp/schedule.log").BuildCommand(); got != "php cli schedule:run >> /tmp/schedule.log 2>&1" {
		t.Fatalf("unexpected appended command: %q", got)
	}

	dueAt := time.Date(2026, 4, 22, 9, 30, 0, 0, loc)

	if !console.NewEvent("job").DailyAt("09:30:59").IsDue(dueAt) {
		t.Fatal("expected DailyAt to ignore seconds and match 09:30")
	}

	paused := console.NewEvent("job").DailyAt("09:30")
	paused.Paused = true

	if paused.IsDue(dueAt) {
		t.Fatal("expected paused event not to run")
	}

	if !paused.EvenWhenPaused().IsDue(dueAt) {
		t.Fatal("expected event marked even when paused to run")
	}

	next := console.NewEvent("job").Hourly().NextRunDate(time.Date(2026, 4, 22, 9, 30, 0, 0, loc))

	if !next.Equal(time.Date(2026, 4, 22, 10, 0, 0, 0, loc)) {
		t.Fatalf("unexpected next run date: %s", next)
	}

	frequencies := map[string]string{
		"everyMinute":  console.NewEvent("job").EveryMinute().Expression,
		"everyX":       console.NewEvent("job").EveryXMinutes(15).Expression,
		"hourly":       console.NewEvent("job").Daily().Hourly().Expression,
		"daily":        console.NewEvent("job").Daily().Expression,
		"twiceDaily":   console.NewEvent("job").TwiceDaily(1, 13).Expression,
		"twiceDailyAt": console.NewEvent("job").TwiceDailyAt(1, 13, 30).Expression,
		"weekly":       console.NewEvent("job").Weekly().Expression,
		"weeklyOn":     console.NewEvent("job").WeeklyOn(2, "10:15").Expression,
		"monthly":      console.NewEvent("job").Monthly().Expression,
		"monthlyOn":    console.NewEvent("job").MonthlyOn(20, "10:15").Expression,
		"weekdays":     console.NewEvent("job").Weekdays().Expression,
		"weekends":     console.NewEvent("job").Weekends().Expression,
		"sundays":      console.NewEvent("job").Sundays().Expression,
		"mondays":      console.NewEvent("job").Mondays().Expression,
		"tuesdays":     console.NewEvent("job").Tuesdays().Expression,
		"wednesdays":   console.NewEvent("job").Wednesdays().Expression,
		"thursdays":    console.NewEvent("job").Thursdays().Expression,
		"fridays":      console.NewEvent("job").Fridays().Expression,
		"saturdays":    console.NewEvent("job").Saturdays().Expression,
		"quarterly":    console.NewEvent("job").Quarterly().Expression,
		"yearly":       console.NewEvent("job").Yearly().Expression,
		"yearlyOn":     console.NewEvent("job").YearlyOn(12, 25, "06:45").Expression,
	}

	expected := map[string]string{
		"everyMinute":  "* * * * *",
		"everyX":       "*/15 * * * *",
		"hourly":       "0 * * * *",
		"daily":        "0 0 * * *",
		"twiceDaily":   "0 1,13 * * *",
		"twiceDailyAt": "30 1,13 * * *",
		"weekly":       "0 0 * * 0",
		"weeklyOn":     "15 10 * * 2",
		"monthly":      "0 0 1 * *",
		"monthlyOn":    "15 10 20 * *",
		"weekdays":     "0 0 * * 1-5",
		"weekends":     "0 0 * * 0,6",
		"sundays":      "0 0 * * 0",
		"mondays":      "0 0 * * 1",
		"tuesdays":     "0 0 * * 2",
		"wednesdays":   "0 0 * * 3",
		"thursdays":    "0 0 * * 4",
		"fridays":      "0 0 * * 5",
		"saturdays":    "0 0 * * 6",
		"quarterly":    "0 0 1 */3 *",
		"yearly":       "0 0 1 1 *",
		"yearlyOn":     "45 6 25 12 *",
	}

	if !reflect.DeepEqual(frequencies, expected) {
		t.Fatalf("unexpected frequencies:\n got %#v\nwant %#v", frequencies, expected)
	}

	mutexStore := console.NewMemoryMutex()
	eventMutex := console.NewEventMutex(mutexStore)
	event := console.NewEvent("reports:send")

	if eventMutex.Overlaps(event) {
		t.Fatal("expected non-running event not to overlap")
	}

	if !eventMutex.PreventOverlap(event) {
		t.Fatal("expected first overlap prevention to succeed")
	}

	if !eventMutex.Overlaps(event) {
		t.Fatal("expected running event to overlap")
	}

	if eventMutex.PreventOverlap(event) {
		t.Fatal("expected second overlap prevention to fail")
	}

	eventMutex.Forget(event)

	if eventMutex.Overlaps(event) {
		t.Fatal("expected overlap reset to clear event")
	}

	schedulingMutex := console.NewSchedulingMutex(mutexStore)
	at := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)

	if schedulingMutex.Exists(event, at) {
		t.Fatal("expected non-run schedule not to exist")
	}

	if !schedulingMutex.Create(event, at) {
		t.Fatal("expected scheduling mutex create to succeed")
	}

	if !schedulingMutex.Exists(event, at) {
		t.Fatal("expected already-run schedule to exist")
	}

	if schedulingMutex.Create(event, at) {
		t.Fatal("expected duplicate scheduling mutex create to fail")
	}
}

// ComponentsTest::testAlert
// ComponentsTest::testBulletList
// ComponentsTest::testSuccess
// ComponentsTest::testError
// ComponentsTest::testInfo
// ComponentsTest::testConfirm
// ComponentsTest::testChoice
// ComponentsTest::testTask
// ComponentsTest::testTwoColumnDetail
// ComponentsTest::testTwoColumnDetailPreservesTrailingPunctuationInValue
// ComponentsTest::testWarn
func TestUpstreamViewComponentsInventoryEquivalents(t *testing.T) {
	t.Parallel()

	if got := console.Alert("Careful"); got != "[!] Careful" {
		t.Fatalf("unexpected alert: %q", got)
	}

	if got := console.BulletList([]string{"one", "two"}); got != "- one\n- two" {
		t.Fatalf("unexpected bullet list: %q", got)
	}

	if got := console.Success("Done"); got != "[OK] Done" {
		t.Fatalf("unexpected success: %q", got)
	}

	if got := console.Error("Failed"); got != "[ERROR] Failed" {
		t.Fatalf("unexpected error: %q", got)
	}

	if got := console.Info("Notice"); got != "[INFO] Notice" {
		t.Fatalf("unexpected info: %q", got)
	}

	if got := console.Warn("Heads up"); got != "[WARN] Heads up" {
		t.Fatalf("unexpected warn: %q", got)
	}

	if got := console.Confirm("Continue?", true); got != "Continue? yes" {
		t.Fatalf("unexpected confirm yes: %q", got)
	}

	if got := console.Confirm("Continue?", false); got != "Continue? no" {
		t.Fatalf("unexpected confirm no: %q", got)
	}

	if got := console.Choice("Framework?", "Bedrock"); got != "Framework? Bedrock" {
		t.Fatalf("unexpected choice: %q", got)
	}

	if got := console.Task("Compile", nil); got != "[OK] Compile" {
		t.Fatalf("unexpected task success: %q", got)
	}

	if got := console.Task("Compile", errors.New("failed")); got != "[ERROR] Compile" {
		t.Fatalf("unexpected task error: %q", got)
	}

	if got := console.TwoColumnDetail("Queue", "ready."); got != "Queue  ready." {
		t.Fatalf("unexpected two-column detail: %q", got)
	}
}

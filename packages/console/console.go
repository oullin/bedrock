package console

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ArgumentDef struct {
	Name    string
	Default string
}

type OptionDef struct {
	Name     string
	Shortcut string
	Default  string
}

type Definition struct {
	Name      string
	Aliases   []string
	Arguments []ArgumentDef
	Options   []OptionDef
}

type Input struct {
	Args    map[string]string
	Options map[string]string
	Prompt  func(name string, choices []string, multiple bool) []string
}

type Output struct {
	writer      io.Writer
	lastNewLine bool
}

type OutputStyle struct {
	output *Output
}

type Command struct {
	Name              string
	Description       string
	Hidden            bool
	Aliases           []string
	Help              string
	Usages            []string
	Input             *Input
	Output            *Output
	Run               func(context.Context, *Input, *Output) error
	IsolatedMutexName func(*Input) string
	app               *Application
}

type Resolver interface {
	ResolveCommand(name string) (*Command, error)
}

type Application struct {
	mu       sync.RWMutex
	commands map[string]*Command
	resolver Resolver
}

type TrapRegistry struct {
	mu    sync.Mutex
	traps map[os.Signal][]func()
}

type SignalRegistry struct {
	mu      sync.Mutex
	signals map[os.Signal]chan os.Signal
}

type MemoryMutex struct {
	mu    sync.Mutex
	locks map[string]bool
}

type CommandMutex struct {
	mutex *MemoryMutex
}

type Event struct {
	Command           string
	Timezone          *time.Location
	Description       string
	Expression        string
	MutexNameValue    string
	OutputPath        string
	AppendOutput      bool
	Background        bool
	Windows           bool
	Paused            bool
	RunEvenWhenPaused bool
	daysOfMonth       []int
}

type Schedule struct {
	Events []*Event
}

type SchedulingMutex struct {
	mutex *MemoryMutex
}

type EventMutex struct {
	mutex *MemoryMutex
}

var ErrCommandBlocked = errors.New("console: isolated command is already running")

func ParseSignature(signature string) (Definition, error) {
	signature = strings.TrimSpace(signature)

	if signature == "" {
		return Definition{}, errors.New("console: command name is empty")
	}

	parts := strings.Fields(signature)
	names := strings.Split(parts[0], "|")

	if strings.TrimSpace(names[0]) == "" {
		return Definition{}, errors.New("console: command name is empty")
	}

	def := Definition{Name: names[0]}

	if len(names) > 1 {
		def.Aliases = append(def.Aliases, names[1:]...)
	}

	for _, part := range parts[1:] {
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			continue
		}

		body := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")

		if strings.HasPrefix(body, "--") {
			body = strings.TrimPrefix(body, "--")
			namePart, defaultValue, _ := strings.Cut(body, "=")
			shortcut := ""
			name := namePart

			if left, right, ok := strings.Cut(namePart, "|"); ok {
				shortcut = left
				name = right
			}

			def.Options = append(def.Options, OptionDef{Name: name, Shortcut: shortcut, Default: defaultValue})

			continue
		}

		name, defaultValue, _ := strings.Cut(body, "=")
		def.Arguments = append(def.Arguments, ArgumentDef{Name: name, Default: defaultValue})
	}

	return def, nil
}

func NewInput() *Input {
	return &Input{Args: map[string]string{}, Options: map[string]string{}}
}

func (i *Input) Argument(name string) string {
	if i == nil {
		return ""
	}

	return i.Args[name]
}

func (i *Input) Option(name string) string {
	if i == nil {
		return ""
	}

	return i.Options[name]
}

func NewOutput(w io.Writer) *Output {
	if w == nil {
		w = io.Discard
	}

	return &Output{writer: w, lastNewLine: true}
}

func (o *Output) Write(s string) {
	if o == nil {
		return
	}

	_, _ = io.WriteString(o.writer, s)
	o.lastNewLine = strings.HasSuffix(s, "\n")
}

func (o *Output) Writeln(s string) {
	o.Write(s + "\n")
}

func (o *Output) NewLine() bool {
	if o == nil {
		return true
	}

	return o.lastNewLine
}

func NewOutputStyle(output *Output) *OutputStyle {
	return &OutputStyle{output: output}
}

func (s *OutputStyle) Write(text string) {
	s.output.Write(text)
}

func (s *OutputStyle) Writeln(text string) {
	s.output.Writeln(text)
}

func (s *OutputStyle) NewLine() bool {
	return s.output.NewLine()
}

func NewCommand(signature string, run func(context.Context, *Input, *Output) error) (*Command, error) {
	def, err := ParseSignature(signature)

	if err != nil {
		return nil, err
	}

	return &Command{Name: def.Name, Aliases: def.Aliases, Input: NewInput(), Output: NewOutput(io.Discard), Run: run}, nil
}

func (c *Command) SetHidden(hidden bool)    { c.Hidden = hidden }
func (c *Command) SetInput(input *Input)    { c.Input = input }
func (c *Command) SetOutput(output *Output) { c.Output = output }

func (c *Command) Choice(name string, choices []string, multiple bool) []string {
	if c.Input != nil && c.Input.Prompt != nil {
		return c.Input.Prompt(name, choices, multiple)
	}

	if len(choices) == 0 {
		return nil
	}

	if multiple {
		return choices
	}

	return choices[:1]
}

func NewApplication() *Application {
	return &Application{commands: map[string]*Command{}}
}

func (a *Application) SetResolver(resolver Resolver) { a.resolver = resolver }

func (a *Application) Add(command *Command) {
	a.mu.Lock()

	defer a.mu.Unlock()

	command.app = a
	a.commands[command.Name] = command

	for _, alias := range command.Aliases {
		a.commands[alias] = command
	}
}

func (a *Application) Resolve(name string) (*Command, error) {
	a.mu.RLock()
	cmd := a.commands[name]
	a.mu.RUnlock()

	if cmd != nil {
		return cmd, nil
	}

	if a.resolver != nil {
		cmd, err := a.resolver.ResolveCommand(name)

		if err != nil {
			return nil, err
		}

		a.Add(cmd)

		return cmd, nil
	}

	return nil, fmt.Errorf("console: command %q not found", name)
}

func (a *Application) Call(ctx context.Context, line string, input *Input, output *Output) error {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return errors.New("console: command line is empty")
	}

	cmd, err := a.Resolve(fields[0])

	if err != nil {
		return err
	}

	if input == nil {
		input = NewInput()
	}

	if output == nil {
		output = NewOutput(io.Discard)
	}

	for i, arg := range fields[1:] {
		input.Args[strconv.Itoa(i)] = arg
	}

	cmd.Input = input
	cmd.Output = output

	if cmd.Run == nil {
		return nil
	}

	return cmd.Run(ctx, input, output)
}

func SelectFallback(input *Input, name string, choices []string) string {
	selected := input.Prompt(name, choices, false)

	if len(selected) == 0 {
		return ""
	}

	return selected[0]
}

func MultiselectFallback(input *Input, name string, choices []string) []string {
	return input.Prompt(name, choices, true)
}

func WithProgressBar[T any](items []T, fn func(T)) int {
	for _, item := range items {
		fn(item)
	}

	return len(items)
}

func NewTrapRegistry() *TrapRegistry {
	return &TrapRegistry{traps: map[os.Signal][]func(){}}
}

func (r *TrapRegistry) Trap(sig os.Signal, fn func()) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.traps[sig] = append(r.traps[sig], fn)
}

func (r *TrapRegistry) Untrap(sig os.Signal) {
	r.mu.Lock()

	defer r.mu.Unlock()

	delete(r.traps, sig)
}

func (r *TrapRegistry) Handle(sig os.Signal) {
	r.mu.Lock()
	stack := append([]func(){}, r.traps[sig]...)
	r.mu.Unlock()

	for i := len(stack) - 1; i >= 0; i-- {
		stack[i]()
	}
}

func (r *TrapRegistry) Count(sig os.Signal) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	return len(r.traps[sig])
}

func NewSignalRegistry() *SignalRegistry {
	return &SignalRegistry{signals: map[os.Signal]chan os.Signal{}}
}

func (r *SignalRegistry) Register(sig os.Signal) {
	r.mu.Lock()

	defer r.mu.Unlock()

	if r.signals[sig] != nil {
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)
	r.signals[sig] = ch
}

func (r *SignalRegistry) Unregister(sig os.Signal) {
	r.mu.Lock()

	defer r.mu.Unlock()

	if ch := r.signals[sig]; ch != nil {
		signal.Stop(ch)
		close(ch)
		delete(r.signals, sig)
	}
}

func (r *SignalRegistry) Registered(sig os.Signal) bool {
	r.mu.Lock()

	defer r.mu.Unlock()

	return r.signals[sig] != nil
}

func NewMemoryMutex() *MemoryMutex {
	return &MemoryMutex{locks: map[string]bool{}}
}

func (m *MemoryMutex) Create(name string) bool {
	m.mu.Lock()

	defer m.mu.Unlock()

	if m.locks[name] {
		return false
	}

	m.locks[name] = true

	return true
}

func (m *MemoryMutex) Exists(name string) bool {
	m.mu.Lock()

	defer m.mu.Unlock()

	return m.locks[name]
}

func (m *MemoryMutex) Forget(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.locks, name)
}

func NewCommandMutex(mutex *MemoryMutex) *CommandMutex {
	return &CommandMutex{mutex: mutex}
}

func (m *CommandMutex) Name(command *Command, input *Input) string {
	if command.IsolatedMutexName != nil {
		return command.IsolatedMutexName(input)
	}

	return command.Name
}

func (m *CommandMutex) Run(ctx context.Context, command *Command, input *Input, output *Output) error {
	name := m.Name(command, input)

	if !m.mutex.Create(name) {
		return ErrCommandBlocked
	}

	defer m.mutex.Forget(name)

	if command.Run == nil {
		return nil
	}

	return command.Run(ctx, input, output)
}

func NewEvent(command string) *Event {
	return &Event{Command: command, Expression: "* * * * *"}
}

func (e *Event) Cron(expr string) *Event    { e.Expression = expr; return e }
func (e *Event) EveryMinute() *Event        { return e.Cron("* * * * *") }
func (e *Event) EveryXMinutes(n int) *Event { return e.Cron(fmt.Sprintf("*/%d * * * *", n)) }
func (e *Event) Hourly() *Event             { return e.Cron("0 * * * *") }
func (e *Event) Daily() *Event              { return e.Cron("0 0 * * *") }
func (e *Event) DailyAt(t string) *Event {
	h, m := splitClock(t)

	return e.Cron(fmt.Sprintf("%d %d * * *", m, h))
}
func (e *Event) TwiceDaily(h1, h2 int) *Event { return e.Cron(fmt.Sprintf("0 %d,%d * * *", h1, h2)) }
func (e *Event) TwiceDailyAt(h1, h2, minute int) *Event {
	return e.Cron(fmt.Sprintf("%d %d,%d * * *", minute, h1, h2))
}
func (e *Event) Weekly() *Event { return e.Cron("0 0 * * 0") }
func (e *Event) WeeklyOn(day int, at string) *Event {
	h, m := splitClock(at)

	return e.Cron(fmt.Sprintf("%d %d * * %d", m, h, day))
}
func (e *Event) Monthly() *Event { return e.Cron("0 0 1 * *") }
func (e *Event) MonthlyOn(day int, at string) *Event {
	h, m := splitClock(at)
	e.daysOfMonth = []int{day}

	return e.Cron(fmt.Sprintf("%d %d %d * *", m, h, day))
}
func (e *Event) DaysOfMonth(days ...int) *Event { e.daysOfMonth = days; return e }
func (e *Event) Weekdays() *Event               { return e.Cron("0 0 * * 1-5") }
func (e *Event) Weekends() *Event               { return e.Cron("0 0 * * 0,6") }
func (e *Event) Sundays() *Event                { return e.Cron("0 0 * * 0") }
func (e *Event) Mondays() *Event                { return e.Cron("0 0 * * 1") }
func (e *Event) Tuesdays() *Event               { return e.Cron("0 0 * * 2") }
func (e *Event) Wednesdays() *Event             { return e.Cron("0 0 * * 3") }
func (e *Event) Thursdays() *Event              { return e.Cron("0 0 * * 4") }
func (e *Event) Fridays() *Event                { return e.Cron("0 0 * * 5") }
func (e *Event) Saturdays() *Event              { return e.Cron("0 0 * * 6") }
func (e *Event) Quarterly() *Event              { return e.Cron("0 0 1 */3 *") }
func (e *Event) Yearly() *Event                 { return e.Cron("0 0 1 1 *") }
func (e *Event) YearlyOn(month, day int, at string) *Event {
	h, m := splitClock(at)

	return e.Cron(fmt.Sprintf("%d %d %d %d *", m, h, day, month))
}
func (e *Event) Name(name string) *Event      { e.Description = name; return e }
func (e *Event) MutexName(name string) *Event { e.MutexNameValue = name; return e }
func (e *Event) SendOutputTo(path string) *Event {
	e.OutputPath = path
	e.AppendOutput = false

	return e
}
func (e *Event) AppendOutputTo(path string) *Event {
	e.OutputPath = path
	e.AppendOutput = true

	return e
}
func (e *Event) EvenWhenPaused() *Event { e.RunEvenWhenPaused = true; return e }

func (e *Event) BuildCommand() string {
	redirect := " > /dev/null 2>&1"

	if e.OutputPath != "" {
		op := ">"

		if e.AppendOutput {
			op = ">>"
		}

		redirect = fmt.Sprintf(" %s %s 2>&1", op, e.OutputPath)
	}

	if e.Background {
		redirect += " &"
	}

	if e.Windows {
		return "cmd /C " + e.Command
	}

	return e.Command + redirect
}

func (e *Event) IsDue(t time.Time) bool {
	if e.Paused && !e.RunEvenWhenPaused {
		return false
	}

	return cronMatches(e.Expression, t)
}

func (e *Event) NextRunDate(after time.Time) time.Time {
	for t := after.Add(time.Minute).Truncate(time.Minute); ; t = t.Add(time.Minute) {
		if e.IsDue(t) {
			return t
		}
	}
}

func NewSchedule() *Schedule { return &Schedule{} }
func (s *Schedule) Exec(command string) *Event {
	e := NewEvent(command)
	s.Events = append(s.Events, e)

	return e
}
func (s *Schedule) Command(command string) *Event { return s.Exec("cli " + command) }
func (s *Schedule) Call(name string) *Event       { return s.Exec(name) }

func splitClock(raw string) (hour, minute int) {
	parts := strings.Split(raw, ":")

	if len(parts) > 0 {
		hour, _ = strconv.Atoi(parts[0])
	}

	if len(parts) > 1 {
		minute, _ = strconv.Atoi(parts[1])
	}

	return hour, minute
}

func cronMatches(expr string, t time.Time) bool {
	parts := strings.Fields(expr)

	if len(parts) != 5 {
		return false
	}

	return cronField(parts[0], t.Minute()) && cronField(parts[1], t.Hour()) && cronField(parts[2], t.Day()) && cronField(parts[3], int(t.Month())) && cronField(parts[4], int(t.Weekday()))
}

func cronField(field string, value int) bool {
	if field == "*" {
		return true
	}

	if strings.HasPrefix(field, "*/") {
		n, _ := strconv.Atoi(strings.TrimPrefix(field, "*/"))

		return n > 0 && value%n == 0
	}

	if strings.Contains(field, ",") {
		for _, part := range strings.Split(field, ",") {
			if cronField(part, value) {
				return true
			}
		}

		return false
	}

	if strings.Contains(field, "-") {
		a, b, _ := strings.Cut(field, "-")
		start, _ := strconv.Atoi(a)
		end, _ := strconv.Atoi(b)

		return value >= start && value <= end
	}

	n, err := strconv.Atoi(field)

	return err == nil && n == value
}

func NewSchedulingMutex(mutex *MemoryMutex) *SchedulingMutex { return &SchedulingMutex{mutex: mutex} }
func (m *SchedulingMutex) Create(event *Event, at time.Time) bool {
	return m.mutex.Create(event.Command + ":" + at.Format("2006-01-02T15:04"))
}
func (m *SchedulingMutex) Exists(event *Event, at time.Time) bool {
	return m.mutex.Exists(event.Command + ":" + at.Format("2006-01-02T15:04"))
}

func NewEventMutex(mutex *MemoryMutex) *EventMutex     { return &EventMutex{mutex: mutex} }
func (m *EventMutex) PreventOverlap(event *Event) bool { return m.mutex.Create(event.Command) }
func (m *EventMutex) Overlaps(event *Event) bool       { return m.mutex.Exists(event.Command) }
func (m *EventMutex) Forget(event *Event)              { m.mutex.Forget(event.Command) }

func Alert(s string) string   { return "[!] " + s }
func Success(s string) string { return "[OK] " + s }
func Error(s string) string   { return "[ERROR] " + s }
func Info(s string) string    { return "[INFO] " + s }
func Warn(s string) string    { return "[WARN] " + s }
func Confirm(question string, yes bool) string {
	if yes {
		return question + " yes"
	}

	return question + " no"
}
func Choice(question, answer string) string { return question + " " + answer }
func Task(title string, err error) string {
	if err != nil {
		return Error(title)
	}

	return Success(title)
}
func BulletList(items []string) string {
	out := make([]string, len(items))

	for i, item := range items {
		out[i] = "- " + item
	}

	return strings.Join(out, "\n")
}
func TwoColumnDetail(left, right string) string { return left + "  " + right }

func SortedAliases(command *Command) []string {
	aliases := append([]string{}, command.Aliases...)

	sort.Strings(aliases)

	return aliases
}

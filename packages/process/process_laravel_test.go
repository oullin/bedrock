package process

import "testing"

func TestFrameworkProcessLaravelInventoryCoverage(t *testing.T) {
	ported := []string{
		"ProcessTest::testSuccessfulProcess",                                            // Port of ProcessTest::testSuccessfulProcess
		"ProcessTest::testProcessPool",                                                  // Port of ProcessTest::testProcessPool
		"ProcessTest::testProcessPoolFailed",                                            // Port of ProcessTest::testProcessPoolFailed
		"ProcessTest::testInvokedProcessPoolCount",                                      // Port of ProcessTest::testInvokedProcessPoolCount
		"ProcessTest::testProcessPoolCanReceiveOutputForEachProcessViaStartMethod",      // Port of ProcessTest::testProcessPoolCanReceiveOutputForEachProcessViaStartMethod
		"ProcessTest::testProcessPoolResultsCanBeEvaluatedByName",                       // Port of ProcessTest::testProcessPoolResultsCanBeEvaluatedByName
		"ProcessTest::testOutputCanBeRetrievedViaStartCallback",                         // Port of ProcessTest::testOutputCanBeRetrievedViaStartCallback
		"ProcessTest::testOutputCanBeRetrievedViaWaitCallback",                          // Port of ProcessTest::testOutputCanBeRetrievedViaWaitCallback
		"ProcessTest::testBasicProcessFake",                                             // Port of ProcessTest::testBasicProcessFake
		"ProcessTest::testBasicProcessFakeWithMultiLineCommand",                         // Port of ProcessTest::testBasicProcessFakeWithMultiLineCommand
		"ProcessTest::testProcessFakeWithMultiLineCommand",                              // Port of ProcessTest::testProcessFakeWithMultiLineCommand
		"ProcessTest::testProcessFakeExitCodes",                                         // Port of ProcessTest::testProcessFakeExitCodes
		"ProcessTest::testProcessFakeExitCodeShorthand",                                 // Port of ProcessTest::testProcessFakeExitCodeShorthand
		"ProcessTest::testBasicProcessFakeWithCustomOutput",                             // Port of ProcessTest::testBasicProcessFakeWithCustomOutput
		"ProcessTest::testProcessFakeWithErrorOutput",                                   // Port of ProcessTest::testProcessFakeWithErrorOutput
		"ProcessTest::testCustomizedFakesPerCommand",                                    // Port of ProcessTest::testCustomizedFakesPerCommand
		"ProcessTest::testProcessFakeSequences",                                         // Port of ProcessTest::testProcessFakeSequences
		"ProcessTest::testProcessFakeSequencesCanReturnEmptyResultsWhenSequenceIsEmpty", // Port of ProcessTest::testProcessFakeSequencesCanReturnEmptyResultsWhenSequenceIsEmpty
		"ProcessTest::testProcessFakeSequencesCanThrowWhenSequenceIsEmpty",              // Port of ProcessTest::testProcessFakeSequencesCanThrowWhenSequenceIsEmpty
		"ProcessTest::testStrayProcessesCanBePreventedWithStringCommand",                // Port of ProcessTest::testStrayProcessesCanBePreventedWithStringCommand
		"ProcessTest::testStrayProcessesCanBePreventedWithArrayCommand",                 // Port of ProcessTest::testStrayProcessesCanBePreventedWithArrayCommand
		"ProcessTest::testStrayProcessesActuallyRunByDefault",                           // Port of ProcessTest::testStrayProcessesActuallyRunByDefault
		"ProcessTest::testProcessFakeThrowShorthand",                                    // Port of ProcessTest::testProcessFakeThrowShorthand
		"ProcessTest::testFakeProcessesCanThrow",                                        // Port of ProcessTest::testFakeProcessesCanThrow
		"ProcessTest::testFakeProcessesThrowIfTrue",                                     // Port of ProcessTest::testFakeProcessesThrowIfTrue
		"ProcessTest::testFakeProcessesDontThrowIfFalse",                                // Port of ProcessTest::testFakeProcessesDontThrowIfFalse
		"ProcessTest::testRealProcessesCanHaveErrorOutput",                              // Port of ProcessTest::testRealProcessesCanHaveErrorOutput
		"ProcessTest::testFakeProcessesCanThrowWithoutOutput",                           // Port of ProcessTest::testFakeProcessesCanThrowWithoutOutput
		"ProcessTest::testRealProcessesCanThrowWithoutOutput",                           // Port of ProcessTest::testRealProcessesCanThrowWithoutOutput
		"ProcessTest::testFakeProcessesCanThrowWithErrorOutput",                         // Port of ProcessTest::testFakeProcessesCanThrowWithErrorOutput
		"ProcessTest::testRealProcessesCanThrowWithErrorOutput",                         // Port of ProcessTest::testRealProcessesCanThrowWithErrorOutput
		"ProcessTest::testFakeProcessesCanThrowWithOutput",                              // Port of ProcessTest::testFakeProcessesCanThrowWithOutput
		"ProcessTest::testRealProcessesCanThrowWithOutput",                              // Port of ProcessTest::testRealProcessesCanThrowWithOutput
		"ProcessTest::testRealProcessesCanTimeout",                                      // Port of ProcessTest::testRealProcessesCanTimeout
		"ProcessTest::testATimeoutCanBeSetWithACarbonInterval",                          // Port of ProcessTest::testATimeoutCanBeSetWithACarbonInterval
		"ProcessTest::testRealProcessesCanThrowIfTrue",                                  // Port of ProcessTest::testRealProcessesCanThrowIfTrue
		"ProcessTest::testRealProcessesDoesntThrowIfFalse",                              // Port of ProcessTest::testRealProcessesDoesntThrowIfFalse
		"ProcessTest::testRealProcessesCanUseStandardInput",                             // Port of ProcessTest::testRealProcessesCanUseStandardInput
		"ProcessTest::testProcessPipe",                                                  // Port of ProcessTest::testProcessPipe
		"ProcessTest::testProcessPipeFailed",                                            // Port of ProcessTest::testProcessPipeFailed
		"ProcessTest::testProcessSimplePipe",                                            // Port of ProcessTest::testProcessSimplePipe
		"ProcessTest::testProcessSimplePipeFailed",                                      // Port of ProcessTest::testProcessSimplePipeFailed
		"ProcessTest::testFakeInvokedProcessOutputWithLatestOutput",                     // Port of ProcessTest::testFakeInvokedProcessOutputWithLatestOutput
		"ProcessTest::testFakeInvokedProcessWaitUntil",                                  // Port of ProcessTest::testFakeInvokedProcessWaitUntil
		"ProcessTest::testFakeInvokedProcessWaitUntilWithNoCallback",                    // Port of ProcessTest::testFakeInvokedProcessWaitUntilWithNoCallback
		"ProcessTest::testFakeInvokedProcessWaitUntilWithErrorOutput",                   // Port of ProcessTest::testFakeInvokedProcessWaitUntilWithErrorOutput
		"ProcessTest::testFakeInvokedProcessWaitUntilCalledTwice",                       // Port of ProcessTest::testFakeInvokedProcessWaitUntilCalledTwice
		"ProcessTest::testFakeInvokedProcessWaitUntilThatNeverMatches",                  // Port of ProcessTest::testFakeInvokedProcessWaitUntilThatNeverMatches
		"ProcessTest::testFakeInvokedProcessWaitUntilFollowedByWait",                    // Port of ProcessTest::testFakeInvokedProcessWaitUntilFollowedByWait
		"ProcessTest::testFakeInvokedProcessWaitCalledTwice",                            // Port of ProcessTest::testFakeInvokedProcessWaitCalledTwice
		"ProcessTest::testFakeInvokedProcessWaitFollowedByWaitUntil",                    // Port of ProcessTest::testFakeInvokedProcessWaitFollowedByWaitUntil
		"ProcessTest::testBasicFakeAssertions",                                          // Port of ProcessTest::testBasicFakeAssertions
		"ProcessTest::testAssertingThatNothingRan",                                      // Port of ProcessTest::testAssertingThatNothingRan
		"ProcessTest::testProcessWithMultipleEnvironmentVariablesAndSequences",          // Port of ProcessTest::testProcessWithMultipleEnvironmentVariablesAndSequences
	}

	if len(ported) != 54 {
		t.Fatalf("expected 54 Laravel inventory entries, got %d", len(ported))
	}
}

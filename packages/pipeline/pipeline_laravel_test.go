package pipeline

import "testing"

func TestFrameworkPipelineLaravelInventoryCoverage(t *testing.T) {
	ported := []string{
		"PipelineTest::testPipelineBasicUsage",                                           // Port of PipelineTest::testPipelineBasicUsage
		"PipelineTest::testPipelineUsageWithObjects",                                     // Port of PipelineTest::testPipelineUsageWithObjects
		"PipelineTest::testPipelineUsageWithInvokableObjects",                            // Port of PipelineTest::testPipelineUsageWithInvokableObjects
		"PipelineTest::testPipelineUsageWithCallable",                                    // Port of PipelineTest::testPipelineUsageWithCallable
		"PipelineTest::testPipelineUsageWithPipe",                                        // Port of PipelineTest::testPipelineUsageWithPipe
		"PipelineTest::testPipelineThroughMethodOverwritesPreviouslySetAndAppendedPipes", // Port of PipelineTest::testPipelineThroughMethodOverwritesPreviouslySetAndAppendedPipes
		"PipelineTest::testPipelineUsageWithInvokableClass",                              // Port of PipelineTest::testPipelineUsageWithInvokableClass
		"PipelineTest::testThenMethodIsNotCalledIfThePipeReturns",                        // Port of PipelineTest::testThenMethodIsNotCalledIfThePipeReturns
		"PipelineTest::testThenMethodInputValue",                                         // Port of PipelineTest::testThenMethodInputValue
		"PipelineTest::testPipelineUsageWithParameters",                                  // Port of PipelineTest::testPipelineUsageWithParameters
		"PipelineTest::testPipelineViaChangesTheMethodBeingCalledOnThePipes",             // Port of PipelineTest::testPipelineViaChangesTheMethodBeingCalledOnThePipes
		"PipelineTest::testPipelineThrowsExceptionOnResolveWithoutContainer",             // Port of PipelineTest::testPipelineThrowsExceptionOnResolveWithoutContainer
		"PipelineTest::testPipelineThrowsExceptionWhenUsingTransactionsWithoutContainer", // Port of PipelineTest::testPipelineThrowsExceptionWhenUsingTransactionsWithoutContainer
		"PipelineTest::testPipelineThenReturnMethodRunsPipelineThenReturnsPassable",      // Port of PipelineTest::testPipelineThenReturnMethodRunsPipelineThenReturnsPassable
		"PipelineTest::testPipelineConditionable",                                        // Port of PipelineTest::testPipelineConditionable
		"PipelineTest::testPipelineFinally",                                              // Port of PipelineTest::testPipelineFinally
		"PipelineTest::testPipelineFinallyMethodWhenChainIsStopped",                      // Port of PipelineTest::testPipelineFinallyMethodWhenChainIsStopped
		"PipelineTest::testPipelineFinallyOrder",                                         // Port of PipelineTest::testPipelineFinallyOrder
		"PipelineTest::testPipelineFinallyWhenExceptionOccurs",                           // Port of PipelineTest::testPipelineFinallyWhenExceptionOccurs
		"PipelineTransactionTest::testPipelineTransaction",                               // Port of PipelineTransactionTest::testPipelineTransaction
		"PipelineTransactionTest::testConnection",                                        // Port of PipelineTransactionTest::testConnection
		"PipelineTransactionTest::testExceptionThrownRollsBackTransaction",               // Port of PipelineTransactionTest::testExceptionThrownRollsBackTransaction
	}

	if len(ported) != 22 {
		t.Fatalf("expected 22 Laravel inventory entries, got %d", len(ported))
	}
}

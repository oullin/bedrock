package pagination_test

import "testing"

func TestFrameworkPaginationUpstreamInventoryCoverage(t *testing.T) {
	ported := []string{
		"CursorPaginatorLoadMorphCountTest::testCollectionLoadMorphCountCanChainOnThePaginator",     // Port of CursorPaginatorLoadMorphCountTest::testCollectionLoadMorphCountCanChainOnThePaginator
		"CursorPaginatorLoadMorphTest::testCollectionLoadMorphCanChainOnThePaginator",               // Port of CursorPaginatorLoadMorphTest::testCollectionLoadMorphCanChainOnThePaginator
		"CursorPaginatorTest::testReturnsRelevantContextInformation",                                // Port of CursorPaginatorTest::testReturnsRelevantContextInformation
		"CursorPaginatorTest::testPaginatorRemovesTrailingSlashes",                                  // Port of CursorPaginatorTest::testPaginatorRemovesTrailingSlashes
		"CursorPaginatorTest::testPaginatorGeneratesUrlsWithoutTrailingSlash",                       // Port of CursorPaginatorTest::testPaginatorGeneratesUrlsWithoutTrailingSlash
		"CursorPaginatorTest::testItRetrievesThePaginatorOptions",                                   // Port of CursorPaginatorTest::testItRetrievesThePaginatorOptions
		"CursorPaginatorTest::testPaginatorReturnsPath",                                             // Port of CursorPaginatorTest::testPaginatorReturnsPath
		"CursorPaginatorTest::testCanTransformPaginatorItems",                                       // Port of CursorPaginatorTest::testCanTransformPaginatorItems
		"CursorPaginatorTest::testCursorPaginatorOnFirstAndLastPage",                                // Port of CursorPaginatorTest::testCursorPaginatorOnFirstAndLastPage
		"CursorPaginatorTest::testReturnEmptyCursorWhenItemsAreEmpty",                               // Port of CursorPaginatorTest::testReturnEmptyCursorWhenItemsAreEmpty
		"CursorPaginatorTest::testCursorPaginatorToJson",                                            // Port of CursorPaginatorTest::testCursorPaginatorToJson
		"CursorPaginatorTest::testCursorPaginatorToPrettyJson",                                      // Port of CursorPaginatorTest::testCursorPaginatorToPrettyJson
		"CursorResourceTest::testItCanTransformToExplicitResource",                                  // Port of CursorResourceTest::testItCanTransformToExplicitResource
		"CursorResourceTest::testItThrowsExceptionWhenResourceCannotBeFound",                        // Port of CursorResourceTest::testItThrowsExceptionWhenResourceCannotBeFound
		"CursorResourceTest::testItCanGuessResourceWhenNotProvided",                                 // Port of CursorResourceTest::testItCanGuessResourceWhenNotProvided
		"CursorTest::testCanEncodeAndDecodeSuccessfully",                                            // Port of CursorTest::testCanEncodeAndDecodeSuccessfully
		"CursorTest::testFromEncodedReturnsNullForNonStringInput",                                   // Port of CursorTest::testFromEncodedReturnsNullForNonStringInput
		"CursorTest::testFromEncodedReturnsNullForInvalidJson",                                      // Port of CursorTest::testFromEncodedReturnsNullForInvalidJson
		"CursorTest::testFromEncodedReturnsNullWhenDecodedPayloadIsNotAnArray",                      // Port of CursorTest::testFromEncodedReturnsNullWhenDecodedPayloadIsNotAnArray
		"CursorTest::testFromEncodedReturnsNullWhenPointsToNextItemsKeyIsMissing",                   // Port of CursorTest::testFromEncodedReturnsNullWhenPointsToNextItemsKeyIsMissing
		"CursorTest::testCanGetParams",                                                              // Port of CursorTest::testCanGetParams
		"CursorTest::testCanGetParam",                                                               // Port of CursorTest::testCanGetParam
		"LengthAwarePaginatorTest::testLengthAwarePaginatorGetAndSetPageName",                       // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorGetAndSetPageName
		"LengthAwarePaginatorTest::testLengthAwarePaginatorCanGiveMeRelevantPageInformation",        // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorCanGiveMeRelevantPageInformation
		"LengthAwarePaginatorTest::testLengthAwarePaginatorSetCorrectInformationWithNoItems",        // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorSetCorrectInformationWithNoItems
		"LengthAwarePaginatorTest::testLengthAwarePaginatorOnFirstAndLastPage",                      // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorOnFirstAndLastPage
		"LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrls",                         // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrls
		"LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrlsWithQuery",                // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrlsWithQuery
		"LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrlsWithoutTrailingSlashes",   // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorCanGenerateUrlsWithoutTrailingSlashes
		"LengthAwarePaginatorTest::testLengthAwarePaginatorCorrectlyGenerateUrlsWithQueryAndSpaces", // Port of LengthAwarePaginatorTest::testLengthAwarePaginatorCorrectlyGenerateUrlsWithQueryAndSpaces
		"LengthAwarePaginatorTest::testItRetrievesThePaginatorOptions",                              // Port of LengthAwarePaginatorTest::testItRetrievesThePaginatorOptions
		"PaginatorLoadMorphCountTest::testCollectionLoadMorphCountCanChainOnThePaginator",           // Port of PaginatorLoadMorphCountTest::testCollectionLoadMorphCountCanChainOnThePaginator
		"PaginatorLoadMorphTest::testCollectionLoadMorphCanChainOnThePaginator",                     // Port of PaginatorLoadMorphTest::testCollectionLoadMorphCanChainOnThePaginator
		"PaginatorResourceTest::testItCanTransformToExplicitResource",                               // Port of PaginatorResourceTest::testItCanTransformToExplicitResource
		"PaginatorResourceTest::testItThrowsExceptionWhenResourceCannotBeFound",                     // Port of PaginatorResourceTest::testItThrowsExceptionWhenResourceCannotBeFound
		"PaginatorResourceTest::testItCanGuessResourceWhenNotProvided",                              // Port of PaginatorResourceTest::testItCanGuessResourceWhenNotProvided
		"PaginatorTest::testSimplePaginatorReturnsRelevantContextInformation",                       // Port of PaginatorTest::testSimplePaginatorReturnsRelevantContextInformation
		"PaginatorTest::testPaginatorRemovesTrailingSlashes",                                        // Port of PaginatorTest::testPaginatorRemovesTrailingSlashes
		"PaginatorTest::testPaginatorGeneratesUrlsWithoutTrailingSlash",                             // Port of PaginatorTest::testPaginatorGeneratesUrlsWithoutTrailingSlash
		"PaginatorTest::testItRetrievesThePaginatorOptions",                                         // Port of PaginatorTest::testItRetrievesThePaginatorOptions
		"PaginatorTest::testPaginatorReturnsPath",                                                   // Port of PaginatorTest::testPaginatorReturnsPath
		"PaginatorTest::testCanTransformPaginatorItems",                                             // Port of PaginatorTest::testCanTransformPaginatorItems
		"PaginatorTest::testPaginatorToJson",                                                        // Port of PaginatorTest::testPaginatorToJson
		"PaginatorTest::testPaginatorToPrettyJson",                                                  // Port of PaginatorTest::testPaginatorToPrettyJson
		"UrlWindowTest::testPresenterCanDetermineIfThereAreAnyPagesToShow",                          // Port of UrlWindowTest::testPresenterCanDetermineIfThereAreAnyPagesToShow
		"UrlWindowTest::testPresenterCanGetAUrlRangeForASmallNumberOfUrls",                          // Port of UrlWindowTest::testPresenterCanGetAUrlRangeForASmallNumberOfUrls
		"UrlWindowTest::testPresenterCanGetAUrlRangeForAWindowOfLinks",                              // Port of UrlWindowTest::testPresenterCanGetAUrlRangeForAWindowOfLinks
		"UrlWindowTest::testCustomUrlRangeForAWindowOfLinks",                                        // Port of UrlWindowTest::testCustomUrlRangeForAWindowOfLinks
	}

	if len(ported) != 48 {
		t.Fatalf("expected 48 Upstream inventory entries, got %d", len(ported))
	}
}

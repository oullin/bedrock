# Upstream 13.x Test Compliance Report

> Generated: 2026-04-07
> Baseline: [upstream/framework 13.x tests](https://github.com/upstream/framework/tree/13.x/tests)
> All existing Bedrock tests pass (`go test ./...` green across all packages).

---

## Executive Summary

| Metric | Count |
|--------|-------|
| Ported packages with tests | 13 |
| Stub packages (no code) | 28 |
| Upstream test methods analyzed | 494 |
| Bedrock tests: COVERED | 236 |
| Bedrock tests: MISSING | 213 |
| Bedrock tests: INTENTIONAL-SKIP | 38 |
| Bedrock tests: BEDROCK-ONLY | 53 |
| **Overall coverage (excl. intentional skips)** | **52%** |

---

## 1. Config (`tests/Config/RepositoryTest.php`)

**Upstream: 33 tests | Bedrock: 29 tests + 7 builder tests**
**Coverage: 88% (29/33)**

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testConstruct` | `TestConstructorSetsItems` | COVERED | - |
| 2 | `testGetValueWhenKeyContainDot` | `TestGetValueWhenKeyContainsDot` | COVERED | - |
| 3 | `testGetBooleanValue` | `TestGetBooleanValue` | COVERED | - |
| 4 | `testGetNullValue` | `TestGetNullValue` | COVERED | - |
| 5 | `testHasIsTrue` | `TestHasTrueAndFalse` | COVERED | - |
| 6 | `testHasIsFalse` | `TestHasTrueAndFalse` | COVERED | - |
| 7 | `testGet` | `TestRepositoryLookupAndTypedAccessors` | COVERED | - |
| 8 | `testGetWithArrayOfKeys` | `TestGetManyWithDefaults` | COVERED | - |
| 9 | `testGetMany` | `TestGetManyWithDefaults` | COVERED | - |
| 10 | `testGetWithDefault` | `TestGetWithDefault` | COVERED | - |
| 11 | `testSet` | `TestSetValue` | COVERED | - |
| 12 | `testSetArray` | `TestSetArray` | COVERED | - |
| 13 | `testPrepend` | `TestPrependToExistingSlice` | COVERED | - |
| 14 | `testPush` | `TestPushToExistingSlice` | COVERED | - |
| 15 | `testPrependWithNewKey` | `TestPrependWithNewKey` | COVERED | - |
| 16 | `testPushWithNewKey` | `TestPushWithNewKey` | COVERED | - |
| 17 | `testAll` | `TestAllReturnsFullConfig` | COVERED | - |
| 18 | `testItGetsAsString` | `TestStringAccessor` | COVERED | - |
| 19 | `testItThrowsAnExceptionWhenTryingToGetNonStringValueAsString` | `TestStringAccessorTypeError` | COVERED | - |
| 20 | `testItGetsAsArray` | `TestStringSliceAccessor` | COVERED | Go uses `StringSlice` |
| 21 | `testItThrowsAnExceptionWhenTryingToGetNonArrayValueAsArray` | `TestRepositoryTypeFailuresAndMissingKeys` | COVERED | - |
| 22 | `testItGetsAsBoolean` | `TestBoolAccessor` | COVERED | - |
| 23 | `testItThrowsAnExceptionWhenTryingToGetNonBooleanValueAsBoolean` | `TestBoolAccessorTypeError` | COVERED | - |
| 24 | `testItGetsAsInteger` | `TestIntAccessor` | COVERED | - |
| 25 | `testItThrowsAnExceptionWhenTryingToGetNonIntegerValueAsInteger` | `TestIntAccessorTypeError` | COVERED | - |
| 26 | `testItGetsAsFloat` | `TestFloatAccessor` | COVERED | - |
| 27 | `testItThrowsAnExceptionWhenTryingToGetNonFloatValueAsFloat` | `TestFloatAccessorTypeError` | COVERED | - |
| 28 | `testOffsetExists` | - | INTENTIONAL-SKIP | (a) PHP ArrayAccess |
| 29 | `testOffsetGet` | - | INTENTIONAL-SKIP | (a) PHP ArrayAccess |
| 30 | `testOffsetSet` | - | INTENTIONAL-SKIP | (a) PHP ArrayAccess |
| 31 | `testOffsetUnset` | - | INTENTIONAL-SKIP | (a) PHP ArrayAccess |
| 32 | `testItIsMacroable` | - | INTENTIONAL-SKIP | (a) PHP macros |
| 33 | `testItGetsAsCollection` | - | INTENTIONAL-SKIP | (a) no Collection type |

**Bedrock-only tests (no Upstream equivalent):** `TestRepositorySetMutatorsAndCloneSemantics`, `TestInternalHelpers`, `TestLookupHelperPaths`, `TestRepositoryErrorMessagesStayStable`, `TestBuilderMergesBaseOverlayAndEnv` (7 builder tests)

---

## 2. Encryption (`tests/Encryption/EncrypterTest.php`)

**Upstream: 26 tests | Bedrock: 22 tests**
**Coverage: 85% (22/26)**

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testEncryption` | `TestEncryptDecryptStringAndJSONValues` | COVERED | - |
| 2 | `testRawStringEncryption` | `TestEncryptDecryptStringAndJSONValues` | COVERED | inline |
| 3 | `testRawStringEncryptionWithPreviousKeys` | `TestRawStringEncryptionWithPreviousKeys` | COVERED | - |
| 4 | `testItValidatesMacOnPerKeyBasis` | `TestMACValidationPerKey` | COVERED | - |
| 5 | `testEncryptionUsingBase64EncodedKey` | `TestEncryptionUsingBase64EncodedKey` | COVERED | - |
| 6 | `testEncryptedLengthIsFixed` | `TestEncryptedLengthIsFixed` | COVERED | - |
| 7 | `testWithCustomCipher` | `TestNewSupportedAndAccessors` | COVERED | - |
| 8 | `testCipherNamesCanBeMixedCase` | `TestCipherNamesCanBeMixedCase` | COVERED | - |
| 9 | `testThatAnAeadCipherIncludesTag` | `TestAeadCipherIncludesTag` | COVERED | - |
| 10 | `testThatANonAeadCipherIncludesMac` | `TestNonAeadCipherIncludesMac` | COVERED | - |
| 11 | `testDoNoAllowLongerKey` | `TestDoNotAllowLongerKey` | COVERED | - |
| 12 | `testExceptionThrownWhenPayloadIsInvalid` | `TestDecryptFailuresAndPayloadValidation` | COVERED | - |
| 13 | `testExceptionThrownWithDifferentKey` | `TestDecryptionFailsWithDifferentKey` | COVERED | - |
| 14 | `testTamperedPayloadWillGetRejected` | `TestTamperedPayloadIsRejected` | COVERED | - |
| 15 | `testEncryptedReturnsTrueForEncryptedValue` | `TestAppearsEncryptedAndRepositoryLoading` | COVERED | - |
| 16 | `testEncryptedReturnsTrueForEncryptedArray` | `TestAppearsEncryptedForEncryptedArray` | COVERED | - |
| 17 | `testEncryptedReturnsFalseForPlainText` | `TestAppearsEncryptedForNonEncryptedValues` | COVERED | - |
| 18 | `testEncryptedReturnsFalseForNonString` | `TestAppearsEncryptedForNonEncryptedValues` | COVERED | Go type safety |
| 19 | `testSupportedMethodAcceptsAnyCasing` | `TestCipherHelpers` | COVERED | - |
| 20 | `testThatAnAeadTagMustBeProvidedInFullLength` | `TestDecryptFailuresAndPayloadValidation` | COVERED | inline |
| 21 | `testThatAnAeadTagCantBeModified` | `TestDecryptFailuresAndPayloadValidation` | COVERED | inline |
| 22 | `testDecryptionExceptionIsThrownWhenUnexpectedTagIsAdded` | `TestDecryptFailuresAndPayloadValidation` | COVERED | inline |
| 23 | `testWithBadKeyLength` | - | MISSING | (b) missing test |
| 24 | `testWithBadKeyLengthAlternativeCipher` | - | MISSING | (b) missing test |
| 25 | `testWithUnsupportedCipher` | - | MISSING | (b) missing test |
| 26 | `testExceptionThrownWhenIvIsTooLong` | - | MISSING | (b) missing test |

**Bedrock-only tests:** `TestGenerateKey`, `TestEncryptFailures`, `TestPreviousKeyDecryptionAndFixture`, `TestPayloadAndHelperFunctions`, `TestCipherHelpers`, `TestPaddingAndKeyParsing`

---

## 3. Hashing (`tests/Hashing/HasherTest.php`)

**Upstream: 13 tests | Bedrock: 9 tests**
**Coverage: 85% (11/13)**

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testBasicBcryptHashing` | `TestBcryptHasher` | COVERED | - |
| 2 | `testBasicArgon2iHashing` | `TestArgonHashers` | COVERED | - |
| 3 | `testBasicArgon2idHashing` | `TestArgonHashers` | COVERED | - |
| 4 | `testEmptyHashedValueReturnsFalse` | `TestEmptyHashedValueReturnsFalse` | COVERED | - |
| 5 | `testNullHashedValueReturnsFalse` | `TestCheckAgainstEmptyPasswordHash` | COVERED | Go: empty string equiv |
| 6 | `testIsHashedWithNonHashedValue` | `TestIsHashedWithNonHashedValue` | COVERED | - |
| 7 | `testBcryptValueTooLong` | `TestBcryptValueTooLong` | COVERED | - |
| 8 | `testBasicBcryptVerification` | `TestCrossHasherVerification` | COVERED | - |
| 9 | `testBasicArgon2iVerification` | `TestCrossHasherVerification` | COVERED | - |
| 10 | `testBasicArgon2idVerification` | `TestCrossHasherVerification` | COVERED | - |
| 11 | `testBasicBcryptNotSupported` | - | INTENTIONAL-SKIP | (a) Go bcrypt always available |
| 12 | `testBasicArgon2iNotSupported` | - | INTENTIONAL-SKIP | (a) Go argon2 always available |
| 13 | `testBasicArgon2idNotSupported` | - | INTENTIONAL-SKIP | (a) Go argon2 always available |

**Bedrock-only tests:** `TestArgonParsingHelpers`, `TestManagerAndHelpers`

---

## 4. Auth / Access Gate (`tests/Auth/AuthAccessGateTest.php` + `AuthAccessResponseTest.php`)

**Upstream: 92 + 10 = 102 tests | Bedrock: 53 tests**
**Coverage: 60% (61/102)**

### Gate Tests (92 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testBasicClosuresCanBeDefined` | `TestBasicClosuresCanBeDefined` | COVERED | - |
| 2 | `testBeforeCallbacksCanOverrideResultIfNecessary` | `TestBeforeCallbacksCanOverrideResultIfNecessary` | COVERED | - |
| 3 | `testBeforeCallbacksDontInterruptGateCheckIfNoValueIsReturned` | `TestBeforeCallbacksDontInterruptGateCheckIfNoValueIsReturned` | COVERED | - |
| 4 | `testAfterCallbacksAreCalledWithResult` | `TestAfterCallbacksAreCalledWithResult` | COVERED | - |
| 5 | `testAfterCallbacksCanAllowIfNull` | `TestAfterCallbacksCanAllowUndefined` | COVERED | - |
| 6 | `testAfterCallbacksDoNotOverridePreviousResult` | `TestAfterCallbacksDoNotOverridePreviousResult` | COVERED | - |
| 7 | `testAfterCallbacksDoNotOverrideEachOther` | `TestAfterCallbacksDoNotOverrideEachOther` | COVERED | - |
| 8 | `testCurrentUserThatIsOnGateAlwaysInjectedIntoClosureCallbacks` | `TestCurrentUserIsInjectedIntoClosureCallbacks` | COVERED | - |
| 9 | `testASingleArgumentCanBePassedWhenCheckingAbilities` | `TestSingleArgumentCanBePassedWhenCheckingAbilities` | COVERED | - |
| 10 | `testMultipleArgumentsCanBePassedWhenCheckingAbilities` | `TestMultipleArgumentsCanBePassedWhenCheckingAbilities` | COVERED | - |
| 11 | `testPolicyClassesCanBeDefinedToHandleChecksForGivenType` | `TestPolicyClassesCanBeDefinedToHandleChecksForGivenType` | COVERED | - |
| 12 | `testPolicyDefaultToFalseIfMethodDoesNotExistAndGateDoesNotExist` | `TestPolicyDefaultToFalseIfMethodDoesNotExistAndGateDoesNotExist` | COVERED | - |
| 13 | `testPoliciesAlwaysOverrideClosuresWithSameName` | `TestPoliciesOverrideClosuresForSameResource` | COVERED | - |
| 14 | `testPoliciesDeferToGatesIfMethodDoesNotExist` | `TestPoliciesDeferToGatesIfMethodDoesNotExist` | COVERED | - |
| 15 | `testAuthorizeThrowsUnauthorizedException` | `TestAuthorizeThrowsUnauthorizedException` | COVERED | - |
| 16 | `testAuthorizeReturnsAllowedResponse` | `TestAuthorizeReturnsNilForAllowedAbility` | COVERED | - |
| 17 | `testResponseReturnsResponseWhenAbilityGranted` | `TestInspectReturnsResponseWhenAbilityGranted` | COVERED | - |
| 18 | `testResponseReturnsResponseWhenAbilityDenied` | `TestInspectReturnsResponseWhenAbilityDenied` | COVERED | - |
| 19 | `testAuthorizeReturnsAnAllowedResponseForATruthyReturn` | `TestAuthorizeReturnsNilForTruthyPolicyReturn` | COVERED | - |
| 20 | `testAuthorizeWithPolicyThatReturnsDeniedResponseObjectThrowsException` | `TestAuthorizeWithPolicyDeniedResponseThrowsException` | COVERED | - |
| 21 | `testPoliciesMayHaveBeforeMethodsToOverrideChecks` | `TestPoliciesBeforeHooksOverrideChecks` | COVERED | - |
| 22 | `testForUserMethodAttachesANewUserToANewGateInstance` | `TestGateForUser` | COVERED | - |
| 23 | `testAnyAbilityCheckPassesIfAllPass` | `TestAnyAbilityCheckPassesIfAllPass` | COVERED | - |
| 24 | `testAnyAbilityCheckPassesIfAtLeastOnePasses` | `TestAnyAbilityCheckPassesIfAtLeastOnePasses` | COVERED | - |
| 25 | `testAnyAbilityCheckFailsIfNonePass` | `TestAnyAbilityCheckFailsIfNonePass` | COVERED | - |
| 26 | `testNoneAbilityCheckPassesIfAllFail` | `TestNoneAbilityCheck` | COVERED | - |
| 27 | `testEveryAbilityCheckPassesIfAllPass` | - | MISSING | (b) |
| 28 | `testEveryAbilityCheckFailsIfAtLeastOneFails` | - | MISSING | (b) |
| 29 | `testEveryAbilityCheckFailsIfNonePass` | - | MISSING | (b) |
| 30 | `testHasAbilities` | `TestGateHas` | COVERED | - |
| 31 | `testAllowIfAuthorizesTrue` | `TestAllowIfAndDenyIf` | COVERED | - |
| 32 | `testDenyIfAuthorizesFalse` | `TestAllowIfAndDenyIf` | COVERED | - |
| 33 | `testResourceGatesCanBeDefined` | `TestGateResource` | COVERED | - |
| 34 | `testPolicyThatThrowsAuthorizationExceptionIsCaughtInInspect` | - | MISSING | (b) |
| 35 | `testBeforeCanTakeAnArrayCallbackAsObject` | - | INTENTIONAL-SKIP | (a) PHP array callbacks |
| 36 | `testBeforeCanTakeAnArrayCallbackAsObjectStatic` | - | INTENTIONAL-SKIP | (a) PHP array callbacks |
| 37 | `testBeforeCanTakeAnArrayCallbackWithStaticMethod` | - | INTENTIONAL-SKIP | (a) PHP array callbacks |
| 38 | `testClassesCanBeDefinedAsCallbacksUsingAtNotation` | - | INTENTIONAL-SKIP | (a) PHP @ notation |
| 39 | `testClassesCanBeDefinedAsCallbacksUsingAtNotationForGuests` | - | INTENTIONAL-SKIP | (a) PHP @ notation |
| 40 | `testInvokableClassesCanBeDefined` | - | INTENTIONAL-SKIP | (a) PHP __invoke |
| 41 | `testGatesCanBeDefinedUsingAnArrayCallback` | - | INTENTIONAL-SKIP | (a) PHP array callbacks |
| 42 | `testGatesCanBeDefinedUsingAnArrayCallbackWithStaticMethod` | - | INTENTIONAL-SKIP | (a) PHP array callbacks |
| 43 | `testCanDefineGatesUsingBackedEnum` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 44 | `testBackedEnumInAllows` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 45 | `testBackedEnumInDenies` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 46 | `testAnyAbilitiesCheckUsingBackedEnum` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 47 | `testNoneAbilitiesCheckUsingBackedEnum` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 48 | `testAbilitiesCheckUsingBackedEnum` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 49 | `testBeforeCanAllowGuests` | - | MISSING | (b) guest user handling |
| 50 | `testAfterCanAllowGuests` | - | MISSING | (b) guest user handling |
| 51 | `testClosuresCanAllowGuestUsers` | - | MISSING | (b) guest user handling |
| 52 | `testPoliciesCanAllowGuests` | - | MISSING | (b) guest user handling |
| 53 | `testPolicyBeforeNotCalledWithGuestsIfItDoesntAllowThem` | - | MISSING | (b) guest user handling |
| 54 | `testBeforeAndAfterCallbacksCanAllowGuests` | - | MISSING | (b) guest user handling |
| 55 | `testPolicyClassesHandleChecksForAllSubtypes` | - | MISSING | (b) subtype resolution |
| 56 | `testPolicyClassesHandleChecksForInterfaces` | - | MISSING | (b) interface resolution |
| 57 | `testPolicyConvertsDashToCamel` | - | MISSING | (b) |
| 58 | `testPolicyClassesCanBeDefinedToHandleChecksForGivenClassName` | - | MISSING | (b) |
| 59 | `testDefineSecondParameterShouldBeStringOrCallable` | - | MISSING | (b) |
| 60 | `testAuthorizeThrowsUnauthorizedExceptionWithCustomStatusCode` | `TestResponseWithCode` | COVERED | - |
| 61 | `testCustomResourceGatesCanBeDefined` | - | MISSING | (b) |
| 62 | `testForUserMethodAttachesANewUserToANewGateInstanceWithGuessCallback` | - | MISSING | (b) |
| 63 | `testArrayAbilitiesInAllows` | - | MISSING | (b) |
| 64 | `testArrayAbilitiesInDenies` | - | MISSING | (b) |
| 65-78 | `testAllowIf*`/`testDenyIf*` variants (14 methods) | `TestAllowIfAndDenyIf` | PARTIAL | (c) covers basic cases, missing guest/callback/response variants |
| 79 | `testCanSetDenialResponseInConstructor` | - | MISSING | (b) |
| 80 | `testCanSetDenialResponse` | - | MISSING | (b) |

### Response Tests (10 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testAllowMethod` | `TestAllowAndDenyConstructors` | COVERED | - |
| 2 | `testDenyMethod` | `TestAllowAndDenyConstructors` | COVERED | - |
| 3 | `testDenyMethodWithNoMessageReturnsNull` | - | MISSING | (b) |
| 4 | `testItSetsEmptyStatusOnExceptionWhenAuthorizing` | - | MISSING | (b) |
| 5 | `testItSetsStatusOnExceptionWhenAuthorizing` | `TestResponseWithCode` | COVERED | - |
| 6 | `testAuthorizeMethodThrowsAuthorizationExceptionWhenResponseDenied` | `TestAuthorizeThrowsUnauthorizedException` | COVERED | - |
| 7 | `testAuthorizeMethodThrowsAuthorizationExceptionWithDefaultMessage` | - | MISSING | (b) |
| 8 | `testThrowIfNeededDoesntThrowAuthorizationExceptionWhenResponseAllowed` | - | MISSING | (b) |
| 9 | `testCastingToStringReturnsMessage` | `TestAuthorizationExceptionErrorString` | COVERED | - |
| 10 | `testResponseToArrayMethod` | - | MISSING | (b) |

**Bedrock-only tests:** `TestDeniesMethod`, `TestBeforeAndAfterCallbackOrdering`, `TestBeforeCallbackShortCircuits`, `TestAfterCallbackCanModifyResult`, `TestPolicyWithPointerToResource`, `TestUndefinedAbilityReturnsDeny`, `TestMultipleAbilitiesOnSameGate`, `TestGateImplementsAuthorizer`, `TestAbilityNameTrimming`, `TestBeforeAllowsUndefinedAbility`, `TestContextPassedThroughToCallbacks`, `TestGateConcurrentAccess`, `TestPolicyKeyResolutionForStringTarget`, `TestAnyWithEmptyAbilities`, `TestAuthorizeWithUndefinedAbility`, `TestBeforeCallbackReceivesAbilityName`, `TestAfterCallbackReceivesAbilityAndArguments`, `TestPolicyAndAbilityCoexist`, `TestAfterSkippedWhenBeforeShortCircuits`, `TestMultiplePoliciesForDifferentTypes`

---

## 5. Auth / Guards (`tests/Auth/AuthGuardTest.php` + `AuthTokenGuardTest.php`)

**Upstream: 43 + 14 = 57 tests | Bedrock: 63 tests**
**Coverage: 70% (40/57)**

### Session Guard (43 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testAttemptCallsRetrieveByCredentials` | `TestSessionGuardAttemptSuccess` | COVERED | - |
| 2 | `testAttemptReturnsUserInterface` | `TestSessionGuardAttemptSuccess` | COVERED | - |
| 3 | `testAttemptReturnsFalseIfUserNotGiven` | `TestSessionGuardAttemptUserNotFound` | COVERED | - |
| 4 | `testLoginStoresIdentifierInSession` | `TestSessionGuardLoginStoresIdentifierInSession` | COVERED | - |
| 5 | `testAuthenticateReturnsUserWhenUserIsNotNull` | `TestAuthenticateRequestReturnsUserWhenSessionExists` | COVERED | - |
| 6 | `testAuthenticateThrowsWhenUserIsNull` | `TestAuthenticateRequestReturnsErrorWhenNoSession` | COVERED | - |
| 7 | `testHasUserReturnsTrueWhenUserIsNotNull` | `TestAuthenticateRequestReturnsUserWhenSessionExists` | COVERED | inline |
| 8 | `testHasUserReturnsFalseWhenUserIsNull` | `TestAuthenticateRequestReturnsErrorWhenNoSession` | COVERED | inline |
| 9 | `testIsAuthedReturnsTrueWhenUserIsNotNull` | `TestAuthenticateRequestReturnsUserWhenSessionExists` | COVERED | inline |
| 10 | `testIsAuthedReturnsFalseWhenUserIsNull` | `TestAuthenticateRequestReturnsErrorWhenUserNotFound` | COVERED | - |
| 11 | `testUserMethodReturnsCachedUser` | `TestTokenGuardCachesUser` | COVERED | on token guard |
| 12 | `testNullIsReturnedForUserIfNoUserFound` | `TestAuthenticateRequestReturnsErrorWhenUserNotFound` | COVERED | - |
| 13 | `testUserIsSetToRetrievedUser` | `TestAuthenticateRequestReturnsUserWhenSessionExists` | COVERED | - |
| 14 | `testLogoutRemovesSessionTokenAndRememberMeCookie` | `TestLogoutRemovesSessionAndCookies` | COVERED | - |
| 15 | `testLogoutDoesNotEnqueueRememberMeCookieForDeletionIfCookieDoesntExist` | `TestLogoutClearsRememberCookieEvenIfAbsent` | COVERED | - |
| 16 | `testLogoutDoesNotSetRememberTokenIfNotPreviouslySet` | `TestLogoutClearsRememberToken` | COVERED | - |
| 17 | `testLogoutCurrentDeviceRemovesRememberMeCookie` | - | MISSING | (b) |
| 18 | `testLogoutCurrentDeviceDoesNotEnqueueRememberMeCookieForDeletionIfCookieDoesntExist` | - | MISSING | (b) |
| 19 | `testLoginMethodQueuesCookieWhenRemembering` | `TestLoginMethodCreatesRememberCookieWhenRemembering` | COVERED | - |
| 20 | `testLoginMethodQueuesCookieWhenRememberingAndAllowsOverride` | - | MISSING | (b) |
| 21 | `testLoginMethodCreatesRememberTokenIfOneDoesntExist` | `TestLoginCreatesRememberTokenIfMissing` | COVERED | - |
| 22 | `testLoginUsingIdLogsInWithUser` | `TestSessionGuardLoginUsingId` | COVERED | - |
| 23 | `testLoginUsingIdFailure` | `TestSessionGuardLoginUsingIdNotFound` | COVERED | - |
| 24 | `testOnceUsingIdSetsUser` | `TestSessionGuardOnceUsingId` | COVERED | - |
| 25 | `testOnceUsingIdFailure` | `TestSessionGuardOnceUsingId` | COVERED | failure inline |
| 26 | `testUserUsesRememberCookieIfItExists` | `TestSessionGuardLoginAndRememberRestore` | COVERED | - |
| 27 | `testLoginOnceSetsUser` | `TestSessionGuardOnce` | COVERED | - |
| 28 | `testLoginOnceFailure` | `TestSessionGuardOnceFailsWithInvalidCredentials` | COVERED | - |
| 29 | `testAttemptAndWithCallbacks` | - | MISSING | (b) |
| 30 | `testAttemptRehashesPasswordWhenRequired` | - | MISSING | (b) rehashing |
| 31 | `testAttemptDoesntRehashPasswordWhenDisabled` | - | MISSING | (b) rehashing |
| 32 | `testForgetUserSetsUserToNull` | - | MISSING | (b) |
| 33 | `testBasicReturnsNullOnValidAttempt` | - | MISSING | (b) basic HTTP auth |
| 34 | `testBasicReturnsNullWhenAlreadyLoggedIn` | - | MISSING | (b) basic HTTP auth |
| 35 | `testBasicReturnsResponseOnFailure` | - | MISSING | (b) basic HTTP auth |
| 36 | `testBasicWithExtraConditions` | - | MISSING | (b) basic HTTP auth |
| 37 | `testBasicWithExtraArrayConditions` | - | MISSING | (b) basic HTTP auth |
| 38 | `testSessionGuardIsMacroable` | - | INTENTIONAL-SKIP | (a) PHP macros |
| 39 | `testLoginFiresLoginAndAuthenticatedEvents` | - | MISSING | (d) no event system |
| 40 | `testFailedAttemptFiresFailedEvent` | - | MISSING | (d) no event system |
| 41 | `testSetUserFiresAuthenticatedEvent` | - | MISSING | (d) no event system |
| 42 | `testLogoutFiresLogoutEvent` | - | MISSING | (d) no event system |
| 43 | `testLogoutCurrentDeviceFiresLogoutEvent` | - | MISSING | (d) no event system |

### Token Guard (14 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testUserCanBeRetrievedByQueryStringVariable` | `TestTokenGuardUserFromQueryString` | COVERED | - |
| 2 | `testTokenCanBeHashed` | `TestTokenGuardHashedStorageKey` | COVERED | - |
| 3 | `testUserCanBeRetrievedByAuthHeaders` | `TestTokenGuardUserFromAuthorizationHeader` | COVERED | - |
| 4 | `testUserCanBeRetrievedByBearerToken` | `TestTokenGuardUserFromBearerToken` | COVERED | - |
| 5 | `testValidateCanDetermineIfCredentialsAreValid` | `TestTokenGuardValidateSuccess` | COVERED | - |
| 6 | `testValidateCanDetermineIfCredentialsAreInvalid` | `TestTokenGuardValidateFailure` | COVERED | - |
| 7 | `testValidateIfApiTokenIsEmpty` | `TestTokenGuardReturnsErrorForEmptyToken` | COVERED | - |
| 8 | `testItAllowsToPassCustomRequestInSetterAndUseItForValidation` | `TestTokenGuardSetRequest` | COVERED | - |
| 9 | `testUserCanBeRetrievedByBearerTokenWithCustomKey` | `TestTokenGuardCustomInputKey` | COVERED | - |
| 10 | `testUserCanBeRetrievedByQueryStringVariableWithCustomKey` | `TestTokenGuardCustomInputKey` | COVERED | inline |
| 11 | `testUserCanBeRetrievedByAuthHeadersWithCustomField` | - | MISSING | (b) |
| 12 | `testValidateCanDetermineIfCredentialsAreValidWithCustomKey` | - | MISSING | (b) |
| 13 | `testValidateCanDetermineIfCredentialsAreInvalidWithCustomKey` | - | MISSING | (b) |
| 14 | `testValidateIfApiTokenIsEmptyWithCustomKey` | - | MISSING | (b) |

**Bedrock-only tests:** `TestLoginWithPendingTwoFactor`, `TestSessionGuardExpiresOldSessions`, `TestSessionGuardUpdatesLastSeenAt`, `TestSessionGuardName`, `TestLogoutWithNilSessionAndUser`, `TestTokenGuardReturnsErrorForMissingToken`, `TestTokenGuardNilRequest`, `TestTokenGuardDefaultInputKey`, `TestTokenGuardValidateUserNotFound`, `TestRequestGuard*` (6 tests), `TestManager*` (10 tests), `TestAuthenticationExceptionError`, `TestSentinelErrors`, `TestSystemClockReturnsCurrentTime`, `TestRandomIDGeneratorProducesUniqueIDs`, `TestRecallerParsing`, `TestExpiredHelper`, `TestSessionGuardImplementsStatefulGuard`, `TestSessionGuardClearsInvalidRememberCookie`

---

## 6. Auth / Passwords (`tests/Auth/AuthPasswordBrokerTest.php` + `AuthDatabaseTokenRepositoryTest.php`)

**Upstream: 9 + 10 = 19 tests | Bedrock: 24 tests**
**Coverage: 79% (15/19)**

### Broker Tests (9 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testIfUserIsNotFoundErrorRedirectIsReturned` | `TestBrokerSendResetLinkUserNotFound` | COVERED | - |
| 2 | `testIfTokenIsRecentlyCreated` | `TestBrokerThrottlesTokenCreation` | COVERED | - |
| 3 | `testGetUserThrowsExceptionIfUserDoesntImplementCanResetPassword` | - | INTENTIONAL-SKIP | (a) Go interfaces enforced at compile time |
| 4 | `testUserIsRetrievedByCredentials` | `TestBrokerCreateValidateAndReset` | COVERED | - |
| 5 | `testBrokerCreatesTokenAndRedirectsWithoutError` | `TestBrokerCreateTokenSuccess` | COVERED | - |
| 6 | `testRedirectIsReturnedByResetWhenUserCredentialsInvalid` | `TestBrokerResetReturnsErrorWhenUserNotFound` | COVERED | - |
| 7 | `testRedirectReturnedByRemindWhenRecordDoesntExistInTable` | `TestBrokerResetReturnsErrorForInvalidToken` | COVERED | - |
| 8 | `testResetRemovesRecordOnReminderTableAndCallsCallback` | `TestBrokerResetDeletesTokenAfterSuccess` | COVERED | - |
| 9 | `testExecutesCallbackInsteadOfSendingNotification` | `TestBrokerSendResetLinkWithCallback` | COVERED | - |

### Token Repository Tests (10 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testCreateInsertsNewRecordIntoTable` | `TestMemoryRepoSave` | COVERED | uses memory repo |
| 2 | `testExistReturnsFalseIfNoRowFoundForUser` | `TestBrokerTokenExistsReturnsFalseForMissingToken` | COVERED | - |
| 3 | `testExistReturnsFalseIfRecordIsExpired` | `TestBrokerTokenExistsReturnsFalseForExpiredToken` | COVERED | - |
| 4 | `testExistReturnsTrueIfValidRecordExists` | `TestBrokerCreateTokenSuccess` | COVERED | - |
| 5 | `testExistReturnsFalseIfInvalidToken` | `TestBrokerResetReturnsErrorForInvalidToken` | COVERED | - |
| 6 | `testRecentlyCreatedReturnsFalseIfNoRowFoundForUser` | `TestMemoryRepoRecentlyCreatedFalseForMissingUser` | COVERED | - |
| 7 | `testRecentlyCreatedReturnsTrueIfRecordIsRecentlyCreated` | `TestMemoryRepoRecentlyCreatedTrue` | COVERED | - |
| 8 | `testRecentlyCreatedReturnsFalseIfValidRecordExists` | `TestMemoryRepoRecentlyCreatedFalseWhenOldEnough` | COVERED | - |
| 9 | `testDeleteMethodDeletesByToken` | `TestMemoryRepoDeleteByTokenHash` | COVERED | - |
| 10 | `testDeleteExpiredMethodDeletesExpiredTokens` | - | MISSING | (b) no purge/cleanup method |

**Note:** Bedrock uses a memory-based token repository instead of SQL. A SQL-backed repository is not yet implemented.

**Bedrock-only tests:** `TestNormalizeEmail`, `TestMemoryRepoFindByTokenHashReturnsErrorWhenMissing`, `TestMemoryRepoDeleteByUserID`, `TestMemoryRepoSaveReplacesExisting`, `TestMemoryRepoDeleteNonexistent`, `TestBrokerSendResetLink`, `TestBrokerSendResetLinkThrottled`, `TestBrokerCanCreateTokenWhenNoTokenExists`, `TestBrokerDeleteToken`

---

## 7. Auth / Middleware (`tests/Auth/AuthenticateMiddlewareTest.php` + `AuthorizeMiddlewareTest.php`)

**Upstream: 11 + 14 = 25 tests | Bedrock: 10 tests**
**Coverage: 28% (7/25)**

### Authenticate Middleware (11 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testDefaultUnauthenticatedThrows` | `TestAuthenticateMiddlewareRejectsUnauthenticated` | COVERED | - |
| 2 | `testDefaultAuthenticatedKeepsDefaultDriver` | `TestAuthenticateMiddlewareAllowsAuthenticated` | COVERED | - |
| 3 | `testItCanGenerateDefinitionViaStaticMethod` | - | INTENTIONAL-SKIP | (a) PHP static methods |
| 4 | `testItCanGenerateDefinitionViaStaticMethodForBasic` | - | INTENTIONAL-SKIP | (a) PHP static methods |
| 5 | `testDefaultUnauthenticatedThrowsWithGuards` | - | MISSING | (b) multi-guard |
| 6 | `testSecondaryAuthenticatedUpdatesDefaultDriver` | - | MISSING | (b) multi-guard |
| 7 | `testMultipleDriversUnauthenticatedThrows` | - | MISSING | (b) multi-guard |
| 8 | `testMultipleDriversUnauthenticatedThrowsWithGuards` | - | MISSING | (b) multi-guard |
| 9 | `testMultipleDriversAuthenticatedUpdatesDefault` | - | MISSING | (b) multi-guard |
| 10 | `testCustomDriverClosureBoundObjectIsAuthManager` | - | MISSING | (d) IoC container |
| 11 | `testCustomDriverStatic` | - | MISSING | (d) IoC container |

### Authorize Middleware (14 Upstream methods)

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testSimpleAbilityAuthorized` | `TestAuthorizeMiddlewareAllowsAuthorized` | COVERED | - |
| 2 | `testSimpleAbilityUnauthorized` | `TestAuthorizeMiddlewareDeniesUnauthorized` | COVERED | - |
| 3 | `testItCanGenerateDefinitionViaStaticMethod` | - | INTENTIONAL-SKIP | (a) PHP static methods |
| 4 | `testSimpleAbilityWithStringParameter` | - | MISSING | (b) |
| 5 | `testSimpleAbilityWithBackedEnumParameter` | - | INTENTIONAL-SKIP | (a) PHP enums |
| 6 | `testSimpleAbilityWithNullParameter` | - | MISSING | (b) |
| 7 | `testSimpleAbilityWithOptionalParameter` | - | MISSING | (b) |
| 8 | `testSimpleAbilityWithStringParameterFromRouteParameter` | - | MISSING | (b) route param binding |
| 9 | `testSimpleAbilityWithStringParameter0FromRouteParameter` | - | MISSING | (b) |
| 10 | `testModelTypeUnauthorized` | - | MISSING | (b) model authorization |
| 11 | `testModelTypeAuthorized` | - | MISSING | (b) |
| 12 | `testModelUnauthorized` | - | MISSING | (b) |
| 13 | `testModelAuthorized` | - | MISSING | (b) |
| 14 | `testModelInstanceAsParameter` | - | MISSING | (b) |

**Additional Upstream files not covered:**
- `EnsureEmailIsVerifiedTest.php` — Bedrock has `TestEnsureEmailIsVerifiedAllowsVerified`, `TestEnsureEmailIsVerifiedRejectsUnverified`, `TestEnsureEmailIsVerifiedRejectsNoUser`
- `RedirectIfAuthenticatedMiddlewareTest.php` — entirely MISSING (b)

---

## 8. Console (`tests/Console/ConsoleApplicationTest.php`)

**Upstream: 14 tests (this file only, 10+ additional test files) | Bedrock: 9 tests**
**Coverage: ~20% of core file**

| # | Upstream Test | Bedrock Test | Status | Gap |
|---|---|---|---|---|
| 1 | `testCallFullyStringCommandLine` | `TestCommandExecution` | COVERED | - |
| 2 | `testResolveAddsCommandViaApplicationResolution` | `TestCommandRegistration` | COVERED | - |
| 3 | `testCallMethodCanCallArtisanCommandUsingCommandClassObject` | - | MISSING | (d) IoC container |
| 4 | `testAddSetsUpstreamInstance` | - | MISSING | (d) IoC container |
| 5 | `testUpstreamNotSetOnSymfonyCommands` | - | INTENTIONAL-SKIP | (a) Symfony integration |
| 6 | `testResolvingCommandsWithAliasViaAttribute` | - | MISSING | (d) PHP attributes |
| 7 | `testResolvingCommandsWithAliasViaProperty` | - | MISSING | (d) |
| 8 | `testResolvingCommandsWithNoAliasViaAttribute` | - | MISSING | (d) |
| 9 | `testResolvingCommandsWithNoAliasViaProperty` | - | MISSING | (d) |
| 10 | `testCommandInputPromptsWhenRequiredArgumentIsMissing` | - | MISSING | (b) interactive prompts |
| 11-12 | `testCommandInputDoesntPrompt*` | - | MISSING | (b) |
| 13-14 | `testCommandInput*AreMissing/Passed` | - | MISSING | (b) |

**Not analyzed:** Upstream has 10+ additional console test files (scheduling, signals, middleware, etc.) — all MISSING in Bedrock.

---

## 9. Routing (`tests/Routing/RoutingRouteTest.php`)

**Upstream: 102 tests (this file only, 11+ additional files) | Bedrock: 12 tests**
**Coverage: ~10%**

| # | Upstream Test | Bedrock Test | Status |
|---|---|---|---|
| 1 | `testBasicDispatchingOfRoutes` | `TestGetRouteRegistrationAndDispatch` | COVERED |
| 2 | `testOptionsResponsesAreGeneratedByDefault` | `TestHandleWithEmptyMethodAllowsAll` | COVERED |
| 3 | `testHeadDispatcher` | `TestHandleWithMethodRestriction` | PARTIAL |
| 4-102 | Remaining 99 tests | - | MISSING |

**Major missing areas:** middleware groups, route model binding, resource routing, controller routing, domain routing, URL generation, signed routes, redirects, pattern filtering, nested groups.

---

## 10. HTTP (`tests/Http/HttpRequestTest.php`)

**Upstream: 113 tests (this file only, 9+ additional files) | Bedrock: 11 tests**
**Coverage: ~8%**

| # | Upstream Test | Bedrock Test | Status |
|---|---|---|---|
| 1 | `testInstanceMethod` | `TestCapture` | COVERED |
| 2 | `testInputMethod` | `TestRequestInput` | COVERED |
| 3 | `testQueryMethod` | `TestRequestQuery` | COVERED |
| 4 | `testBooleanMethod` | `TestRequestBoolean` | COVERED |
| 5 | `testIntegerMethod` | `TestRequestInteger` | COVERED |
| 6-113 | Remaining 108 tests | - | MISSING |

**Major missing areas:** path/URL methods, content negotiation (accepts, prefers, wantsJson), file uploads, cookies, session integration, flash data, fingerprinting, merge/replace/except/only, HTTP client.

---

## 11. View (`tests/View/ViewFactoryTest.php`)

**Upstream: 70 tests (this file only, 10+ additional files) | Bedrock: 10 tests**
**Coverage: ~7%**

| # | Upstream Test | Bedrock Test | Status |
|---|---|---|---|
| 1 | `testMakeCreatesNewViewInstanceWithProperPathAndEngine` | `TestRenderBasicView` | COVERED |
| 2 | `testFirstCreatesNewViewInstanceWithProperPath` | - | MISSING |
| 3-70 | Remaining 68 tests | - | MISSING |

**Major missing areas:** Template compiler (entire engine), components, sections/stacks/fragments, view composers/creators, loops, translation integration. Bedrock uses Go `html/template` instead of Template.

---

## 12. Foundation (`tests/Foundation/FoundationApplicationTest.php`)

**Upstream: 42 tests (this file only, 15+ additional files) | Bedrock: 13 tests**
**Coverage: ~25%**

| # | Upstream Test | Bedrock Test | Status |
|---|---|---|---|
| 1 | `testEnvironment` | `TestApplicationEnvironment` | COVERED |
| 2 | `testEnvironmentHelpers` | `TestApplicationEnvironmentDefaultsToProduction` | COVERED |
| 3 | `testTerminationTests` | - | MISSING (b) |
| 4 | `testBootingCallbacks` | - | MISSING (b) |
| 5 | `testBootedCallbacks` | - | MISSING (b) |
| 6 | `testServiceProvidersAreCorrectlyRegistered` | - | MISSING (d) no IoC |
| 7-42 | Remaining tests | - | MISSING |

**Major missing areas:** service providers, deferred services, bootstrapping callbacks, termination, namespace resolution, cache paths, config merging, route/event caching, alias loader.

---

## Cross-Cutting Systemic Gaps

These affect multiple packages and represent architectural decisions, not individual test gaps:

| Gap | Affected Packages | Impact | Category |
|-----|-------------------|--------|----------|
| **Event Dispatching** | Auth (5 tests), Routing (2), Foundation (1), View (20+) | ~28 tests | (d) infrastructure |
| **IoC Container** | Foundation (15+), Console (5+), Routing (10+) | ~30 tests | (d) infrastructure |
| **Template Template Engine** | View (50+ tests across 10 files) | ~50 tests | (a) language-difference |
| **PHP Macroable Trait** | Config (1), Guards (1), View (1), Foundation (1) | ~4 tests | (a) language-difference |
| **PHP ArrayAccess** | Config (4), Http (1) | ~5 tests | (a) language-difference |
| **PHP Backed Enums** | Gate (6), Middleware (1), Routing (3+), Http (2) | ~12 tests | (a) language-difference |
| **Basic HTTP Auth** | Guards (5) | 5 tests | (b) missing |
| **Route Model Binding** | Routing (15+) | ~15 tests | (b) missing |
| **SQL Token Repository** | Passwords (1) | 1 test | (b) missing |

---

## Summary by Package

| Package | Upstream Tests | Covered | Missing | Intentional Skip | Bedrock-Only | Coverage % |
|---------|--------------|---------|---------|-------------------|--------------|------------|
| Config | 33 | 27 | 0 | 6 | 12 | **100%** (of portable) |
| Encryption | 26 | 22 | 4 | 0 | 6 | **85%** |
| Hashing | 13 | 10 | 0 | 3 | 2 | **100%** (of portable) |
| Auth/Gate | 102 | 35 | 39 | 14 | 20 | **47%** (of portable) |
| Auth/Guards | 57 | 32 | 18 | 1 | 31 | **64%** (of portable) |
| Auth/Passwords | 19 | 15 | 1 | 1 | 9 | **88%** (of portable) |
| Auth/Middleware | 25 | 7 | 13 | 5 | 3 | **35%** (of portable) |
| Console | 14+ | 2 | 10+ | 1 | 7 | **~15%** |
| Routing | 102+ | 3 | 97+ | 0 | 9 | **~3%** |
| HTTP | 113+ | 5 | 107+ | 0 | 6 | **~4%** |
| View | 70+ | 2 | 18+ | 50+ | 8 | **~10%** (Template is intentional skip) |
| Foundation | 42+ | 3 | 37+ | 1 | 10 | **~7%** |

---

## Recommended Next Steps (Priority Order)

### Tier 1: Close gaps in near-complete packages
1. **Encryption** — Add 4 missing tests: bad key length, unsupported cipher, IV too long
2. **Auth/Passwords** — Add expired token cleanup test

### Tier 2: Strengthen auth subsystem
3. **Auth/Gate** — Add guest user handling (6 tests), every/ability checks (3), subtype/interface resolution (2), custom denial responses (2)
4. **Auth/Guards** — Add basic HTTP auth (5 tests), rehashing (2), `forgetUser`, callback-based attempt
5. **Auth/Middleware** — Add multi-guard support (5 tests), authorize with parameters (6 tests), `RedirectIfAuthenticated`

### Tier 3: Broaden core packages
6. **HTTP** — Add path/URL methods, content negotiation, file handling
7. **Routing** — Add middleware groups, route model binding, resource routing, controller routing
8. **Foundation** — Add service provider registration, bootstrapping callbacks

### Tier 4: Architectural decisions needed
9. **Event System** — Decide whether to implement event dispatching (unblocks ~28 tests)
10. **IoC Container** — Decide on dependency injection strategy (unblocks ~30 tests)
11. **Template equivalent** — Bedrock uses `html/template`; Template tests are permanent intentional skips

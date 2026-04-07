package container

// Laravel's AfterResolvingAttributeCallbackTest relies entirely on PHP
// attributes (#[Attribute]) and reflection-based dependency injection.
// Go has no equivalent of PHP attributes, so all 4 Laravel tests in this
// file are intentional skips.
//
// INTENTIONAL-SKIP: testCallbackIsCalledAfterDependencyResolutionWithAttribute
// INTENTIONAL-SKIP: testCallbackIsCalledAfterClassWithAttributeIsResolved
// INTENTIONAL-SKIP: testCallbackIsCalledAfterClassWithConstructorAndAttributeIsResolved
// INTENTIONAL-SKIP: testCallbackIsCalledOnAppCall

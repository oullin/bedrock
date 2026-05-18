package workflow

// Event-name helpers for registering listeners on the typed dispatcher.

func EventNameGuard(workflow string) string { return "workflow." + workflow + ".guard" }

func EventNameGuardNamed(workflow string, transition string) string {
	return "workflow." + workflow + ".guard." + transition
}

func EventNameLeave(workflow string) string { return "workflow." + workflow + ".leave" }

func EventNameLeavePlace(workflow string, place string) string {
	return "workflow." + workflow + ".leave." + place
}

func EventNameTransition(workflow string) string { return "workflow." + workflow + ".transition" }

func EventNameTransitionNamed(workflow string, transition string) string {
	return "workflow." + workflow + ".transition." + transition
}

func EventNameEnter(workflow string) string { return "workflow." + workflow + ".enter" }

func EventNameEnterPlace(workflow string, place string) string {
	return "workflow." + workflow + ".enter." + place
}

func EventNameEntered(workflow string) string { return "workflow." + workflow + ".entered" }

func EventNameEnteredPlace(workflow string, place string) string {
	return "workflow." + workflow + ".entered." + place
}

func EventNameCompleted(workflow string) string { return "workflow." + workflow + ".completed" }

func EventNameCompletedNamed(workflow string, transition string) string {
	return "workflow." + workflow + ".completed." + transition
}

func EventNameAnnounce(workflow string) string { return "workflow." + workflow + ".announce" }

func EventNameAnnounceNamed(workflow string, transition string) string {
	return "workflow." + workflow + ".announce." + transition
}

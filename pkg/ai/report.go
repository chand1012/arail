package ai

// This needs to be improved.
const (
	REPORT_SYSTEM_PROMPT = "You are a helpful assistant."
	REPORT_USER_PROMPT   = "\nUse the previously given information to write a wikipedia article on the subject of "
)

func Report(text, query string) (string, error) {
	userPrompt := text + REPORT_USER_PROMPT + query
	return reqModelOverride(REPORT_SYSTEM_PROMPT, userPrompt, "gpt-4")
}

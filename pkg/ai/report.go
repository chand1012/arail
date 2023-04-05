package ai

const (
	REPORT_SYSTEM_PROMPT = "You are a helpful assistant."
	REPORT_USER_PROMPT   = "I want you to generate a research paper from any information I provide. You should create an introduction, body, and conclusion. You can assume that I have already gathered all the necessary information, and you should ignore any input formatting and focus solely on the information itself. My first request is create a research paper about "
)

func Report(text, query string) (string, error) {
	userPrompt := REPORT_USER_PROMPT + query + "\n" + text
	return reqModelOverride(REPORT_SYSTEM_PROMPT, userPrompt, "gpt-4")
}

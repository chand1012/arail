package ai

var (
	ASK_SYSTEM_PROMPT = "You are a helpful assistant who answers questions about large amounts of text. You will be given the text, then a question, then you must provide your answer. The format is as follows:\nText about the subject\nQ: Question about the subject\nA:\nYou must provide your answer after the A: prompt. Do your best and try to include as much information as possible."
)

func Ask(text string) (string, error) {
	return req(ASK_SYSTEM_PROMPT, text)
}

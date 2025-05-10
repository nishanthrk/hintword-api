package agent_service

type Response struct {
	Choices []struct {
		ContentFilterResults struct {
			Hate struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"hate"`
			ProtectedMaterialCode struct {
				Filtered bool `json:"filtered"`
				Detected bool `json:"detected"`
			} `json:"protected_material_code"`
			ProtectedMaterialText struct {
				Filtered bool `json:"filtered"`
				Detected bool `json:"detected"`
			} `json:"protected_material_text"`
			SelfHarm struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"self_harm"`
			Sexual struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"sexual"`
			Violence struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"violence"`
		} `json:"content_filter_results"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		Message      struct {
			Annotations []interface{} `json:"annotations"`
			Content     string        `json:"content"`
			Refusal     interface{}   `json:"refusal"`
			Role        string        `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Created             int    `json:"created"`
	Id                  string `json:"id"`
	Model               string `json:"model"`
	Object              string `json:"object"`
	PromptFilterResults []struct {
		PromptIndex          int `json:"prompt_index"`
		ContentFilterResults struct {
			Hate struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"hate"`
			Jailbreak struct {
				Filtered bool `json:"filtered"`
				Detected bool `json:"detected"`
			} `json:"jailbreak"`
			SelfHarm struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"self_harm"`
			Sexual struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"sexual"`
			Violence struct {
				Filtered bool   `json:"filtered"`
				Severity string `json:"severity"`
			} `json:"violence"`
		} `json:"content_filter_results"`
	} `json:"prompt_filter_results"`
	SystemFingerprint string `json:"system_fingerprint"`
	Usage             struct {
		CompletionTokens        int `json:"completion_tokens"`
		CompletionTokensDetails struct {
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			AudioTokens              int `json:"audio_tokens"`
			ReasoningTokens          int `json:"reasoning_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
		PromptTokens        int `json:"prompt_tokens"`
		PromptTokensDetails struct {
			AudioTokens  int `json:"audio_tokens"`
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

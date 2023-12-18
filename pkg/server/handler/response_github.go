package handler

type gitHubfileContentResponse struct {
	Content *string `json:"content"`
}

func newGitHubfileContentResponse(content *string) *gitHubfileContentResponse {
	return &gitHubfileContentResponse{
		Content: content,
	}
}
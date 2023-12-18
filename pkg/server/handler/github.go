package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/labstack/echo/v4"
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/rs/zerolog/log"
)

type gitHubUrlParts struct {
	Owner string
	Repo string
	Ref string
	Path string
}

func getGitHubUrlParts(url string) (*gitHubUrlParts, error) {
	ret := &gitHubUrlParts{}

	if !strings.HasPrefix(url, "https://github.com") {
		return nil, errors.New("URL should begin with https://github.com")
	}

	parts := strings.Split(url[19:], "/")
	if len(parts) < 5 {
		return nil, errors.New("URL should be of format: https://github.com/OWNER/REPO/blob/REF/path/to/filename.extension")
	}

	ret.Owner = parts[0]
	ret.Repo = parts[1]
	ret.Ref = parts[3]
	ret.Path = strings.Join(parts[4:], "/")

	return ret, nil
}

func (h *Handler) RetrieveGitHubFileContent(c echo.Context) error {
	r := &gitHubfileContentRetrieveRequest{}
	if err := r.bind(c); err != nil {
		return unprocessableEntity(c, err)
	}

	client := github.NewClient(nil)
	if r.CredentialId != nil {
		credential, err := h.credentialStore.GetById(*r.CredentialId)
		if err != nil {
			return unprocessableEntity(c, errors.New("Credentials not found"))
		}
		
		decryptedSecret, err := ske.Decrypt(credential.Secret)
		if err != nil {
			panic(err)
		}
		client = github.NewClient(nil).WithAuthToken(decryptedSecret)
	}

	p, err := getGitHubUrlParts(r.Url)
	if err != nil {
		return unprocessableEntity(c, err)
	}

	content, _, _, err := client.Repositories.GetContents(context.Background(), p.Owner, p.Repo, p.Path, &github.RepositoryContentGetOptions{Ref: p.Ref})
	if err != nil {
		if r.CredentialId != nil {
			log.Error().Err(err).Str("url", r.Url).Uint("credentialId", *r.CredentialId).Msg("Error while retriveing file content from GitHub")
		} else {
			log.Error().Err(err).Str("url", r.Url).Msg("Error while retriveing file content from GitHub")
		}
		return unprocessableEntity(c, errors.New("Error while retrieving file content from provided GitHub URL"))
	}

	text, err := content.GetContent()
	if err != nil {
		panic(err)
	}

	return ok(c, newGitHubfileContentResponse(&text))
}
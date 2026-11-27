package webhook

import (
	"encoding/json"
	"strings"
)

type githubCommit struct {
	ID       string   `json:"id"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type githubPayload struct {
	Ref        string         `json:"ref"`
	After      string         `json:"after"`
	HeadCommit *githubCommit  `json:"head_commit"`
	Commits    []githubCommit `json:"commits"`
	Sender     struct {
		Login string `json:"login"`
	} `json:"sender"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		CloneURL string `json:"clone_url"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Action      string `json:"action"`
	PullRequest *struct {
		Head struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"head"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
	Release *struct {
		TagName string `json:"tag_name"`
	} `json:"release"`
}

// parseGitHub maps a GitHub event and payload to a delivery. handled is false
// for an event or action this receiver does not start runs for.
func parseGitHub(event string, body []byte) (delivery Delivery, handled bool) {
	var payload githubPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return Delivery{}, false
	}
	delivery.Repository = payload.Repository.CloneURL
	if delivery.Repository == "" {
		delivery.Repository = payload.Repository.HTMLURL
	}
	delivery.Actor = payload.Sender.Login
	if delivery.Actor == "" {
		delivery.Actor = payload.Pusher.Name
	}

	switch event {
	case "ping":
		delivery.Event = eventPing
		return delivery, true
	case "push":
		delivery.Event = eventPush
		delivery.Branch, delivery.ReleaseTag = branchFromRef(payload.Ref)
		if payload.HeadCommit != nil {
			delivery.Commit = payload.HeadCommit.ID
		} else {
			delivery.Commit = payload.After
		}
		for _, commit := range payload.Commits {
			delivery.Files = append(delivery.Files, commit.Added...)
			delivery.Files = append(delivery.Files, commit.Removed...)
			delivery.Files = append(delivery.Files, commit.Modified...)
		}
		return delivery, true
	case "pull_request":
		switch payload.Action {
		case "opened", "synchronize", "reopened":
		default:
			return Delivery{}, false
		}
		if payload.PullRequest == nil {
			return Delivery{}, false
		}
		delivery.Event = eventPullRequest
		delivery.Branch = payload.PullRequest.Head.Ref
		delivery.Commit = payload.PullRequest.Head.SHA
		if payload.PullRequest.User.Login != "" {
			delivery.Actor = payload.PullRequest.User.Login
		}
		// GitHub only reports the changed-file count in the payload; the file
		// list needs an API call, so a pull request matches on its branch.
		return delivery, true
	case "release":
		if payload.Action != "created" || payload.Release == nil {
			return Delivery{}, false
		}
		delivery.Event = eventRelease
		delivery.ReleaseTag = payload.Release.TagName
		delivery.Branch = payload.Release.TagName
		return delivery, true
	default:
		return Delivery{}, false
	}
}

type gitlabCommit struct {
	ID       string   `json:"id"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type gitlabPayload struct {
	ObjectKind   string         `json:"object_kind"`
	Ref          string         `json:"ref"`
	CheckoutSHA  string         `json:"checkout_sha"`
	UserUsername string         `json:"user_username"`
	UserName     string         `json:"user_name"`
	Commits      []gitlabCommit `json:"commits"`
	Project      struct {
		GitHTTPURL        string `json:"git_http_url"`
		WebURL            string `json:"web_url"`
		PathWithNamespace string `json:"path_with_namespace"`
	} `json:"project"`
	Action           string `json:"action"`
	Tag              string `json:"tag"`
	ObjectAttributes *struct {
		Action       string `json:"action"`
		SourceBranch string `json:"source_branch"`
		LastCommit   *struct {
			ID string `json:"id"`
		} `json:"last_commit"`
	} `json:"object_attributes"`
}

// parseGitLab maps a GitLab event and payload to a delivery. GitLab names the
// event in the body instead of a header.
func parseGitLab(body []byte) (delivery Delivery, handled bool) {
	var payload gitlabPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return Delivery{}, false
	}
	delivery.Repository = payload.Project.GitHTTPURL
	if delivery.Repository == "" {
		delivery.Repository = payload.Project.WebURL
	}
	delivery.Actor = payload.UserUsername
	if delivery.Actor == "" {
		delivery.Actor = payload.UserName
	}

	switch payload.ObjectKind {
	case "push", "tag_push":
		delivery.Event = eventPush
		delivery.Branch, delivery.ReleaseTag = branchFromRef(payload.Ref)
		delivery.Commit = payload.CheckoutSHA
		for _, commit := range payload.Commits {
			delivery.Files = append(delivery.Files, commit.Added...)
			delivery.Files = append(delivery.Files, commit.Removed...)
			delivery.Files = append(delivery.Files, commit.Modified...)
		}
		return delivery, true
	case "merge_request":
		if payload.ObjectAttributes == nil {
			return Delivery{}, false
		}
		switch payload.ObjectAttributes.Action {
		case "open", "update", "reopen":
		default:
			return Delivery{}, false
		}
		delivery.Event = eventPullRequest
		delivery.Branch = payload.ObjectAttributes.SourceBranch
		if payload.ObjectAttributes.LastCommit != nil {
			delivery.Commit = payload.ObjectAttributes.LastCommit.ID
		}
		// GitLab reports merge request changes only through its API.
		return delivery, true
	case "release":
		if payload.Action != "" && payload.Action != "create" {
			return Delivery{}, false
		}
		delivery.Event = eventRelease
		delivery.ReleaseTag = payload.Tag
		delivery.Branch = payload.Tag
		return delivery, true
	case "ping":
		delivery.Event = eventPing
		return delivery, true
	default:
		return Delivery{}, false
	}
}

func trimDedupe(files []string) []string {
	seen := make(map[string]bool, len(files))
	kept := files[:0]
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" || seen[file] {
			continue
		}
		seen[file] = true
		kept = append(kept, file)
	}
	return kept
}

package github

// FilterByOrganizations keeps the pull requests whose repository owner is one
// of orgs, preserving their order.
func FilterByOrganizations(prs []PullRequest, orgs []string) []PullRequest {
	allowed := make(map[string]struct{}, len(orgs))
	for _, org := range orgs {
		allowed[org] = struct{}{}
	}

	var filtered []PullRequest
	for _, pr := range prs {
		if _, ok := allowed[pr.Owner()]; ok {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

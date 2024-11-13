package ghclient

func calcTotalForksCount(repos []RepoData) int {
	count := 0
	for _, r := range repos {
		count = count + r.ForksCount
	}
	return count
}

func calcUserActivity(repos []RepoData) UserActivity {
	userActivity := make(map[int]int)
	for _, r := range repos {
		pushedAt := r.PushedAt.Year()
		createdAt := r.CreatedAt.Year()
		userActivity[pushedAt] += 1
		if createdAt < pushedAt {
			for y := createdAt; y < pushedAt; y++ {
				userActivity[y] += 1
			}
		}
	}
	return userActivity
}

package i18n

// Key identifies a translatable message.
type Key string

const (
	NotAGitRepo   Key = "not_a_git_repo"
	NothingStaged Key = "nothing_staged"
)

var catalog = map[Locale]map[Key]string{
	En: {
		NotAGitRepo:   "nice try, buddy, but you're not in a git repo",
		NothingStaged: "nothing staged; run 'git add' first",
	},
	De: {
		NotAGitRepo:   "kollege, netter Versuch, aber du bist in keinem Git-Repo",
		NothingStaged: "nichts vorgemerkt; erst 'git add' ausführen",
	},
}

// T returns the message for key in the given locale, falling back to
// English if the locale or key is missing from the catalog.
func T(locale Locale, key Key) string {
	if msgs, ok := catalog[locale]; ok {
		if m, ok := msgs[key]; ok {
			return m
		}
	}
	return catalog[En][key]
}

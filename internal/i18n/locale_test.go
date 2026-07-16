package i18n

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name        string
		cfgLanguage string
		env         map[string]string
		want        Locale
	}{
		{"config german", "german", nil, De},
		{"config deutsch", "deutsch", nil, De},
		{"config de", "de", nil, De},
		{"config english", "english", nil, En},
		{"config overrides env", "german", map[string]string{"LANG": "en_US.UTF-8"}, De},
		{"env LANG de", "", map[string]string{"LANG": "de_DE.UTF-8"}, De},
		{"env LC_ALL takes priority over LANG", "", map[string]string{"LC_ALL": "de_DE.UTF-8", "LANG": "en_US.UTF-8"}, De},
		{"unrecognized config falls back to env", "conventional", map[string]string{"LANG": "de_DE.UTF-8"}, De},
		{"no config no env defaults english", "", nil, En},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
				t.Setenv(key, "")
			}
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			if got := Detect(tt.cfgLanguage); got != tt.want {
				t.Errorf("Detect(%q) = %q, want %q", tt.cfgLanguage, got, tt.want)
			}
		})
	}
}

func TestT(t *testing.T) {
	if got := T(De, NotAGitRepo); got == "" {
		t.Error("T(De, NotAGitRepo) = empty")
	}
	if got := T(En, NotAGitRepo); got == "" {
		t.Error("T(En, NotAGitRepo) = empty")
	}
	if got := T(Locale("fr"), NothingStaged); got != T(En, NothingStaged) {
		t.Errorf("T(unknown locale) = %q, want English fallback %q", got, T(En, NothingStaged))
	}
}

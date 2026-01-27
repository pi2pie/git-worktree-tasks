package cli

import "testing"

func TestApplyMergeMode(t *testing.T) {
	cases := []struct {
		name   string
		mode   string
		noFF   bool
		squash bool
		rebase bool
	}{
		{"default", "ff", false, false, false},
		{"no-ff", "no-ff", true, false, false},
		{"squash", "squash", false, true, false},
		{"rebase", "rebase", false, false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &finishOptions{}
			if err := applyMergeMode(opts, tc.mode); err != nil {
				t.Fatalf("applyMergeMode() error = %v", err)
			}
			if opts.noFF != tc.noFF || opts.squash != tc.squash || opts.rebase != tc.rebase {
				t.Fatalf("applyMergeMode() flags = noFF:%v squash:%v rebase:%v, want noFF:%v squash:%v rebase:%v",
					opts.noFF, opts.squash, opts.rebase, tc.noFF, tc.squash, tc.rebase,
				)
			}
		})
	}
}

func TestApplyMergeMode_Invalid(t *testing.T) {
	opts := &finishOptions{}
	if err := applyMergeMode(opts, "invalid"); err == nil {
		t.Fatalf("applyMergeMode() expected error")
	}
}

func TestValidateMergeStrategy(t *testing.T) {
	cases := []struct {
		name    string
		noFF    bool
		squash  bool
		rebase  bool
		wantErr bool
	}{
		{"none", false, false, false, false},
		{"no-ff", true, false, false, false},
		{"squash", false, true, false, false},
		{"rebase", false, false, true, false},
		{"two", true, true, false, true},
		{"three", true, true, true, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &finishOptions{noFF: tc.noFF, squash: tc.squash, rebase: tc.rebase}
			err := validateMergeStrategy(opts)
			if tc.wantErr && err == nil {
				t.Fatalf("validateMergeStrategy() expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("validateMergeStrategy() error = %v", err)
			}
		})
	}
}

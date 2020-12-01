const progressPrinterDiff = `diff --git README.md README.md
new file mode 100644
index 0000000..3363c39
--- /dev/null
+++ README.md
@@ -0,0 +1,3 @@
+# README
+
+This is the readme
diff --git a/b/c/c.txt a/b/c/c.txt
deleted file mode 100644
index 5da75cf..0000000
--- a/b/c/c.txt
+++ /dev/null
@@ -1 +0,0 @@
-this is c
diff --git x/x.txt x/x.txt
index 627c2ae..88f1836 100644
--- x/x.txt
+++ x/x.txt
@@ -1 +1 @@
-this is x
+this is x (or is it?)
`

	statuses := []*campaigns.TaskStatus{
		{
			RepoName:           "github.com/sourcegraph/automation-testing",
			StartedAt:          time.Now().Add(time.Duration(-5) * time.Second),
			CurrentlyExecuting: "echo Hello World > README.md",
		},
	}

	printer := newCampaignProgressPrinter(out, true, 4)
	printer.forceNoSpinner = true
	// Print with all three tasks running
	printer.PrintStatuses(statuses)
		"⠋  Executing... (0/3, 0 errored)  ",
		"├── github.com/sourcegraph/src-cli      Downloading archive                   0s",
		"└── github.com/sourcegraph/automati...  echo Hello World > README.md          0s",
	}
	if !cmp.Equal(want, have) {
		t.Fatalf("wrong output:\n%s", cmp.Diff(want, have))
	}

	// Now mark the last task as completed
	statuses[len(statuses)-1] = &campaigns.TaskStatus{
		RepoName:           "github.com/sourcegraph/automation-testing",
		StartedAt:          time.Now().Add(time.Duration(-5) * time.Second),
		FinishedAt:         time.Now().Add(time.Duration(5) * time.Second),
		CurrentlyExecuting: "",
		Err:                nil,
		ChangesetSpec: &campaigns.ChangesetSpec{
			BaseRepository: "graphql-id",
			CreatedChangeset: &campaigns.CreatedChangeset{
				BaseRef:        "refs/heads/main",
				BaseRev:        "d34db33f",
				HeadRepository: "graphql-id",
				HeadRef:        "refs/heads/my-campaign",
				Title:          "This is my campaign",
				Body:           "This is my campaign",
				Commits: []campaigns.GitCommitDescription{
					{
						Message: "This is my campaign",
						Diff:    progressPrinterDiff,
					},
				},
				Published: false,
			},
		},
	}

	printer.PrintStatuses(statuses)
	have = buf.Lines()
	want = []string{
		"github.com/sourcegraph/automation-testing",
		"\tREADME.md   | 3 +++",
		"\ta/b/c/c.txt | 1 -",
		"\tx/x.txt     | 2 +-",
		"  3 files changed, 4 insertions, 2 deletions",
		"  Execution took 10s",
		"",
		"⠋  Executing... (1/3, 0 errored)  ███████████████▍",
		"│                                                                                 0s", // Not sure why this is here, bug?
		"├── github.com/sourcegraph/sourcegraph  echo Hello World > README.md          0s",
		"├── github.com/sourcegraph/src-cli      Downloading archive                   0s",
		"└── github.com/sourcegraph/automati...  3 files changed ++++               0s",
	}
	if !cmp.Equal(want, have) {
		t.Fatalf("wrong output:\n%s", cmp.Diff(want, have))
	// Print again to make sure we get the same result
	printer.PrintStatuses(statuses)
	have = buf.Lines()
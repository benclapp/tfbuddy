package gitlab_hooks

import (
	"fmt"
	"os"
	"testing"

	"github.com/zapier/tfbuddy/pkg/allow_list"
	"github.com/zapier/tfbuddy/pkg/comment_actions"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zapier/tfbuddy/pkg/mocks"
	"github.com/zapier/tfbuddy/pkg/runstream"
	"github.com/zapier/tfbuddy/pkg/tfc_api"
	"github.com/zapier/tfbuddy/pkg/tfc_trigger"
	"github.com/zapier/tfbuddy/pkg/vcs"
)

func Test_parseCommentCommand(t *testing.T) {
	type args struct {
		noteBody string
	}
	tests := []struct {
		name    string
		args    args
		want    *comment_actions.CommentOpts
		wantErr bool
	}{
		{
			name: "tfc plan (all workspaces)",
			args: args{"tfc plan"},
			want: &comment_actions.CommentOpts{
				Args: comment_actions.CommentArgs{
					Agent:   "tfc",
					Command: "plan",
					Rest:    nil,
				},
				Workspace: "",
			},
			wantErr: false,
		},
		{
			name: "tfc apply (all workspaces)",
			args: args{"tfc apply"},
			want: &comment_actions.CommentOpts{
				Args: comment_actions.CommentArgs{
					Agent:   "tfc",
					Command: "apply",
					Rest:    nil,
				},
				Workspace: "",
			},
			wantErr: false,
		},
		{
			name: "tfc plan (single workspaces)",
			args: args{"tfc plan -w service-foo"},
			want: &comment_actions.CommentOpts{
				Args: comment_actions.CommentArgs{
					Agent:   "tfc",
					Command: "plan",
					Rest:    nil,
				},
				Workspace: "service-foo",
			},
			wantErr: false,
		},
		{
			name: "tfc apply (single workspaces)",
			args: args{"tfc apply -w service-foo"},
			want: &comment_actions.CommentOpts{
				Args: comment_actions.CommentArgs{
					Agent:   "tfc",
					Command: "apply",
					Rest:    nil,
				},
				Workspace: "service-foo",
			},
			wantErr: false,
		},
		{
			name:    "not tfc command",
			args:    args{"amazing gitlab review comment"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := comment_actions.ParseCommentCommand(tt.args.noteBody)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommentCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestProcessNoteEventPlanError(t *testing.T) {
	os.Setenv(allow_list.GitlabProjectAllowListEnv, "zapier/")
	defer os.Unsetenv(allow_list.GitlabProjectAllowListEnv)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockGitClient := mocks.NewMockGitClient(mockCtrl)

	mockApiClient := mocks.NewMockApiClient(mockCtrl)
	mockStreamClient := mocks.NewMockStreamClient(mockCtrl)
	mockProject := mocks.NewMockProject(mockCtrl)
	mockProject.EXPECT().GetPathWithNamespace().Return("zapier/service-tf-buddy")

	mockLastCommit := mocks.NewMockCommit(mockCtrl)
	mockLastCommit.EXPECT().GetSHA().Return("abvc12345")

	mockAttributes := mocks.NewMockMRAttributes(mockCtrl)

	mockAttributes.EXPECT().GetNote().Return("tfc plan -w service-tf-buddy")
	mockAttributes.EXPECT().GetType().Return("SomeNote")

	mockMREvent := mocks.NewMockMRCommentEvent(mockCtrl)
	mockMREvent.EXPECT().GetProject().Return(mockProject)
	mockMREvent.EXPECT().GetAttributes().Return(mockAttributes).Times(2)
	mockMREvent.EXPECT().GetLastCommit().Return(mockLastCommit)

	mockSimpleMR := mocks.NewMockMR(mockCtrl)
	mockSimpleMR.EXPECT().GetSourceBranch().Return("DTA-2009")

	mockSimpleMR.EXPECT().GetInternalID().Return(101)
	mockMREvent.EXPECT().GetMR().Return(mockSimpleMR).Times(2)

	mockTFCTrigger := mocks.NewMockTrigger(mockCtrl)
	mockTFCConfig := mocks.NewMockTriggerConfig(mockCtrl)
	mockTFCConfig.EXPECT().SetAction(tfc_trigger.PlanAction)
	mockTFCConfig.EXPECT().SetWorkspace("service-tf-buddy")
	mockTFCTrigger.EXPECT().GetConfig().Return(mockTFCConfig).Times(2)
	mockTFCTrigger.EXPECT().TriggerTFCEvents().Return(nil, fmt.Errorf("something went wrong"))

	client := &GitlabEventWorker{
		gl:        mockGitClient,
		tfc:       mockApiClient,
		runstream: mockStreamClient,
		triggerCreation: func(gl vcs.GitClient, tfc tfc_api.ApiClient, runstream runstream.StreamClient, cfg tfc_trigger.TriggerConfig) tfc_trigger.Trigger {
			return mockTFCTrigger
		},
	}

	proj, err := client.processNoteEvent(mockMREvent)
	if err == nil {
		t.Error("expected error")
		return
	}
	if proj != "zapier/service-tf-buddy" {
		t.Error("unexpected project")
		return
	}
}
func TestProcessNoteEventPlan(t *testing.T) {
	os.Setenv(allow_list.GitlabProjectAllowListEnv, "zapier/")
	defer os.Unsetenv(allow_list.GitlabProjectAllowListEnv)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockGitClient := mocks.NewMockGitClient(mockCtrl)

	mockApiClient := mocks.NewMockApiClient(mockCtrl)
	mockStreamClient := mocks.NewMockStreamClient(mockCtrl)
	mockProject := mocks.NewMockProject(mockCtrl)
	mockProject.EXPECT().GetPathWithNamespace().Return("zapier/service-tf-buddy")

	mockLastCommit := mocks.NewMockCommit(mockCtrl)
	mockLastCommit.EXPECT().GetSHA().Return("abvc12345")

	mockAttributes := mocks.NewMockMRAttributes(mockCtrl)

	mockAttributes.EXPECT().GetNote().Return("tfc plan -w service-tf-buddy")
	mockAttributes.EXPECT().GetType().Return("SomeNote")

	mockMREvent := mocks.NewMockMRCommentEvent(mockCtrl)
	mockMREvent.EXPECT().GetProject().Return(mockProject)
	mockMREvent.EXPECT().GetAttributes().Return(mockAttributes).Times(2)
	mockMREvent.EXPECT().GetLastCommit().Return(mockLastCommit)

	mockSimpleMR := mocks.NewMockMR(mockCtrl)
	mockSimpleMR.EXPECT().GetSourceBranch().Return("DTA-2009")

	mockSimpleMR.EXPECT().GetInternalID().Return(101)
	mockMREvent.EXPECT().GetMR().Return(mockSimpleMR).Times(2)

	mockTFCTrigger := mocks.NewMockTrigger(mockCtrl)
	mockTFCConfig := mocks.NewMockTriggerConfig(mockCtrl)
	mockTFCConfig.EXPECT().SetAction(tfc_trigger.PlanAction)
	mockTFCConfig.EXPECT().SetWorkspace("service-tf-buddy")
	mockTFCTrigger.EXPECT().GetConfig().Return(mockTFCConfig).Times(2)
	mockTFCTrigger.EXPECT().TriggerTFCEvents().Return(&tfc_trigger.TriggeredTFCWorkspaces{
		Executed: []string{"service-tf-buddy"},
	}, nil)

	client := &GitlabEventWorker{
		gl:        mockGitClient,
		tfc:       mockApiClient,
		runstream: mockStreamClient,
		triggerCreation: func(gl vcs.GitClient, tfc tfc_api.ApiClient, runstream runstream.StreamClient, cfg tfc_trigger.TriggerConfig) tfc_trigger.Trigger {
			return mockTFCTrigger
		},
	}

	proj, err := client.processNoteEvent(mockMREvent)
	if err != nil {
		t.Fatal(err)
	}
	if proj != "zapier/service-tf-buddy" {
		t.Fatal("expected a project name to be returned")
	}
}

func TestProcessNoteEventPlanFailedWorkspace(t *testing.T) {
	os.Setenv(allow_list.GitlabProjectAllowListEnv, "zapier/")
	defer os.Unsetenv(allow_list.GitlabProjectAllowListEnv)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockGitClient := mocks.NewMockGitClient(mockCtrl)

	mockGitClient.EXPECT().CreateMergeRequestComment(101, "zapier/service-tf-buddy", ":no_entry: service-tf-buddy could not be run because: could not fetch upstream\n").Return(nil)

	mockApiClient := mocks.NewMockApiClient(mockCtrl)
	mockStreamClient := mocks.NewMockStreamClient(mockCtrl)
	mockProject := mocks.NewMockProject(mockCtrl)
	mockProject.EXPECT().GetPathWithNamespace().Return("zapier/service-tf-buddy").Times(2)

	mockLastCommit := mocks.NewMockCommit(mockCtrl)
	mockLastCommit.EXPECT().GetSHA().Return("abvc12345")

	mockAttributes := mocks.NewMockMRAttributes(mockCtrl)

	mockAttributes.EXPECT().GetNote().Return("tfc plan -w service-tf-buddy")
	mockAttributes.EXPECT().GetType().Return("SomeNote")

	mockMREvent := mocks.NewMockMRCommentEvent(mockCtrl)
	mockMREvent.EXPECT().GetProject().Return(mockProject).Times(2)
	mockMREvent.EXPECT().GetAttributes().Return(mockAttributes).Times(2)
	mockMREvent.EXPECT().GetLastCommit().Return(mockLastCommit)

	mockSimpleMR := mocks.NewMockMR(mockCtrl)
	mockSimpleMR.EXPECT().GetSourceBranch().Return("DTA-2009")

	mockSimpleMR.EXPECT().GetInternalID().Return(101).Times(2)
	mockMREvent.EXPECT().GetMR().Return(mockSimpleMR).Times(3)

	mockTFCTrigger := mocks.NewMockTrigger(mockCtrl)
	mockTFCConfig := mocks.NewMockTriggerConfig(mockCtrl)
	mockTFCConfig.EXPECT().SetAction(tfc_trigger.PlanAction)
	mockTFCConfig.EXPECT().SetWorkspace("service-tf-buddy")
	mockTFCTrigger.EXPECT().GetConfig().Return(mockTFCConfig).Times(2)
	mockTFCTrigger.EXPECT().TriggerTFCEvents().Return(&tfc_trigger.TriggeredTFCWorkspaces{
		Errored: []*tfc_trigger.ErroredWorkspace{{
			Name:  "service-tf-buddy",
			Error: "could not fetch upstream",
		},
		},
	}, nil)

	worker := &GitlabEventWorker{
		tfc:       mockApiClient,
		gl:        mockGitClient,
		runstream: mockStreamClient,
		triggerCreation: func(gl vcs.GitClient, tfc tfc_api.ApiClient, runstream runstream.StreamClient, cfg tfc_trigger.TriggerConfig) tfc_trigger.Trigger {
			return mockTFCTrigger
		},
	}

	proj, err := worker.processNoteEvent(mockMREvent)
	if err != nil {
		t.Fatal(err)
	}
	if proj != "zapier/service-tf-buddy" {
		t.Fatal("expected a project name to be returned")
	}
}

func TestProcessNoteEventPlanFailedMultipleWorkspaces(t *testing.T) {
	os.Setenv(allow_list.GitlabProjectAllowListEnv, "zapier/")
	defer os.Unsetenv(allow_list.GitlabProjectAllowListEnv)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockGitClient := mocks.NewMockGitClient(mockCtrl)

	mockGitClient.EXPECT().CreateMergeRequestComment(101, "zapier/service-tf-buddy", ":no_entry: service-tf-buddy could not be run because: could not fetch upstream\nservice-tf-buddy-staging could not be run because: workspace has been modified on target branch\n").Return(nil)

	mockApiClient := mocks.NewMockApiClient(mockCtrl)
	mockStreamClient := mocks.NewMockStreamClient(mockCtrl)
	mockProject := mocks.NewMockProject(mockCtrl)
	mockProject.EXPECT().GetPathWithNamespace().Return("zapier/service-tf-buddy").Times(2)

	mockLastCommit := mocks.NewMockCommit(mockCtrl)
	mockLastCommit.EXPECT().GetSHA().Return("abvc12345")

	mockAttributes := mocks.NewMockMRAttributes(mockCtrl)

	mockAttributes.EXPECT().GetNote().Return("tfc plan")
	mockAttributes.EXPECT().GetType().Return("SomeNote")

	mockMREvent := mocks.NewMockMRCommentEvent(mockCtrl)
	mockMREvent.EXPECT().GetProject().Return(mockProject).Times(2)
	mockMREvent.EXPECT().GetAttributes().Return(mockAttributes).Times(2)
	mockMREvent.EXPECT().GetLastCommit().Return(mockLastCommit)

	mockSimpleMR := mocks.NewMockMR(mockCtrl)
	mockSimpleMR.EXPECT().GetSourceBranch().Return("DTA-2009")

	mockSimpleMR.EXPECT().GetInternalID().Return(101).Times(2)
	mockMREvent.EXPECT().GetMR().Return(mockSimpleMR).Times(3)

	mockTFCTrigger := mocks.NewMockTrigger(mockCtrl)
	mockTFCConfig := mocks.NewMockTriggerConfig(mockCtrl)
	mockTFCConfig.EXPECT().SetAction(tfc_trigger.PlanAction)
	mockTFCConfig.EXPECT().SetWorkspace("")
	mockTFCTrigger.EXPECT().GetConfig().Return(mockTFCConfig).Times(2)
	mockTFCTrigger.EXPECT().TriggerTFCEvents().Return(&tfc_trigger.TriggeredTFCWorkspaces{
		Errored: []*tfc_trigger.ErroredWorkspace{{
			Name:  "service-tf-buddy",
			Error: "could not fetch upstream",
		},
			{
				Name:  "service-tf-buddy-staging",
				Error: "workspace has been modified on target branch",
			},
		},
	}, nil)

	client := &GitlabEventWorker{
		gl:        mockGitClient,
		tfc:       mockApiClient,
		runstream: mockStreamClient,
		triggerCreation: func(gl vcs.GitClient, tfc tfc_api.ApiClient, runstream runstream.StreamClient, cfg tfc_trigger.TriggerConfig) tfc_trigger.Trigger {
			return mockTFCTrigger
		},
	}

	proj, err := client.processNoteEvent(mockMREvent)
	if err != nil {
		t.Fatal(err)
	}
	if proj != "zapier/service-tf-buddy" {
		t.Fatal("expected a project name to be returned")
	}
}

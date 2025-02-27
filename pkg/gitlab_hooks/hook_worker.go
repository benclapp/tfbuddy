package gitlab_hooks

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/zapier/tfbuddy/pkg/runstream"
	"github.com/zapier/tfbuddy/pkg/tfc_api"
	"github.com/zapier/tfbuddy/pkg/tfc_trigger"
	"github.com/zapier/tfbuddy/pkg/vcs"
)

type GitlabEventWorker struct {
	tfc             tfc_api.ApiClient
	gl              vcs.GitClient
	runstream       runstream.StreamClient
	triggerCreation TriggerCreationFunc
}

func NewGitlabEventWorker(h *GitlabHooksHandler, js nats.JetStreamContext) *GitlabEventWorker {
	w := &GitlabEventWorker{
		tfc:             h.tfc,
		gl:              h.gl,
		runstream:       h.runstream,
		triggerCreation: tfc_trigger.NewTFCTrigger,
	}

	_, err := h.mrStream.QueueSubscribe("gitlab_mr_event_worker", w.processMREventStreamMsg)
	if err != nil {
		log.Error().Err(err).Msg("could not subscribe to hook stream")
	}

	_, err = h.notesStream.QueueSubscribe("gitlab_note_event_worker", w.processNoteEventStreamMsg)
	if err != nil {
		log.Error().Err(err).Msg("could not subscribe to hook stream")
	}

	return w
}

func (w *GitlabEventWorker) processNoteEventStreamMsg(msg *NoteEventMsg) error {

	w.processNoteEvent(msg)

	return nil
}

func (w *GitlabEventWorker) processMREventStreamMsg(msg *MergeRequestEventMsg) error {

	w.processMergeRequestEvent(msg)

	return nil
}

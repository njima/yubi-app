package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
	"github.com/rs/zerolog"
)

// stubEpisodeGradeRepo is a hand-rolled fake for repository.EpisodeGrade.
// It captures the IDs requested and returns a pre-canned aggregate map.
type stubEpisodeGradeRepo struct {
	aggMap     map[string]repository.GradeAggregate
	err        error
	calledWith []string

	upsertResult model.EpisodeGrade
	upsertErr    error
	// lastUpsertArg records the argument passed to Upsert so callers can
	// verify the usecase forwards the validated domain model unchanged.
	lastUpsertArg *model.EpisodeGrade

	myGrade    *model.EpisodeGrade
	myGradeErr error

	listItems       []repository.EpisodeGradeListItem
	listTotal       int
	listErr         error
	lastListEpisode string
	lastListLimit   int
	lastListOffset  int
}

func (s *stubEpisodeGradeRepo) GetAverageMap(_ context.Context, _ repository.Conn, ids []string) (map[string]repository.GradeAggregate, error) {
	s.calledWith = ids
	return s.aggMap, s.err
}

func (s *stubEpisodeGradeRepo) Upsert(_ context.Context, _ repository.Conn, grade model.EpisodeGrade) (model.EpisodeGrade, error) {
	s.lastUpsertArg = &grade
	return s.upsertResult, s.upsertErr
}

func (s *stubEpisodeGradeRepo) GetMyGrade(_ context.Context, _ repository.Conn, _, _ string) (*model.EpisodeGrade, error) {
	return s.myGrade, s.myGradeErr
}

func (s *stubEpisodeGradeRepo) ListByEpisodeID(_ context.Context, _ repository.Conn, episodeID string, limit, offset int) ([]repository.EpisodeGradeListItem, int, error) {
	s.lastListEpisode = episodeID
	s.lastListLimit = limit
	s.lastListOffset = offset
	return s.listItems, s.listTotal, s.listErr
}

func newEpisodeUsecaseWithGradeStub(gradeRepo repository.EpisodeGrade) *episode {
	return &episode{gradeRepo: gradeRepo, logger: zerolog.Nop()}
}

func TestMergeGradeAggregates_PopulatesEpisodesWithMatchingGrades(t *testing.T) {
	stub := &stubEpisodeGradeRepo{
		aggMap: map[string]repository.GradeAggregate{
			"ep-1": {Average: 0.70, Count: 2},
			"ep-2": {Average: 0.85, Count: 3},
		},
	}
	uc := newEpisodeUsecaseWithGradeStub(stub)

	episodes := model.Episodes{
		{IDNatural: "ep-1"},
		{IDNatural: "ep-2"},
		{IDNatural: "ep-3"}, // no grades — should remain nil/0
	}

	uc.mergeGradeAggregates(context.Background(), episodes)

	if episodes[0].AverageGrade == nil || *episodes[0].AverageGrade != 0.70 {
		t.Errorf("episodes[0].AverageGrade = %v, want 0.70", episodes[0].AverageGrade)
	}
	if episodes[0].GradeCount != 2 {
		t.Errorf("episodes[0].GradeCount = %v, want 2", episodes[0].GradeCount)
	}
	if episodes[1].AverageGrade == nil || *episodes[1].AverageGrade != 0.85 {
		t.Errorf("episodes[1].AverageGrade = %v, want 0.85", episodes[1].AverageGrade)
	}
	if episodes[1].GradeCount != 3 {
		t.Errorf("episodes[1].GradeCount = %v, want 3", episodes[1].GradeCount)
	}
	if episodes[2].AverageGrade != nil {
		t.Errorf("episodes[2].AverageGrade = %v, want nil (no grades)", episodes[2].AverageGrade)
	}
	if episodes[2].GradeCount != 0 {
		t.Errorf("episodes[2].GradeCount = %v, want 0", episodes[2].GradeCount)
	}
}

func TestMergeGradeAggregates_NoEpisodesShortCircuits(t *testing.T) {
	stub := &stubEpisodeGradeRepo{}
	uc := newEpisodeUsecaseWithGradeStub(stub)

	uc.mergeGradeAggregates(context.Background(), model.Episodes{})

	if stub.calledWith != nil {
		t.Errorf("gradeRepo should not be called with empty episode list, got calledWith=%v", stub.calledWith)
	}
}

func TestMergeGradeAggregates_PassesAllEpisodeIDsToRepo(t *testing.T) {
	stub := &stubEpisodeGradeRepo{aggMap: map[string]repository.GradeAggregate{}}
	uc := newEpisodeUsecaseWithGradeStub(stub)

	episodes := model.Episodes{
		{IDNatural: "ep-a"},
		{IDNatural: "ep-b"},
		{IDNatural: "ep-c"},
	}

	uc.mergeGradeAggregates(context.Background(), episodes)

	wantIDs := []string{"ep-a", "ep-b", "ep-c"}
	if len(stub.calledWith) != len(wantIDs) {
		t.Fatalf("calledWith length = %d, want %d", len(stub.calledWith), len(wantIDs))
	}
	for i, want := range wantIDs {
		if stub.calledWith[i] != want {
			t.Errorf("calledWith[%d] = %v, want %v", i, stub.calledWith[i], want)
		}
	}
}

func TestMergeGradeAggregates_RepoErrorLeavesEpisodesUntouched(t *testing.T) {
	wantErr := errors.New("db down")
	stub := &stubEpisodeGradeRepo{err: wantErr}
	uc := newEpisodeUsecaseWithGradeStub(stub)

	episodes := model.Episodes{{IDNatural: "ep-1"}}

	// mergeGradeAggregates must NOT propagate the error; the core read path
	// keeps working even when the grade aggregation dependency is unavailable
	// (e.g. migration skew during rolling deploys).
	uc.mergeGradeAggregates(context.Background(), episodes)

	if episodes[0].AverageGrade != nil {
		t.Errorf("AverageGrade should remain nil on repo error, got %v", episodes[0].AverageGrade)
	}
	if episodes[0].GradeCount != 0 {
		t.Errorf("GradeCount should remain 0 on repo error, got %v", episodes[0].GradeCount)
	}
}

func TestEpisodeUpdateRejectsDirectStatusTransition(t *testing.T) {
	uc := &episode{}

	t.Run("allows no-op status update", func(t *testing.T) {
		status := model.EpisodeStatusRecording
		ep := model.Episode{
			IDNatural:      "ep-1",
			OrganizationID: "org-1",
			TaskID:         "task-1",
			LocationID:     "loc-1",
			RobotID:        "robot-1",
			UserID:         "user-1",
			Status:         model.EpisodeStatusRecording,
		}

		got, err := uc.update(context.Background(), ep, EpisodeUpdateInput{Status: &status})
		if err != nil {
			t.Fatalf("episode.update() error = %v, want nil", err)
		}
		if got.Status != model.EpisodeStatusRecording {
			t.Fatalf("episode.update() Status = %v, want %v", got.Status, model.EpisodeStatusRecording)
		}
	})

	t.Run("rejects lifecycle status change", func(t *testing.T) {
		status := model.EpisodeStatusRecording
		ep := model.Episode{
			IDNatural:      "ep-1",
			OrganizationID: "org-1",
			TaskID:         "task-1",
			LocationID:     "loc-1",
			RobotID:        "robot-1",
			UserID:         "user-1",
			Status:         model.EpisodeStatusReady,
		}

		if _, err := uc.update(context.Background(), ep, EpisodeUpdateInput{Status: &status}); err == nil {
			t.Fatal("episode.update() error = nil, want error")
		}
	})
}

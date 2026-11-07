package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
)

type applicationRepositoryStub struct {
	gotPostingID, gotCandidateID string
	result                       domain.Application
	err                          error
	calls                        int
}

func (s *applicationRepositoryStub) Apply(_ context.Context, postingID, candidateID string) (domain.Application, error) {
	s.calls++
	s.gotPostingID = postingID
	s.gotCandidateID = candidateID
	return s.result, s.err
}
func TestApplicationAdmissionApplyDelegates(t *testing.T) {
	submitted := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	want := domain.Application{ID: "app-id", PostingID: "posting-id", Status: domain.ApplicationAdmitted, SubmittedAt: submitted}
	repo := &applicationRepositoryStub{result: want}
	service := NewApplicationAdmission(repo)
	got, err := service.Apply(context.Background(), "posting-id", "candidate-id")
	if err != nil {
		t.Fatalf("Apply() error=%v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Apply()=%+v, want %+v", got, want)
	}
	if repo.calls != 1 || repo.gotPostingID != "posting-id" || repo.gotCandidateID != "candidate-id" {
		t.Fatalf("repository call=%+v", repo)
	}
}
func TestApplicationAdmissionPropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	service := NewApplicationAdmission(&applicationRepositoryStub{err: wantErr})
	_, err := service.Apply(context.Background(), "posting-id", "candidate-id")
	if !errors.Is(err, wantErr) {
		t.Fatalf("Apply() error=%v, want %v", err, wantErr)
	}
}
func TestApplicationAdmissionRejectsMissingServiceOrRepository(t *testing.T) {
	var nilService *ApplicationAdmission
	if _, err := nilService.Apply(context.Background(), "posting", "candidate"); err == nil {
		t.Fatal("nil service should fail")
	}
	if _, err := NewApplicationAdmission(nil).Apply(context.Background(), "posting", "candidate"); err == nil {
		t.Fatal("nil repository should fail")
	}
}

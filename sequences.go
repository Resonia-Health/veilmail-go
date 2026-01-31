package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// SequencesService handles automation sequence operations.
type SequencesService struct {
	client *httpClient
}

// Create creates a new automation sequence.
func (s *SequencesService) Create(ctx context.Context, params *CreateSequenceParams) (*Sequence, error) {
	var resp Sequence
	err := s.client.post(ctx, "/v1/sequences", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// List returns a paginated list of sequences.
func (s *SequencesService) List(ctx context.Context, params *ListParams) (*ListResponse[Sequence], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Sequence]
	err := s.client.get(ctx, "/v1/sequences", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single sequence by ID.
func (s *SequencesService) Get(ctx context.Context, id string) (*Sequence, error) {
	var resp Sequence
	err := s.client.get(ctx, fmt.Sprintf("/v1/sequences/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a sequence (only DRAFT or PAUSED).
func (s *SequencesService) Update(ctx context.Context, id string, params *UpdateSequenceParams) (*Sequence, error) {
	var resp Sequence
	err := s.client.put(ctx, fmt.Sprintf("/v1/sequences/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete removes a sequence (only DRAFT).
func (s *SequencesService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/sequences/%s", id))
}

// Activate activates a sequence.
func (s *SequencesService) Activate(ctx context.Context, id string) (*Sequence, error) {
	var resp Sequence
	err := s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/activate", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses an active sequence.
func (s *SequencesService) Pause(ctx context.Context, id string) (*Sequence, error) {
	var resp Sequence
	err := s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/pause", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Archive archives a sequence.
func (s *SequencesService) Archive(ctx context.Context, id string) (*Sequence, error) {
	var resp Sequence
	err := s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/archive", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddStep adds a step to a sequence.
func (s *SequencesService) AddStep(ctx context.Context, sequenceID string, params *CreateSequenceStepParams) (*SequenceStep, error) {
	var resp SequenceStep
	err := s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/steps", sequenceID), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateStep updates a sequence step.
func (s *SequencesService) UpdateStep(ctx context.Context, sequenceID, stepID string, params *UpdateSequenceStepParams) (*SequenceStep, error) {
	var resp SequenceStep
	err := s.client.put(ctx, fmt.Sprintf("/v1/sequences/%s/steps/%s", sequenceID, stepID), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteStep removes a step from a sequence.
func (s *SequencesService) DeleteStep(ctx context.Context, sequenceID, stepID string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/sequences/%s/steps/%s", sequenceID, stepID))
}

// ReorderSteps reorders the steps of a sequence.
func (s *SequencesService) ReorderSteps(ctx context.Context, sequenceID string, steps []StepPosition) error {
	body := map[string]interface{}{
		"steps": steps,
	}
	return s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/steps/reorder", sequenceID), body, nil)
}

// Enroll manually enrolls subscribers into a sequence.
func (s *SequencesService) Enroll(ctx context.Context, sequenceID string, subscriberIDs []string) (*EnrollResult, error) {
	body := map[string]interface{}{
		"subscriberIds": subscriberIDs,
	}
	var resp EnrollResult
	err := s.client.post(ctx, fmt.Sprintf("/v1/sequences/%s/enroll", sequenceID), body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListEnrollments returns a paginated list of enrollments.
func (s *SequencesService) ListEnrollments(ctx context.Context, sequenceID string, params *ListParams) (*ListResponse[SequenceEnrollment], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[SequenceEnrollment]
	err := s.client.get(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments", sequenceID), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveEnrollment removes an enrollment from a sequence.
func (s *SequencesService) RemoveEnrollment(ctx context.Context, sequenceID, enrollmentID string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/sequences/%s/enrollments/%s", sequenceID, enrollmentID))
}

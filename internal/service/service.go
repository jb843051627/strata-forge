package service

import (
	"context"
	"sync"

	"github.com/jb843051627/strata-forge/internal/clock"
	"github.com/jb843051627/strata-forge/internal/engine"
	"github.com/jb843051627/strata-forge/internal/store"
)

type LabService struct {
	store    *store.Store
	workflow *engine.Workflow
	quality  *engine.QualityEngine
	age      *engine.AgeEstimator
	clock    clock.Clock

	queueMu       sync.Mutex
	queue         chan int64
	workers       sync.WaitGroup
	closed        bool
	runCreationMu sync.Mutex
	runStateMu    sync.Mutex
	pendingRuns   map[int64]bool

	reviewAuditMu      sync.Mutex
	runAuditMu         sync.Mutex
	measurementAuditMu sync.Mutex
	sampleAuditMu      sync.Mutex
	alertAuditMu       sync.Mutex
	reportAuditMu      sync.Mutex
	reviewAudit        *engine.ReviewAudit
	runAudit           *engine.RunAudit
	measurementAudit   *engine.MeasurementAudit
	sampleAudit        *engine.SampleAudit
	alertAudit         *engine.AlertAudit
	reportAudit        *engine.ReportAudit
	instruments        *engine.InstrumentCatalog
}

func New(st *store.Store, c clock.Clock) *LabService {
	return &LabService{
		store:            st,
		workflow:         engine.NewWorkflow(),
		quality:          engine.NewQualityEngine(),
		age:              engine.NewAgeEstimator(),
		clock:            c,
		queue:            make(chan int64, 32),
		pendingRuns:      make(map[int64]bool),
		reviewAudit:      engine.NewReviewAudit(),
		runAudit:         engine.NewRunAudit(),
		measurementAudit: engine.NewMeasurementAudit(),
		sampleAudit:      engine.NewSampleAudit(),
		alertAudit:       engine.NewAlertAudit(),
		reportAudit:      engine.NewReportAudit(),
		instruments:      engine.NewInstrumentCatalog(),
	}
}

func (s *LabService) Close() {
	s.queueMu.Lock()
	s.closed = true
	s.queueMu.Unlock()
	s.workers.Wait()
}

func (s *LabService) Store() *store.Store {
	return s.store
}

func (s *LabService) Quality() *engine.QualityEngine {
	return s.quality
}

func (s *LabService) Clock() clock.Clock {
	return s.clock
}

func (s *LabService) contextError(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

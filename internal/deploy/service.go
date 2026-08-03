package deploy

import "context"

type Service struct {
	engine *Engine
}

func NewService(engine *Engine) *Service {
	return &Service{
		engine: engine,
	}
}

func (s *Service) StartProject(ctx context.Context, p Project) error {
	err := s.engine.Start(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) StopProject(ctx context.Context, p Project) error {
	err := s.engine.Stop(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RedeployProject(ctx context.Context, req RedeployRequest) error {
	err := s.engine.Redeploy(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RollbackProject(ctx context.Context, req RollbackRequest) error {
	err := s.engine.Rollback(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

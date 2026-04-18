package deploy

type RedeployRequest struct {
	Project          Project
	TargetRevision   string
	PreviousRevision string
}

type RollbackRequest struct {
	Project  Project
	Revision string
}

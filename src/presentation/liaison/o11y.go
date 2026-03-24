package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/goinfinite/os/src/domain/useCase"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
)

type O11yLiaison struct {
	transientDbSvc *internalDbInfra.TransientDatabaseService
}

func NewO11yLiaison(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *O11yLiaison {
	return &O11yLiaison{
		transientDbSvc: transientDbSvc,
	}
}

func (liaison *O11yLiaison) ReadOverview() tkPresentation.LiaisonResponse {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(liaison.transientDbSvc)

	o11yOverview, err := useCase.ReadO11yOverview(o11yQueryRepo, true)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, o11yOverview,
	)
}

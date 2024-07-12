package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func DeleteCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	deleteDto dto.DeleteCron,
) error {
	if deleteDto.Id == nil && deleteDto.Comment == nil {
		return errors.New("CronIdOrCommentRequired")
	}

	if deleteDto.Id != nil {
		err := cronCmdRepo.Delete(*deleteDto.Id)
		if err != nil {
			log.Printf("DeleteCronError: %s", err)
			return errors.New("DeleteCronInfraError")
		}
		return nil
	}

	err := cronCmdRepo.DeleteByComment(*deleteDto.Comment)
	if err != nil {
		log.Printf("DeleteCronByCommentError: %s", err)
		return errors.New("DeleteCronByCommentInfraError")
	}

	return nil
}

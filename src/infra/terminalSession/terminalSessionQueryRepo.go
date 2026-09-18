package terminalSessionInfra

import (
	"cmp"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
)

type TerminalSessionQueryRepo struct {
	accountQueryRepo repository.AccountQueryRepo
}

func NewTerminalSessionQueryRepo(
	accountQueryRepo repository.AccountQueryRepo,
) *TerminalSessionQueryRepo {
	return &TerminalSessionQueryRepo{accountQueryRepo: accountQueryRepo}
}

func (repo *TerminalSessionQueryRepo) resolveTargetAccounts(
	accountIdPtr *tkValueObject.AccountId,
) (accounts []entity.Account, err error) {
	if accountIdPtr != nil {
		accountEntity, err := repo.accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
			AccountId: accountIdPtr,
		})
		if err != nil {
			return accounts, err
		}

		return []entity.Account{accountEntity}, nil
	}

	accountsResponse, err := repo.accountQueryRepo.Read(dto.ReadAccountsRequest{
		Pagination: tkDto.PaginationUnpaginated,
	})
	if err != nil {
		return accounts, err
	}

	return accountsResponse.Accounts, nil
}

func (repo *TerminalSessionQueryRepo) readAccountTerminalSessions(
	accountEntity entity.Account,
	terminalSessionIdPtr *valueObject.TerminalSessionId,
) (terminalSessions []entity.TerminalSession, err error) {
	client, err := NewTerminalMultiplexerClient(accountEntity.Username)
	if err != nil {
		return terminalSessions, err
	}

	multiplexerSessions, err := client.ListSessions()
	if err != nil {
		return terminalSessions, err
	}

	terminalSessions = []entity.TerminalSession{}
	for _, multiplexerSession := range multiplexerSessions {
		if terminalSessionIdPtr != nil && multiplexerSession.Id != *terminalSessionIdPtr {
			continue
		}

		terminalSessions = append(terminalSessions, entity.NewTerminalSession(
			multiplexerSession.Id, accountEntity.Id, accountEntity.Username,
			multiplexerSession.WorkingDir, multiplexerSession.Command,
			multiplexerSession.CreatedAt, multiplexerSession.AttachedClients,
		))
	}

	return terminalSessions, nil
}

func (repo *TerminalSessionQueryRepo) sortTerminalSessions(
	terminalSessions []entity.TerminalSession,
	sortByPtr *tkValueObject.PaginationSortBy,
	sortDirectionPtr *tkValueObject.PaginationSortDirection,
) []entity.TerminalSession {
	sortByStr := "createdAt"
	if sortByPtr != nil {
		sortByStr = sortByPtr.String()
	}

	sortDirectionStr := "asc"
	if sortDirectionPtr != nil {
		sortDirectionStr = sortDirectionPtr.String()
	}

	slices.SortStableFunc(terminalSessions, func(
		first, second entity.TerminalSession,
	) int {
		if sortDirectionStr != "asc" {
			first, second = second, first
		}

		switch sortByStr {
		case "id", "terminalSessionId":
			return strings.Compare(first.Id.String(), second.Id.String())
		case "accountUsername":
			return strings.Compare(
				first.AccountUsername.String(), second.AccountUsername.String(),
			)
		case "createdAt":
			return cmp.Compare(first.CreatedAt, second.CreatedAt)
		default:
			return 0
		}
	})

	return terminalSessions
}

func (repo *TerminalSessionQueryRepo) paginateTerminalSessions(
	terminalSessions []entity.TerminalSession,
	paginationDto tkDto.Pagination,
) (responseDto dto.ReadTerminalSessionsResponse, err error) {
	paginationStartIndex := int(paginationDto.PageNumber) *
		int(paginationDto.ItemsPerPage)
	paginationEndIndex := paginationStartIndex + int(paginationDto.ItemsPerPage)

	paginatedTerminalSessions := []entity.TerminalSession{}
	if paginationStartIndex < len(terminalSessions) {
		pageEndIndex := min(paginationEndIndex, len(terminalSessions))
		paginatedTerminalSessions = terminalSessions[paginationStartIndex:pageEndIndex]
	}

	itemsTotal := uint64(len(terminalSessions))
	pagesTotal, err := tkInfraDb.PaginationPagesTotalResolver(
		itemsTotal, paginationDto.ItemsPerPage,
	)
	if err != nil {
		return responseDto, err
	}

	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadTerminalSessionsResponse{
		Pagination:       paginationDto,
		TerminalSessions: paginatedTerminalSessions,
	}, nil
}

func (repo *TerminalSessionQueryRepo) Read(
	requestDto dto.ReadTerminalSessionsRequest,
) (responseDto dto.ReadTerminalSessionsResponse, err error) {
	targetAccounts, err := repo.resolveTargetAccounts(requestDto.AccountId)
	if err != nil {
		return responseDto, err
	}

	terminalSessions := []entity.TerminalSession{}
	readableAccountsCount := 0

	for _, accountEntity := range targetAccounts {
		accountTerminalSessions, err := repo.readAccountTerminalSessions(
			accountEntity, requestDto.TerminalSessionId,
		)
		if err != nil {
			slog.Error(
				"ReadAccountTerminalSessionsError",
				slog.String("accountUsername", accountEntity.Username.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		readableAccountsCount++
		terminalSessions = append(terminalSessions, accountTerminalSessions...)
	}

	if len(targetAccounts) > 0 && readableAccountsCount == 0 {
		return responseDto, errors.New("ReadTargetAccountsTerminalSessionsError")
	}

	sortedTerminalSessions := repo.sortTerminalSessions(
		terminalSessions,
		requestDto.Pagination.SortBy,
		requestDto.Pagination.SortDirection,
	)

	return repo.paginateTerminalSessions(sortedTerminalSessions, requestDto.Pagination)
}

func (repo *TerminalSessionQueryRepo) ReadFirst(
	requestDto dto.ReadTerminalSessionsRequest,
) (terminalSessionEntity entity.TerminalSession, err error) {
	requestDto.Pagination = tkDto.PaginationSingleItem
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return terminalSessionEntity, err
	}
	if len(responseDto.TerminalSessions) == 0 {
		return terminalSessionEntity, repository.ErrTerminalSessionNotFound
	}

	return responseDto.TerminalSessions[0], nil
}

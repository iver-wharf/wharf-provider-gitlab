package main

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type gitLabPaging struct {
	totalItems   int
	totalPages   int
	itemsPerPage int
	currentPage  int
	nextPage     int
	previousPage int
}

func mapToPaging(resp *gitlab.Response) gitLabPaging {
	return gitLabPaging{
		totalItems:   resp.TotalItems,
		totalPages:   resp.TotalPages,
		itemsPerPage: resp.ItemsPerPage,
		currentPage:  resp.CurrentPage,
		nextPage:     resp.NextPage,
		previousPage: resp.PreviousPage,
	}
}

func (p gitLabPaging) next() int {
	if p.currentPage >= p.totalPages {
		log.Debug().WithInt("page", p.totalPages).Message("Found end of collection.")
		return -1
	}

	log.Debug().
		WithInt("currentPage", p.currentPage).
		WithInt("nextPage", p.nextPage).
		WithInt("totalPages", p.totalPages).
		Message("Fetching next page.")
	return p.nextPage
}

type getProjects func(int) ([]*gitlab.Project, gitLabPaging, error)
type putProjects func([]*gitlab.Project) string

func importPaginatedProjects(get getProjects, put putProjects) error {
	errMessage := ""
	page := 0
	for page >= 0 {
		projects, paging, err := get(page)
		if err != nil {
			log.Error().WithError(err).Message("Failed to get projects.")
			return err
		}

		if len(projects) > 0 {
			errMessage += put(projects)
		}

		page = paging.next()
	}

	if errMessage != "" {
		return fmt.Errorf(errMessage)
	}

	return nil
}

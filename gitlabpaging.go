package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
		log.WithField("page", p.totalPages).Debugln("found end of collection")
		return -1
	}

	log.WithFields(log.Fields{
		"current page": p.currentPage,
		"next page":    p.nextPage,
		"total pages":  p.totalPages}).
		Debugln("fetching next page")
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
			log.WithError(err).Errorln("failed to get projects")
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

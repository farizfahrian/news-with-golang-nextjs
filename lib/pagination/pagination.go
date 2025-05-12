package pagination

import (
	"math"
	"news-with-golang/internal/core/domain/entity"
)

type PaginationInterface interface {
	ApplyPagination(totalData, page, perPage int) (*entity.Page, error)
}

type Options struct{}

// ApplyPagination implements PaginationInterface.
func (o *Options) ApplyPagination(totalData int, page int, perPage int) (*entity.Page, error) {
	newPage := page
	if newPage <= 0 {
		return nil, ErrorPage
	}

	limitData := 10
	if perPage > 0 {
		limitData = perPage
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limitData)))

	lastPage := (newPage * limitData)
	firstPage := lastPage - limitData
	if totalData < lastPage {
		lastPage = totalData
	}

	zeroPage := &entity.Page{PageCount: 1, Page: newPage}
	if totalData == 0 && newPage == 1 {
		return zeroPage, nil
	}

	if newPage > totalPage {
		return nil, ErrorMaxPage
	}

	if newPage < 1 {
		return nil, ErrorPage
	}

	return &entity.Page{
		Page:       newPage,
		PerPage:    limitData,
		PageCount:  totalPage,
		TotalCount: totalData,
		FirstPage:  firstPage,
		LastPage:   lastPage,
	}, nil
}

func NewPagination() PaginationInterface {
	pagination := new(Options)

	return pagination
}

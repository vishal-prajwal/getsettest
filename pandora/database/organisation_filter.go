package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func InjectOrganisationFilter(filter *OrganisationContentFilter) func(*bun.SelectQuery) *bun.SelectQuery {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if filter == nil {
			return q
		}

		tm, ok := q.Table().GetModel().(bun.TableModel)
		if !ok {
			panic("cannot access cast %v to bun.TableModel")
		}

		// Use the alias of the model
		alias := tm.Table().Alias

		// Not all cards actually support filtering on organisation groups.
		// So we use this to check that the model actually supports it.
		supportsOrganisationGroups := false

		for _, field := range tm.Table().DataFields {
			if field.Name == "organisation_groups" {
				supportsOrganisationGroups = true
			}
		}

		switch {
		case len(filter.Organisations) > 0 && filter.IncludeOriginalContent:
			q.Where(fmt.Sprintf("%s.organisation_id IN (?) OR %s.organisation_id IS NULL", alias, alias), bun.In(filter.Organisations))

		case len(filter.Organisations) > 0 && !filter.IncludeOriginalContent:
			q.Where(fmt.Sprintf("%s.organisation_id IN (?)", alias), bun.In(filter.Organisations))

		case filter.HideOrganisationContent:
			q.Where(fmt.Sprintf("%s.organisation_id IS NULL", alias))
		}

		if len(filter.OrganisationGroups) > 0 && supportsOrganisationGroups {
			if filter.IncludeOriginalContent {
				q.Where(
					fmt.Sprintf("%s.organisation_groups && ? OR (%s.organisation_groups IS NULL OR array_length(%s.organisation_groups, 1) IS NULL)", alias, alias, alias),
					pgdialect.Array(filter.OrganisationGroups),
				)
			} else {
				q.WhereGroup(" AND ", func(query *bun.SelectQuery) *bun.SelectQuery {
					query.WhereOr(fmt.Sprintf("%s.organisation_groups && ?", alias), pgdialect.Array(filter.OrganisationGroups))

					if !filter.EnsureContentNotAssociateToAGroup {
						query.WhereOr(
							fmt.Sprintf("%s.organisation_groups IS NULL OR array_length(%s.organisation_groups, 1) IS NULL", alias, alias),
						)
					}

					return query
				})
			}
		}

		if filter.EnsureContentNotAssociateToAGroup && supportsOrganisationGroups {
			q.Where(fmt.Sprintf("%s.organisation_groups IS NULL OR array_length(%s.organisation_groups, 1) IS NULL", alias, alias))
		}

		return q
	}
}

type OrganisationContentFilter struct {
	// Organisations which organisations we want to filter on
	Organisations []uuid.UUID

	// OrganisationGroups which groups the content needs to belong too to view the content
	OrganisationGroups []uuid.UUID

	// IncludeOriginalContent ensures that we also allow original content to be displayed
	IncludeOriginalContent bool

	// HideOrganisationContent ensures that no content from any organisation is visible
	HideOrganisationContent bool

	// EnsureContentNotAssociateToAGroup ensures that any content returned is not associated to a group
	EnsureContentNotAssociateToAGroup bool
}

package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/lib/pq"
)

var (
	orgMembersQuery = regexp.QuoteMeta("SELECT" +
		" zitadel.projections.org_members.creation_date" +
		", zitadel.projections.org_members.change_date" +
		", zitadel.projections.org_members.sequence" +
		", zitadel.projections.org_members.resource_owner" +
		", zitadel.projections.org_members.user_id" +
		", zitadel.projections.org_members.roles" +
		", zitadel.projections.login_names.login_name" +
		", zitadel.projections.users_humans.email" +
		", zitadel.projections.users_humans.first_name" +
		", zitadel.projections.users_humans.last_name" +
		", zitadel.projections.users_humans.display_name" +
		", zitadel.projections.users_machines.name" +
		", zitadel.projections.users_humans.avater_key" +
		", COUNT(*) OVER () " +
		"FROM zitadel.projections.org_members " +
		"LEFT JOIN zitadel.projections.users_humans " +
		"ON zitadel.projections.org_members.user_id = zitadel.projections.users_humans.user_id " +
		"LEFT JOIN zitadel.projections.users_machines " +
		"ON zitadel.projections.org_members.user_id = zitadel.projections.users_machines.user_id " +
		"LEFT JOIN zitadel.projections.login_names " +
		"ON zitadel.projections.org_members.user_id = zitadel.projections.login_names.user_id " +
		"WHERE zitadel.projections.login_names.is_primary = $1")
	orgMembersColumns = []string{
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"user_id",
		"roles",
		"login_name",
		"email",
		"first_name",
		"last_name",
		"display_name",
		"name",
		"avater_key",
		"count",
	}
)

func Test_OrgMemberPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareOrgMembersQuery no result",
			prepare: prepareOrgMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					orgMembersQuery,
					nil,
					nil,
				),
			},
			object: &Members{
				Members: []*Member{},
			},
		},
		{
			name:    "prepareOrgMembersQuery human found",
			prepare: prepareOrgMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					orgMembersQuery,
					orgMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							pq.StringArray{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						Roles:              []string{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareOrgMembersQuery machine found",
			prepare: prepareOrgMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					orgMembersQuery,
					orgMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							pq.StringArray{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						Roles:              []string{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareOrgMembersQuery multiple users",
			prepare: prepareOrgMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					orgMembersQuery,
					orgMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-1",
							pq.StringArray{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
						},
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-2",
							pq.StringArray{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-1",
						Roles:              []string{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
					},
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-2",
						Roles:              []string{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareOrgMembersQuery sql err",
			prepare: prepareOrgMembersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					orgMembersQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
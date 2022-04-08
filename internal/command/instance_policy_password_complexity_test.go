package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddDefaultPasswordComplexityPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordComplexityPolicy
	}
	type res struct {
		want *domain.PasswordComplexityPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid policy, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    0,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "password complexity policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								8,
								true, true, true, true,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    8,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewPasswordComplexityPolicyAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									8,
									true, true, true, true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    8,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				want: &domain.PasswordComplexityPolicy{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					MinLength:    8,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultPasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeDefaultPasswordComplexityPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordComplexityPolicy
	}
	type res struct {
		want *domain.PasswordComplexityPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid policy, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    0,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "password complexity policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    8,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								8,
								true, true, true, true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    8,
					HasUppercase: true,
					HasLowercase: true,
					HasNumber:    true,
					HasSymbol:    true,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								8,
								true, true, true, true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultPasswordComplexityPolicyChangedEvent(context.Background(), 10, false, false, false, false),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordComplexityPolicy{
					MinLength:    10,
					HasUppercase: false,
					HasLowercase: false,
					HasNumber:    false,
					HasSymbol:    false,
				},
			},
			res: res{
				want: &domain.PasswordComplexityPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					MinLength:    10,
					HasUppercase: false,
					HasLowercase: false,
					HasNumber:    false,
					HasSymbol:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultPasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func newDefaultPasswordComplexityPolicyChangedEvent(ctx context.Context, minLength uint64, hasUpper, hasLower, hasNumber, hasSymbol bool) *instance.PasswordComplexityPolicyChangedEvent {
	event, _ := instance.NewPasswordComplexityPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.PasswordComplexityPolicyChanges{
			policy.ChangeMinLength(minLength),
			policy.ChangeHasUppercase(hasUpper),
			policy.ChangeHasLowercase(hasLower),
			policy.ChangeHasSymbol(hasNumber),
			policy.ChangeHasNumber(hasSymbol),
		},
	)
	return event
}
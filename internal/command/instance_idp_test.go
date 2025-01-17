package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOAuthProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"clientID",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("clientSecret"),
									},
									"auth",
									"token",
									"user",
									nil,
									idp.Options{},
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"clientID",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("clientSecret"),
									},
									"auth",
									"token",
									"user",
									[]string{"user"},
									idp.Options{
										IsCreationAllowed: true,
										IsLinkingAllowed:  true,
										IsAutoCreation:    true,
										IsAutoUpdate:      true,
									},
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGenericOAuthProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOAuthProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewOAuthIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
										"id1",
										[]idp.OAuthIDPChanges{
											idp.ChangeOAuthName("new name"),
											idp.ChangeOAuthClientID("clientID2"),
											idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
												CryptoType: crypto.TypeEncryption,
												Algorithm:  "enc",
												KeyID:      "id",
												Crypted:    []byte("newSecret"),
											}),
											idp.ChangeOAuthAuthorizationEndpoint("new auth"),
											idp.ChangeOAuthTokenEndpoint("new token"),
											idp.ChangeOAuthUserEndpoint("new user"),
											idp.ChangeOAuthScopes([]string{"openid", "profile"}),
											idp.ChangeOAuthOptions(idp.OptionChanges{
												IsCreationAllowed: &t,
												IsLinkingAllowed:  &t,
												IsAutoCreation:    &t,
												IsAutoUpdate:      &t,
											}),
										},
									)
									return event
								}(),
							),
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOAuthProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GoogleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"",
									"clientID",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("clientSecret"),
									},
									nil,
									idp.Options{},
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"",
									"clientID",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("clientSecret"),
									},
									[]string{"openid"},
									idp.Options{
										IsCreationAllowed: true,
										IsLinkingAllowed:  true,
										IsAutoCreation:    true,
										IsAutoUpdate:      true,
									},
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGoogleProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGoogleIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
										"id1",
										[]idp.GoogleIDPChanges{
											idp.ChangeGoogleClientID("clientID2"),
											idp.ChangeGoogleClientSecret(&crypto.CryptoValue{
												CryptoType: crypto.TypeEncryption,
												Algorithm:  "enc",
												KeyID:      "id",
												Crypted:    []byte("newSecret"),
											}),
											idp.ChangeGoogleScopes([]string{"openid", "profile"}),
											idp.ChangeGoogleOptions(idp.OptionChanges{
												IsCreationAllowed: &t,
												IsLinkingAllowed:  &t,
												IsAutoCreation:    &t,
												IsAutoUpdate:      &t,
											}),
										},
									)
									return event
								}(),
							),
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGoogleProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider LDAPProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid host",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid admin",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid password",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"host",
									"",
									false,
									"baseDN",
									"userObjectClass",
									"userUniqueAttribute",
									"admin",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("password"),
									},
									idp.LDAPAttributes{},
									idp.Options{},
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
					Password:            "password",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"host",
									"port",
									true,
									"baseDN",
									"userObjectClass",
									"userUniqueAttribute",
									"admin",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("password"),
									},
									idp.LDAPAttributes{
										IDAttribute:                "id",
										FirstNameAttribute:         "firstName",
										LastNameAttribute:          "lastName",
										DisplayNameAttribute:       "displayName",
										NickNameAttribute:          "nickName",
										PreferredUsernameAttribute: "preferredUsername",
										EmailAttribute:             "email",
										EmailVerifiedAttribute:     "emailVerified",
										PhoneAttribute:             "phone",
										PhoneVerifiedAttribute:     "phoneVerified",
										PreferredLanguageAttribute: "preferredLanguage",
										AvatarURLAttribute:         "avatarURL",
										ProfileAttribute:           "profile",
									},
									idp.Options{
										IsCreationAllowed: true,
										IsLinkingAllowed:  true,
										IsAutoCreation:    true,
										IsAutoUpdate:      true,
									},
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					Port:                "port",
					TLS:                 true,
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
					Password:            "password",
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "id",
						FirstNameAttribute:         "firstName",
						LastNameAttribute:          "lastName",
						DisplayNameAttribute:       "displayName",
						NickNameAttribute:          "nickName",
						PreferredUsernameAttribute: "preferredUsername",
						EmailAttribute:             "email",
						EmailVerifiedAttribute:     "emailVerified",
						PhoneAttribute:             "phone",
						PhoneVerifiedAttribute:     "phoneVerified",
						PreferredLanguageAttribute: "preferredLanguage",
						AvatarURLAttribute:         "avatarURL",
						ProfileAttribute:           "profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceLDAPProvider(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider LDAPProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: LDAPProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid host",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid admin",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"host",
								"",
								false,
								"baseDN",
								"userObjectClass",
								"userUniqueAttribute",
								"admin",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"host",
								"port",
								false,
								"baseDN",
								"userObjectClass",
								"userUniqueAttribute",
								"admin",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewLDAPIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
										"id1",
										"name",
										[]idp.LDAPIDPChanges{
											idp.ChangeLDAPName("new name"),
											idp.ChangeLDAPHost("new host"),
											idp.ChangeLDAPPort("new port"),
											idp.ChangeLDAPTLS(true),
											idp.ChangeLDAPBaseDN("new baseDN"),
											idp.ChangeLDAPUserObjectClass("new userObjectClass"),
											idp.ChangeLDAPUserUniqueAttribute("new userUniqueAttribute"),
											idp.ChangeLDAPAdmin("new admin"),
											idp.ChangeLDAPPassword(&crypto.CryptoValue{
												CryptoType: crypto.TypeEncryption,
												Algorithm:  "enc",
												KeyID:      "id",
												Crypted:    []byte("new password"),
											}),
											idp.ChangeLDAPAttributes(idp.LDAPAttributeChanges{
												IDAttribute:                stringPointer("new id"),
												FirstNameAttribute:         stringPointer("new firstName"),
												LastNameAttribute:          stringPointer("new lastName"),
												DisplayNameAttribute:       stringPointer("new displayName"),
												NickNameAttribute:          stringPointer("new nickName"),
												PreferredUsernameAttribute: stringPointer("new preferredUsername"),
												EmailAttribute:             stringPointer("new email"),
												EmailVerifiedAttribute:     stringPointer("new emailVerified"),
												PhoneAttribute:             stringPointer("new phone"),
												PhoneVerifiedAttribute:     stringPointer("new phoneVerified"),
												PreferredLanguageAttribute: stringPointer("new preferredLanguage"),
												AvatarURLAttribute:         stringPointer("new avatarURL"),
												ProfileAttribute:           stringPointer("new profile"),
											}),
											idp.ChangeLDAPOptions(idp.OptionChanges{
												IsCreationAllowed: &t,
												IsLinkingAllowed:  &t,
												IsAutoCreation:    &t,
												IsAutoUpdate:      &t,
											}),
										},
									)
									return event
								}(),
							),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "new name",
					Host:                "new host",
					Port:                "new port",
					TLS:                 true,
					BaseDN:              "new baseDN",
					UserObjectClass:     "new userObjectClass",
					UserUniqueAttribute: "new userUniqueAttribute",
					Admin:               "new admin",
					Password:            "new password",
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "new id",
						FirstNameAttribute:         "new firstName",
						LastNameAttribute:          "new lastName",
						DisplayNameAttribute:       "new displayName",
						NickNameAttribute:          "new nickName",
						PreferredUsernameAttribute: "new preferredUsername",
						EmailAttribute:             "new email",
						EmailVerifiedAttribute:     "new emailVerified",
						PhoneAttribute:             "new phone",
						PhoneVerifiedAttribute:     "new phoneVerified",
						PreferredLanguageAttribute: "new preferredLanguage",
						AvatarURLAttribute:         "new avatarURL",
						ProfileAttribute:           "new profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceLDAPProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

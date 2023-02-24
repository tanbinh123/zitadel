package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgOIDCIDPWriteModel struct {
	OIDCIDPWriteModel
}

func NewOIDCOrgIDPWriteModel(orgID, id string) *OrgOIDCIDPWriteModel {
	return &OrgOIDCIDPWriteModel{
		OIDCIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgOIDCIDPWriteModel) Reduce() error {
	return wm.OIDCIDPWriteModel.Reduce()
}

func (wm *OrgOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OIDCIDPAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *org.OIDCIDPChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *org.IDPRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *org.IDPConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *org.IDPConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *org.IDPOIDCConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *org.IDPOIDCConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *org.IDPConfigRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgOIDCIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OIDCIDPAddedEventType,
			org.OIDCIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.IDPConfigAddedEventType,
			org.IDPConfigChangedEventType,
			org.IDPOIDCConfigAddedEventType,
			org.IDPOIDCConfigChangedEventType,
			org.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *OrgOIDCIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.OIDCIDPChangedEvent, error) {

	changes, err := wm.OIDCIDPWriteModel.NewChanges(
		name,
		issuer,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewOIDCIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgJWTIDPWriteModel struct {
	JWTIDPWriteModel
}

func NewJWTOrgIDPWriteModel(orgID, id string) *OrgJWTIDPWriteModel {
	return &OrgJWTIDPWriteModel{
		JWTIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgJWTIDPWriteModel) Reduce() error {
	return wm.JWTIDPWriteModel.Reduce()
}

func (wm *OrgJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.JWTIDPAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *org.JWTIDPChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *org.IDPRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *org.IDPConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *org.IDPConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *org.IDPJWTConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *org.IDPJWTConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *org.IDPConfigRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgJWTIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.JWTIDPAddedEventType,
			org.JWTIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.IDPConfigAddedEventType,
			org.IDPConfigChangedEventType,
			org.IDPJWTConfigAddedEventType,
			org.IDPJWTConfigChangedEventType,
			org.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *OrgJWTIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) (*org.JWTIDPChangedEvent, error) {

	changes, err := wm.JWTIDPWriteModel.NewChanges(
		name,
		issuer,
		jwtEndpoint,
		keysEndpoint,
		headerName,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewJWTIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgGoogleIDPWriteModel struct {
	GoogleIDPWriteModel
}

func NewGoogleOrgIDPWriteModel(orgID, id string) *OrgGoogleIDPWriteModel {
	return &OrgGoogleIDPWriteModel{
		GoogleIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgGoogleIDPWriteModel) Reduce() error {
	return wm.GoogleIDPWriteModel.Reduce()
}

func (wm *OrgGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.GoogleIDPAddedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *org.GoogleIDPChangedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *org.IDPRemovedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.GoogleIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgGoogleIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GoogleIDPAddedEventType,
			org.GoogleIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *OrgGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*org.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgLDAPIDPWriteModel struct {
	LDAPIDPWriteModel
}

func NewLDAPOrgIDPWriteModel(orgID, id string) *OrgLDAPIDPWriteModel {
	return &OrgLDAPIDPWriteModel{
		LDAPIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgLDAPIDPWriteModel) Reduce() error {
	return wm.LDAPIDPWriteModel.Reduce()
}

func (wm *OrgLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.LDAPIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgLDAPIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgLDAPIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password string,
	secretCrypto crypto.Crypto,
	attributes idp.LDAPAttributes,
	options idp.Options,
) (*org.LDAPIDPChangedEvent, error) {

	changes, err := wm.LDAPIDPWriteModel.NewChanges(
		name,
		host,
		port,
		tls,
		baseDN,
		userObjectClass,
		userUniqueAttribute,
		admin,
		password,
		secretCrypto,
		attributes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := org.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgIDPRemoveWriteModel struct {
	IDPRemoveWriteModel
}

func NewOrgIDPRemoveWriteModel(orgID, id string) *OrgIDPRemoveWriteModel {
	return &OrgIDPRemoveWriteModel{
		IDPRemoveWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgIDPRemoveWriteModel) Reduce() error {
	return wm.IDPRemoveWriteModel.Reduce()
}

func (wm *OrgIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.JWTIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.GoogleIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *instance.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.IDPRemoveWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgIDPRemoveWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.GoogleIDPAddedEventType,
			org.GoogleIDPChangedEventType,
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

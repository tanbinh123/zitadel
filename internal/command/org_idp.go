package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func (c *Commands) AddOrgGenericOAuthProvider(ctx context.Context, resourceOwner string, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGenericOAuthProvider(ctx context.Context, resourceOwner, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGoogleProvider(ctx context.Context, resourceOwner string, provider GoogleProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGoogleProvider(ctx context.Context, resourceOwner, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgLDAPProvider(ctx context.Context, resourceOwner string, provider LDAPProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgLDAPProvider(ctx context.Context, resourceOwner, id string, provider LDAPProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) DeleteOrgProvider(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareDeleteOrgProvider(orgAgg, resourceOwner, id))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) prepareAddOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-DF4ga", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOAuthIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-asfsa", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-JNsd3", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-W2vqs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGoogleIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-S32t1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-ds432", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Dqrg1", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SAfdd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SDVg2", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sv31s", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdgf4", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-AEG2w", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SAD5n", "Errors.Invalid.Argument")
		}
		if provider.Password = strings.TrimSpace(provider.Password); provider.Password == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdf5h", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.Password), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewLDAPIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Host,
					provider.Port,
					provider.TLS,
					provider.BaseDN,
					provider.UserObjectClass,
					provider.UserUniqueAttribute,
					provider.Admin,
					secret,
					provider.LDAPAttributes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dgdbs", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Sffgd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dz62d", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-vb3ss", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-hbere", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-ASFt6", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-DG45z", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-ASF3F", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				writeModel.Name,
				provider.Name,
				provider.Host,
				provider.Port,
				provider.TLS,
				provider.BaseDN,
				provider.UserObjectClass,
				provider.UserUniqueAttribute,
				provider.Admin,
				provider.Password,
				c.idpConfigEncryption,
				provider.LDAPAttributes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareDeleteOrgProvider(a *org.Aggregate, resourceOwner, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgIDPRemoveWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Se3tg", "Errors.Org.IDPConfig.NotExisting")
			}
			return []eventstore.Command{org.NewIDPRemovedEvent(ctx, &a.Aggregate, id, writeModel.name)}, nil
		}, nil
	}
}

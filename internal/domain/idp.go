package domain

type IDPState int32

const (
	IDPStateUnspecified IDPState = iota
	IDPStateActive
	IDPStateInactive
	IDPStateRemoved

	idpStateCount
)

func (s IDPState) Valid() bool {
	return s >= 0 && s < idpStateCount
}

func (s IDPState) Exists() bool {
	return s != IDPStateUnspecified && s != IDPStateRemoved
}

type IDPType int32

const (
	IDPTypeUnspecified IDPType = iota
	IDPTypeOIDC
	IDPTypeJWT
	IDPTypeOAuth
	IDPTypeLDAP
	IDPTypeAzureAD
	IDPTypeGitHub
	IDPTypeGitHubEE
	IDPTypeGitLab
	IDPTypeGitLabSelfHosted
	IDPTypeGoogle
)

package types

type GoogleIDTokenClaims struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Avatar     string `json:"picture"`
}

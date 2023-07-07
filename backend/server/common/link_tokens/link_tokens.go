package link_tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/errors"
)

type TokenInfo struct {
	OrganizationID int64   `json:"organization_id"`
	EndCustomerID  string  `json:"end_customer_id"`
	DestinationIDs []int64 `json:"destination_ids"`
}

func (t TokenInfo) HasDestination(destinationID int64) bool {
	for _, id := range t.DestinationIDs {
		if id == destinationID {
			return true
		}
	}

	return false
}

type LinkTokenClaims struct {
	TokenInfo `json:"token_info"`
	jwt.RegisteredClaims
}

func CreateLinkToken(tokenInfo TokenInfo) (*string, error) {
	rawToken := jwt.NewWithClaims(crypto.SigningMethodKMSHS256, LinkTokenClaims{
		tokenInfo,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour * 24)),
		},
	})

	signedToken, err := rawToken.SignedString(nil)
	if err != nil {
		return nil, errors.Wrap(err, "(link_tokens.CreateLinkToken) signing token")
	}

	return &signedToken, nil
}

func ValidateLinkToken(linkTokenStr string) (*TokenInfo, error) {
	token, err := jwt.ParseWithClaims(linkTokenStr, &LinkTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil // no key needs to be fetched— we just call the GCP KMS endpoint
	})

	if err != nil {
		return nil, errors.Wrap(err, "(link_tokens.ValidateLinkToken) parsing token")
	}

	if !token.Valid {
		return nil, errors.Wrapf(jwt.ErrTokenInvalidClaims, "(link_tokens.ValidateLinkToken) token invalid: %v", token.Raw)
	}

	claims, ok := token.Claims.(*LinkTokenClaims)
	if !ok {
		return nil, errors.Wrapf(jwt.ErrTokenInvalidClaims, "(link_tokens.ValidateLinkToken) invalid claims: %v", token.Raw)
	}

	return &claims.TokenInfo, nil
}

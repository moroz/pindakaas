package config

import "time"

const UserTokenLength = 32

const AccessTokenSessionKey = "access_token"
const AccessTokenValidity = 24 * time.Hour

const SessionCookieName = "_pindakaas_session"

const OIDCStateSessionKey = "oauth_state"

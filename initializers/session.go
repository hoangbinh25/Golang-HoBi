package initializers

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("mylongsecretkeywith32byteslength!!!"))

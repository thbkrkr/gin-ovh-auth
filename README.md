## OVH Authentication

Secure your Go [gin](https://github.com/gin-gonic/gin#gin-web-framework) API with the OVH Authentication.

### Getting started

#### Environment variables

```sh
# OVH Application credentials
# https://api.ovh.com/createApp/
OVH_ENDPOINT=ovh-eu
OVH_APPLICATION_KEY=??
OVH_APPLICATION_SECRET=??

# Login URL of the page with the login button
AUTH_LOGIN_URL=https://??/s/login.html
# Secret to sign the auth data
AUTH_SECRET=??
```

#### API

```go
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/thbkrkr/ovh-auth"
)

func main() {
  router := gin.Default()

  authRouter := ovhauth.Secure(r)

  // Public endpoint
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "ping"})
  })

  // Secured endpoint
  authRouter.GET("/api/me", func(c *gin.Context) {
    me := ovhauth.GetAuthUser(c)
    c.JSON(200, me)
  })

  router.Run(":8080")
}

```

#### login.html ($AUTH_LOGIN_URL)

Add a 'clickable' element with a `login` id in the `$AUTH_LOGIN_URL` page.

```html
<button id="login">Login</button>
<script src="//thbkrkr.github.io/ovh-auth/js/auth.1.min.js"></script>
```

#### index.html ($AUTH_HOME_URL)

Call your API with `$get` and `$post`.

```html
<script src="//thbkrkr.github.io/ovh-auth/js/auth.1.min.js"></script>
<script>
$get('/api/me', function(me) {
  $('.me').html(JSON.stringify(me))
})
</script>
```
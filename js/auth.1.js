(function(){

  var self = this

  AUTH_LOADER_ELEM_ID = 'loader'
  AUTH_LOGIN_ELEM_ID = 'login'
  AUTH_LOGOUT_ELEM_ID = 'logout'

  AUTH_CREDS_URL = '/auth/credential'
  AUTH_VALIDATION_URL = '/auth/validate/'
  AUTH_LOCAL_STORAGE_KEY = 'auth'
  TOKEN_LOCAL_STORAGE_KEY = 'token'

  DEFAULT_AUTH_HOME_PATH = '/s/'
  DEFAULT_AUTH_LOGIN_PATH = '/s/login.html'

  if (typeof(AUTH_HOME_PATH) === 'undefined') {
    AUTH_HOME_PATH = DEFAULT_AUTH_HOME_PATH
  }
  if (typeof(AUTH_LOGIN_PATH) === 'undefined') {
    AUTH_LOGIN_PATH = DEFAULT_AUTH_LOGIN_PATH
  }

  window.addEventListener('load', function() {
    loaderElem = document.getElementById(AUTH_LOADER_ELEM_ID)
    loginElem = document.getElementById(AUTH_LOGIN_ELEM_ID)
    logoutElem = document.getElementById(AUTH_LOGOUT_ELEM_ID)

    // If already authenticated, go to the home
    if (localStorage.getItem(AUTH_LOCAL_STORAGE_KEY)) {
      if (window.location.pathname != AUTH_HOME_PATH) {
        window.location = AUTH_HOME_PATH
        return
      }
    }

    // Validate the token query parameter, set auth data in the local storage
    // and redirect to the home
    validateToken = false
    token = localStorage.getItem(TOKEN_LOCAL_STORAGE_KEY) || $param('token')
    if (token) {
      if (loaderElem) loaderElem.style.display = 'inline-block'
      validateToken = true

      // validate a token and store auth data in the local storage
      $get(AUTH_VALIDATION_URL + token, function(creds) {
        // ok go to the home
        localStorage.removeItem(TOKEN_LOCAL_STORAGE_KEY)
        localStorage.setItem(AUTH_LOCAL_STORAGE_KEY, creds.auth)
        window.location = AUTH_HOME_PATH
      }, function(err) {
        localStorage.removeItem(TOKEN_LOCAL_STORAGE_KEY)
        window.location = AUTH_LOGIN_PATH
      })
      return
    }

    // Install login click
    if (loginElem && !validateToken) {
      if (loaderElem) loaderElem.style.display = 'none'

      loginElem.innerHTML = 'Login'
      loginElem.style.display = 'inline-block'

      loginElem.addEventListener('click', function() {
        loginElem.style.display = 'none'
        if (loaderElem) loaderElem.style.display = 'inline-block'

        // get an url to associate a consumer key to a user
        home = window.location.origin + AUTH_LOGIN_PATH
        $get(AUTH_CREDS_URL+'?redirect='+home, function(creds) {
          localStorage.setItem(TOKEN_LOCAL_STORAGE_KEY, creds.token)
          window.location = creds.url
        }, function(err) {
          window.location = AUTH_LOGIN_PATH
        })
      })

    }

    // Install logout click
    if (logoutElem) {
      logoutElem.addEventListener('click', function() {
        localStorage.removeItem(AUTH_LOCAL_STORAGE_KEY)
        window.location = AUTH_LOGIN_PATH
      })
    }
  })

})()
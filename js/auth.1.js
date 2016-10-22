(function(){

  var self = this

  AUTH_LOGIN_ELEM_ID = 'login'
  AUTH_LOGOUT_ELEM_ID = 'logout'
  AUTH_CREDS_URL = '/auth/credential'
  AUTH_VALIDATION_URL = '/auth/validate/'
  AUTH_LOCAL_STORAGE_KEY = 'auth'

  DEFAULT_AUTH_HOME_PATH = '/s/'
  DEFAULT_AUTH_LOGIN_PATH = '/s/login.html'

  if (typeof(AUTH_HOME_PATH) === 'undefined') {
    AUTH_HOME_PATH = DEFAULT_AUTH_HOME_PATH
  }

  if (typeof(AUTH_LOGIN_PATH) === 'undefined') {
    AUTH_LOGIN_PATH = DEFAULT_AUTH_LOGIN_PATH
  }

  window.addEventListener('load', function() {
    // Ask a validation URL and redirect to it
    loginElem = document.getElementById(AUTH_LOGIN_ELEM_ID)
    if (loginElem) {
      loginElem.addEventListener('click', function() {
        loginElem.innerHTML = '...'
        $get(AUTH_CREDS_URL, function(creds) {
          window.location = creds.url
        })
      })
    }

    // Already authenticated, go to the home
    if (localStorage.getItem(AUTH_LOCAL_STORAGE_KEY)) {
      if (window.location.pathname != AUTH_HOME_PATH) {
        window.location.pathname = AUTH_HOME_PATH
      }
    } else {
      $('.container').show()
    }

    // Validate the token query parameter to set auth data in the local storage
    token = $param('token')
    if (token) {
      $get(AUTH_VALIDATION_URL + token, function(creds) {
        localStorage.setItem(AUTH_LOCAL_STORAGE_KEY, creds.auth)
        window.location = AUTH_HOME_PATH
      })
    }

    logoutElem = document.getElementById(AUTH_LOGOUT_ELEM_ID)
    if (logoutElem) {
      logoutElem.addEventListener('click', function() {
        $logout()
      })
    }

  })

  this.$logout = function $logout() {
    window.location = AUTH_LOGIN_PATH
    localStorage.removeItem(AUTH_LOCAL_STORAGE_KEY)
  }

  // --

})()
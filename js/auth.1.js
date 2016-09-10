(function(){

  var self = this

  AUTH_LOGIN_ELEM_ID = 'login'
  AUTH_LOGOUT_ELEM_ID = 'logout'
  AUTH_CREDS_URL = '/auth/credential'
  AUTH_VALIDATION_URL = '/auth/validate/'
  AUTH_LOCAL_STORAGE_KEY = 'auth'

  DEFAULT_AUTH_HOME_PATH = '/s/'
  DEFAULT_AUTH_LOGIN_PATH = '/s/login.html'
  DEFAULT_AUTH_HEADER = 'X-Auth'

  if (typeof(AUTH_HOME_PATH) === 'undefined') {
    AUTH_HOME_PATH = DEFAULT_AUTH_HOME_PATH
  }

  if (typeof(AUTH_LOGIN_PATH) === 'undefined') {
    AUTH_LOGIN_PATH = DEFAULT_AUTH_LOGIN_PATH
  }
  if (typeof(AUTH_HEADER) === 'undefined') {
    AUTH_HEADER = DEFAULT_AUTH_HEADER
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

  //@ Extract query param from the URL
  this.$param = function $param(name) {
      var url = window.location.href
      name = name.replace(/[\[\]]/g, "\\$&")
      var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
          results = regex.exec(url)
      if (!results) return null
      if (!results[2]) return ''
      return decodeURIComponent(results[2].replace(/\+/g, " "))
  }

  //@ HTTP GET
  this.$get = function $get(url, onSuccess, onFailure, noauth) {
    return $ajax(url, undefined, onSuccess, onFailure, noauth)
  }

  //@ HTTP POST
  this.$post = function $post(url, data, onSuccess, onFailure, noauth) {
    return $ajax(url, data, onSuccess, onFailure, noauth)
  }

  // $.ajax with the auth header and a redirection to the login page
  // with a  local storage clean if an HTTP call return 401.
  this.$ajax = function $ajax(url, data, onSuccess, onFailure, noauth) {
    var headers = {}
    var auth = localStorage.getItem(AUTH_LOCAL_STORAGE_KEY)
    if (auth) {
      headers[AUTH_HEADER] = auth
    }
    type = 'GET'
    if (data) {
      type = 'POST'
    }
    req = {
      type: type,
      url: url,
      contentType: "application/json",
      headers: headers,
      success: function(response) {
        if (typeof response === 'object') {
          onSuccess(response)
        } else {
          onSuccess(JSON.parse(response))
        }
      },
      error: function(xhr, errorType, error) {
        // Redirect to home if authentication error
        if (xhr.status === 401 && !noauth) {
          $logout()
        }
        if (!onFailure) {
          console.error("error on " + url + " status="+ xhr.status)
          return
        }
        onFailure(error)
      }
    }
    if (data) {
      req.data = JSON.stringify(data)
    }
    $.ajax(req)
  }

})()
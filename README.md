# [emailauth](https://emailauth.herokuapp.com) API Documentation

**API platform for authenticating users without requiring a password.**

## Installation

### Heroku

```bash
git clone https://github.com/scottmotte/emailauth.git
cd emailauth
heroku create emailauth
heroku addons:add sendgrid
git push heroku master
heroku run rake db:migrate
heroku config:set FROM=login@yourapp.com
```

Next, create your first app.

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/apps/create.json \
-d "email=you@email.com" \
-d "app_name=myappname"
```

Nice, that's all it takes to get your authentication system running. Now let's plug that into our app using the embeddable JavaScript.

Place a script tag wherever you want the login form displayed.  

```html
<script src='/path/to/emailauth.js' 
        data-app_name="your_app_name" 
        data-root_url="https://emailauth.herokuapp.com"></script>
```

Get the latest [emailauth.js here](https://github.com/scottmotte/emailauth-js/blob/master/build/emailauth.js). Replace the `data-app_name` with your own.

Next, bind to the emailauth:login_confirm event to get the successful login data. This is where you would make an internal request to your application to set the session for the user.

```html
<script>
  emailauth.script.addEventListener('emailauth:login_confirm', function(e) {
    console.log(e.data);
    $.post("/login/success", {email: e.data.identity.email}, function(data) {
      window.location.href = "/dashboard";
    });    
  }, false); 
</script>
```

## API Overview

The [emailauth.herokuapp.com](https://emailauth.herokuapp.com) API is based around REST. It uses standard HTTP authentication. [JSON](https://www.json.org/) is returned in all responses from the API, including errors.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

* [Summary](#summary)
* [Apps](#apps)
* [Logins](#logins)
* [Identity](#identities) 

## Summary

### API Endpoint

* https://emailauth.herokuapp.com/api/v0

## Apps

To start using the emailauth.io API, you must first create an app.

### POST /apps/create

Pass an email and app_name to create your app at emailauth.herokuapp.com.

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/v0/apps/create.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/apps/create.json -d "email=test@example.com" -d "app_name=myapp" -d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  app: {
    id: "APP_123453423784",
    email: "test@example.com",
    app_name: "myapp"
  }
}
```

## Logins

### POST /login/request

Request a login.

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/v0/login/request.json
```

#### Required Parameters

* email
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/login/request.json \ 
-d "email=test@example.com" \
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  login: {
    email: "test@example.com",
    app_name: "your_app_name"
  }
  identity: {
    id: "IDNT_1234348347834",  
    email: "test@example.com",
  }
}
```

### POST /login/confirm

Confirm a login. Email and authcode must match to get a success response back. 

#### Definition

```bash
POST https://emailauth.herokuapp.com/api/v0/login/confirm.json
```

#### Required Parameters

* email
* authcode
* app_name

#### Example Request

```bash
curl -X POST https://emailauth.herokuapp.com/api/v0/login/confirm.json \
-d "email=test@example.com" \
-d "authcode=7389" \ 
-d "app_name=your_app_name"
```

#### Example Response
```javascript
{
  success: true,
  login: {
    email: "test@example.com",
    app_name: "your_app_name"
  }
  identity: {
    id: "IDNT_1234348347834",  
    email: "test@example.com",
  }
}
```

## Identities

No endpoints yet. Some identities will be shown when making a `/login/confirm.json` call.

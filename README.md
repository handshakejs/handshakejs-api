# [handshakejs](https://handshakejs.herokuapp.com) API Documentation

![](https://raw.githubusercontent.com/handshakejs/handshakejs-api/sm-go/handshakejs-logo.png)

**API platform for authenticating users without requiring a password.**

## Installation

### Heroku

```bash
git clone https://github.com/scottmotte/handshakejs-api.git
cd handshakejs-api
heroku create handshakejs-api -b https://github.com/kr/heroku-buildpack-go.git
heroku addons:add sendgrid
heroku addons:add redistogo
heroku config
```

Note the REDISTOGOURL, SENGRID_PASSWORD, and SENDGRID_USERNAME.

```
heroku config:set DB_ENCRYPTION_SALT="somesecretsaltthatis32characters"
heroku config:set REDIS_URL=REDISTOGOURL
heroku config:set FROM=you@youremail.com
heroku config:set SMTP_ADDRESS=smtp.sendgrid.net 
heroku config:set SMTP_PORT=587
heroku config:set SMTP_USERNAME=SENDGRID_USERNAME
heroku config:set SMTP_PASSWORD=SENDGRID_PASSWORD
heroku config:set SUBJECT_TEMPLATE="Your code: {{authcode}}. Please enter it to login."
heroku config:set TEXT_TEMPLATE="Your code: {{authcode}}. Please enter it to login."
heroku config:set HTML_TEMPLATE="Your code: <b>{{authcode}}</b>. Please enter it to login."
heroku config
```

Finally, deploy it.

```
git push heroku master
```

Next, create your first app. Replace email and app_name with your information. 

<https://handshakejs-api.herokuapp.com/api/v1/apps/create.json?email=[email]&app_name=[app_name]>

Nice, that's all it takes to get your authentication system running. Now let's plug that into our app using the embeddable JavaScript.

### Old, this will move to a separate spot probably with instructions.

Place a script tag wherever you want the login form displayed.  

```html
<script src='/path/to/handshake.js' 
        data-app_name="your_app_name" 
        data-root_url="https://handshakejs-api.herokuapp.com"></script>
```

Get the latest [handshake.js here](https://github.com/sendgrid/handshakejs-script/blob/master/build/handshake.js). Replace the `data-app_name` with your own.

Next, bind to the handshake:login_confirm event to get the successful login data. This is where you would make an internal request to your application to set the session for the user.

```html
<script>
  handshake.script.addEventListener('handshake:login_confirm', function(e) {
    console.log(e.data);
    $.post("/login/success", {email: e.data.identity.email, hash: e.data.identity.hash}, function(data) {
      window.location.href = "/dashboard";
    });    
  }, false); 
</script>
```

Then you'd setup a route in your app at /login/success to do something like this (setting the session). Here's an example in ruby and there is also a [full example ruby app](https://github.com/handshakejs/handshakejs-example-ruby).

```ruby
  post "/login/success" do
    salt    = "the_secret_salt_when_you_created_an_app_that_only_you_should_know"
    pbkdf2  = PBKDF2.new(:password=>params[:email], :salt=>salt, :iterations=>1000, :key_length => 16, :hash_function => "sha1")

    session[:user] = params[:email] if pbkdf2.hex_string == params[:hash]
    redirect "/dashboard"
  end
```

## API Overview

The [handshakejs.herokuapp.com](https://handshakejs.herokuapp.com) API is based around REST. It uses standard HTTP authentication. [JSON](https://www.json.org/) is returned in all responses from the API, including errors.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

* [Summary](#summary)
* [Apps](#apps)
* [Login](#login)

## Summary

### API Endpoint

* https://handshakejs-api.herokuapp.com/api/v1

## Apps

To start using the handshake API, you must first create an app.

### POST /apps/create

Pass an email and app_name to create your app at handshakejs.herokuapp.com.

#### Definition

```bash
ANY https://handshakejs-api.herokuapp.com/api/v1/apps/create.json?app_name=[app_name]&email=[email]&salt=[salt]
```

#### Required Parameters

* app_name

#### Optional Parameters

* email
* salt

#### Example Request

<https://handshakejs-api.herokuapp.com/api/v1/apps/create.json?app_name=[app_name]&email=[email]&salt=[salt]>

#### Example Response
```javascript
{
  "apps": [{
    "email": "test@example.com",
    "app_name": "myapp",
    "salt": "the_default_generated_salt_that_you_should_keep_secret"
  }]
}
```

#### Example Error
```javascript
{
  errors: [{
    "code": "not_unique",
    "field": "app_name",
    "message": "app_name must be unique"
  }]
}
```

## Logins

### POST /login/request

Request a login.

#### Definition

```bash
ANY https://handshakejs-api.herokuapp.com/api/v0/login/request.json?email=[email]&app_name=[app_name]
```

#### Required Parameters

* email
* app_name

#### Example Request

<https://handshakejs-api.herokuapp.com/api/v0/login/request.json?email=[email]&app_name=[app_name]>

#### Example Response
```javascript
{
  "identities": [{
    "email": "test@example.com",
    "app_name": "your_app_name",
    "authcode_expired_at": "1382833591309"
  }]
}
```

#### Example Error
```javascript
{
  "errors": [{
    "code": "required",
    "field": "email",
    "message": "email cannot be blank"
  }]
}
```

### POST /login/confirm

Confirm a login. Email and authcode must match to get a success response back. 

#### Definition

```bash
ANY https://handshakejs-api.herokuapp.com/api/v1/login/confirm.json?email=[email]&authcode=[authcode]&app_name=[app_name]
```

#### Required Parameters

* email
* authcode
* app_name

#### Example Request

<https://handshakejs-api.herokuapp.com/api/v1/login/confirm.json?email=[email]&authcode=[authcode]&app_name=[app_name]>

#### Example Response
```javascript
{
  "identities": [{
    "email": "test@example.com",
    "app_name": "your_app_name",
    "hash": "523f325279fd3446a78894b55cf4d777"
  }]
}
```

#### Example Error
```javascript
{
  "errors": [{
    "code": "incorrect",
    "field": "authcode",
    "message": "the authcode was incorrect"
  }]
}
```

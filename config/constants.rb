RACK_ENV                  = ENV['RACK_ENV']
DATABASE_URL              = ENV['DATABASE_URL'] 
MAX_POOL_HEROKU_DEV_LIMIT = 500
DEFAULT_COUNT             = 10
MAX_COUNT                 = 100
DEFAULT_OFFSET            = 0
FROM                      = ENV['FROM'] || "login@emailauth.io"
SMTP_ADDRESS              = ENV['SMTP_ADDRESS'] || "smtp.sendgrid.net"
SMTP_PORT                 = ENV['SMTP_PORT'] || 25
SMTP_USERNAME             = ENV['SMTP_USERNAME'] || ENV['SENDGRID_USERNAME']
SMTP_PASSWORD             = ENV['SMTP_PASSWORD'] || ENV['SENDGRID_PASSWORD']

case RACK_ENV
when "production"
when "test"
else  
  # defined above 
end

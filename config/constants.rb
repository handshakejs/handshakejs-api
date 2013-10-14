RACK_ENV                  = ENV['RACK_ENV']
DATABASE_URL              = ENV['DATABASE_URL'] 
MAX_POOL_HEROKU_DEV_LIMIT = 500
DEFAULT_COUNT             = 10
MAX_COUNT                 = 100
DEFAULT_OFFSET            = 0

case RACK_ENV
when "production"
when "test"
else  
  # defined above 
end

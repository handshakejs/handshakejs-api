class Application < Sinatra::Base
  register Sinatra::ActiveRecordExtension
  
  enable :raise_errors
  #disable :show_exceptions
  #disable :raise_errors

  configure do
    set :database, DATABASE_URL
    ActiveRecord::Base.logger = nil unless RACK_ENV == "development"
  end

  configure :development do
    register Sinatra::Reloader
  end

  helpers do
    def jsonify(*args)
      render(:jsonify, *args)
    end

    def req_basic_auth
      header = request.env["HTTP_AUTHORIZATION"]
      return nil unless header
      token   = header.split(/\s+/).pop()
      auth    = Base64.decode64 token
      auth.split(/:/)[0]
    end
  end

  before do
    response['Access-Control-Allow-Origin'] = "*"
  end
  
  get "/" do
    redirect "/api/v0"
  end

  get "/exception" do
    raise "Exception"
  end
end

class ApiV0 < Application
  use Rack::PostBodyContentTypeParser

  layout false

  before do
    headers['X-Frame-Options']  = 'GOFORIT'
    response['Cache-Control']   = "private, no-cache=true, no-store=true, max-age=0"
    content_type :json
  end

  def set_count_and_offset
    @count  = params[:count] ||= DEFAULT_COUNT
    @count  = @count.to_i
    @count  = [@count, MAX_COUNT].min
    @offset = params[:offset] ||= DEFAULT_OFFSET
    @offset = [0, @offset.to_i].max
  end

  def return_order(order)
    case order.try(:upcase)
    when "ASC"
      "ASC"
    else
      "DESC"
    end
  end

  not_found do
    @message = "API endpoint not found."
    jsonify :"/api_v0/error"
  end

  get "/" do
    jsonify :"/api_v0/index"
  end

  post "/apps/create.json" do
    params_to_save = { 
      email:      params[:email],
      app_name:   params[:app_name]
    }

    @app = App.new(params_to_save)

    if @app.save
      jsonify :"/api_v0/apps/create"
    else
      @message = @app.errors.full_messages.to_sentence
      jsonify :"/api_v0/error"
    end
  end

  post "/login/request.json" do
    @app = App.where(app_name: params[:app_name]).first

    if !!@app
      params_to_save = { email: params[:email] }
      @login  = @app.logins.build(params_to_save)

      if @login.save
        # SEND THE EMAIL HERE WITH THE AUTHCODE!!!!
        jsonify :"/api_v0/login/request"
      else
        @message = @login.errors.full_messages.to_sentence
        jsonify :"/api_v0/error"
      end
    else
      @message = "Sorry, we couldn't find an app by that app_name."
      jsonify :"/api_v0/error"
    end
  end

  post "/login/confirm.json" do
    params_to_confirm = { 
      email:    params[:email], 
      authcode: params[:authcode], 
      app_name: params[:app_name] 
    }

    @login_confirm = Login.confirm(params_to_confirm)

    if @login_confirm[0]
      @login = @login_confirm[0]
      jsonify :"/api_v0/login/confirm"
    else
      @message = @login_confirm[1] 
      jsonify :"/api_v0/error"
    end
  end
end

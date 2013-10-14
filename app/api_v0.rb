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
    jsonify :"/api/error"
  end

  post "/apps.json" do
    params_to_save = { 
      email:      params[:email],
      app_name:   params[:app_name]
    }

    @app = App.new(params_to_save)

    if @app.save
      jsonify :"/api/apps/create"
    else
      @message = @app.errors.full_messages.to_sentence
      jsonify :"/api/error"
    end
  end

  








  # crusty from copying in
  post "/sessions.json" do
    if params[:email].blank?
      @message = "Please provide an email address."
      return jsonify :"/api_v0/error"
    elsif params[:password].blank?
      @message = "Please provide a password."
      return jsonify :"/api_v0/error"
    end

    @person = Person.authenticate(params[:email], params[:password])
    if @person
      jsonify :"/api_v0/sessions/create"
    else
      @message = "The account information provided is incorrect."
      jsonify :"/api_v0/error"
    end
  end

  post "/documents.json" do
    require_api_person!

    params_to_save = { 
      url:              params[:url],
      output_formats:   params[:output_formats].try(:upcase).try(:split, ",")
    }

    @document = current_api_person.documents.build(params_to_save)

    if @document.save
      if QUEUE_INLINE
        Document.process!(@document.id)
      else
        QC.enqueue "Document.process!", @document.id
      end

      jsonify :"/api_v0/documents/create"
    else
      @message = @document.errors.full_messages.to_sentence
      jsonify :"/api_v0/error"
    end
  end

  get "/documents/:id.json" do
    require_api_person!

    @document = current_api_person.documents.where(short_id: params[:id]).first

    if !@document
      @message = "Document not found."
      return jsonify :"/api_v0/error"
    end

    jsonify :"/api_v0/documents/show"
  end

  get "/documents/:id/pages.json" do
    require_api_person!

    @document = current_api_person.documents.where(short_id: params[:id]).first

    if !@document
      @message = "Document not found."
      return jsonify :"/api_v0/error"
    end

    @pages          = @document.pages.processed
    @total_count    = @pages.count
    @count          = @total_count
    @offset         = 0

    jsonify :"/api_v0/documents/pages/index"
  end

  get "/events.json" do
    require_api_person!
    set_count_and_offset
    @events       = current_api_person.events.order("created_at #{return_order(params[:order])}").limit(@count).offset(@offset)
    @total_count  = @events

    if params[:type]
      type_strings  = [params[:type]]
      @total_count  = @total_count.for_type_strings(type_strings)
      @events       = @events.for_type_strings(type_strings)
    end

    @total_count  = @total_count.count
    @events       = @events

    jsonify :"/api_v0/events/index"
  end

  post "/webhooks.json" do
    require_api_person!

    params_to_save = { 
      url:              params[:url]
    }

    @webhook = current_api_person.webhooks.build(params_to_save)

    if @webhook.save
      jsonify :"/api_v0/webhooks/create"
    else
      @message = @webhook.errors.full_messages.to_sentence
      jsonify :"/api_v0/error"
    end
  end

  get "/webhooks.json" do
    require_api_person!

    @webhooks = current_api_person.webhooks

    jsonify :"/api_v0/webhooks/index"
  end

  post "/webhooks/:id/delete.json" do
    require_api_person!

    @webhook = current_api_person.webhooks.where(short_id: params[:id]).first
    
    if @webhook && @webhook.destroy
      jsonify :"/api_v0/webhooks/delete"
    else
      @message = "Webhook not found."
      jsonify :"/api_v0/error"
    end
  end
end

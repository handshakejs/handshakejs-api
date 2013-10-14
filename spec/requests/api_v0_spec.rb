require 'spec_helper'

describe "/api/v0" do
  include Rack::Test::Methods

  def app
    ApiV0
  end

  let(:body)  { last_response.body }
  let(:json)  { JSON.parse(body) }

  context "not requiring authorization" do
    describe "POST /apps" do
      let(:email)     { "scott@mailinator.com" }
      let(:app_name)  { "mailinator" }
      let(:params)    { {email: email, app_name: app_name} }

      before do
        post "/apps.json", params, format: 'json'
      end

      context "valid email and password" do
        it do
          json['success'].should eq true
          json['app']['email'].should eq email
          json['app']['app_name'].should eq app_name
        end
      end

      context "missing email" do
        let(:email)   { " " }

        it do
          json['success'].should eq false
          json['error']['message'].should_not be_blank 
        end
      end

      context "missing app_name" do
        let(:app_name)   { "" }

        it do
          json['success'].should eq false
          json['error']['message'].should_not be_blank
        end
      end
    end
  end

  # context "requiring authorization" do
  #   let(:email)                 { "scott@mailinator.com" }
  #   let(:password)              { "password" }
  #   let(:output_formats)        { nil }
  #   let!(:person)               { Person.create(email: email, password: password) }
  #   let(:secret_api_key)        { person.secret_api_key }
  #   let(:url)                   { "http://scottmotte.com/assets/resume.pdf" }

  #   before { basic_authorize(secret_api_key, nil) }

  #   describe "POST /documents" do
  #     let(:params) { {url: url, output_formats: output_formats} }

  #     before do
  #       post "/documents.json", params, format: 'json'
  #     end

  #     it do
  #       body['success'].should eq true
  #       body['document']['url'].should eq url
  #       Document.count.should eq 1
  #     end

  #     context "secret_api_key incorrect" do
  #       let(:secret_api_key) { "invalid" }

  #       it { body['success'].should eq false }
  #     end

  #     context "missing url" do
  #       let(:url) { "" }

  #       it do 
  #         body['success'].should eq false
  #         body['error']['message'].should eq "Url can't be blank"
  #       end
  #     end

  #     context "not real url" do
  #       let(:url) { "curl -i -X HEAD http://scottmotte.com/assets/noexiste.pdf" }

  #       it do
  #         body['success'].should eq false
  #         body['error']['message'].should eq "Url is invalid"
  #       end
  #     end

  #     context "output_format is pdf" do
  #       let(:output_formats)         { "pdf" }

  #       it do
  #         body['success'].should eq true
  #         @document_id  = body['document']['id']
  #         response2     = get "/documents/#{@document_id}.json", {}, format: 'json'
  #         body2         = JSON.parse(response2.body)
  #         body2['document']['pdf'].should_not be_blank
  #       end
  #     end
  #   end

  #   describe "GET /documents/:id" do
  #     let!(:document) { FactoryGirl.build(:document, person: person) }

  #     context "does not yet exist" do
  #       before do
  #         get "/documents/1.json", {}, format: 'json'
  #       end

  #       it { body['success'].should eq false }
  #     end

  #     context "default" do
  #       before do
  #         document.save!
  #         get "/documents/#{document.short_id}.json", format: 'json'
  #       end

  #       it do
  #         body['success'].should eq true
  #         body['document']['id'].should eq document.short_id
  #       end
  #     end

  #     context "exists but belongs to a different person" do
  #       let(:person2)     { FactoryGirl.build(:person, email: "differentperson@mailinator.com") }

  #       before do
  #         document.save
  #         person2.save
  #         basic_authorize(person2.secret_api_key, nil)
  #         get "/documents/#{document.short_id}.json", {}, format: 'json'
  #       end

  #       it { body['success'].should eq false }
  #     end
  #   end

  #   describe "GET /documents/:id/pages" do
  #     let!(:document) { FactoryGirl.build(:document, person: person) }

  #     context "default" do
  #       before do
  #         document.save!
  #         document.process!
  #         get "/documents/#{document.short_id}/pages.json", format: 'json'
  #       end

  #       it do
  #         body['success'].should eq true
  #         body['pages'].count.should >= 1
  #       end
  #     end
  #   end

  #   describe "POST /webhooks" do
  #     before { post "/webhooks.json", {url: "http://webhookurl.com"}, format: 'json' }

  #     it do
  #       Webhook.count.should eq 1
  #       body['success'].should eq true
  #       body['webhook']['id'].should_not be_blank
  #     end
  #   end

  #   describe "GET /webhooks" do
  #     let(:webhook) { FactoryGirl.create(:webhook, person: person) }

  #     before do
  #       webhook.save
  #       get "/webhooks.json", {}, format: 'json'
  #     end 

  #     it do
  #       body['success'].should eq true
  #       body['webhooks'].count.should eq 1
  #     end
  #   end

  #   describe "POST /webhooks/:id/delete" do
  #     let(:webhook) { FactoryGirl.create(:webhook, person: person) }

  #     before do
  #       webhook.save!
  #       post "/webhooks/#{webhook.short_id}/delete.json", {}, format: 'json'
  #     end  

  #     it do
  #       Webhook.count.should eq 0
  #       body['success'].should eq true
  #     end
  #   end

  #   describe "GET /events" do
  #     let(:document)  { FactoryGirl.create(:document, person: person) }
  #     let(:event)     { Event.last }

  #     before { document.save! }
  #     
  #     context "default" do
  #       before { get "/events.json", {}, format: 'json' }

  #       it do
  #         Event.count.should >= 1
  #         body['success'].should eq true
  #         body['events'].count.should >= 1
  #         body['events'][0]['id'].should eq event.short_id
  #       end
  #     end

  #     context "limit to 1" do
  #       before { get "/events.json?count=1", {}, format: 'json' }

  #       it { body['events'].count.should eq 1 }
  #     end

  #     context "filter by event_type" do
  #       before { get "/events.json?type=document.created" }

  #       it { body['events'].count.should eq 1 }
  #     end
  #   end
  # end
end

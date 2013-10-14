require 'spec_helper' 

describe App do
  let(:app) { FactoryGirl.build(:app) }

  it { app.should be_valid }

  context "missing email" do
    before { app.email = "" }

    it { app.should_not be_valid }
  end

  context "missing app_name" do
    before { app.app_name = "" }

    it { app.should_not be_valid }
  end

  context "app_name already exists" do
    let(:app2) { FactoryGirl.build(:app) }

    before do
      app.save!
      app2.app_name = app.app_name
    end

    it { app2.should_not be_valid }
  end

  context "email already exists" do
    let(:app2) { FactoryGirl.build(:app) }

    before do
      app.save!
      app2.email = app.email
    end

    it { app2.should_not be_valid }
  end

  context "successful save" do
    before { app.save! }

    it { app.secret_api_key.should_not be_blank }
    it { app.public_api_key.should_not be_blank }
  end
end

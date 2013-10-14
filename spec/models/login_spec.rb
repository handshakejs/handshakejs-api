require 'spec_helper' 

describe Login do
  let(:login) { FactoryGirl.build(:login) }

  it { login.should be_valid }

  context "missing email" do
    before { login.email = "" }

    it { login.should_not be_valid }
  end

  context "successful save" do
    before { login.save! }

    it { login.authcode.should_not be_blank }
    it { login.authcode.length.should eq 4 }
    it { login.should be_requested }
    it { login.should_not be_confirmed }
    it { login.identity.should_not be_blank }
  end

  describe "#mark_confirmed!" do
    let!(:login) { FactoryGirl.create(:login) }

    before { login.mark_confirmed! }

    it { login.should be_confirmed }
    it { login.should_not be_requested }
  end

  describe ".confirm" do
    let!(:login)    { FactoryGirl.create(:login) }
    let(:email)     { login.email }
    let(:authcode)  { login.authcode }
    let(:app_name)  { login.app.app_name }
    let(:params)    { {email: email, authcode: authcode, app_name: app_name } }
    let(:confirm)   { Login.confirm(params) }

    context "valid" do
      before { confirm }

      it { confirm[0].should be_true }
    end
    
    context "incorrect authcode" do
      let(:authcode) { "incorrect" }
      
      before { confirm }
      
      it { confirm[0].should_not be_true }
    end

    context "incorrect email" do
      let(:email) { "incorrect" }
      
      before { confirm }
      
      it { confirm[0].should_not be_true }
    end

    context "incorrect app_name" do
      let(:app_name) { "incorrect" }
      
      before { confirm }
      
      it { confirm[0].should_not be_true }
    end
  end
end

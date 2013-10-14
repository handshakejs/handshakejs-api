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
  end

  describe "#mark_confirmed!" do
    let!(:login) { FactoryGirl.create(:login) }

    before { login.mark_confirmed! }

    it { login.should be_confirmed }
    it { login.should_not be_requested }
  end
end

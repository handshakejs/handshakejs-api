require 'spec_helper' 

describe Identity do
  let(:identity) { FactoryGirl.build(:identity) }
  
  it { identity.should be_valid }

  context "missing email" do
    before { identity.email = "" }

    it { identity.should_not be_valid }
  end

  context "email already exists" do
    let(:identity2) { FactoryGirl.build(:identity, app: identity.app) }

    before do
      identity.save!
      identity2.email = identity.email 
    end

    it { identity2.should_not be_valid }
  end

  context "successful save" do
    before do
      identity.save!
    end

    it { identity.short_id.should_not be_blank }
  end

end

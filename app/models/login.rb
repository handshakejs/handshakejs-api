class Login < ActiveRecord::Base
  belongs_to :app

  validates :email,    presence: true

  before_save :generate_authcode

  state_machine :status, :initial => :requested do
    event :mark_confirmed! do
      transition any => :confirmed
    end
  end

  private

  def generate_authcode
    self.authcode = random_authcode 
  end

  def random_authcode
    '%04d' % rand(10 ** 4)
  end
end
